package apple

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequestTestNotification(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"testNotificationToken":"token123"}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/notifications/test" && req.Method == "POST"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.RequestTestNotification()
	assert.NoError(t, err)
	assert.Equal(t, "token123", result.TestNotificationToken)
}

func TestGetTestNotificationStatus(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{
		"signedPayload":"signed.payload.jws",
		"sendAttempts":[
			{"attemptDate":1700000000000,"sendAttemptResult":"SUCCESS"}
		],
		"firstSendAttemptResult":"SUCCESS"
	}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/notifications/test/token123" && req.Method == "GET"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.GetTestNotificationStatus("token123")
	assert.NoError(t, err)
	assert.Equal(t, "signed.payload.jws", result.SignedPayload)
	assert.Len(t, result.SendAttempts, 1)
	assert.Equal(t, "SUCCESS", result.SendAttempts[0].SendAttemptResult)
	assert.Equal(t, "SUCCESS", result.FirstSendAttemptResult)
}

func TestGetNotificationHistory(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{
		"notificationHistory":[
			{
				"signedPayload":"signed.payload.jws",
				"sendAttempts":[{"attemptDate":1700000000000,"sendAttemptResult":"SUCCESS"}],
				"firstSendAttemptResult":"SUCCESS"
			}
		],
		"hasMore":false,
		"paginationToken":"page-token"
	}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/notifications/history" && req.Method == "POST"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.GetNotificationHistory(&ASNotificationHistoryRequest{
		StartDate: 1700000000000,
		EndDate:   1700100000000,
	})
	assert.NoError(t, err)
	assert.Len(t, result.NotificationHistory, 1)
	assert.False(t, result.HasMore)
	assert.Equal(t, "page-token", result.PaginationToken)
	assert.Equal(t, "signed.payload.jws", result.NotificationHistory[0].SignedPayload)
}
