package apple

import (
	"fmt"
	"net/url"
)

// LookUpOrderID looks up an order by its order ID.
func (s *appStoreServer) LookUpOrderID(orderID string) (*ASOrderLookupResponse, error) {
	path := fmt.Sprintf("/inApps/v1/lookup/%s", orderID)
	var result ASOrderLookupResponse
	if err := s.doRequest("GET", path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRefundHistory gets the refund history for a transaction.
func (s *appStoreServer) GetRefundHistory(transactionID string, revision string) (*ASRefundHistoryResponse, error) {
	path := fmt.Sprintf("/inApps/v2/refund-history/%s", transactionID)

	var q url.Values
	if revision != "" {
		q = make(url.Values)
		q.Set("revision", revision)
	}

	var result ASRefundHistoryResponse
	if err := s.doRequest("GET", path, q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
