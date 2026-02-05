package apple

// ASEnvironment represents the App Store server environment.
type ASEnvironment string

const (
	ASEnvironmentSandbox    ASEnvironment = "Sandbox"
	ASEnvironmentProduction ASEnvironment = "Production"
)

// ASStatus represents the subscription status.
type ASStatus int32

const (
	ASStatusActive       ASStatus = 1
	ASStatusExpired      ASStatus = 2
	ASStatusBillingRetry ASStatus = 3
	ASStatusGracePeriod  ASStatus = 4
	ASStatusRevoked      ASStatus = 5
)

// --- V1 Notification Types ---

// ASNotificationTypeV1 represents a V1 notification type.
type ASNotificationTypeV1 string

const (
	ASNotificationTypeV1InitialBuy              ASNotificationTypeV1 = "INITIAL_BUY"
	ASNotificationTypeV1Cancel                  ASNotificationTypeV1 = "CANCEL"
	ASNotificationTypeV1DidChangeRenewalPref    ASNotificationTypeV1 = "DID_CHANGE_RENEWAL_PREF"
	ASNotificationTypeV1DidChangeRenewalStatus  ASNotificationTypeV1 = "DID_CHANGE_RENEWAL_STATUS"
	ASNotificationTypeV1DidFailToRenew          ASNotificationTypeV1 = "DID_FAIL_TO_RENEW"
	ASNotificationTypeV1DidRecover              ASNotificationTypeV1 = "DID_RECOVER"
	ASNotificationTypeV1DidRenew                ASNotificationTypeV1 = "DID_RENEW"
	ASNotificationTypeV1InteractiveRenewal      ASNotificationTypeV1 = "INTERACTIVE_RENEWAL"
	ASNotificationTypeV1PriceIncreaseConsent    ASNotificationTypeV1 = "PRICE_INCREASE_CONSENT"
	ASNotificationTypeV1Refund                  ASNotificationTypeV1 = "REFUND"
	ASNotificationTypeV1Revoke                  ASNotificationTypeV1 = "REVOKE"
	ASNotificationTypeV1ConsumptionRequest      ASNotificationTypeV1 = "CONSUMPTION_REQUEST"
)

// --- V2 Notification Types ---

// ASNotificationType represents a V2 notification type.
type ASNotificationType string

const (
	ASNotificationTypeSubscribed            ASNotificationType = "SUBSCRIBED"
	ASNotificationTypeDidRenew              ASNotificationType = "DID_RENEW"
	ASNotificationTypeDidFailToRenew        ASNotificationType = "DID_FAIL_TO_RENEW"
	ASNotificationTypeDidChangeRenewalPref  ASNotificationType = "DID_CHANGE_RENEWAL_PREF"
	ASNotificationTypeDidChangeRenewalStat  ASNotificationType = "DID_CHANGE_RENEWAL_STATUS"
	ASNotificationTypeExpired               ASNotificationType = "EXPIRED"
	ASNotificationTypeGracePeriodExpired    ASNotificationType = "GRACE_PERIOD_EXPIRED"
	ASNotificationTypePriceIncrease         ASNotificationType = "PRICE_INCREASE"
	ASNotificationTypeOfferRedeemed         ASNotificationType = "OFFER_REDEEMED"
	ASNotificationTypeRefund                ASNotificationType = "REFUND"
	ASNotificationTypeRefundDeclined        ASNotificationType = "REFUND_DECLINED"
	ASNotificationTypeRefundReversed        ASNotificationType = "REFUND_REVERSED"
	ASNotificationTypeConsumptionRequest    ASNotificationType = "CONSUMPTION_REQUEST"
	ASNotificationTypeRenewalExtended       ASNotificationType = "RENEWAL_EXTENDED"
	ASNotificationTypeRenewalExtension      ASNotificationType = "RENEWAL_EXTENSION"
	ASNotificationTypeRevoke                ASNotificationType = "REVOKE"
	ASNotificationTypeExternalPurchaseToken ASNotificationType = "EXTERNAL_PURCHASE_TOKEN"
	ASNotificationTypeOneTimeCharge         ASNotificationType = "ONE_TIME_CHARGE"
	ASNotificationTypeTest                  ASNotificationType = "TEST"
)

// ASNotificationSubtype represents a V2 notification subtype.
type ASNotificationSubtype string

