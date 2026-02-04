package apple

// GetCurrentUser gets the current user's identity from a CloudKit database.
func (c *cloudKit) GetCurrentUser(db CKDatabase) (*CKUserInfo, error) {
	var resp CKUserInfo
	if err := c.doRequest("POST", db, "users/current", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DiscoverAllUsers discovers all users in a CloudKit database.
func (c *cloudKit) DiscoverAllUsers(db CKDatabase) (*CKUsersResponse, error) {
	var resp CKUsersResponse
	if err := c.doRequest("POST", db, "users/discover", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupUsers looks up users by email addresses or phone numbers.
func (c *cloudKit) LookupUsers(db CKDatabase, req *CKUserLookupRequest) (*CKUsersResponse, error) {
	var resp CKUsersResponse

	if len(req.EmailAddresses) > 0 {
		body := struct {
			EmailAddresses []string `json:"emailAddresses"`
		}{
			EmailAddresses: req.EmailAddresses,
		}
		if err := c.doRequest("POST", db, "users/lookup/email", body, &resp); err != nil {
			return nil, err
		}
	}

	if len(req.PhoneNumbers) > 0 {
		var phoneResp CKUsersResponse
		body := struct {
			PhoneNumbers []string `json:"phoneNumbers"`
		}{
			PhoneNumbers: req.PhoneNumbers,
		}
		if err := c.doRequest("POST", db, "users/lookup/phone", body, &phoneResp); err != nil {
			return nil, err
		}
		resp.Users = append(resp.Users, phoneResp.Users...)
	}

	return &resp, nil
}
