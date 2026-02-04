package apple

// QueryRecords queries records in a CloudKit database.
func (c *cloudKit) QueryRecords(db CKDatabase, req *CKQueryRequest) (*CKQueryResponse, error) {
	var resp CKQueryResponse
	if err := c.doRequest("POST", db, "records/query", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifyRecords modifies records in a CloudKit database.
func (c *cloudKit) ModifyRecords(db CKDatabase, req *CKRecordsModifyRequest) (*CKRecordsModifyResponse, error) {
	var resp CKRecordsModifyResponse
	if err := c.doRequest("POST", db, "records/modify", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupRecords looks up records in a CloudKit database.
func (c *cloudKit) LookupRecords(db CKDatabase, req *CKRecordsLookupRequest) (*CKRecordsLookupResponse, error) {
	var resp CKRecordsLookupResponse
	if err := c.doRequest("POST", db, "records/lookup", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RecordChanges gets record changes in a CloudKit database.
func (c *cloudKit) RecordChanges(db CKDatabase, req *CKRecordChangesRequest) (*CKRecordChangesResponse, error) {
	var resp CKRecordChangesResponse
	if err := c.doRequest("POST", db, "records/changes", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
