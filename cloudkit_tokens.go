package apple

// CreateTokens registers APNs tokens for CloudKit push notifications.
func (c *cloudKit) CreateTokens(db CKDatabase, req *CKTokensCreateRequest) (*CKTokensCreateResponse, error) {
	var resp CKTokensCreateResponse
	if err := c.doRequest("POST", db, "tokens/create", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
