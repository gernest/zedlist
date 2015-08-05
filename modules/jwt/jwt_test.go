package jwt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gernest/zedlist/modules/query"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type JWTKeys struct {
	RSAKeyHolder
	pub  []byte
	priv []byte
}

func NewJWTKeys(priv, pub []byte) *JWTKeys {
	return &JWTKeys{pub: pub, priv: priv}
}

func (j *JWTKeys) GetPublicBytes() []byte {
	return j.pub
}

func (j *JWTKeys) GetPrivateBytes() []byte {
	return j.priv
}

func home(ctx *echo.Context) error {
	v := fmt.Sprintf("%v", ctx.Get("_claims"))
	return ctx.String(http.StatusOK, v)
}

func TestNewJWTAuth(t *testing.T) {
	e := echo.New()
	privateKey, err := ioutil.ReadFile("test/sample_key")
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := ioutil.ReadFile("test/sample_key.pub")
	if err != nil {
		t.Fatal(err)
	}
	holder := NewJWTKeys(privateKey, publicKey)
	e.Use(middleware.JWTAuth(NewJWTAuth(holder)))
	e.Get("/", home)

	claims := map[string]string{
		"user":       "gernest",
		"occupation": "dreamer",
	}
	tok, err := NewToken(holder, claims)
	if err != nil {
		t.Error(err)
	}
	err = query.Create(tok)
	if err != nil {
		t.Error(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tok.Key))

	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	if !strings.Contains(w.Body.String(), "dreamer") {
		t.Errorf("expected %s to contain dreamer", w.Body.String())
	}
}