const (
	ASSubtypeInitialBuy         ASNotificationSubtype = "INITIAL_BUY"
	ASSubtypeResubscribe        ASNotificationSubtype = "RESUBSCRIBE"
	ASSubtypeDowngrade          ASNotificationSubtype = "DOWNGRADE"
	ASSubtypeUpgrade            ASNotificationSubtype = "UPGRADE"
	ASSubtypeAutoRenewEnabled   ASNotificationSubtype = "AUTO_RENEW_ENABLED"
	ASSubtypeAutoRenewDisabled  ASNotificationSubtype = "AUTO_RENEW_DISABLED"
	ASSubtypeVoluntary          ASNotificationSubtype = "VOLUNTARY"
	ASSubtypeBillingRetry       ASNotificationSubtype = "BILLING_RETRY"
	ASSubtypePriceIncrease      ASNotificationSubtype = "PRICE_INCREASE"
	ASSubtypeAccepted           ASNotificationSubtype = "ACCEPTED"
	ASSubtypePending            ASNotificationSubtype = "PENDING"
	ASSubtypeBillingRecovery    ASNotificationSubtype = "BILLING_RECOVERY"
	ASSubtypeProductNotForSale  ASNotificationSubtype = "PRODUCT_NOT_FOR_SALE"
	ASSubtypeFailure            ASNotificationSubtype = "FAILURE"
	ASSubtypeGracePeriod        ASNotificationSubtype = "GRACE_PERIOD"
	ASSubtypeSummary            ASNotificationSubtype = "SUMMARY"
	ASSubtypeUnreported         ASNotificationSubtype = "UNREPORTED"
)

// ASTransactionType represents the type of in-app purchase transaction.
type ASTransactionType string

const (
	ASTransactionTypeAutoRenewable  ASTransactionType = "Auto-Renewable Subscription"
	ASTransactionTypeNonConsumable  ASTransactionType = "Non-Consumable"
	ASTransactionTypeConsumable     ASTransactionType = "Consumable"
	ASTransactionTypeNonRenewing    ASTransactionType = "Non-Renewing Subscription"
)

// ASInAppOwnershipType represents the ownership type of an in-app purchase.
type ASInAppOwnershipType string

const (
	ASOwnershipTypePurchased   ASInAppOwnershipType = "PURCHASED"
	ASOwnershipTypeFamilyShared ASInAppOwnershipType = "FAMILY_SHARED"
)

// ASTransactionReason represents the reason for a transaction.
type ASTransactionReason string

const (
	ASTransactionReasonPurchase ASTransactionReason = "PURCHASE"
	ASTransactionReasonRenewal  ASTransactionReason = "RENEWAL"
)

// --- V1 Structs ---

// ASNotificationV1 represents a V1 App Store Server Notification.
type ASNotificationV1 struct {
	NotificationType    ASNotificationTypeV1 `json:"notification_type"`
	Password            string               `json:"password,omitempty"`
	Environment         ASEnvironment        `json:"environment,omitempty"`
	AutoRenewAdamID     string               `json:"auto_renew_adam_id,omitempty"`
	AutoRenewProductID  string               `json:"auto_renew_product_id,omitempty"`
	AutoRenewStatus     string               `json:"auto_renew_status,omitempty"`
	AutoRenewStatusDate string               `json:"auto_renew_status_change_date,omitempty"`
	AutoRenewStatusMs   string               `json:"auto_renew_status_change_date_ms,omitempty"`
	AutoRenewStatusPst  string               `json:"auto_renew_status_change_date_pst,omitempty"`
	UnifiedReceipt      *ASUnifiedReceipt    `json:"unified_receipt,omitempty"`
	BID                 string               `json:"bid,omitempty"`
	BVRS                string               `json:"bvrs,omitempty"`
}

// ASUnifiedReceipt represents the unified receipt in a V1 notification.
type ASUnifiedReceipt struct {
	Status             int                    `json:"status"`
	Environment        ASEnvironment          `json:"environment,omitempty"`
	LatestReceipt      string                 `json:"latest_receipt,omitempty"`
	LatestReceiptInfo  []ASLatestReceiptInfo  `json:"latest_receipt_info,omitempty"`
	PendingRenewalInfo []ASPendingRenewalInfo `json:"pending_renewal_info,omitempty"`
}

