package apple

import "encoding/json"

// CKPushNotification represents a CloudKit push notification payload from APNs.
type CKPushNotification struct {
	CK *CKNotificationPayload `json:"ck"`
}

// CKNotificationPayload represents the CloudKit-specific notification data.
type CKNotificationPayload struct {
	ContainerIdentifier    string                      `json:"cid"`
	NotificationID         string                      `json:"nid"`
	QueryNotification      *CKQueryNotification        `json:"qry,omitempty"`
	RecordZoneNotification *CKRecordZoneNotification   `json:"zry,omitempty"`
}

// CKQueryNotification represents a query-based subscription notification.
type CKQueryNotification struct {
	RecordName              string         `json:"rid,omitempty"`
	RecordType              string         `json:"rct,omitempty"`
	SubscriptionID          string         `json:"sid,omitempty"`
	ZoneID                  *CKZoneID      `json:"zid,omitempty"`
	QueryNotificationReason int            `json:"fo,omitempty"` // 1=created, 2=updated, 3=deleted
	RecordFields            map[string]any `json:"af,omitempty"`
}

// CKRecordZoneNotification represents a record zone subscription notification.
type CKRecordZoneNotification struct {
	ZoneID         *CKZoneID `json:"zid,omitempty"`
	SubscriptionID string    `json:"sid,omitempty"`
}

// ParseCKPushNotification parses a CloudKit push notification from raw APNs payload bytes.
func ParseCKPushNotification(payload []byte) (*CKPushNotification, error) {
	var notification CKPushNotification
	if err := json.Unmarshal(payload, &notification); err != nil {
		return nil, &CKError{Code: CKErrorBadRequest, Reason: err.Error()}
	}
	if notification.CK == nil {
		return nil, &CKError{Code: CKErrorBadRequest, Reason: "missing ck field"}
	}
	return &notification, nil
}
