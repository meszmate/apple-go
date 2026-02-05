package apple

// --- Transaction Types ---

// ASTransactionInfoResponse represents the response for GetTransactionInfo.
type ASTransactionInfoResponse struct {
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}

// ASTransactionHistoryParams represents query parameters for GetTransactionHistory.
type ASTransactionHistoryParams struct {
	Revision           string `url:"revision,omitempty"`
	StartDate          int64  `url:"startDate,omitempty"`
	EndDate            int64  `url:"endDate,omitempty"`
	ProductID          string `url:"productId,omitempty"`
	ProductType        string `url:"productType,omitempty"`
	Sort               string `url:"sort,omitempty"`
	InAppOwnershipType string `url:"inAppOwnershipType,omitempty"`
	Revoked            *bool  `url:"revoked,omitempty"`
}

// ASTransactionHistoryResponse represents the response for GetTransactionHistory.
type ASTransactionHistoryResponse struct {
	Revision           string   `json:"revision"`
	HasMore            bool     `json:"hasMore"`
	BundleID           string   `json:"bundleId"`
	AppAppleID         int64    `json:"appAppleId"`
	Environment        string   `json:"environment"`
	SignedTransactions []string `json:"signedTransactions"`
}

// --- Subscription Types ---

// ASSubscriptionStatusesResponse represents the response for GetAllSubscriptionStatuses.
type ASSubscriptionStatusesResponse struct {
	Environment string                    `json:"environment"`
	AppAppleID  int64                     `json:"appAppleId"`
	BundleID    string                    `json:"bundleId"`
	Data        []ASSubscriptionGroupStatus `json:"data"`
}

// ASSubscriptionGroupStatus represents a subscription group's status.
type ASSubscriptionGroupStatus struct {
	SubscriptionGroupIdentifier string                  `json:"subscriptionGroupIdentifier"`
	LastTransactions            []ASLastTransactionItem `json:"lastTransactions"`
}

// ASLastTransactionItem represents the last transaction in a subscription group.
type ASLastTransactionItem struct {
	Status                int    `json:"status"`
	OriginalTransactionID string `json:"originalTransactionId"`
	SignedTransactionInfo string `json:"signedTransactionInfo"`
	SignedRenewalInfo     string `json:"signedRenewalInfo"`
}

// ASExtendSubscriptionRequest represents a request to extend a subscription.
type ASExtendSubscriptionRequest struct {
	ExtendByDays      int    `json:"extendByDays"`
	ExtendReasonCode  int    `json:"extendReasonCode"`
	RequestIdentifier string `json:"requestIdentifier"`
}

// ASExtendSubscriptionResponse represents the response for extending a subscription.
type ASExtendSubscriptionResponse struct {
	EffectiveDate         int64  `json:"effectiveDate"`
	OriginalTransactionID string `json:"originalTransactionId"`
	Success               bool   `json:"success"`
	WebOrderLineItemID    string `json:"webOrderLineItemId"`
}

// ASMassExtendRequest represents a request to mass-extend subscriptions.
type ASMassExtendRequest struct {
	ExtendByDays            int      `json:"extendByDays"`
	ExtendReasonCode        int      `json:"extendReasonCode"`
	RequestIdentifier       string   `json:"requestIdentifier"`
	ProductID               string   `json:"productId"`
	StorefrontCountryCodes  []string `json:"storefrontCountryCodes,omitempty"`
}

// ASMassExtendResponse represents the response for mass-extending subscriptions.
type ASMassExtendResponse struct {
	RequestIdentifier string `json:"requestIdentifier"`
}

// ASExtensionStatusResponse represents the response for GetExtensionStatus.
type ASExtensionStatusResponse struct {
	RequestIdentifier string `json:"requestIdentifier"`
	Complete          bool   `json:"complete"`
	CompleteDate      int64  `json:"completeDate,omitempty"`
	SucceededCount    int64  `json:"succeededCount,omitempty"`
	FailedCount       int64  `json:"failedCount,omitempty"`
}

// --- Order / Refund Types ---

// ASOrderLookupResponse represents the response for LookUpOrderID.
type ASOrderLookupResponse struct {
	Status             int      `json:"status"`
	SignedTransactions []string `json:"signedTransactions"`
}

// ASRefundHistoryResponse represents the response for GetRefundHistory.
type ASRefundHistoryResponse struct {
	HasMore            bool     `json:"hasMore"`
	Revision           string   `json:"revision"`
	SignedTransactions []string `json:"signedTransactions"`
}

// --- Consumption Types ---

// ASConsumptionRequest represents a request to send consumption info.
type ASConsumptionRequest struct {
	AccountTenure             int    `json:"accountTenure"`
	AppAccountToken           string `json:"appAccountToken,omitempty"`
	ConsumptionStatus         int    `json:"consumptionStatus"`
	CustomerConsented         bool   `json:"customerConsented"`
	DeliveryStatus            int    `json:"deliveryStatus"`
	LifetimeDollarsPurchased  int    `json:"lifetimeDollarsPurchased"`
	LifetimeDollarsRefunded   int    `json:"lifetimeDollarsRefunded"`
	Platform                  int    `json:"platform"`
	PlayTime                  int    `json:"playTime"`
	SampleContentProvided     bool   `json:"sampleContentProvided"`
	UserStatus                int    `json:"userStatus"`
	RefundPreference          int    `json:"refundPreference,omitempty"`
}

// --- Notification Types ---

// ASTestNotificationResponse represents the response for RequestTestNotification.
type ASTestNotificationResponse struct {
	TestNotificationToken string `json:"testNotificationToken"`
}

// ASTestNotificationStatusResponse represents the response for GetTestNotificationStatus.
type ASTestNotificationStatusResponse struct {
	SignedPayload          string                     `json:"signedPayload"`
	SendAttempts           []ASNotificationSendAttempt `json:"sendAttempts"`
	FirstSendAttemptResult string                     `json:"firstSendAttemptResult,omitempty"`
}

// ASNotificationSendAttempt represents a send attempt for a notification.
type ASNotificationSendAttempt struct {
	AttemptDate       int64  `json:"attemptDate"`
	SendAttemptResult string `json:"sendAttemptResult"`
}

// ASNotificationHistoryRequest represents a request to get notification history.
type ASNotificationHistoryRequest struct {
	StartDate           int64  `json:"startDate"`
	EndDate             int64  `json:"endDate"`
	NotificationType    string `json:"notificationType,omitempty"`
	NotificationSubtype string `json:"notificationSubtype,omitempty"`
	TransactionID       string `json:"transactionId,omitempty"`
	OnlyFailures        bool   `json:"onlyFailures,omitempty"`
	PaginationToken     string `json:"paginationToken,omitempty"`
}

// ASNotificationHistoryResponse represents the response for GetNotificationHistory.
type ASNotificationHistoryResponse struct {
	NotificationHistory []ASNotificationHistoryItem `json:"notificationHistory"`
	HasMore             bool                        `json:"hasMore"`
	PaginationToken     string                      `json:"paginationToken,omitempty"`
}

// ASNotificationHistoryItem represents a single notification in the history.
type ASNotificationHistoryItem struct {
	SignedPayload          string                     `json:"signedPayload"`
	SendAttempts           []ASNotificationSendAttempt `json:"sendAttempts"`
	FirstSendAttemptResult string                     `json:"firstSendAttemptResult,omitempty"`
}
