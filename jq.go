package jq

// #cgo LDFLAGS: -ljq
// #include <jq.h>
// typedef struct jv_parser jv_parser_struct;
import "C"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"
)

const (
	BUFSIZE = 4096
)

type JQ struct {
	src string
}

func New(src string) *JQ {
	jq := new(JQ)
	jq.src = src

	return jq
}

func toJson(str string) interface{} {
	var data interface{}
	json.Unmarshal([]byte(fmt.Sprintf("[%s]", str)), &data)
	result := data.([]interface{})
	return result[0]
}

func (self *JQ) Search(pattern string) (r interface{}, err error) {
	jq_state := C.jq_init()
	compiled := C.jq_compile(jq_state, C.CString(pattern))

	if compiled != 1 {
		C.jq_teardown(&jq_state)
		err = errors.New("compile error")
		return
	}

	jv_parser := C.jv_parser_new(0)

	processed := false
	src := bytes.NewBufferString(self.src)
	buf := make([]byte, BUFSIZE)
	for n, _ := src.Read(buf); n > 0; n, _ = src.Read(buf) {
		C.jv_parser_set_buf(jv_parser, (*C.char)(unsafe.Pointer(&buf[0])), C.int(n), 1)
		var value C.jv
		for value = C.jv_parser_next(jv_parser); C.jv_is_valid(value) != 0 && processed == false; value = C.jv_parser_next(jv_parser) {
			C.jq_start(jq_state, value, 0)

			var result C.jv
			for result = C.jq_next(jq_state); C.jv_is_valid(result) != 0 && processed == false; result = C.jq_next(jq_state) {
				dumped := C.jv_dump_string(result, 0)
				gostring := C.GoString(C.jv_string_value(dumped))
				r = toJson(gostring)
				processed = true
			}

			C.jv_free(result)
		}

		if C.jv_invalid_has_msg(C.jv_copy(value)) != 0 {
			msg := C.jv_invalid_get_msg(value)
			gomsg := C.GoString(C.jv_string_value(msg))
			C.jv_free(msg)
			err = errors.New(gomsg)
			goto end
		} else {
			C.jv_free(value)
		}
	}

end:

	C.jv_parser_free(jv_parser)
	C.jq_teardown(&jq_state)

	return
}
