package apple

import "fmt"

// CKErrorCode represents a CloudKit error code.
type CKErrorCode string

const (
	CKErrorAccessDenied           CKErrorCode = "ACCESS_DENIED"
	CKErrorAtomicError            CKErrorCode = "ATOMIC_ERROR"
	CKErrorAuthenticationFailed   CKErrorCode = "AUTHENTICATION_FAILED"
	CKErrorAuthenticationRequired CKErrorCode = "AUTHENTICATION_REQUIRED"
	CKErrorBadRequest             CKErrorCode = "BAD_REQUEST"
	CKErrorConflict               CKErrorCode = "CONFLICT"
	CKErrorExists                 CKErrorCode = "EXISTS"
	CKErrorInternalError          CKErrorCode = "INTERNAL_ERROR"
	CKErrorNotFound               CKErrorCode = "NOT_FOUND"
	CKErrorQuotaExceeded          CKErrorCode = "QUOTA_EXCEEDED"
	CKErrorThrottled              CKErrorCode = "THROTTLED"
	CKErrorTryAgainLater          CKErrorCode = "TRY_AGAIN_LATER"
	CKErrorZoneNotFound           CKErrorCode = "ZONE_NOT_FOUND"
	CKErrorUnknownError           CKErrorCode = "UNKNOWN_ERROR"
	CKErrorSignatureError         CKErrorCode = "SIGNATURE_ERROR"
	CKErrorNetworkError           CKErrorCode = "NETWORK_ERROR"
)

// CKError represents a CloudKit API error.
type CKError struct {
	Code            CKErrorCode `json:"ckErrorCode,omitempty"`
	Reason          string      `json:"reason,omitempty"`
	ServerErrorCode string      `json:"serverErrorCode,omitempty"`
	RetryAfter      int         `json:"retryAfter,omitempty"`
	UUID            string      `json:"uuid,omitempty"`
}

// Error implements the error interface.
func (e *CKError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("cloudkit: %s: %s", e.Code, e.Reason)
	}
	return fmt.Sprintf("cloudkit: %s", e.Code)
}

// CKErrorResponse represents the top-level error response from the CloudKit API.
type CKErrorResponse struct {
	UUID            string      `json:"uuid,omitempty"`
	ServerErrorCode CKErrorCode `json:"serverErrorCode,omitempty"`
	Reason          string      `json:"reason,omitempty"`
	RetryAfter      int         `json:"retryAfter,omitempty"`
	Records         []struct {
		ServerErrorCode CKErrorCode `json:"serverErrorCode,omitempty"`
		Reason          string      `json:"reason,omitempty"`
		RetryAfter      int         `json:"retryAfter,omitempty"`
	} `json:"records,omitempty"`
}
