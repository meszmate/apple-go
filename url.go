package apple

import (
	"net/url"
	"strings"
)

type ResponseMode string

const (
	// Not Supported with id_token
	ResponseModeQuery ResponseMode = "query"

	ResponseModeFragment ResponseMode = "fragment"

	// Required mode if you requested any scope
	ResponseModeFormPost ResponseMode = "form_post"
)

type ResponseType string

const (
	ResponseTypeCode   ResponseType = "code"
	ResponseTypeCodeID ResponseType = "code id_token"
)

// AuthorizeURLConfig collects every optional parameter that can be put
// into the authorization request.
type AuthorizeURLConfig struct {
	// Required
	ClientID    string // "Services ID" (NOT the App ID)
	RedirectURI string

	// Optional
	State        string
	Scope        []string // {"email","name"} stb.
	Nonce        string
	ResponseMode ResponseMode // "form_post" | "fragment"
	ResponseType ResponseType // code or code and id_token
}

func AuthorizeURL(cfg AuthorizeURLConfig) string {
	u, _ := url.Parse("https://appleid.apple.com/auth/authorize")

	var responseMode string
	if cfg.ResponseMode != "" {
		responseMode = string(cfg.ResponseMode)
	} else {
		responseMode = string(ResponseModeFormPost)
	}

	var responseType string
	if cfg.ResponseType != "" {
		responseType = string(cfg.ResponseType)
	} else {
		responseType = string(ResponseTypeCodeID)
	}

	q := url.Values{}
	q.Add("response_type", responseType)
	q.Add("response_mode", responseMode)
	q.Add("client_id", cfg.ClientID)
	q.Add("redirect_uri", cfg.RedirectURI)

	if cfg.State != "" {
		q.Add("state", cfg.State)
	}
	if cfg.Nonce != "" {
		q.Add("nonce", cfg.Nonce)
	}
	if len(cfg.Scope) > 0 {
		q.Add("scope", strings.Join(cfg.Scope, " "))
	}

	u.RawQuery = q.Encode()
	return u.String()
}
