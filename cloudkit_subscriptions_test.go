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

func TestListSubscriptions(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKSubscriptionsResponse{
		Subscriptions: []CKSubscription{
			{SubscriptionID: "sub-1", SubscriptionType: "query"},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/public/subscriptions/list"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.ListSubscriptions(CKDatabasePublic)
	assert.NoError(t, err)
	assert.Len(t, result.Subscriptions, 1)
	assert.Equal(t, "sub-1", result.Subscriptions[0].SubscriptionID)
}

func TestListSubscriptionsError(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	errResp := CKErrorResponse{ServerErrorCode: CKErrorAccessDenied, Reason: "denied"}
	errBytes, _ := json.Marshal(errResp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{StatusCode: 403, Body: io.NopCloser(bytes.NewReader(errBytes))},
		nil,
	)

	result, err := ck.ListSubscriptions(CKDatabasePublic)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestModifySubscriptions(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKSubscriptionsModifyResponse{
		Subscriptions: []CKSubscription{
			{SubscriptionID: "new-sub", SubscriptionType: "query"},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/public/subscriptions/modify"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.ModifySubscriptions(CKDatabasePublic, &CKSubscriptionsModifyRequest{
		Operations: []CKSubscriptionOperation{
			{
				OperationType: CKOperationCreate,
				Subscription: CKSubscription{
					SubscriptionType: "query",
					Query:            &CKQuery{RecordType: "TestType"},
					FiresOn:          []string{"create", "update"},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, result.Subscriptions, 1)
}
