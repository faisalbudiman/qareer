package utils

import (
	"encoding/json"
	"net/http"
)

type httpSuccess struct {
	Data interface{} `json:"data"`
}

type httpError struct {
	Error httpErrorData `json:"error"`
}

type httpErrorData struct {
	Code    int   `json:"code"`
	Message error `json:"message"`
}

type HTTPJSONResponse struct {
}

// Resp ia a success json response
func (h *HTTPJSONResponse) Resp(w http.ResponseWriter, code int, data interface{}) {
	b, err := json.Marshal(httpSuccess{
		Data: data,
	})
	if err != nil {
		h.Err(w, code, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

// Err is a error json response
func (h *HTTPJSONResponse) Err(w http.ResponseWriter, code int, err error) {
	b, err := json.Marshal(httpError{
		Error: httpErrorData{
			Code:    code,
			Message: err,
		},
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}
