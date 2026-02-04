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

func TestGetCurrentUser(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKUserInfo{
		UserRecordName: "_abc123",
		FirstName:      "John",
		LastName:       "Doe",
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/public/users/current"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.GetCurrentUser(CKDatabasePublic)
	assert.NoError(t, err)
	assert.Equal(t, "_abc123", result.UserRecordName)
	assert.Equal(t, "John", result.FirstName)
}

func TestGetCurrentUserError(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	errResp := CKErrorResponse{ServerErrorCode: CKErrorAuthenticationFailed, Reason: "bad key"}
	errBytes, _ := json.Marshal(errResp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{StatusCode: 401, Body: io.NopCloser(bytes.NewReader(errBytes))},
		nil,
	)

	result, err := ck.GetCurrentUser(CKDatabasePublic)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDiscoverAllUsers(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKUsersResponse{
		Users: []CKUserInfo{
			{UserRecordName: "_user1", FirstName: "Alice"},
			{UserRecordName: "_user2", FirstName: "Bob"},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/public/users/discover"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.DiscoverAllUsers(CKDatabasePublic)
	assert.NoError(t, err)
	assert.Len(t, result.Users, 2)
}

func TestLookupUsersByEmail(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKUsersResponse{
		Users: []CKUserInfo{
			{UserRecordName: "_user1", EmailAddress: "alice@example.com"},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/public/users/lookup/email"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.LookupUsers(CKDatabasePublic, &CKUserLookupRequest{
		EmailAddresses: []string{"alice@example.com"},
	})
	assert.NoError(t, err)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "alice@example.com", result.Users[0].EmailAddress)
}

func TestLookupUsersByPhone(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKUsersResponse{
		Users: []CKUserInfo{
			{UserRecordName: "_user2", FirstName: "Bob"},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/public/users/lookup/phone"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.LookupUsers(CKDatabasePublic, &CKUserLookupRequest{
		PhoneNumbers: []string{"+1234567890"},
	})
	assert.NoError(t, err)
	assert.Len(t, result.Users, 1)
}

func TestLookupUsersEmpty(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	result, err := ck.LookupUsers(CKDatabasePublic, &CKUserLookupRequest{})
	assert.NoError(t, err)
	assert.Len(t, result.Users, 0)
}
