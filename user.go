package apple

import (
	"time"

	"github.com/tideland/gorest/jwt"
)

var (
	// RealUserStatusUnsupported unsupported, only works in iOS >= 14.
	RealUserStatusUnsupported RealUserStatus = 0
	// RealUserStatusUnknown cannot determine if the user is real.
	RealUserStatusUnknown RealUserStatus = 1
	// RealUserStatusLikelyReal user is likely real.
	RealUserStatusLikelyReal RealUserStatus = 2
)

// RealUserStatus an integer value that indicates whether the user appears to be
// a real person.
type RealUserStatus int

// AppleUser is the model to hold information about the user.
type AppleUser struct {
	// Standard JWT claims
	Issuer   string    `json:"iss"`   // Issuer (always https://appleid.apple.com)
	Audience string    `json:"aud"`   // Your app's client ID
	Subject  string    `json:"sub"`   // Unique user identifier (UID)
	IssuedAt time.Time `json:"iat"`   // Token issue time
	Expiry   time.Time `json:"exp"`   // Token expiration time
	Nonce    string    `json:"nonce"` // Replayed nonce if provided in request

	// Apple-specific claims
	Email          string         `json:"email"`            // User's email (real or proxy)
	EmailVerified  bool           `json:"email_verified"`   // Whether email is verified
	IsPrivateEmail bool           `json:"is_private_email"` // Whether email is Apple's private relay
	RealUserStatus RealUserStatus `json:"real_user_status"` // Fraud detection status

	// Optional Apple claims
	AuthTime       *time.Time `json:"auth_time"`       // Time of user authentication
	NonceSupported *bool      `json:"nonce_supported"` // Whether nonce is supported
	TransferSub    string     `json:"transfer_sub"`    // App transfer identifier
	OrgID          string     `json:"org_id"`          // Organization ID (for managed accounts)
}

// GetUserInfoFromIDToken retrieves the user info from the JWT id token.
// It maps every documented Apple id_token claim into the AppleUser struct.
func GetUserInfoFromIDToken(idToken string) (*AppleUser, error) {
	token, err := jwt.Decode(idToken)
	if err != nil {
		return nil, err
	}

	u := AppleUser{}
	claims := token.Claims()

	/* ---------- standard JWT claims ---------- */
	if v, ok := claims["iss"].(string); ok {
		u.Issuer = v
	}
	if v, ok := claims["aud"].(string); ok {
		u.Audience = v
	}
	if v, ok := claims["sub"].(string); ok {
		u.Subject = v
	}
	if v, ok := claims["nonce"].(string); ok {
		u.Nonce = v
	}
	if v, ok := claims["iat"].(float64); ok {
		u.IssuedAt = time.Unix(int64(v), 0)
	}
	if v, ok := claims["exp"].(float64); ok {
		u.Expiry = time.Unix(int64(v), 0)
	}

	/* ---------- Apple-specific claims ---------- */
	if v, ok := claims["email"].(string); ok {
		u.Email = v
	}
	u.EmailVerified = parseBool(claims, "email_verified")
	u.IsPrivateEmail = parseBool(claims, "is_private_email")

	/* real_user_status */
	switch v := claims["real_user_status"].(type) {
	case float64:
		switch RealUserStatus(int(v)) {
		case RealUserStatusLikelyReal:
			u.RealUserStatus = RealUserStatusLikelyReal
		case RealUserStatusUnknown:
			u.RealUserStatus = RealUserStatusUnknown
		default:
			u.RealUserStatus = RealUserStatusUnsupported
		}
	case int:
		switch RealUserStatus(v) {
		case RealUserStatusLikelyReal:
			u.RealUserStatus = RealUserStatusLikelyReal
		case RealUserStatusUnknown:
			u.RealUserStatus = RealUserStatusUnknown
		default:
			u.RealUserStatus = RealUserStatusUnsupported
		}
	default:
		u.RealUserStatus = RealUserStatusUnsupported
	}

	/* optional Apple claims */
	if v, ok := claims["auth_time"].(float64); ok {
		t := time.Unix(int64(v), 0)
		u.AuthTime = &t
	}
	if v, ok := claims["nonce_supported"].(bool); ok {
		u.NonceSupported = &v
	}
	if v, ok := claims["transfer_sub"].(string); ok {
		u.TransferSub = v
	}
	if v, ok := claims["org_id"].(string); ok {
		u.OrgID = v
	}

	return &u, nil
}

/*
helper: JSON numbers arrive as float64, and Apple sometimes

	sends "true"/"false" strings for booleans.
*/
func parseBool(claims map[string]interface{}, key string) bool {
	switch v := claims[key].(type) {
	case bool:
		return v
	case string:
		return v == "true"
	case float64:
		return v != 0
	default:
		return false
	}
}
