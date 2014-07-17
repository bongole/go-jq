package jq

/*
#cgo LDFLAGS: -ljq
#include <jq.h>
typedef struct jv_parser jv_parser_struct;
extern void jq_err_callback(void *,jv);
*/
import "C"

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"unsafe"
)

const (
	BUFSIZE = 4096
)

type handler func(interface{})

type JQ struct {
	src interface{}
    err error
}

//export jq_err_callback
func jq_err_callback(p unsafe.Pointer, value C.jv) {
    self := (*JQ)(p)
    self.err = errors.New(C.GoString(C.jv_string_value(value)))
}

func New(src interface{}) *JQ {
	jq := new(JQ)
	jq.src = src
    jq.err = nil

	return jq
}

func toJson(str string) interface{} {
	var data interface{}
	json.Unmarshal([]byte(fmt.Sprintf("[%s]", str)), &data)
	result := data.([]interface{})
	return result[0]
}

func (self *JQ) Search(pattern string, block handler) (err error) {
	var src io.ReadSeeker
	switch t := self.src.(type) {
	case string:
		src = strings.NewReader(t)
	case io.ReadSeeker:
		src = t
	default:
		err = errors.New("src is not io.ReadSeeker or string")
		return
	}

	jq_state := C.jq_init()
    C.jq_set_error_cb(jq_state, (C.jq_err_cb)(C.jq_err_callback), unsafe.Pointer(self))

	compiled := C.jq_compile(jq_state, C.CString(pattern))

	if compiled != 1 {
		C.jq_teardown(&jq_state)
		err = errors.New("compile error")
		return
	}
	jv_parser := C.jv_parser_new(0)

	src.Seek(0, 0)
	buf := make([]byte, BUFSIZE)
	for n, _ := src.Read(buf); n > 0; n, _ = src.Read(buf) {
		C.jv_parser_set_buf(jv_parser, (*C.char)(unsafe.Pointer(&buf[0])), C.int(n), 1)
		var value C.jv
		for value = C.jv_parser_next(jv_parser); C.jv_is_valid(value) != 0; value = C.jv_parser_next(jv_parser) {
			C.jq_start(jq_state, value, 0)

			var result C.jv
			for result = C.jq_next(jq_state); C.jv_is_valid(result) != 0; result = C.jq_next(jq_state) {
				dumped := C.jv_dump_string(result, 0)
				gostring := C.GoString(C.jv_string_value(dumped))
				block(toJson(gostring))
			}

			C.jv_free(result)
		}

		if C.jv_invalid_has_msg(C.jv_copy(value)) != 0 {
			msg := C.jv_invalid_get_msg(value)
			self.err = errors.New(C.GoString(C.jv_string_value(msg)))
			C.jv_free(msg)
		} else {
			C.jv_free(value)
		}
	}

    if self.err != nil {
        err = self.err
    }

	C.jv_parser_free(jv_parser)
	C.jq_teardown(&jq_state)

	return
}
