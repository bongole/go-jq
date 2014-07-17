package jq

import (
	"testing"
)

func TestJQ(t *testing.T) {
	jq := New("{ \"foo\": 1 }")

	err := jq.Search(".", func(d interface{}) {
		m := d.(map[string]interface{})
		if m["foo"].(float64) != 1 {
			t.Errorf("invalid result")
		}
	})

	err = jq.Search(".foo", func(d interface{}) {
		i := d.(float64)
		if i != 1 {
			t.Errorf("not a number")
		}
	})

	if err != nil {
		t.Error(err)
	}

	jq = New("{ foo: 2 }")
	err = jq.Search(".", func(d interface{}) {})

	if err == nil {
		t.Error(err)
	}

	jq = New("{ foo: 2 }")
	err = jq.Search(".bar", func(d interface{}) {})

	if err == nil {
		t.Error(err)
	}

}
