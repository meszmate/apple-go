package apple

// ListSubscriptions lists all subscriptions in a CloudKit database.
func (c *cloudKit) ListSubscriptions(db CKDatabase) (*CKSubscriptionsResponse, error) {
	var resp CKSubscriptionsResponse
	if err := c.doRequest("POST", db, "subscriptions/list", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifySubscriptions modifies subscriptions in a CloudKit database.
func (c *cloudKit) ModifySubscriptions(db CKDatabase, req *CKSubscriptionsModifyRequest) (*CKSubscriptionsModifyResponse, error) {
	var resp CKSubscriptionsModifyResponse
	if err := c.doRequest("POST", db, "subscriptions/modify", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
