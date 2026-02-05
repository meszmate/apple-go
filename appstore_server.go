package apple

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	asProductionBaseURL = "https://api.storekit.itunes.apple.com"
	asSandboxBaseURL    = "https://api.storekit-sandbox.itunes.apple.com"
)

// AppStoreServerAPI provides methods for communicating with the App Store Server API v2.
type AppStoreServerAPI interface {
	// Transactions
	GetTransactionInfo(transactionID string) (*ASTransactionInfoResponse, error)
	GetTransactionHistory(originalTransactionID string, params *ASTransactionHistoryParams) (*ASTransactionHistoryResponse, error)

	// Subscriptions
	GetAllSubscriptionStatuses(originalTransactionID string) (*ASSubscriptionStatusesResponse, error)
	ExtendSubscription(originalTransactionID string, req *ASExtendSubscriptionRequest) (*ASExtendSubscriptionResponse, error)
	MassExtendSubscriptions(req *ASMassExtendRequest) (*ASMassExtendResponse, error)
	GetExtensionStatus(productID, requestIdentifier string) (*ASExtensionStatusResponse, error)

	// Order / Refunds
	LookUpOrderID(orderID string) (*ASOrderLookupResponse, error)
	GetRefundHistory(transactionID string, revision string) (*ASRefundHistoryResponse, error)

	// Consumption
	SendConsumptionInfo(originalTransactionID string, req *ASConsumptionRequest) error

	// Notifications
	RequestTestNotification() (*ASTestNotificationResponse, error)
	GetTestNotificationStatus(testNotificationToken string) (*ASTestNotificationStatusResponse, error)
	GetNotificationHistory(req *ASNotificationHistoryRequest) (*ASNotificationHistoryResponse, error)
}

type asHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type appStoreServer struct {
	issuerID     string
	keyID        string
	bundleID     string
	keyContent   []byte
	baseURL      string
	httpClient   asHTTPClient
	rootCertPool *x509.CertPool
}

// NewAppStoreServerAPI creates a new App Store Server API client using a key file path.
func NewAppStoreServerAPI(issuerID, keyID, bundleID, keyPath string, sandbox bool) (AppStoreServerAPI, error) {
	keyContent, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return newAppStoreServer(issuerID, keyID, bundleID, keyContent, sandbox), nil
}

// NewAppStoreServerAPIB64 creates a new App Store Server API client using a base64-encoded key.
func NewAppStoreServerAPIB64(issuerID, keyID, bundleID, b64Key string, sandbox bool) (AppStoreServerAPI, error) {
	keyContent, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return nil, err
	}
	return newAppStoreServer(issuerID, keyID, bundleID, keyContent, sandbox), nil
}

func newAppStoreServer(issuerID, keyID, bundleID string, keyContent []byte, sandbox bool) *appStoreServer {
	baseURL := asProductionBaseURL
	if sandbox {
		baseURL = asSandboxBaseURL
	}
	return &appStoreServer{
		issuerID:     issuerID,
		keyID:        keyID,
		bundleID:     bundleID,
		keyContent:   keyContent,
		baseURL:      baseURL,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		rootCertPool: appleRootCertPool(),
	}
}

type appStoreServerClaims struct {
	jwt.StandardClaims
	BundleID string `json:"bid"`
}

func (s *appStoreServer) generateToken() (string, error) {
	privateKey, err := parseECPrivateKey(s.keyContent)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := appStoreServerClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    s.issuerID,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(15 * time.Minute).Unix(),
			Audience:  "appstoreconnect-v1",
		},
		BundleID: s.bundleID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, &claims)
	token.Header["kid"] = s.keyID
	token.Header["typ"] = "JWT"

	return token.SignedString(privateKey)
}

func (s *appStoreServer) doRequest(method, path string, queryParams url.Values, body, result any) error {
	token, err := s.generateToken()
	if err != nil {
		return err
	}

	fullURL := s.baseURL + path
	if queryParams != nil {
		encoded := queryParams.Encode()
		if encoded != "" {
			fullURL += "?" + encoded
		}
	}

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr ASAPIError
		if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil && apiErr.ErrorCode != 0 {
			apiErr.HTTPStatus = resp.StatusCode
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if v, err := strconv.Atoi(ra); err == nil {
					apiErr.RetryAfter = v
				}
			}
			return &apiErr
		}
		return fmt.Errorf("appstore api: unexpected status code: %d", resp.StatusCode)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return err
		}
	}

	return nil
}
