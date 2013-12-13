package jq

import (
    "testing"
    "bytes"
)

func TestJQ(t *testing.T) {
    jq, _ := New(".")
    json := bytes.NewBufferString("{ hoge: 1 }")
    jq.Search(json, func(d interface{}){
        t.Errorf("%v", d)
    })
}
