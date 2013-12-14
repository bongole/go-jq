package jq

import (
	"testing"
)

func TestJQ(t *testing.T) {
	jq := New("{ \"foo\": 1 }")

	d, err := jq.Search(".")
	m := d.(map[string]interface{})
	if m["foo"].(float64) != 1 {
		t.Errorf("invalid result")
	}

	d, err = jq.Search(".foo")
	i := d.(float64)
	if i != 1 {
		t.Errorf("not a number")
	}

	if err != nil {
		t.Error(err)
	}

	jq = New("{ foo: 2 }")
	d, err = jq.Search(".")

	if err == nil {
		t.Error(err)
	}
}
