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

	"github.com/netbox-community/go-netbox/v3/netbox/models"
)

// DcimModuleBayTemplatesCreateReader is a Reader for the DcimModuleBayTemplatesCreate structure.
type DcimModuleBayTemplatesCreateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DcimModuleBayTemplatesCreateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewDcimModuleBayTemplatesCreateCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDcimModuleBayTemplatesCreateDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDcimModuleBayTemplatesCreateCreated creates a DcimModuleBayTemplatesCreateCreated with default headers values
func NewDcimModuleBayTemplatesCreateCreated() *DcimModuleBayTemplatesCreateCreated {
	return &DcimModuleBayTemplatesCreateCreated{}
}

/*
DcimModuleBayTemplatesCreateCreated describes a response with status code 201, with default header values.

DcimModuleBayTemplatesCreateCreated dcim module bay templates create created
*/
type DcimModuleBayTemplatesCreateCreated struct {
	Payload *models.ModuleBayTemplate
}

// IsSuccess returns true when this dcim module bay templates create created response has a 2xx status code
func (o *DcimModuleBayTemplatesCreateCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this dcim module bay templates create created response has a 3xx status code
func (o *DcimModuleBayTemplatesCreateCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this dcim module bay templates create created response has a 4xx status code
func (o *DcimModuleBayTemplatesCreateCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this dcim module bay templates create created response has a 5xx status code
func (o *DcimModuleBayTemplatesCreateCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this dcim module bay templates create created response a status code equal to that given
func (o *DcimModuleBayTemplatesCreateCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the dcim module bay templates create created response
func (o *DcimModuleBayTemplatesCreateCreated) Code() int {
	return 201
}

func (o *DcimModuleBayTemplatesCreateCreated) Error() string {
	return fmt.Sprintf("[POST /dcim/module-bay-templates/][%d] dcimModuleBayTemplatesCreateCreated  %+v", 201, o.Payload)
}

func (o *DcimModuleBayTemplatesCreateCreated) String() string {
	return fmt.Sprintf("[POST /dcim/module-bay-templates/][%d] dcimModuleBayTemplatesCreateCreated  %+v", 201, o.Payload)
}

func (o *DcimModuleBayTemplatesCreateCreated) GetPayload() *models.ModuleBayTemplate {
	return o.Payload
}

func (o *DcimModuleBayTemplatesCreateCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ModuleBayTemplate)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDcimModuleBayTemplatesCreateDefault creates a DcimModuleBayTemplatesCreateDefault with default headers values
func NewDcimModuleBayTemplatesCreateDefault(code int) *DcimModuleBayTemplatesCreateDefault {
	return &DcimModuleBayTemplatesCreateDefault{
		_statusCode: code,
	}
}

/*
DcimModuleBayTemplatesCreateDefault describes a response with status code -1, with default header values.

DcimModuleBayTemplatesCreateDefault dcim module bay templates create default
*/
type DcimModuleBayTemplatesCreateDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this dcim module bay templates create default response has a 2xx status code
func (o *DcimModuleBayTemplatesCreateDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this dcim module bay templates create default response has a 3xx status code
func (o *DcimModuleBayTemplatesCreateDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this dcim module bay templates create default response has a 4xx status code
func (o *DcimModuleBayTemplatesCreateDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this dcim module bay templates create default response has a 5xx status code
func (o *DcimModuleBayTemplatesCreateDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this dcim module bay templates create default response a status code equal to that given
func (o *DcimModuleBayTemplatesCreateDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the dcim module bay templates create default response
func (o *DcimModuleBayTemplatesCreateDefault) Code() int {
	return o._statusCode
}

func (o *DcimModuleBayTemplatesCreateDefault) Error() string {
	return fmt.Sprintf("[POST /dcim/module-bay-templates/][%d] dcim_module-bay-templates_create default  %+v", o._statusCode, o.Payload)
}

func (o *DcimModuleBayTemplatesCreateDefault) String() string {
	return fmt.Sprintf("[POST /dcim/module-bay-templates/][%d] dcim_module-bay-templates_create default  %+v", o._statusCode, o.Payload)
}

func (o *DcimModuleBayTemplatesCreateDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *DcimModuleBayTemplatesCreateDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}