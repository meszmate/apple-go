# apple-go

[![Go Reference](https://pkg.go.dev/badge/github.com/meszmate/apple-go.svg)](https://pkg.go.dev/github.com/meszmate/apple-go)
[![CI](https://github.com/meszmate/apple-go/actions/workflows/ci.yml/badge.svg)](https://github.com/meszmate/apple-go/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Go package for Apple Sign-In token validation, App Store Server API v2, App Store Server Notifications, and CloudKit server-to-server API.

## Installation

```
go get github.com/meszmate/apple-go
```

---

## Apple Sign-In

### Setup

From a key file:

```go
auth, err := apple.New("com.example.app", "TEAM123456", "KEYID12345", "/path/to/AuthKey.p8")
```

From a base64-encoded key (environment variable):

```go
auth, err := apple.NewB64("com.example.app", "TEAM123456", "KEYID12345", os.Getenv("APPLE_KEY"))
```

### Authorization URL

```go
cfg := apple.AuthorizeURLConfig{
    ClientID:    "com.example.app.login",
    RedirectURI: "https://example.com/auth/apple/callback",
    Scope:       []string{"email", "name"},
    State:       "csrf-123",
    Nonce:       "nonce-abc",
}

loginURL := apple.AuthorizeURL(cfg)
// Redirect the user to loginURL
```

### Validate Authorization Code

```go
// From a mobile app
tokenResponse, err := auth.ValidateCode("<AUTHORIZATION-CODE>")

// From a web app with redirect URI
tokenResponse, err := auth.ValidateCodeWithRedirectURI("<AUTHORIZATION-CODE>", "https://example.com/callback")
```

### Validate Refresh Token

```go
tokenResponse, err := auth.ValidateRefreshToken("<REFRESH-TOKEN>")
```

### Get User Info from ID Token

```go
user, err := apple.GetUserInfoFromIDToken(tokenResponse.IDToken)
if err != nil {
    log.Fatal(err)
}

fmt.Println(user.Subject) // Unique user identifier
fmt.Println(user.Email)   // User's email
```

---

## App Store Server Notifications

Parse and verify App Store Server Notifications (V1 and V2) with JWS signature verification against the Apple Root CA.

### Setup

```go
asn := apple.NewAppStoreNotifications()
```

### Parse V1 Notification

```go
notification, err := asn.ParseV1(requestBody)
if err != nil {
    log.Fatal(err)
}
fmt.Println(notification.NotificationType)
fmt.Println(notification.UnifiedReceipt.LatestReceiptInfo[0].TransactionID)
```

### Parse V2 Notification

```go
notification, err := asn.ParseV2(requestBody)
if err != nil {
    log.Fatal(err)
}
fmt.Println(notification.NotificationType) // e.g. "SUBSCRIBED"
fmt.Println(notification.Subtype)          // e.g. "INITIAL_BUY"
```

### Decode Signed Transaction / Renewal Info

```go
txn, err := asn.DecodeTransactionInfo(notification.Data.SignedTransactionInfo)
renewal, err := asn.DecodeRenewalInfo(notification.Data.SignedRenewalInfo)
```

---

## App Store Server API v2

Full HTTP client for the App Store Server API v2. Supports all 12 endpoints with JWT Bearer authentication.

### Setup

From a key file:

```go
api, err := apple.NewAppStoreServerAPI(
    "issuer-id",
    "key-id",
    "com.example.app",
    "/path/to/SubscriptionKey.p8",
    false, // true for sandbox
)
```

From a base64-encoded key:

```go
api, err := apple.NewAppStoreServerAPIB64(
    "issuer-id",
    "key-id",
    "com.example.app",
    os.Getenv("APPSTORE_KEY"),
    false,
)
```

### Transactions

```go
// Get transaction info
txnResp, err := api.GetTransactionInfo("transaction-id")

// Get transaction history
historyResp, err := api.GetTransactionHistory("original-transaction-id", &apple.ASTransactionHistoryParams{
    ProductID: "com.example.sub.monthly",
    Sort:      "DESCENDING",
})

// Send consumption info
err = api.SendConsumptionInfo("original-transaction-id", &apple.ASConsumptionRequest{
    AccountTenure:     3,
    ConsumptionStatus: 0,
    CustomerConsented: true,
    DeliveryStatus:    0,
    Platform:          1,
    PlayTime:          1,
    UserStatus:        0,
})
```

### Subscriptions

```go
// Get all subscription statuses
statuses, err := api.GetAllSubscriptionStatuses("original-transaction-id")

// Extend a subscription
extResp, err := api.ExtendSubscription("original-transaction-id", &apple.ASExtendSubscriptionRequest{
    ExtendByDays:      30,
    ExtendReasonCode:  0,
    RequestIdentifier: "unique-request-id",
})

// Mass extend subscriptions
massResp, err := api.MassExtendSubscriptions(&apple.ASMassExtendRequest{
    ExtendByDays:           30,
    ExtendReasonCode:       0,
    RequestIdentifier:      "unique-request-id",
    ProductID:              "com.example.sub.monthly",
    StorefrontCountryCodes: []string{"USA", "GBR"},
})

// Check mass extension status
statusResp, err := api.GetExtensionStatus("com.example.sub.monthly", "unique-request-id")
```

### Orders & Refunds

```go
// Look up an order
orderResp, err := api.LookUpOrderID("order-id")

// Get refund history
refundResp, err := api.GetRefundHistory("transaction-id", "")
// Paginate with revision token
nextPage, err := api.GetRefundHistory("transaction-id", refundResp.Revision)
```

### Notifications

```go
// Request a test notification
testResp, err := api.RequestTestNotification()

// Check test notification status
statusResp, err := api.GetTestNotificationStatus(testResp.TestNotificationToken)

// Get notification history
historyResp, err := api.GetNotificationHistory(&apple.ASNotificationHistoryRequest{
    StartDate: 1700000000000,
    EndDate:   1700100000000,
})
```

### Error Handling

API methods return `*ASAPIError` for server-side errors:

```go
resp, err := api.GetTransactionInfo("invalid-id")
if err != nil {
    var apiErr *apple.ASAPIError
    if errors.As(err, &apiErr) {
        fmt.Println("Error code:", apiErr.ErrorCode)
        fmt.Println("Message:", apiErr.ErrorMessage)
        fmt.Println("HTTP status:", apiErr.HTTPStatus)
        if apiErr.RetryAfter > 0 {
            fmt.Println("Retry after:", apiErr.RetryAfter, "seconds")
        }
    }
}
```

---

## CloudKit

Server-to-server API for Apple CloudKit. Requires a CloudKit server-to-server key from the Apple Developer portal.

### Setup

From a key file:

```go
ck, err := apple.NewCloudKit(
    "CloudKitKeyID",
    "iCloud.com.example.app",
    apple.CKEnvironmentDevelopment,
    "/path/to/CloudKitKey.p8",
)
```

From a base64-encoded key:

```go
ck, err := apple.NewCloudKitB64(
    "CloudKitKeyID",
    "iCloud.com.example.app",
    apple.CKEnvironmentProduction,
    os.Getenv("CLOUDKIT_KEY"),
)
```

### Records

**Query records:**

```go
resp, err := ck.QueryRecords(apple.CKDatabasePublic, &apple.CKQueryRequest{
    Query: apple.CKQuery{
        RecordType: "Todo",
        FilterBy: []apple.CKFilter{
            {
                Comparator: "EQUALS",
                FieldName:  "status",
                FieldValue: &apple.CKField{Value: "active", Type: "STRING"},
            },
        },
        SortBy: []apple.CKSort{
            {FieldName: "createdAt", Ascending: false},
        },
    },
    ResultsLimit: 50,
})
```

**Create a record:**

```go
resp, err := ck.ModifyRecords(apple.CKDatabasePublic, &apple.CKRecordsModifyRequest{
    Operations: []apple.CKRecordOperation{
        {
            OperationType: apple.CKOperationCreate,
            Record: apple.CKRecord{
                RecordType: "Todo",
                Fields: map[string]*apple.CKField{
                    "title":  {Value: "Buy groceries", Type: "STRING"},
                    "done":   {Value: 0, Type: "INT64"},
                },
            },
        },
    },
})
```

**Update a record:**

```go
resp, err := ck.ModifyRecords(apple.CKDatabasePublic, &apple.CKRecordsModifyRequest{
    Operations: []apple.CKRecordOperation{
        {
            OperationType: apple.CKOperationUpdate,
            Record: apple.CKRecord{
                RecordName:      "record-uuid",
                RecordType:      "Todo",
                RecordChangeTag: "existing-change-tag",
                Fields: map[string]*apple.CKField{
                    "done": {Value: 1, Type: "INT64"},
                },
            },
        },
    },
})
```

**Delete a record:**

```go
resp, err := ck.ModifyRecords(apple.CKDatabasePublic, &apple.CKRecordsModifyRequest{
    Operations: []apple.CKRecordOperation{
        {
            OperationType: apple.CKOperationDelete,
            Record: apple.CKRecord{
                RecordName:      "record-uuid",
                RecordType:      "Todo",
                RecordChangeTag: "existing-change-tag",
            },
        },
    },
})
```

**Lookup records by name:**

```go
resp, err := ck.LookupRecords(apple.CKDatabasePublic, &apple.CKRecordsLookupRequest{
    Records: []apple.CKRecord{
        {RecordName: "record-uuid-1"},
        {RecordName: "record-uuid-2"},
    },
})
```

**Get record changes:**

```go
resp, err := ck.RecordChanges(apple.CKDatabasePublic, &apple.CKRecordChangesRequest{
    ZoneID:    apple.CKZoneID{ZoneName: "_defaultZone"},
    SyncToken: "previous-sync-token",
})
// Use resp.SyncToken for the next call
```

### Zones

```go
// List all zones
zones, err := ck.ListZones(apple.CKDatabasePrivate)

// Create a zone
resp, err := ck.ModifyZones(apple.CKDatabasePrivate, []apple.CKZone{
    {ZoneID: apple.CKZoneID{ZoneName: "MyCustomZone"}},
}, apple.CKOperationCreate)

// Delete a zone
resp, err := ck.ModifyZones(apple.CKDatabasePrivate, []apple.CKZone{
    {ZoneID: apple.CKZoneID{ZoneName: "MyCustomZone"}},
}, apple.CKOperationDelete)

// Get zone changes
changes, err := ck.ZoneChanges(apple.CKDatabasePrivate, &apple.CKZoneChangesRequest{
    ZoneIDs: []apple.CKZoneID{{ZoneName: "_defaultZone"}},
})
```

### Subscriptions

```go
// List subscriptions
subs, err := ck.ListSubscriptions(apple.CKDatabasePublic)

// Create a subscription with push notifications
resp, err := ck.ModifySubscriptions(apple.CKDatabasePublic, &apple.CKSubscriptionsModifyRequest{
    Operations: []apple.CKSubscriptionOperation{
        {
            OperationType: apple.CKOperationCreate,
            Subscription: apple.CKSubscription{
                SubscriptionType: "query",
                Query:            &apple.CKQuery{RecordType: "Todo"},
                FiresOn:          []string{"create", "update", "delete"},
                NotificationInfo: &apple.CKNotificationInfo{
                    AlertBody: "A todo was changed",
                    ShouldSendContentAvailable: true,
                },
            },
        },
    },
})

// Delete a subscription
resp, err := ck.ModifySubscriptions(apple.CKDatabasePublic, &apple.CKSubscriptionsModifyRequest{
    Operations: []apple.CKSubscriptionOperation{
        {
            OperationType: apple.CKOperationDelete,
            Subscription:  apple.CKSubscription{SubscriptionID: "sub-id"},
        },
    },
})
```

### Assets

```go
// Request an upload URL
uploadResp, err := ck.UploadAssets(apple.CKDatabasePublic, &apple.CKAssetsUploadRequest{
    Tokens: []apple.CKAssetUploadRequest{
        {RecordType: "Photo", FieldName: "image"},
    },
})

// Upload the file data to the returned URL
file, _ := os.Open("photo.jpg")
defer file.Close()
// Use the cloudKit struct's UploadAssetData method or PUT directly
```

### Users

```go
// Get current user
user, err := ck.GetCurrentUser(apple.CKDatabasePublic)

// Discover all users
users, err := ck.DiscoverAllUsers(apple.CKDatabasePublic)

// Lookup users by email
users, err := ck.LookupUsers(apple.CKDatabasePublic, &apple.CKUserLookupRequest{
    EmailAddresses: []string{"user@example.com"},
})
```

### APNs Tokens

```go
resp, err := ck.CreateTokens(apple.CKDatabasePublic, &apple.CKTokensCreateRequest{
    Tokens: []apple.CKAPNsToken{
        {APNsToken: "device-token-hex", APNsEnvironment: "production"},
    },
})
```

### Push Notifications

Parse incoming APNs payloads from CloudKit subscriptions:

```go
notification, err := apple.ParseCKPushNotification(apnsPayload)
if err != nil {
    log.Fatal(err)
}

fmt.Println(notification.CK.ContainerIdentifier)
fmt.Println(notification.CK.NotificationID)

// Query subscription notification
if qry := notification.CK.QueryNotification; qry != nil {
    fmt.Println("Record:", qry.RecordName)
    fmt.Println("Type:", qry.RecordType)
    fmt.Println("Reason:", qry.QueryNotificationReason) // 1=created, 2=updated, 3=deleted
}

// Record zone subscription notification
if zry := notification.CK.RecordZoneNotification; zry != nil {
    fmt.Println("Zone:", zry.ZoneID.ZoneName)
    fmt.Println("Subscription:", zry.SubscriptionID)
}
```

### Error Handling

All CloudKit methods return `*CKError` on failure:

```go
resp, err := ck.QueryRecords(apple.CKDatabasePublic, req)
if err != nil {
    var ckErr *apple.CKError
    if errors.As(err, &ckErr) {
        fmt.Println("CloudKit error:", ckErr.Code)
        fmt.Println("Reason:", ckErr.Reason)
        if ckErr.RetryAfter > 0 {
            fmt.Println("Retry after:", ckErr.RetryAfter, "seconds")
        }
    }
}
```

---

## License

MIT
