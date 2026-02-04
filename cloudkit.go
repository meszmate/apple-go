package apple

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	cloudKitBaseURL = "https://api.apple-cloudkit.com"
	cloudKitVersion = "1"
)

// CloudKit is the interface for the CloudKit server-to-server API.
type CloudKit interface {
	// Records
	QueryRecords(db CKDatabase, req *CKQueryRequest) (*CKQueryResponse, error)
	ModifyRecords(db CKDatabase, req *CKRecordsModifyRequest) (*CKRecordsModifyResponse, error)
	LookupRecords(db CKDatabase, req *CKRecordsLookupRequest) (*CKRecordsLookupResponse, error)
	RecordChanges(db CKDatabase, req *CKRecordChangesRequest) (*CKRecordChangesResponse, error)
	// Zones
	ListZones(db CKDatabase) (*CKZonesResponse, error)
	LookupZones(db CKDatabase, zoneIDs []CKZoneID) (*CKZonesResponse, error)
	ModifyZones(db CKDatabase, zones []CKZone, operationType CKOperationType) (*CKZonesResponse, error)
	ZoneChanges(db CKDatabase, req *CKZoneChangesRequest) (*CKZoneChangesResponse, error)
	// Subscriptions
	ListSubscriptions(db CKDatabase) (*CKSubscriptionsResponse, error)
	ModifySubscriptions(db CKDatabase, req *CKSubscriptionsModifyRequest) (*CKSubscriptionsModifyResponse, error)
	// Assets
	UploadAssets(db CKDatabase, req *CKAssetsUploadRequest) (*CKAssetsUploadResponse, error)
	// Users
	GetCurrentUser(db CKDatabase) (*CKUserInfo, error)
	DiscoverAllUsers(db CKDatabase) (*CKUsersResponse, error)
	LookupUsers(db CKDatabase, req *CKUserLookupRequest) (*CKUsersResponse, error)
	// Tokens
	CreateTokens(db CKDatabase, req *CKTokensCreateRequest) (*CKTokensCreateResponse, error)
}

type ckHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type cloudKit struct {
	KeyID       string
	Container   string
	Environment CKEnvironment
	KeyContent  []byte
	httpClient  ckHTTPClient
}

// NewCloudKit creates a new CloudKit client using a key file path.
func NewCloudKit(keyID, container string, environment CKEnvironment, keyPath string) (CloudKit, error) {
	keyContent, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return &cloudKit{
		KeyID:       keyID,
		Container:   container,
		Environment: environment,
		KeyContent:  keyContent,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// NewCloudKitB64 creates a new CloudKit client using a base64-encoded key.
func NewCloudKitB64(keyID, container string, environment CKEnvironment, b64Key string) (CloudKit, error) {
	keyContent, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return nil, err
	}
	return &cloudKit{
		KeyID:       keyID,
		Container:   container,
		Environment: environment,
		KeyContent:  keyContent,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// buildURL constructs the full URL and returns the subpath used for signing.
func (c *cloudKit) buildURL(db CKDatabase, subpath string) (fullURL, urlSubpath string) {
	urlSubpath = fmt.Sprintf("/database/%s/%s/%s/%s/%s",
		cloudKitVersion, c.Container, string(c.Environment), string(db), subpath)
	fullURL = cloudKitBaseURL + urlSubpath
	return fullURL, urlSubpath
}

// sign creates the CloudKit ECDSA signature for a request.
func (c *cloudKit) sign(body []byte, subpath string) (date, signature string, err error) {
	privateKey, err := parseECPrivateKey(c.KeyContent)
	if err != nil {
		return "", "", err
	}

	date = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	bodyHash := sha256.Sum256(body)
	bodyHashB64 := base64.StdEncoding.EncodeToString(bodyHash[:])

	message := fmt.Sprintf("%s:%s:%s", date, bodyHashB64, subpath)
	messageHash := sha256.Sum256([]byte(message))

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, messageHash[:])
	if err != nil {
		return "", "", err
	}

	signature = base64.StdEncoding.EncodeToString(sig)
	return date, signature, nil
}

// doRequest executes a signed HTTP request to the CloudKit API.
func (c *cloudKit) doRequest(method string, db CKDatabase, subpath string, body, result any) error {
	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return err
		}
	} else {
		bodyBytes = []byte("{}")
	}

	fullURL, urlSubpath := c.buildURL(db, subpath)

	date, signature, err := c.sign(bodyBytes, urlSubpath)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Apple-CloudKit-Request-KeyID", c.KeyID)
	req.Header.Set("X-Apple-CloudKit-Request-ISO8601Date", date)
	req.Header.Set("X-Apple-CloudKit-Request-SignatureV1", signature)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &CKError{Code: CKErrorNetworkError, Reason: err.Error()}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &CKError{Code: CKErrorNetworkError, Reason: err.Error()}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp CKErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errResp); jsonErr == nil && errResp.ServerErrorCode != "" {
			return &CKError{
				Code:            errResp.ServerErrorCode,
				Reason:          errResp.Reason,
				ServerErrorCode: string(errResp.ServerErrorCode),
				RetryAfter:      errResp.RetryAfter,
				UUID:            errResp.UUID,
			}
		}
		return &CKError{
			Code:   CKErrorUnknownError,
			Reason: fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
		}
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return err
		}
	}

	return nil
}
