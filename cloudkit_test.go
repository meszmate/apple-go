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

const testECPrivateKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQggrixEnG1iCk5UR5i
g3Rf1492YSmMMFNv2qYiKCEkOVahRANCAAQYi8kq0gamP3++N6oUAm7AkNeqnYzR
2gmKCXQTZEzCtJMQPeXqt04OnY+XBEm4kvJHdQx37cETW21xSfoeE1jr
-----END PRIVATE KEY-----`

// MockedCKHTTPClient implements ckHTTPClient for testing.
type MockedCKHTTPClient struct {
	mock.Mock
}

func (m *MockedCKHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	resArg := args.Get(0)
	resp, ok := resArg.(*http.Response)
	if !ok {
		return nil, errors.New("first parameter should be of type *http.Response")
	}
	return resp, args.Error(1)
}

func newTestCloudKit(client *MockedCKHTTPClient) *cloudKit {
	return &cloudKit{
		KeyID:       "test-key-id",
		Container:   "iCloud.com.example.app",
		Environment: CKEnvironmentDevelopment,
		KeyContent:  []byte(testECPrivateKey),
		httpClient:  client,
	}
}

func TestBuildURL(t *testing.T) {
	ck := newTestCloudKit(nil)

	fullURL, subpath := ck.buildURL(CKDatabasePublic, "records/query")
	assert.Equal(t,
		"https://api.apple-cloudkit.com/database/1/iCloud.com.example.app/development/public/records/query",
		fullURL,
	)
	assert.Equal(t,
		"/database/1/iCloud.com.example.app/development/public/records/query",
		subpath,
	)
}

func TestBuildURLPrivateProduction(t *testing.T) {
	ck := &cloudKit{
		KeyID:       "key",
		Container:   "iCloud.com.example.prod",
		Environment: CKEnvironmentProduction,
		KeyContent:  []byte(testECPrivateKey),
	}

	fullURL, subpath := ck.buildURL(CKDatabasePrivate, "zones/list")
	assert.Equal(t,
		"https://api.apple-cloudkit.com/database/1/iCloud.com.example.prod/production/private/zones/list",
		fullURL,
	)
	assert.Equal(t,
		"/database/1/iCloud.com.example.prod/production/private/zones/list",
		subpath,
	)
}

func TestSign(t *testing.T) {
	ck := newTestCloudKit(nil)

	date, signature, err := ck.sign([]byte(`{"test":"data"}`), "/database/1/test/development/public/records/query")
	assert.NoError(t, err)
	assert.NotEmpty(t, date)
	assert.NotEmpty(t, signature)
	assert.Contains(t, date, "T")
	assert.Contains(t, date, "Z")
}

func TestSignInvalidKey(t *testing.T) {
	ck := &cloudKit{
		KeyID:       "key",
		Container:   "container",
		Environment: CKEnvironmentDevelopment,
		KeyContent:  []byte("not-a-valid-key"),
	}

	_, _, err := ck.sign([]byte(`{}`), "/test")
	assert.Error(t, err)
}

func TestDoRequestSuccess(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	respBody := `{"records":[{"recordName":"test-record"}]}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return strings.Contains(req.URL.String(), "records/query") &&
			req.Header.Get("X-Apple-CloudKit-Request-KeyID") == "test-key-id" &&
			req.Header.Get("X-Apple-CloudKit-Request-ISO8601Date") != "" &&
			req.Header.Get("X-Apple-CloudKit-Request-SignatureV1") != ""
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	var result CKQueryResponse
	err := ck.doRequest("POST", CKDatabasePublic, "records/query",
		&CKQueryRequest{Query: CKQuery{RecordType: "TestRecord"}}, &result)
	assert.NoError(t, err)
	assert.Len(t, result.Records, 1)
	assert.Equal(t, "test-record", result.Records[0].RecordName)
}

func TestDoRequestErrorResponse(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	errResp := CKErrorResponse{
		UUID:            "test-uuid",
		ServerErrorCode: CKErrorAccessDenied,
		Reason:          "Access denied",
	}
	errRespBytes, _ := json.Marshal(errResp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{
			StatusCode: 403,
			Body:       io.NopCloser(bytes.NewReader(errRespBytes)),
		},
		nil,
	)

	var result CKQueryResponse
	err := ck.doRequest("POST", CKDatabasePublic, "records/query", nil, &result)
	assert.Error(t, err)

	var ckErr *CKError
	assert.True(t, errors.As(err, &ckErr))
	assert.Equal(t, CKErrorAccessDenied, ckErr.Code)
	assert.Equal(t, "Access denied", ckErr.Reason)
	assert.Equal(t, "test-uuid", ckErr.UUID)
}

func TestDoRequestNetworkError(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		(*http.Response)(nil),
		errors.New("connection refused"),
	)

	var result CKQueryResponse
	err := ck.doRequest("POST", CKDatabasePublic, "records/query", nil, &result)
	assert.Error(t, err)

	var ckErr *CKError
	assert.True(t, errors.As(err, &ckErr))
	assert.Equal(t, CKErrorNetworkError, ckErr.Code)
}

func TestDoRequestNilBody(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"zones":[]}`))),
		},
		nil,
	)

	var result CKZonesResponse
	err := ck.doRequest("POST", CKDatabasePublic, "zones/list", nil, &result)
	assert.NoError(t, err)
}

func TestParseECPrivateKey(t *testing.T) {
	key, err := parseECPrivateKey([]byte(testECPrivateKey))
	assert.NoError(t, err)
	assert.NotNil(t, key)
}

func TestParseECPrivateKeyInvalid(t *testing.T) {
	_, err := parseECPrivateKey([]byte("not-a-key"))
	assert.Error(t, err)
}

func TestCKErrorString(t *testing.T) {
	err := &CKError{Code: CKErrorAccessDenied, Reason: "not allowed"}
	assert.Equal(t, "cloudkit: ACCESS_DENIED: not allowed", err.Error())
}

func TestCKErrorStringNoReason(t *testing.T) {
	err := &CKError{Code: CKErrorInternalError}
	assert.Equal(t, "cloudkit: INTERNAL_ERROR", err.Error())
}

func TestDoRequestUnknownErrorResponse(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte("not json"))),
		},
		nil,
	)

	var result CKQueryResponse
	err := ck.doRequest("POST", CKDatabasePublic, "records/query", nil, &result)
	assert.Error(t, err)

	var ckErr *CKError
	assert.True(t, errors.As(err, &ckErr))
	assert.Equal(t, CKErrorUnknownError, ckErr.Code)
}
