package apple

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadAssets(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKAssetsUploadResponse{
		Tokens: []CKAssetUploadResponse{
			{RecordName: "rec-1", FieldName: "photo", URL: "https://upload.example.com/token123"},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/public/assets/upload"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.UploadAssets(CKDatabasePublic, &CKAssetsUploadRequest{
		Tokens: []CKAssetUploadRequest{
			{RecordType: "Photo", FieldName: "photo"},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, result.Tokens, 1)
	assert.Equal(t, "https://upload.example.com/token123", result.Tokens[0].URL)
}

func TestUploadAssetData(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == "PUT" &&
			req.URL.String() == "https://upload.example.com/token123" &&
			req.Header.Get("Content-Type") == "application/octet-stream"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader([]byte("{}")))},
		nil,
	)

	err := ck.UploadAssetData("https://upload.example.com/token123", strings.NewReader("file-content"))
	assert.NoError(t, err)
}

func TestUploadAssetDataError(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == "PUT"
	})).Return(
		&http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte("error")))},
		nil,
	)

	err := ck.UploadAssetData("https://upload.example.com/token123", strings.NewReader("data"))
	assert.Error(t, err)
}
