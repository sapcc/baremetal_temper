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

// DcimConsolePortTemplatesBulkDeleteReader is a Reader for the DcimConsolePortTemplatesBulkDelete structure.
type DcimConsolePortTemplatesBulkDeleteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DcimConsolePortTemplatesBulkDeleteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDcimConsolePortTemplatesBulkDeleteNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDcimConsolePortTemplatesBulkDeleteDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDcimConsolePortTemplatesBulkDeleteNoContent creates a DcimConsolePortTemplatesBulkDeleteNoContent with default headers values
func NewDcimConsolePortTemplatesBulkDeleteNoContent() *DcimConsolePortTemplatesBulkDeleteNoContent {
	return &DcimConsolePortTemplatesBulkDeleteNoContent{}
}

/*
DcimConsolePortTemplatesBulkDeleteNoContent describes a response with status code 204, with default header values.

DcimConsolePortTemplatesBulkDeleteNoContent dcim console port templates bulk delete no content
*/
type DcimConsolePortTemplatesBulkDeleteNoContent struct {
}

// IsSuccess returns true when this dcim console port templates bulk delete no content response has a 2xx status code
func (o *DcimConsolePortTemplatesBulkDeleteNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this dcim console port templates bulk delete no content response has a 3xx status code
func (o *DcimConsolePortTemplatesBulkDeleteNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this dcim console port templates bulk delete no content response has a 4xx status code
func (o *DcimConsolePortTemplatesBulkDeleteNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this dcim console port templates bulk delete no content response has a 5xx status code
func (o *DcimConsolePortTemplatesBulkDeleteNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this dcim console port templates bulk delete no content response a status code equal to that given
func (o *DcimConsolePortTemplatesBulkDeleteNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the dcim console port templates bulk delete no content response
func (o *DcimConsolePortTemplatesBulkDeleteNoContent) Code() int {
	return 204
}

func (o *DcimConsolePortTemplatesBulkDeleteNoContent) Error() string {
	return fmt.Sprintf("[DELETE /dcim/console-port-templates/][%d] dcimConsolePortTemplatesBulkDeleteNoContent ", 204)
}

func (o *DcimConsolePortTemplatesBulkDeleteNoContent) String() string {
	return fmt.Sprintf("[DELETE /dcim/console-port-templates/][%d] dcimConsolePortTemplatesBulkDeleteNoContent ", 204)
}

func (o *DcimConsolePortTemplatesBulkDeleteNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDcimConsolePortTemplatesBulkDeleteDefault creates a DcimConsolePortTemplatesBulkDeleteDefault with default headers values
func NewDcimConsolePortTemplatesBulkDeleteDefault(code int) *DcimConsolePortTemplatesBulkDeleteDefault {
	return &DcimConsolePortTemplatesBulkDeleteDefault{
		_statusCode: code,
	}
}

/*
DcimConsolePortTemplatesBulkDeleteDefault describes a response with status code -1, with default header values.

DcimConsolePortTemplatesBulkDeleteDefault dcim console port templates bulk delete default
*/
type DcimConsolePortTemplatesBulkDeleteDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this dcim console port templates bulk delete default response has a 2xx status code
func (o *DcimConsolePortTemplatesBulkDeleteDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this dcim console port templates bulk delete default response has a 3xx status code
func (o *DcimConsolePortTemplatesBulkDeleteDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this dcim console port templates bulk delete default response has a 4xx status code
func (o *DcimConsolePortTemplatesBulkDeleteDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this dcim console port templates bulk delete default response has a 5xx status code
func (o *DcimConsolePortTemplatesBulkDeleteDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this dcim console port templates bulk delete default response a status code equal to that given
func (o *DcimConsolePortTemplatesBulkDeleteDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the dcim console port templates bulk delete default response
func (o *DcimConsolePortTemplatesBulkDeleteDefault) Code() int {
	return o._statusCode
}

func (o *DcimConsolePortTemplatesBulkDeleteDefault) Error() string {
	return fmt.Sprintf("[DELETE /dcim/console-port-templates/][%d] dcim_console-port-templates_bulk_delete default  %+v", o._statusCode, o.Payload)
}

func (o *DcimConsolePortTemplatesBulkDeleteDefault) String() string {
	return fmt.Sprintf("[DELETE /dcim/console-port-templates/][%d] dcim_console-port-templates_bulk_delete default  %+v", o._statusCode, o.Payload)
}

func (o *DcimConsolePortTemplatesBulkDeleteDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *DcimConsolePortTemplatesBulkDeleteDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