// ASLatestReceiptInfo represents a transaction in the latest receipt info array.
type ASLatestReceiptInfo struct {
	AppAccountToken             string `json:"app_account_token,omitempty"`
	CancellationDate            string `json:"cancellation_date,omitempty"`
	CancellationDateMs          string `json:"cancellation_date_ms,omitempty"`
	CancellationDatePst         string `json:"cancellation_date_pst,omitempty"`
	CancellationReason          string `json:"cancellation_reason,omitempty"`
	ExpiresDate                 string `json:"expires_date,omitempty"`
	ExpiresDateMs               string `json:"expires_date_ms,omitempty"`
	ExpiresDatePst              string `json:"expires_date_pst,omitempty"`
	InAppOwnershipType          string `json:"in_app_ownership_type,omitempty"`
	IsInIntroOfferPeriod        string `json:"is_in_intro_offer_period,omitempty"`
	IsTrialPeriod               string `json:"is_trial_period,omitempty"`
	IsUpgraded                  string `json:"is_upgraded,omitempty"`
	OfferCodeRefName            string `json:"offer_code_ref_name,omitempty"`
	OriginalPurchaseDate        string `json:"original_purchase_date,omitempty"`
	OriginalPurchaseDateMs      string `json:"original_purchase_date_ms,omitempty"`
	OriginalPurchaseDatePst     string `json:"original_purchase_date_pst,omitempty"`
	OriginalTransactionID       string `json:"original_transaction_id,omitempty"`
	ProductID                   string `json:"product_id,omitempty"`
	PromotionalOfferID          string `json:"promotional_offer_id,omitempty"`
	PurchaseDate                string `json:"purchase_date,omitempty"`
	PurchaseDateMs              string `json:"purchase_date_ms,omitempty"`
	PurchaseDatePst             string `json:"purchase_date_pst,omitempty"`
	Quantity                    string `json:"quantity,omitempty"`
	SubscriptionGroupIdentifier string `json:"subscription_group_identifier,omitempty"`
	TransactionID               string `json:"transaction_id,omitempty"`
	WebOrderLineItemID          string `json:"web_order_line_item_id,omitempty"`
}

// ASPendingRenewalInfo represents a pending renewal in a V1 notification.
type ASPendingRenewalInfo struct {
	AutoRenewProductID          string `json:"auto_renew_product_id,omitempty"`
	AutoRenewStatus             string `json:"auto_renew_status,omitempty"`
	ExpirationIntent            string `json:"expiration_intent,omitempty"`
	GracePeriodExpiresDate      string `json:"grace_period_expires_date,omitempty"`
	GracePeriodExpiresDateMs    string `json:"grace_period_expires_date_ms,omitempty"`
	GracePeriodExpiresDatePst   string `json:"grace_period_expires_date_pst,omitempty"`
	IsInBillingRetryPeriod      string `json:"is_in_billing_retry_period,omitempty"`
	OfferCodeRefName            string `json:"offer_code_ref_name,omitempty"`
	OriginalTransactionID       string `json:"original_transaction_id,omitempty"`
	PriceConsentStatus          string `json:"price_consent_status,omitempty"`
	ProductID                   string `json:"product_id,omitempty"`
	PromotionalOfferID          string `json:"promotional_offer_id,omitempty"`
	SubscriptionGroupIdentifier string `json:"subscription_group_identifier,omitempty"`
}

// --- V2 Structs ---

// ASSignedPayload represents the outer envelope of a V2 notification.
type ASSignedPayload struct {
	SignedPayload string `json:"signedPayload"`
}

// ASNotificationV2 represents a decoded V2 App Store Server Notification.
type ASNotificationV2 struct {
	NotificationType     ASNotificationType     `json:"notificationType"`
	Subtype              ASNotificationSubtype  `json:"subtype,omitempty"`
	NotificationUUID     string                 `json:"notificationUUID"`
	Data                 *ASNotificationData    `json:"data,omitempty"`
	Summary              *ASNotificationSummary `json:"summary,omitempty"`
	ExternalPurchaseToken *ASExternalPurchaseToken `json:"externalPurchaseToken,omitempty"`
	Version              string                 `json:"version"`
	SignedDate           int64                  `json:"signedDate"`
}

// ASNotificationData represents the data field of a V2 notification.
type ASNotificationData struct {
	AppAppleID               int64         `json:"appAppleId,omitempty"`
	BundleID                 string        `json:"bundleId,omitempty"`
	BundleVersion            string        `json:"bundleVersion,omitempty"`
	Environment              ASEnvironment `json:"environment,omitempty"`
	SignedTransactionInfo    string        `json:"signedTransactionInfo,omitempty"`
	SignedRenewalInfo        string        `json:"signedRenewalInfo,omitempty"`
	Status                   ASStatus      `json:"status,omitempty"`
	ConsumptionRequestReason string        `json:"consumptionRequestReason,omitempty"`
}

