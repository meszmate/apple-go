package apple

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllSubscriptionStatuses(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{
		"environment":"Sandbox",
		"appAppleId":123,
		"bundleId":"com.example.app",
		"data":[{
			"subscriptionGroupIdentifier":"group1",
			"lastTransactions":[{
				"status":1,
				"originalTransactionId":"orig123",
				"signedTransactionInfo":"signed.txn",
				"signedRenewalInfo":"signed.renewal"
			}]
		}]
	}`

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/subscriptions/orig123" && req.Method == "GET"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.GetAllSubscriptionStatuses("orig123")
	assert.NoError(t, err)
	assert.Equal(t, "Sandbox", result.Environment)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "group1", result.Data[0].SubscriptionGroupIdentifier)
	assert.Len(t, result.Data[0].LastTransactions, 1)
	assert.Equal(t, 1, result.Data[0].LastTransactions[0].Status)
	assert.Equal(t, "signed.txn", result.Data[0].LastTransactions[0].SignedTransactionInfo)
}

func TestExtendSubscription(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"effectiveDate":1700000000000,"originalTransactionId":"orig123","success":true,"webOrderLineItemId":"woli123"}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/subscriptions/extend/orig123" && req.Method == "PUT"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.ExtendSubscription("orig123", &ASExtendSubscriptionRequest{
		ExtendByDays:      30,
		ExtendReasonCode:  0,
		RequestIdentifier: "req-123",
	})
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "orig123", result.OriginalTransactionID)
	assert.Equal(t, int64(1700000000000), result.EffectiveDate)
}

func TestMassExtendSubscriptions(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"requestIdentifier":"req-456"}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/subscriptions/extend/mass" && req.Method == "POST"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.MassExtendSubscriptions(&ASMassExtendRequest{
		ExtendByDays:           30,
		ExtendReasonCode:       0,
		RequestIdentifier:      "req-456",
		ProductID:              "com.example.sub",
		StorefrontCountryCodes: []string{"USA"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "req-456", result.RequestIdentifier)
}

func TestGetExtensionStatus(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"requestIdentifier":"req-456","complete":true,"completeDate":1700000000000,"succeededCount":100,"failedCount":5}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/subscriptions/extend/mass/com.example.sub/req-456" && req.Method == "GET"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.GetExtensionStatus("com.example.sub", "req-456")
	assert.NoError(t, err)
	assert.True(t, result.Complete)
	assert.Equal(t, int64(100), result.SucceededCount)
	assert.Equal(t, int64(5), result.FailedCount)
}
