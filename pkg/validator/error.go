package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type (
	Errors map[string]error
)

var (
	ErrRequired     = errors.New("is required")
	ErrLatLong      = errors.New("invalid lat long")
	ErrIPError      = errors.New("invalid IP address")
	ErrInvalidValue = errors.New("invalid value")
)

func MergeError(errs ...error) error {
	m := make(Errors)

	for _, err := range errs {
		if err != nil {
			for k, v := range err.(Errors) {
				m[k] = v
			}
		}

	}

	if len(m) == 0 {
		return nil
	}

	return m
}

func (es Errors) Error() string {
	if len(es) == 0 {
		return ""
	}

	keys := []string{}
	for key := range es {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var s strings.Builder
	for i, key := range keys {
		if i > 0 {
			s.WriteString("; ")
		}
		if errs, ok := es[key].(Errors); ok {
			fmt.Fprintf(&s, "%v: (%v)", key, errs)
		} else {
			fmt.Fprintf(&s, "%v: %v", key, es[key].Error())
		}
	}
	s.WriteString(".")
	return s.String()
}

func (es Errors) MarshalJSON() ([]byte, error) {
	errs := map[string]interface{}{}
	for key, err := range es {
		if ms, ok := err.(json.Marshaler); ok {
			errs[key] = ms
		} else {
			errs[key] = err.Error()
		}
	}
	return json.Marshal(errs)
}
