package apple

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockedASHTTPClient implements asHTTPClient for testing.
type MockedASHTTPClient struct {
	mock.Mock
}

func (m *MockedASHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	resArg := args.Get(0)
	resp, ok := resArg.(*http.Response)
	if !ok {
		return nil, errors.New("first parameter should be of type *http.Response")
	}
	return resp, args.Error(1)
}

func newTestAppStoreServer(client *MockedASHTTPClient) *appStoreServer {
	return &appStoreServer{
		issuerID:     "test-issuer",
		keyID:        "test-key-id",
		bundleID:     "com.example.app",
		keyContent:   []byte(testECPrivateKey),
		baseURL:      "https://api.storekit-sandbox.itunes.apple.com",
		httpClient:   client,
		rootCertPool: appleRootCertPool(),
	}
}

func TestNewAppStoreServerAPIB64(t *testing.T) {
	// Valid base64 key (won't actually work for signing, but tests constructor)
	_, err := NewAppStoreServerAPIB64("issuer", "kid", "com.example", "aW52YWxpZA==", true)
	// This will succeed for construction (key parse happens at sign time)
	assert.NoError(t, err)
}

func TestNewAppStoreServerAPIB64_InvalidBase64(t *testing.T) {
	_, err := NewAppStoreServerAPIB64("issuer", "kid", "com.example", "not-valid-base64!!!", true)
	assert.Error(t, err)
}

func TestAppStoreServerBaseURLs(t *testing.T) {
	s := newAppStoreServer("issuer", "kid", "com.example", []byte("key"), false)
	assert.Equal(t, asProductionBaseURL, s.baseURL)

	s = newAppStoreServer("issuer", "kid", "com.example", []byte("key"), true)
	assert.Equal(t, asSandboxBaseURL, s.baseURL)
}

func TestGenerateToken(t *testing.T) {
	s := &appStoreServer{
		issuerID:   "test-issuer",
		keyID:      "test-key-id",
		bundleID:   "com.example.app",
		keyContent: []byte(testECPrivateKey),
	}

	token, err := s.generateToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Token should have 3 parts
	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3)
}

func TestGenerateTokenInvalidKey(t *testing.T) {
	s := &appStoreServer{
		issuerID:   "test-issuer",
		keyID:      "test-key-id",
		bundleID:   "com.example.app",
		keyContent: []byte("not-a-valid-key"),
	}

	_, err := s.generateToken()
	assert.Error(t, err)
}

func TestASDoRequestSuccess(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"testNotificationToken":"abc123"}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/notifications/test" &&
			req.Method == "POST" &&
			strings.HasPrefix(req.Header.Get("Authorization"), "Bearer ")
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	var result ASTestNotificationResponse
	err := s.doRequest("POST", "/inApps/v1/notifications/test", nil, nil, &result)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", result.TestNotificationToken)
}

func TestASDoRequestAPIError(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	errResp := ASAPIError{
		ErrorCode:    ASAPIErrorTransactionNotFound,
		ErrorMessage: "Transaction not found",
	}
	errBytes, _ := json.Marshal(errResp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{
			StatusCode: 404,
			Header:     http.Header{},
			Body:       io.NopCloser(bytes.NewReader(errBytes)),
		},
		nil,
	)

	var result ASTransactionInfoResponse
	err := s.doRequest("GET", "/inApps/v1/transactions/12345", nil, nil, &result)
	assert.Error(t, err)

	var apiErr *ASAPIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, ASAPIErrorTransactionNotFound, apiErr.ErrorCode)
	assert.Equal(t, "Transaction not found", apiErr.ErrorMessage)
	assert.Equal(t, 404, apiErr.HTTPStatus)
}

func TestASDoRequestRateLimitError(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	errResp := ASAPIError{
		ErrorCode:    ASAPIErrorRateLimitExceeded,
		ErrorMessage: "Rate limit exceeded",
	}
	errBytes, _ := json.Marshal(errResp)

	header := http.Header{}
	header.Set("Retry-After", "30")

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{
			StatusCode: 429,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader(errBytes)),
		},
		nil,
	)

	var result ASTransactionInfoResponse
	err := s.doRequest("GET", "/inApps/v1/transactions/12345", nil, nil, &result)
	assert.Error(t, err)

	var apiErr *ASAPIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, ASAPIErrorRateLimitExceeded, apiErr.ErrorCode)
	assert.Equal(t, 30, apiErr.RetryAfter)
}

func TestASDoRequestNetworkError(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		(*http.Response)(nil),
		errors.New("connection refused"),
	)

	var result ASTransactionInfoResponse
	err := s.doRequest("GET", "/inApps/v1/transactions/12345", nil, nil, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}

func TestASDoRequestUnknownErrorResponse(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{
			StatusCode: 500,
			Header:     http.Header{},
			Body:       io.NopCloser(bytes.NewReader([]byte("not json"))),
		},
		nil,
	)

	var result ASTransactionInfoResponse
	err := s.doRequest("GET", "/inApps/v1/transactions/12345", nil, nil, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}

func TestASAPIError_Error(t *testing.T) {
	t.Run("with message", func(t *testing.T) {
		err := &ASAPIError{ErrorCode: ASAPIErrorTransactionNotFound, ErrorMessage: "not found"}
		assert.Equal(t, "appstore api: 4040010: not found", err.Error())
	})

	t.Run("without message", func(t *testing.T) {
		err := &ASAPIError{ErrorCode: ASAPIErrorGeneralInternal}
		assert.Equal(t, "appstore api: 5000000", err.Error())
	})
}
