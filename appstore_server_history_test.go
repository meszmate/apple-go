package apple

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLookUpOrderID(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"status":0,"signedTransactions":["jws1","jws2"]}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v1/lookup/ORDER123" && req.Method == "GET"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.LookUpOrderID("ORDER123")
	assert.NoError(t, err)
	assert.Equal(t, 0, result.Status)
	assert.Len(t, result.SignedTransactions, 2)
}

func TestGetRefundHistory(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"hasMore":false,"revision":"rev1","signedTransactions":["jws1"]}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v2/refund-history/txn123" && req.Method == "GET"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.GetRefundHistory("txn123", "")
	assert.NoError(t, err)
	assert.False(t, result.HasMore)
	assert.Equal(t, "rev1", result.Revision)
	assert.Len(t, result.SignedTransactions, 1)
}

func TestGetRefundHistory_WithRevision(t *testing.T) {
	mockedClient := new(MockedASHTTPClient)
	s := newTestAppStoreServer(mockedClient)

	respBody := `{"hasMore":false,"revision":"rev2","signedTransactions":[]}`
	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/inApps/v2/refund-history/txn123" &&
			req.URL.Query().Get("revision") == "rev1"
	})).Return(
		&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(respBody))),
		},
		nil,
	)

	result, err := s.GetRefundHistory("txn123", "rev1")
	assert.NoError(t, err)
	assert.Equal(t, "rev2", result.Revision)
}
