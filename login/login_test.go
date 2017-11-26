package main

import "testing"

func TestRun(t *testing.T) {
	doc := aesEn("127.0.0.1:8080")
	t.Logf("doc decode %s", string(doc))
	t.Logf("doc encode %s", aesDe(doc))
}
