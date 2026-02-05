package apple

import "fmt"

// ASErrorCode represents an App Store notification error code.
type ASErrorCode string

const (
	ASErrorInvalidPayload     ASErrorCode = "INVALID_PAYLOAD"
	ASErrorInvalidJWS         ASErrorCode = "INVALID_JWS"
	ASErrorCertificateInvalid ASErrorCode = "CERTIFICATE_INVALID"
	ASErrorSignatureInvalid   ASErrorCode = "SIGNATURE_INVALID"
	ASErrorUnsupportedAlgo    ASErrorCode = "UNSUPPORTED_ALGORITHM"
	ASErrorInvalidCertChain   ASErrorCode = "INVALID_CERT_CHAIN"
	ASErrorDecodeError        ASErrorCode = "DECODE_ERROR"
)

// ASError represents an App Store notification processing error.
type ASError struct {
	Code   ASErrorCode `json:"code,omitempty"`
	Reason string      `json:"reason,omitempty"`
}

// Error implements the error interface.
func (e *ASError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("appstore: %s: %s", e.Code, e.Reason)
	}
	return fmt.Sprintf("appstore: %s", e.Code)
}

// ASAPIErrorCode represents an App Store Server API error code.
type ASAPIErrorCode int

const (
	ASAPIErrorAccountNotFound                        ASAPIErrorCode = 4040001
	ASAPIErrorAccountNotFoundRetryable               ASAPIErrorCode = 4040002
	ASAPIErrorAppNotFound                            ASAPIErrorCode = 4040003
	ASAPIErrorAppNotFoundRetryable                   ASAPIErrorCode = 4040004
	ASAPIErrorInvalidRequest                         ASAPIErrorCode = 4000000
	ASAPIErrorInvalidEmptyStorefrontCountryCodeList  ASAPIErrorCode = 4000002
	ASAPIErrorInvalidStorefrontCountryCode           ASAPIErrorCode = 4000003
	ASAPIErrorOriginalTransactionIDNotFound          ASAPIErrorCode = 4040005
	ASAPIErrorOriginalTransactionIDNotFoundRetryable ASAPIErrorCode = 4040006
	ASAPIErrorServerNotificationURLNotFound          ASAPIErrorCode = 4040007
	ASAPIErrorTestNotificationNotFound               ASAPIErrorCode = 4040008
	ASAPIErrorStatusRequestNotFound                  ASAPIErrorCode = 4040009
	ASAPIErrorTransactionNotFound                    ASAPIErrorCode = 4040010
	ASAPIErrorRateLimitExceeded                      ASAPIErrorCode = 4290000
	ASAPIErrorGeneralInternal                        ASAPIErrorCode = 5000000
	ASAPIErrorGeneralInternalRetryable               ASAPIErrorCode = 5000001
	ASAPIErrorInvalidExtendByDays                    ASAPIErrorCode = 4000014
	ASAPIErrorInvalidExtendReasonCode                ASAPIErrorCode = 4000015
	ASAPIErrorInvalidRequestIdentifier               ASAPIErrorCode = 4000016
	ASAPIErrorSubscriptionExtensionIneligible        ASAPIErrorCode = 4030004
	ASAPIErrorSubscriptionMaxExtension               ASAPIErrorCode = 4030005
	ASAPIErrorFamilySharedSubscriptionExtIneligible  ASAPIErrorCode = 4030007
	ASAPIErrorInvalidStartDate                       ASAPIErrorCode = 4000006
	ASAPIErrorInvalidEndDate                         ASAPIErrorCode = 4000007
	ASAPIErrorInvalidNotificationType                ASAPIErrorCode = 4000008
	ASAPIErrorMultipleFiltersSupplied                ASAPIErrorCode = 4000009
	ASAPIErrorInvalidSort                            ASAPIErrorCode = 4000010
	ASAPIErrorInvalidProductType                     ASAPIErrorCode = 4000011
	ASAPIErrorInvalidProductID                       ASAPIErrorCode = 4000012
	ASAPIErrorInvalidSubscriptionGroupID             ASAPIErrorCode = 4000013
	ASAPIErrorInvalidInAppOwnershipType              ASAPIErrorCode = 4000026
	ASAPIErrorInvalidRevoked                         ASAPIErrorCode = 4000030
	ASAPIErrorInvalidPaginationToken                 ASAPIErrorCode = 4000005
	ASAPIErrorInvalidConsumptionStatus               ASAPIErrorCode = 4000017
	ASAPIErrorInvalidPlatform                        ASAPIErrorCode = 4000018
	ASAPIErrorInvalidPlayTime                        ASAPIErrorCode = 4000019
	ASAPIErrorInvalidSampleContentProvided           ASAPIErrorCode = 4000020
	ASAPIErrorInvalidDeliveryStatus                  ASAPIErrorCode = 4000021
	ASAPIErrorInvalidAppAccountToken                 ASAPIErrorCode = 4000022
	ASAPIErrorInvalidAccountTenure                   ASAPIErrorCode = 4000023
	ASAPIErrorInvalidLifetimeDollarsPurchased        ASAPIErrorCode = 4000024
	ASAPIErrorInvalidLifetimeDollarsRefunded         ASAPIErrorCode = 4000025
	ASAPIErrorInvalidUserStatus                      ASAPIErrorCode = 4000027
	ASAPIErrorInvalidRefundPreference                ASAPIErrorCode = 4000028
)

// ASAPIError represents an error returned by the App Store Server API.
type ASAPIError struct {
	ErrorCode    ASAPIErrorCode `json:"errorCode"`
	ErrorMessage string         `json:"errorMessage"`
	HTTPStatus   int            `json:"-"`
	RetryAfter   int            `json:"-"`
}

// Error implements the error interface.
func (e *ASAPIError) Error() string {
	if e.ErrorMessage != "" {
		return fmt.Sprintf("appstore api: %d: %s", e.ErrorCode, e.ErrorMessage)
	}
	return fmt.Sprintf("appstore api: %d", e.ErrorCode)
}
