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

package dcim

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// DcimModuleBaysDeleteReader is a Reader for the DcimModuleBaysDelete structure.
type DcimModuleBaysDeleteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DcimModuleBaysDeleteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDcimModuleBaysDeleteNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDcimModuleBaysDeleteDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDcimModuleBaysDeleteNoContent creates a DcimModuleBaysDeleteNoContent with default headers values
func NewDcimModuleBaysDeleteNoContent() *DcimModuleBaysDeleteNoContent {
	return &DcimModuleBaysDeleteNoContent{}
}

/*
DcimModuleBaysDeleteNoContent describes a response with status code 204, with default header values.

DcimModuleBaysDeleteNoContent dcim module bays delete no content
*/
type DcimModuleBaysDeleteNoContent struct {
}

// IsSuccess returns true when this dcim module bays delete no content response has a 2xx status code
func (o *DcimModuleBaysDeleteNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this dcim module bays delete no content response has a 3xx status code
func (o *DcimModuleBaysDeleteNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this dcim module bays delete no content response has a 4xx status code
func (o *DcimModuleBaysDeleteNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this dcim module bays delete no content response has a 5xx status code
func (o *DcimModuleBaysDeleteNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this dcim module bays delete no content response a status code equal to that given
func (o *DcimModuleBaysDeleteNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the dcim module bays delete no content response
func (o *DcimModuleBaysDeleteNoContent) Code() int {
	return 204
}

func (o *DcimModuleBaysDeleteNoContent) Error() string {
	return fmt.Sprintf("[DELETE /dcim/module-bays/{id}/][%d] dcimModuleBaysDeleteNoContent ", 204)
}

func (o *DcimModuleBaysDeleteNoContent) String() string {
	return fmt.Sprintf("[DELETE /dcim/module-bays/{id}/][%d] dcimModuleBaysDeleteNoContent ", 204)
}

func (o *DcimModuleBaysDeleteNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDcimModuleBaysDeleteDefault creates a DcimModuleBaysDeleteDefault with default headers values
func NewDcimModuleBaysDeleteDefault(code int) *DcimModuleBaysDeleteDefault {
	return &DcimModuleBaysDeleteDefault{
		_statusCode: code,
	}
}

/*
DcimModuleBaysDeleteDefault describes a response with status code -1, with default header values.

DcimModuleBaysDeleteDefault dcim module bays delete default
*/
type DcimModuleBaysDeleteDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this dcim module bays delete default response has a 2xx status code
func (o *DcimModuleBaysDeleteDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this dcim module bays delete default response has a 3xx status code
func (o *DcimModuleBaysDeleteDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this dcim module bays delete default response has a 4xx status code
func (o *DcimModuleBaysDeleteDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this dcim module bays delete default response has a 5xx status code
func (o *DcimModuleBaysDeleteDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this dcim module bays delete default response a status code equal to that given
func (o *DcimModuleBaysDeleteDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the dcim module bays delete default response
func (o *DcimModuleBaysDeleteDefault) Code() int {
	return o._statusCode
}

func (o *DcimModuleBaysDeleteDefault) Error() string {
	return fmt.Sprintf("[DELETE /dcim/module-bays/{id}/][%d] dcim_module-bays_delete default  %+v", o._statusCode, o.Payload)
}

func (o *DcimModuleBaysDeleteDefault) String() string {
	return fmt.Sprintf("[DELETE /dcim/module-bays/{id}/][%d] dcim_module-bays_delete default  %+v", o._statusCode, o.Payload)
}

func (o *DcimModuleBaysDeleteDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *DcimModuleBaysDeleteDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}