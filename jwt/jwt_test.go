package jwt_test

import (
	"bytes"
	"testing"

	"github.com/zhixian001/jwt-cracker/jwt"
)

func TestParseJWT(t *testing.T) {
	testJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIyMDE5MDczMSIsIm5hbWUiOiJLb25nYWwiLCJpYXQiOjE1NjQ1NzA4MDB9.GfXFIG2fNzUKOC050zUdniguhAEAY1UCeSBHy9VJF28"

	err, parsedJwt := jwt.ParseJWT(testJWT)

	if err != nil {
		t.Errorf("JWT parse failed.\n")
	}

	t.Log(*parsedJwt.Header)
	t.Log(parsedJwt.Payload)
	t.Log(bytes.NewBuffer(parsedJwt.Message).String())
	t.Log(bytes.NewBuffer(parsedJwt.Signature).String())
}
