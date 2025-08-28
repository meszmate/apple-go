package apple

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetUserInfoFromIDToken(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZW1haWwiOiJhbmVtYWlsQHlvdXJkb21haW4iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiaXNfcHJpdmF0ZV9lbWFpbCI6ZmFsc2UsInJlYWxfdXNlcl9zdGF0dXMiOjIsImlhdCI6MTUxNjIzOTAyMn0.M9kQlH4Ybi3-iEYZ3TxdouU8Gjgt2IyfMvATsnFisP4"
	expectedUser := &AppleUser{
		Subject:        "1234567890",
		Email:          "anemail@yourdomain",
		EmailVerified:  true,
		IsPrivateEmail: false,
		RealUserStatus: RealUserStatusLikelyReal,
		IssuedAt:       time.Unix(1516239022, 0),
	}

	au, err := GetUserInfoFromIDToken(jwt)
	assert.Equal(t, nil, err)
	if err == nil && au != nil {
		if ok := reflect.DeepEqual(*expectedUser, *au); !ok {
			t.Errorf("Expected user is different from received. Expected user is %+v, received user is %+v", expectedUser, au)
		}
	}
}

func TestGetUserInfoFromIDToken_DecodeFail(t *testing.T) {
	jwt := "this-is-definetly-not-a-jwt"
	_, err := GetUserInfoFromIDToken(jwt)
	assert.NotEqual(t, nil, err)
}
