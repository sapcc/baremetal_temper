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

// ExtrasExportTemplatesBulkDeleteReader is a Reader for the ExtrasExportTemplatesBulkDelete structure.
type ExtrasExportTemplatesBulkDeleteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ExtrasExportTemplatesBulkDeleteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewExtrasExportTemplatesBulkDeleteNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewExtrasExportTemplatesBulkDeleteDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewExtrasExportTemplatesBulkDeleteNoContent creates a ExtrasExportTemplatesBulkDeleteNoContent with default headers values
func NewExtrasExportTemplatesBulkDeleteNoContent() *ExtrasExportTemplatesBulkDeleteNoContent {
	return &ExtrasExportTemplatesBulkDeleteNoContent{}
}

/*
ExtrasExportTemplatesBulkDeleteNoContent describes a response with status code 204, with default header values.

ExtrasExportTemplatesBulkDeleteNoContent extras export templates bulk delete no content
*/
type ExtrasExportTemplatesBulkDeleteNoContent struct {
}

// IsSuccess returns true when this extras export templates bulk delete no content response has a 2xx status code
func (o *ExtrasExportTemplatesBulkDeleteNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this extras export templates bulk delete no content response has a 3xx status code
func (o *ExtrasExportTemplatesBulkDeleteNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this extras export templates bulk delete no content response has a 4xx status code
func (o *ExtrasExportTemplatesBulkDeleteNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this extras export templates bulk delete no content response has a 5xx status code
func (o *ExtrasExportTemplatesBulkDeleteNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this extras export templates bulk delete no content response a status code equal to that given
func (o *ExtrasExportTemplatesBulkDeleteNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the extras export templates bulk delete no content response
func (o *ExtrasExportTemplatesBulkDeleteNoContent) Code() int {
	return 204
}

func (o *ExtrasExportTemplatesBulkDeleteNoContent) Error() string {
	return fmt.Sprintf("[DELETE /extras/export-templates/][%d] extrasExportTemplatesBulkDeleteNoContent ", 204)
}

func (o *ExtrasExportTemplatesBulkDeleteNoContent) String() string {
	return fmt.Sprintf("[DELETE /extras/export-templates/][%d] extrasExportTemplatesBulkDeleteNoContent ", 204)
}

func (o *ExtrasExportTemplatesBulkDeleteNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewExtrasExportTemplatesBulkDeleteDefault creates a ExtrasExportTemplatesBulkDeleteDefault with default headers values
func NewExtrasExportTemplatesBulkDeleteDefault(code int) *ExtrasExportTemplatesBulkDeleteDefault {
	return &ExtrasExportTemplatesBulkDeleteDefault{
		_statusCode: code,
	}
}

/*
ExtrasExportTemplatesBulkDeleteDefault describes a response with status code -1, with default header values.

ExtrasExportTemplatesBulkDeleteDefault extras export templates bulk delete default
*/
type ExtrasExportTemplatesBulkDeleteDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this extras export templates bulk delete default response has a 2xx status code
func (o *ExtrasExportTemplatesBulkDeleteDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this extras export templates bulk delete default response has a 3xx status code
func (o *ExtrasExportTemplatesBulkDeleteDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this extras export templates bulk delete default response has a 4xx status code
func (o *ExtrasExportTemplatesBulkDeleteDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this extras export templates bulk delete default response has a 5xx status code
func (o *ExtrasExportTemplatesBulkDeleteDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this extras export templates bulk delete default response a status code equal to that given
func (o *ExtrasExportTemplatesBulkDeleteDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the extras export templates bulk delete default response
func (o *ExtrasExportTemplatesBulkDeleteDefault) Code() int {
	return o._statusCode
}

func (o *ExtrasExportTemplatesBulkDeleteDefault) Error() string {
	return fmt.Sprintf("[DELETE /extras/export-templates/][%d] extras_export-templates_bulk_delete default  %+v", o._statusCode, o.Payload)
}

func (o *ExtrasExportTemplatesBulkDeleteDefault) String() string {
	return fmt.Sprintf("[DELETE /extras/export-templates/][%d] extras_export-templates_bulk_delete default  %+v", o._statusCode, o.Payload)
}

func (o *ExtrasExportTemplatesBulkDeleteDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *ExtrasExportTemplatesBulkDeleteDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
