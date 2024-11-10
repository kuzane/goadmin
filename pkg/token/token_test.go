package token

import (
    "testing"
)

func TestGetRole(t *testing.T) {
    tk := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjkwNzAyOTZ9.H_ZRXNvyXceQ1A4CbX_yzBQYsiZcwAWPKNXE1WtAgDg"
    jtc := "VUCVLHABTXGJ6GAA675XE5WN4DDWZWUGBD6OWB36PJJH2LRG35RA===="
    token, err := parseToken(tk, jtc)
    if err != nil {
        t.Fatal(err)
    }
    t.Log(token)
}
