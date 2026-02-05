package apple

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCKPushNotification_QueryNotification(t *testing.T) {
	payload := []byte(`{
		"ck": {
			"cid": "iCloud.com.example.app",
			"nid": "notification-123",
			"qry": {
				"rid": "record-abc",
				"rct": "MyRecordType",
				"sid": "subscription-1",
				"fo": 1,
				"af": {
					"name": "test-value"
				},
				"zid": {
					"zoneName": "_defaultZone"
				}
			}
		}
	}`)

	n, err := ParseCKPushNotification(payload)
	assert.NoError(t, err)
	assert.NotNil(t, n.CK)
	assert.Equal(t, "iCloud.com.example.app", n.CK.ContainerIdentifier)
	assert.Equal(t, "notification-123", n.CK.NotificationID)
	assert.NotNil(t, n.CK.QueryNotification)
	assert.Equal(t, "record-abc", n.CK.QueryNotification.RecordName)
	assert.Equal(t, "MyRecordType", n.CK.QueryNotification.RecordType)
	assert.Equal(t, "subscription-1", n.CK.QueryNotification.SubscriptionID)
	assert.Equal(t, 1, n.CK.QueryNotification.QueryNotificationReason)
	assert.Equal(t, "test-value", n.CK.QueryNotification.RecordFields["name"])
	assert.Equal(t, "_defaultZone", n.CK.QueryNotification.ZoneID.ZoneName)
	assert.Nil(t, n.CK.RecordZoneNotification)
}

func TestParseCKPushNotification_RecordZoneNotification(t *testing.T) {
	payload := []byte(`{
		"ck": {
			"cid": "iCloud.com.example.app",
			"nid": "notification-456",
			"zry": {
				"sid": "zone-sub-1",
				"zid": {
					"zoneName": "CustomZone",
					"ownerRecordName": "_currentUser"
				}
			}
		}
	}`)

	n, err := ParseCKPushNotification(payload)
	assert.NoError(t, err)
	assert.NotNil(t, n.CK)
	assert.Equal(t, "iCloud.com.example.app", n.CK.ContainerIdentifier)
	assert.Nil(t, n.CK.QueryNotification)
	assert.NotNil(t, n.CK.RecordZoneNotification)
	assert.Equal(t, "zone-sub-1", n.CK.RecordZoneNotification.SubscriptionID)
	assert.Equal(t, "CustomZone", n.CK.RecordZoneNotification.ZoneID.ZoneName)
	assert.Equal(t, "_currentUser", n.CK.RecordZoneNotification.ZoneID.OwnerName)
}

func TestParseCKPushNotification_InvalidJSON(t *testing.T) {
	_, err := ParseCKPushNotification([]byte(`{invalid`))
	assert.Error(t, err)

	var ckErr *CKError
	assert.True(t, errors.As(err, &ckErr))
	assert.Equal(t, CKErrorBadRequest, ckErr.Code)
}

func TestParseCKPushNotification_MissingCKField(t *testing.T) {
	_, err := ParseCKPushNotification([]byte(`{"aps":{"content-available":1}}`))
	assert.Error(t, err)

	var ckErr *CKError
	assert.True(t, errors.As(err, &ckErr))
	assert.Equal(t, CKErrorBadRequest, ckErr.Code)
	assert.Contains(t, ckErr.Reason, "missing ck field")
}

func TestParseCKPushNotification_MinimalPayload(t *testing.T) {
	payload := []byte(`{
		"ck": {
			"cid": "iCloud.com.example.app",
			"nid": "notification-789"
		}
	}`)

	n, err := ParseCKPushNotification(payload)
	assert.NoError(t, err)
	assert.Equal(t, "iCloud.com.example.app", n.CK.ContainerIdentifier)
	assert.Equal(t, "notification-789", n.CK.NotificationID)
	assert.Nil(t, n.CK.QueryNotification)
	assert.Nil(t, n.CK.RecordZoneNotification)
}
