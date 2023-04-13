// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2020 The go-netbox Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package extras

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// ExtrasCustomLinksBulkDeleteReader is a Reader for the ExtrasCustomLinksBulkDelete structure.
type ExtrasCustomLinksBulkDeleteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ExtrasCustomLinksBulkDeleteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewExtrasCustomLinksBulkDeleteNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewExtrasCustomLinksBulkDeleteDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewExtrasCustomLinksBulkDeleteNoContent creates a ExtrasCustomLinksBulkDeleteNoContent with default headers values
func NewExtrasCustomLinksBulkDeleteNoContent() *ExtrasCustomLinksBulkDeleteNoContent {
	return &ExtrasCustomLinksBulkDeleteNoContent{}
}

/*
ExtrasCustomLinksBulkDeleteNoContent describes a response with status code 204, with default header values.

ExtrasCustomLinksBulkDeleteNoContent extras custom links bulk delete no content
*/
type ExtrasCustomLinksBulkDeleteNoContent struct {
}

// IsSuccess returns true when this extras custom links bulk delete no content response has a 2xx status code
func (o *ExtrasCustomLinksBulkDeleteNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this extras custom links bulk delete no content response has a 3xx status code
func (o *ExtrasCustomLinksBulkDeleteNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this extras custom links bulk delete no content response has a 4xx status code
func (o *ExtrasCustomLinksBulkDeleteNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this extras custom links bulk delete no content response has a 5xx status code
func (o *ExtrasCustomLinksBulkDeleteNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this extras custom links bulk delete no content response a status code equal to that given
func (o *ExtrasCustomLinksBulkDeleteNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the extras custom links bulk delete no content response
func (o *ExtrasCustomLinksBulkDeleteNoContent) Code() int {
	return 204
}

func (o *ExtrasCustomLinksBulkDeleteNoContent) Error() string {
	return fmt.Sprintf("[DELETE /extras/custom-links/][%d] extrasCustomLinksBulkDeleteNoContent ", 204)
}

func (o *ExtrasCustomLinksBulkDeleteNoContent) String() string {
	return fmt.Sprintf("[DELETE /extras/custom-links/][%d] extrasCustomLinksBulkDeleteNoContent ", 204)
}

func (o *ExtrasCustomLinksBulkDeleteNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewExtrasCustomLinksBulkDeleteDefault creates a ExtrasCustomLinksBulkDeleteDefault with default headers values
func NewExtrasCustomLinksBulkDeleteDefault(code int) *ExtrasCustomLinksBulkDeleteDefault {
	return &ExtrasCustomLinksBulkDeleteDefault{
		_statusCode: code,
	}
}

/*
ExtrasCustomLinksBulkDeleteDefault describes a response with status code -1, with default header values.

ExtrasCustomLinksBulkDeleteDefault extras custom links bulk delete default
*/
type ExtrasCustomLinksBulkDeleteDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this extras custom links bulk delete default response has a 2xx status code
func (o *ExtrasCustomLinksBulkDeleteDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this extras custom links bulk delete default response has a 3xx status code
func (o *ExtrasCustomLinksBulkDeleteDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this extras custom links bulk delete default response has a 4xx status code
func (o *ExtrasCustomLinksBulkDeleteDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this extras custom links bulk delete default response has a 5xx status code
func (o *ExtrasCustomLinksBulkDeleteDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this extras custom links bulk delete default response a status code equal to that given
func (o *ExtrasCustomLinksBulkDeleteDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the extras custom links bulk delete default response
func (o *ExtrasCustomLinksBulkDeleteDefault) Code() int {
	return o._statusCode
}

func (o *ExtrasCustomLinksBulkDeleteDefault) Error() string {
	return fmt.Sprintf("[DELETE /extras/custom-links/][%d] extras_custom-links_bulk_delete default  %+v", o._statusCode, o.Payload)
}

func (o *ExtrasCustomLinksBulkDeleteDefault) String() string {
	return fmt.Sprintf("[DELETE /extras/custom-links/][%d] extras_custom-links_bulk_delete default  %+v", o._statusCode, o.Payload)
}

func (o *ExtrasCustomLinksBulkDeleteDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *ExtrasCustomLinksBulkDeleteDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