// ASNotificationSummary represents the summary field for RENEWAL_EXTENSION/SUMMARY notifications.
type ASNotificationSummary struct {
	RequestIdentifier      string        `json:"requestIdentifier,omitempty"`
	Environment            ASEnvironment `json:"environment,omitempty"`
	AppAppleID             int64         `json:"appAppleId,omitempty"`
	BundleID               string        `json:"bundleId,omitempty"`
	ProductID              string        `json:"productId,omitempty"`
	StorefrontCountryCodes []string      `json:"storefrontCountryCodes,omitempty"`
	FailedCount            int64         `json:"failedCount,omitempty"`
	SucceededCount         int64         `json:"succeededCount,omitempty"`
}

// ASExternalPurchaseToken represents the external purchase token field.
type ASExternalPurchaseToken struct {
	ExternalPurchaseID string        `json:"externalPurchaseId,omitempty"`
	TokenCreationDate  int64         `json:"tokenCreationDate,omitempty"`
	AppAppleID         int64         `json:"appAppleId,omitempty"`
	BundleID           string        `json:"bundleId,omitempty"`
	Environment        ASEnvironment `json:"environment,omitempty"`
}

// ASTransactionInfo represents a decoded JWS transaction info.
type ASTransactionInfo struct {
	TransactionID               string               `json:"transactionId,omitempty"`
	OriginalTransactionID       string               `json:"originalTransactionId,omitempty"`
	WebOrderLineItemID          string               `json:"webOrderLineItemId,omitempty"`
	BundleID                    string               `json:"bundleId,omitempty"`
	ProductID                   string               `json:"productId,omitempty"`
	SubscriptionGroupIdentifier string               `json:"subscriptionGroupIdentifier,omitempty"`
	PurchaseDate                int64                `json:"purchaseDate,omitempty"`
	OriginalPurchaseDate        int64                `json:"originalPurchaseDate,omitempty"`
	ExpiresDate                 int64                `json:"expiresDate,omitempty"`
	Quantity                    int32                `json:"quantity,omitempty"`
	Type                        ASTransactionType    `json:"type,omitempty"`
	InAppOwnershipType          ASInAppOwnershipType `json:"inAppOwnershipType,omitempty"`
	SignedDate                  int64                `json:"signedDate,omitempty"`
	TransactionReason           ASTransactionReason  `json:"transactionReason,omitempty"`
	Storefront                  string               `json:"storefront,omitempty"`
	StorefrontID                string               `json:"storefrontId,omitempty"`
	Environment                 ASEnvironment        `json:"environment,omitempty"`
	Price                       int64                `json:"price,omitempty"`
	Currency                    string               `json:"currency,omitempty"`
	OfferIdentifier             string               `json:"offerIdentifier,omitempty"`
	OfferType                   int32                `json:"offerType,omitempty"`
	OfferDiscountType           string               `json:"offerDiscountType,omitempty"`
	AppAccountToken             string               `json:"appAccountToken,omitempty"`
	IsUpgraded                  bool                 `json:"isUpgraded,omitempty"`
	RevocationDate              int64                `json:"revocationDate,omitempty"`
	RevocationReason            int32                `json:"revocationReason,omitempty"`
}

// ASRenewalInfo represents decoded JWS renewal info.
type ASRenewalInfo struct {
	AutoRenewProductID          string        `json:"autoRenewProductId,omitempty"`
	AutoRenewStatus             int32         `json:"autoRenewStatus,omitempty"`
	Environment                 ASEnvironment `json:"environment,omitempty"`
	ExpirationIntent            int32         `json:"expirationIntent,omitempty"`
	GracePeriodExpiresDate      int64         `json:"gracePeriodExpiresDate,omitempty"`
	IsInBillingRetryPeriod      bool          `json:"isInBillingRetryPeriod,omitempty"`
	OfferIdentifier             string        `json:"offerIdentifier,omitempty"`
	OfferType                   int32         `json:"offerType,omitempty"`
	OriginalTransactionID       string        `json:"originalTransactionId,omitempty"`
	PriceIncreaseStatus         int32         `json:"priceIncreaseStatus,omitempty"`
	ProductID                   string        `json:"productId,omitempty"`
	RecentSubscriptionStartDate int64         `json:"recentSubscriptionStartDate,omitempty"`
	RenewalDate                 int64         `json:"renewalDate,omitempty"`
	RenewalPrice                int64         `json:"renewalPrice,omitempty"`
	Currency                    string        `json:"currency,omitempty"`
	SignedDate                  int64         `json:"signedDate,omitempty"`
	EligibleWinBackOfferIDs     []string      `json:"eligibleWinBackOfferIds,omitempty"`
}
