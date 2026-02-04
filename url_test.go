package apple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizeURLDefaults(t *testing.T) {
	cfg := AuthorizeURLConfig{
		ClientID:    "com.example.app",
		RedirectURI: "https://example.com/callback",
	}

	u := AuthorizeURL(cfg)
	assert.Contains(t, u, "response_type=code+id_token")
	assert.Contains(t, u, "response_mode=form_post")
	assert.Contains(t, u, "client_id=com.example.app")
	assert.Contains(t, u, "redirect_uri=https")
}

func TestAuthorizeURLCustomResponseType(t *testing.T) {
	cfg := AuthorizeURLConfig{
		ClientID:     "com.example.app",
		RedirectURI:  "https://example.com/callback",
		ResponseType: ResponseTypeCode,
		ResponseMode: ResponseModeFragment,
	}

	u := AuthorizeURL(cfg)
	assert.Contains(t, u, "response_type=code")
	assert.Contains(t, u, "response_mode=fragment")
}

func TestAuthorizeURLWithScopeAndState(t *testing.T) {
	cfg := AuthorizeURLConfig{
		ClientID:    "com.example.app",
		RedirectURI: "https://example.com/callback",
		Scope:       []string{"email", "name"},
		State:       "csrf-token",
		Nonce:       "nonce-abc",
	}

	u := AuthorizeURL(cfg)
	assert.Contains(t, u, "scope=email+name")
	assert.Contains(t, u, "state=csrf-token")
	assert.Contains(t, u, "nonce=nonce-abc")
}

func TestAuthorizeURLBugFixResponseTypeNotOverwritingResponseMode(t *testing.T) {
	// This tests the bug fix: when ResponseType is empty, the default
	// should set responseType (not responseMode).
	cfg := AuthorizeURLConfig{
		ClientID:    "com.example.app",
		RedirectURI: "https://example.com/callback",
	}

	u := AuthorizeURL(cfg)
	// response_mode should be form_post (default), not overwritten by ResponseTypeCodeID
	assert.Contains(t, u, "response_mode=form_post")
	// response_type should be "code id_token" (default)
	assert.Contains(t, u, "response_type=code+id_token")
}
