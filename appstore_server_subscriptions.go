package apple

import "fmt"

// GetAllSubscriptionStatuses gets subscription statuses for an original transaction ID.
func (s *appStoreServer) GetAllSubscriptionStatuses(originalTransactionID string) (*ASSubscriptionStatusesResponse, error) {
	path := fmt.Sprintf("/inApps/v1/subscriptions/%s", originalTransactionID)
	var result ASSubscriptionStatusesResponse
	if err := s.doRequest("GET", path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExtendSubscription extends a subscription renewal date.
func (s *appStoreServer) ExtendSubscription(originalTransactionID string, req *ASExtendSubscriptionRequest) (*ASExtendSubscriptionResponse, error) {
	path := fmt.Sprintf("/inApps/v1/subscriptions/extend/%s", originalTransactionID)
	var result ASExtendSubscriptionResponse
	if err := s.doRequest("PUT", path, nil, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// MassExtendSubscriptions mass-extends subscriptions.
func (s *appStoreServer) MassExtendSubscriptions(req *ASMassExtendRequest) (*ASMassExtendResponse, error) {
	var result ASMassExtendResponse
	if err := s.doRequest("POST", "/inApps/v1/subscriptions/extend/mass", nil, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetExtensionStatus gets the status of a mass extension request.
func (s *appStoreServer) GetExtensionStatus(productID, requestIdentifier string) (*ASExtensionStatusResponse, error) {
	path := fmt.Sprintf("/inApps/v1/subscriptions/extend/mass/%s/%s", productID, requestIdentifier)
	var result ASExtensionStatusResponse
	if err := s.doRequest("GET", path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
