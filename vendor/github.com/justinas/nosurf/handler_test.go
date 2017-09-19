package nosurf

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefaultFailureHandler(t *testing.T) {
	writer := httptest.NewRecorder()
	req := dummyGet()

	defaultFailureHandler(writer, req)

	if writer.Code != FailureCode {
		t.Errorf("Wrong status code for defaultFailure Handler: "+
			"expected %d, got %d", FailureCode, writer.Code)
	}
}

func TestRegenerateToken(t *testing.T) {
	hand := New(nil)
	writer := httptest.NewRecorder()

	req := dummyGet()
	token := hand.RegenerateToken(writer, req)

	header := writer.Header().Get("Set-Cookie")
	expected_part := fmt.Sprintf("csrf_token=%s;", token)

	if !strings.Contains(header, expected_part) {
		t.Errorf("Expected header to contain %v, it doesn't. The header is %v.",
			expected_part, header)
	}

}

// Kind of a duplication of TestRegenerateToken,
// but it's still good to test this too.
func TestsetTokenCookie(t *testing.T) {
	hand := New(nil)

	writer := httptest.NewRecorder()
	req := dummyGet()

	token := "dummy"
	hand.setTokenCookie(writer, req, token)

	header := writer.Header().Get("Set-Cookie")
	expected_part := fmt.Sprintf("csrf_token=%s;", token)

	if !strings.Contains(header, expected_part) {
		t.Errorf("Expected header to contain %v, it doesn't. The header is %v.",
			expected_part, header)
	}

	tokenInContext := Token(req)
	if tokenInContext != token {
		t.Errorf("RegenerateToken didn't set the token in the context map!"+
			" Expected %v, got %v", token, tokenInContext)
	}
}

func TestSafeMethodsPass(t *testing.T) {
	handler := New(http.HandlerFunc(succHand))

	for _, method := range safeMethods {
		req, err := http.NewRequest(method, "http://dummy.us", nil)

		if err != nil {
			t.Fatal(err)
		}

		writer := httptest.NewRecorder()
		handler.ServeHTTP(writer, req)

		expected := 200

		if writer.Code != expected {
			t.Errorf("A safe method didn't pass the CSRF check."+
				"Expected HTTP status %d, got %d", expected, writer.Code)
		}

		writer.Flush()
	}
}

func TestExemptedPass(t *testing.T) {
	handler := New(http.HandlerFunc(succHand))
	handler.ExemptPath("/faq")

	req, err := http.NewRequest("POST", "http://dummy.us/faq", strings.NewReader("a=b"))
	if err != nil {
		t.Fatal(err)
	}

	writer := httptest.NewRecorder()
	handler.ServeHTTP(writer, req)

	expected := 200

	if writer.Code != expected {
		t.Errorf("An exempted URL didn't pass the CSRF check."+
			"Expected HTTP status %d, got %d", expected, writer.Code)
	}

	writer.Flush()
}

// Tests that the token/reason context is accessible
// in the success/failure handlers
func TestContextIsAccessible(t *testing.T) {
	// case 1: success
	succHand := func(w http.ResponseWriter, r *http.Request) {
		token := Token(r)
		if token == "" {
			t.Errorf("Token is inaccessible in the success handler")
		}
	}

	hand := New(http.HandlerFunc(succHand))

	// we need a request that passes. Let's just use a safe method for that.
	req := dummyGet()
	writer := httptest.NewRecorder()

	hand.ServeHTTP(writer, req)

	// I'll do the failure case when there is actual logic for failures
}

func TestEmptyRefererFails(t *testing.T) {
	hand := New(http.HandlerFunc(succHand))
	fhand := correctReason(t, ErrNoReferer)
	hand.SetFailureHandler(fhand)

	req, err := http.NewRequest("POST", "https://dummy.us/", strings.NewReader("a=b"))
	if err != nil {
		t.Fatal(err)
	}
	writer := httptest.NewRecorder()

	hand.ServeHTTP(writer, req)

	if writer.Code != FailureCode {
		t.Errorf("A POST request with no Referer should have failed with the code %d, but it didn't.",
			writer.Code)
	}
}

