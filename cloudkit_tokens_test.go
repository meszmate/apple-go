package apple

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTokens(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKTokensCreateResponse{
		Tokens: []CKAPNsToken{
			{APNsToken: "device-token-abc", APNsEnvironment: "development"},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/public/tokens/create"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.CreateTokens(CKDatabasePublic, &CKTokensCreateRequest{
		Tokens: []CKAPNsToken{
			{APNsToken: "device-token-abc", APNsEnvironment: "development"},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, result.Tokens, 1)
	assert.Equal(t, "device-token-abc", result.Tokens[0].APNsToken)
}

func TestCreateTokensError(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	errResp := CKErrorResponse{ServerErrorCode: CKErrorBadRequest, Reason: "invalid token"}
	errBytes, _ := json.Marshal(errResp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{StatusCode: 400, Body: io.NopCloser(bytes.NewReader(errBytes))},
		nil,
	)

	result, err := ck.CreateTokens(CKDatabasePublic, &CKTokensCreateRequest{
		Tokens: []CKAPNsToken{{APNsToken: "bad-token"}},
	})
	assert.Error(t, err)
	assert.Nil(t, result)
}
