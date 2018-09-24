package router

import "testing"

func TestSigningSecret(t *testing.T) {
	r, _ := New(SigningSecret("dummy secret"))
	if r.signingSecret != "dummy secret" {
		t.Error("unexpected signing secret:", r.signingSecret)
	}
}
