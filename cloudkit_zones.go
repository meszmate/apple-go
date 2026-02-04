package apple

// ListZones lists all zones in a CloudKit database.
func (c *cloudKit) ListZones(db CKDatabase) (*CKZonesResponse, error) {
	var resp CKZonesResponse
	if err := c.doRequest("POST", db, "zones/list", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupZones looks up specific zones in a CloudKit database.
func (c *cloudKit) LookupZones(db CKDatabase, zoneIDs []CKZoneID) (*CKZonesResponse, error) {
	body := struct {
		Zones []CKZoneID `json:"zones"`
	}{
		Zones: zoneIDs,
	}
	var resp CKZonesResponse
	if err := c.doRequest("POST", db, "zones/lookup", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifyZones modifies zones in a CloudKit database.
func (c *cloudKit) ModifyZones(db CKDatabase, zones []CKZone, operationType CKOperationType) (*CKZonesResponse, error) {
	type zoneOp struct {
		OperationType CKOperationType `json:"operationType"`
		Zone          CKZone          `json:"zone"`
	}
	operations := make([]zoneOp, len(zones))
	for i, z := range zones {
		operations[i] = zoneOp{
			OperationType: operationType,
			Zone:          z,
		}
	}
	body := struct {
		Operations []zoneOp `json:"operations"`
	}{
		Operations: operations,
	}
	var resp CKZonesResponse
	if err := c.doRequest("POST", db, "zones/modify", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ZoneChanges gets zone changes in a CloudKit database.
func (c *cloudKit) ZoneChanges(db CKDatabase, req *CKZoneChangesRequest) (*CKZoneChangesResponse, error) {
	var resp CKZoneChangesResponse
	if err := c.doRequest("POST", db, "zones/changes", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
