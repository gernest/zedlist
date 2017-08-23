// Package jwt implement JSON Web Token authentication.
package jwt

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/query"

	"github.com/jinzhu/now"
)

var (

	// TokenExpireDate is the expiry time for a generated token
	// defaults to at the end of the current year.
	TokenExpireDate = now.EndOfYear()
)

// RSAKeyHolder is an interface for public and private RSA keys.
type RSAKeyHolder interface {

	// GetPublicBytes returns rsa public key in []byte
	GetPublicBytes() []byte

	// GetPrivateBytes returns private rsa key in []byte
	GetPrivateBytes() []byte
}

// NewJWTAuth returns a JWTValidateFunc for the echo.s JWTAuth middleware.
func NewJWTAuth(keys RSAKeyHolder) func(string, jwt.SigningMethod) ([]byte, error) {
	return func(token string, method jwt.SigningMethod) ([]byte, error) {
		tk, err := query.GetTokenByKey(token)
		k := keys.GetPublicBytes()
		if err != nil {
			return k, err
		}
		if !tk.Valid {
			return k, errors.New("invarid token")
		}
		return k, nil
	}
}

// NewToken returns a new signed JWT token.
func NewToken(rsaKeys RSAKeyHolder, claims map[string]string) (*models.Token, error) {
	tk := &models.Token{}
	token := jwt.New(jwt.SigningMethodRS256)
	mc := make(jwt.MapClaims)
	for k, v := range claims {
		mc[k] = v
	}
	// add expiration
	mc["exp"] = TokenExpireDate.Unix()
	token.Claims = mc
	parsedKey := rsaKeys.GetPrivateBytes()
	rp, err := jwt.ParseRSAPrivateKeyFromPEM(parsedKey)
	if err != nil {
		return nil, err
	}
	tokenStr, err := token.SignedString(rp)
	if err != nil {
		return nil, err
	}
	tk.Key = tokenStr
	var clms []models.Claim
	for k, v := range claims {
		c := models.Claim{
			Key:   k,
			Value: v,
		}
		clms = append(clms, c)
	}
	tk.Claims = clms
	tk.ExpireOn = TokenExpireDate
	tk.Valid = true
	return tk, nil
}
