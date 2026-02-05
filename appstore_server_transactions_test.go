package apple

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTransactionInfo(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"signedTransactionInfo":"signed.txn.jws"}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/transactions/12345" && req.Method == "GET"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.GetTransactionInfo("12345")
	assert.NoError(t, err)
	assert.Equal(t, "signed.txn.jws", result.SignedTransactionInfo)
}

func TestGetTransactionHistory(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"revision":"rev1","hasMore":true,"bundleId":"com.example","appAppleId":123,"environment":"Sandbox","signedTransactions":["jws1","jws2"]}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v2/history/orig123" &&
			req.Method == "GET" &&
			req.URL.Query().Get("productId") == "com.example.sub"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.GetTransactionHistory("orig123", &ASTransactionHistoryParams{
		ProductID: "com.example.sub",
	})
	assert.NoError(t, err)
	assert.Equal(t, "rev1", result.Revision)
	assert.True(t, result.HasMore)
	assert.Len(t, result.SignedTransactions, 2)
}

func TestGetTransactionHistory_NilParams(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"revision":"","hasMore":false,"bundleId":"com.example","appAppleId":123,"environment":"Sandbox","signedTransactions":[]}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v2/history/orig123" &&
			req.URL.RawQuery == ""
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.GetTransactionHistory("orig123", nil)
	assert.NoError(t, err)
	assert.False(t, result.HasMore)
}

func TestSendConsumptionInfo(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/transactions/consumption/orig123" &&
			req.Method == "PUT"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		},
		nil,
	)

	err := s.SendConsumptionInfo("orig123", &ASConsumptionRequest{
		AccountTenure:     3,
		ConsumptionStatus: 0,
		CustomerConsented: true,
		DeliveryStatus:    0,
		Platform:          1,
		PlayTime:          1,
		UserStatus:        0,
	})
	assert.NoError(t, err)
}