func TestDifferentOriginRefererFails(t *testing.T) {
	hand := New(http.HandlerFunc(succHand))
	fhand := correctReason(t, ErrBadReferer)
	hand.SetFailureHandler(fhand)

	req, err := http.NewRequest("POST", "https://dummy.us/", strings.NewReader("a=b"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Referer", "http://attack-on-golang.com")
	writer := httptest.NewRecorder()

	hand.ServeHTTP(writer, req)

	if writer.Code != FailureCode {
		t.Errorf("A POST request with a Referer from a different origin"+
			"should have failed with the code %d, but it didn't.", writer.Code)
	}
}

func TestNoTokenFails(t *testing.T) {
	hand := New(http.HandlerFunc(succHand))
	fhand := correctReason(t, ErrBadToken)
	hand.SetFailureHandler(fhand)

	vals := [][]string{
		{"name", "Jolene"},
	}

	req, err := http.NewRequest("POST", "http://dummy.us", formBodyR(vals))
	if err != nil {
		panic(err)
	}
	writer := httptest.NewRecorder()

	hand.ServeHTTP(writer, req)

	if writer.Code != FailureCode {
		t.Errorf("The check should've failed with the code %d, but instead, it"+
			" returned code %d", FailureCode, writer.Code)
	}
}

func TestWrongTokenFails(t *testing.T) {
	hand := New(http.HandlerFunc(succHand))
	fhand := correctReason(t, ErrBadToken)
	hand.SetFailureHandler(fhand)

	vals := [][]string{
		{"name", "Jolene"},
		// this won't EVER be a valid value with the current scheme
		{FormFieldName, "$#%^&"},
	}

	req, err := http.NewRequest("POST", "http://dummy.us", formBodyR(vals))
	if err != nil {
		panic(err)
	}
	writer := httptest.NewRecorder()

	hand.ServeHTTP(writer, req)

	if writer.Code != FailureCode {
		t.Errorf("The check should've failed with the code %d, but instead, it"+
			" returned code %d", FailureCode, writer.Code)
	}
}

// For this and similar tests we start a test server
// Since it's much easier to get the cookie
// from a normal http.Response than from the recorder
func TestCorrectTokenPasses(t *testing.T) {
	hand := New(http.HandlerFunc(succHand))

	server := httptest.NewServer(hand)
	defer server.Close()

	// issue the first request to get the token
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	cookie := getRespCookie(resp, CookieName)
	if cookie == nil {
		t.Fatal("Cookie was not found in the response.")
	}

	vals := [][]string{
		{"name", "Jolene"},
		{FormFieldName, cookie.Value},
	}

	// Constructing a custom request is suffering
	req, err := http.NewRequest("POST", server.URL, formBodyR(vals))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("The request should have succeeded, but it didn't. Instead, the code was %d",
			resp.StatusCode)
	}
}

func TestPrefersHeaderOverFormValue(t *testing.T) {
	// Let's do a nice trick to find out this:
	// We'll set the correct token in the header
	// And a wrong one in the form.
	// That way, if it succeeds,
	// it will mean that it prefered the header.

	hand := New(http.HandlerFunc(succHand))

	server := httptest.NewServer(hand)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	cookie := getRespCookie(resp, CookieName)
	if cookie == nil {
		t.Fatal("Cookie was not found in the response.")
	}

	vals := [][]string{
		{"name", "Jolene"},
		{FormFieldName, "a very wrong value"},
	}

	req, err := http.NewRequest("POST", server.URL, formBodyR(vals))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(HeaderName, cookie.Value)
	req.AddCookie(cookie)

	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("The request should have succeeded, but it didn't. Instead, the code was %d",
			resp.StatusCode)
	}
}
