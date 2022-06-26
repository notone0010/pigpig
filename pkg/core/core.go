// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package core

import (
	"net/http"

	"encoding/json"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/log"
)

type RawHttpHandler func(w http.ResponseWriter, r *http.Request)

// ErrResponse defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponse(w http.ResponseWriter, r *http.Request, err error, data interface{}) {
	var msg []byte
	if err != nil {
		log.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		msg, _ = json.Marshal(ErrResponse{
			Code:      coder.Code(),
			Message:   coder.String(),
			Reference: coder.Reference(),
		})
		w.WriteHeader(coder.HTTPStatus())
		_, _ = w.Write(msg)
		return
	}
	if data != nil{
		msg, _ = json.Marshal(data)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(msg)
}


// func GetErrorResponse(w http.ResponseWriter, r *http.Request, err error, data interface{}) {
// 	var msg []byte
// 	if err != nil {
// 		log.Errorf("%#+v", err)
// 		coder := errors.ParseCoder(err)
// 		msg, _ = json.Marshal(ErrResponse{
// 			Code:      coder.Code(),
// 			Message:   coder.String(),
// 			Reference: coder.Reference(),
// 		})
// 		w.WriteHeader(coder.HTTPStatus())
// 		_, _ = w.Write(msg)
// 		return
// 	}
// 	if data != nil{
// 		msg, _ = json.Marshal(data)
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	_, _ = w.Write(msg)
// }