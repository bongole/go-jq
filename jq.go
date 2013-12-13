package jq

// #cgo LDFLAGS: -ljq
// #include <jq.h>
// typedef struct jv_parser jv_parser_struct;
import "C"

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
//	"runtime"
	"unsafe"
)

const (
	BUFSIZE = 4096
)

type JQ struct {
	jq_state  *C.jq_state
	jv_parser *C.jv_parser_struct
}

type handler func(interface{})

func New(pattern string) (*JQ, error) {
	jq := new(JQ)
	jq_state := C.jq_init()
	compiled := C.jq_compile(jq_state, C.CString(pattern))

	if compiled != 1 {
		C.jq_teardown(&jq_state)
		return nil, errors.New("compile error")
	}

	jq.jq_state = jq_state
	jq.jv_parser = nil

    /*
	runtime.SetFinalizer(&jq, func(x *JQ) {
		if x.jq_state != nil {
			C.jq_teardown(&x.jq_state)
		}

		if x.jv_parser != nil {
			C.jv_parser_free(x.jv_parser)
		}
	})
    */

	return jq, nil
}

func toJson(str string, block handler) {
	var data interface{}
	json.Unmarshal([]byte(fmt.Sprintf("[%s]", str)), &data)
	result := data.([]interface{})
	block(result[0])
}

func (self *JQ) Search(src io.Reader, block handler) error {
	if self.jv_parser == nil {
		self.jv_parser = C.jv_parser_new(0)
	}

	buf := make([]byte, BUFSIZE)
	for n, _ := src.Read(buf); n > 0; {
		C.jv_parser_set_buf(self.jv_parser, (*C.char)(unsafe.Pointer(&buf[0])), C.int(n), 1)
		var value C.jv
		for value = C.jv_parser_next(self.jv_parser); C.jv_is_valid(value) != 0; value = C.jv_parser_next(self.jv_parser) {
			C.jq_start(self.jq_state, value, 0)

			var result C.jv
			for result = C.jq_next(self.jq_state); C.jv_is_valid(result) != 0; result = C.jq_next(self.jq_state) {

				dumped := C.jv_dump_string(result, 0)
				gostring := C.GoString(C.jv_string_value(dumped))
				toJson(gostring, block)
			}

			C.jv_free(result)
		}

		if C.jv_invalid_has_msg(C.jv_copy(value)) != 0 {
			msg := C.jv_invalid_get_msg(value)
			gomsg := C.GoString(C.jv_string_value(msg))
			C.jv_free(msg)
			C.jv_free(value)
			return errors.New(gomsg)
		} else {
			C.jv_free(value)
		}
	}

	if self.jv_parser != nil {
        C.jv_parser_free(self.jv_parser);
        self.jv_parser = nil
	}

	return nil
}
