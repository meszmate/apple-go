package apple

import (
	"fmt"
	"net/url"
	"strconv"
)

// GetTransactionInfo gets transaction info for a specific transaction ID.
func (s *appStoreServer) GetTransactionInfo(transactionID string) (*ASTransactionInfoResponse, error) {
	path := fmt.Sprintf("/inApps/v1/transactions/%s", transactionID)
	var result ASTransactionInfoResponse
	if err := s.doRequest("GET", path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTransactionHistory gets the transaction history for an original transaction ID.
func (s *appStoreServer) GetTransactionHistory(originalTransactionID string, params *ASTransactionHistoryParams) (*ASTransactionHistoryResponse, error) {
	path := fmt.Sprintf("/inApps/v2/history/%s", originalTransactionID)

	var q url.Values
	if params != nil {
		q = make(url.Values)
		if params.Revision != "" {
			q.Set("revision", params.Revision)
		}
		if params.StartDate != 0 {
			q.Set("startDate", strconv.FormatInt(params.StartDate, 10))
		}
		if params.EndDate != 0 {
			q.Set("endDate", strconv.FormatInt(params.EndDate, 10))
		}
		if params.ProductID != "" {
			q.Set("productId", params.ProductID)
		}
		if params.ProductType != "" {
			q.Set("productType", params.ProductType)
		}
		if params.Sort != "" {
			q.Set("sort", params.Sort)
		}
		if params.InAppOwnershipType != "" {
			q.Set("inAppOwnershipType", params.InAppOwnershipType)
		}
		if params.Revoked != nil {
			q.Set("revoked", strconv.FormatBool(*params.Revoked))
		}
	}

	var result ASTransactionHistoryResponse
	if err := s.doRequest("GET", path, q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SendConsumptionInfo sends consumption info for an original transaction ID.
func (s *appStoreServer) SendConsumptionInfo(originalTransactionID string, req *ASConsumptionRequest) error {
	path := fmt.Sprintf("/inApps/v1/transactions/consumption/%s", originalTransactionID)
	return s.doRequest("PUT", path, nil, req, nil)
}
