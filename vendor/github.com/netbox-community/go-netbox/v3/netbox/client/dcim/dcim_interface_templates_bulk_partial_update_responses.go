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

// DcimInterfaceTemplatesBulkPartialUpdateReader is a Reader for the DcimInterfaceTemplatesBulkPartialUpdate structure.
type DcimInterfaceTemplatesBulkPartialUpdateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DcimInterfaceTemplatesBulkPartialUpdateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDcimInterfaceTemplatesBulkPartialUpdateOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDcimInterfaceTemplatesBulkPartialUpdateDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDcimInterfaceTemplatesBulkPartialUpdateOK creates a DcimInterfaceTemplatesBulkPartialUpdateOK with default headers values
func NewDcimInterfaceTemplatesBulkPartialUpdateOK() *DcimInterfaceTemplatesBulkPartialUpdateOK {
	return &DcimInterfaceTemplatesBulkPartialUpdateOK{}
}

/*
DcimInterfaceTemplatesBulkPartialUpdateOK describes a response with status code 200, with default header values.

DcimInterfaceTemplatesBulkPartialUpdateOK dcim interface templates bulk partial update o k
*/
type DcimInterfaceTemplatesBulkPartialUpdateOK struct {
	Payload *models.InterfaceTemplate
}

// IsSuccess returns true when this dcim interface templates bulk partial update o k response has a 2xx status code
func (o *DcimInterfaceTemplatesBulkPartialUpdateOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this dcim interface templates bulk partial update o k response has a 3xx status code
func (o *DcimInterfaceTemplatesBulkPartialUpdateOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this dcim interface templates bulk partial update o k response has a 4xx status code
func (o *DcimInterfaceTemplatesBulkPartialUpdateOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this dcim interface templates bulk partial update o k response has a 5xx status code
func (o *DcimInterfaceTemplatesBulkPartialUpdateOK) IsServerError() bool {
	return false
}

// IsCode returns true when this dcim interface templates bulk partial update o k response a status code equal to that given
func (o *DcimInterfaceTemplatesBulkPartialUpdateOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the dcim interface templates bulk partial update o k response
func (o *DcimInterfaceTemplatesBulkPartialUpdateOK) Code() int {
	return 200
}

func (o *DcimInterfaceTemplatesBulkPartialUpdateOK) Error() string {
	return fmt.Sprintf("[PATCH /dcim/interface-templates/][%d] dcimInterfaceTemplatesBulkPartialUpdateOK  %+v", 200, o.Payload)
}

func (o *DcimInterfaceTemplatesBulkPartialUpdateOK) String() string {
	return fmt.Sprintf("[PATCH /dcim/interface-templates/][%d] dcimInterfaceTemplatesBulkPartialUpdateOK  %+v", 200, o.Payload)
}

func (o *DcimInterfaceTemplatesBulkPartialUpdateOK) GetPayload() *models.InterfaceTemplate {
	return o.Payload
}

func (o *DcimInterfaceTemplatesBulkPartialUpdateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.InterfaceTemplate)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDcimInterfaceTemplatesBulkPartialUpdateDefault creates a DcimInterfaceTemplatesBulkPartialUpdateDefault with default headers values
func NewDcimInterfaceTemplatesBulkPartialUpdateDefault(code int) *DcimInterfaceTemplatesBulkPartialUpdateDefault {
	return &DcimInterfaceTemplatesBulkPartialUpdateDefault{
		_statusCode: code,
	}
}

/*
DcimInterfaceTemplatesBulkPartialUpdateDefault describes a response with status code -1, with default header values.

DcimInterfaceTemplatesBulkPartialUpdateDefault dcim interface templates bulk partial update default
*/
type DcimInterfaceTemplatesBulkPartialUpdateDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this dcim interface templates bulk partial update default response has a 2xx status code
func (o *DcimInterfaceTemplatesBulkPartialUpdateDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this dcim interface templates bulk partial update default response has a 3xx status code
func (o *DcimInterfaceTemplatesBulkPartialUpdateDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this dcim interface templates bulk partial update default response has a 4xx status code
func (o *DcimInterfaceTemplatesBulkPartialUpdateDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this dcim interface templates bulk partial update default response has a 5xx status code
func (o *DcimInterfaceTemplatesBulkPartialUpdateDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this dcim interface templates bulk partial update default response a status code equal to that given
func (o *DcimInterfaceTemplatesBulkPartialUpdateDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the dcim interface templates bulk partial update default response
func (o *DcimInterfaceTemplatesBulkPartialUpdateDefault) Code() int {
	return o._statusCode
}

func (o *DcimInterfaceTemplatesBulkPartialUpdateDefault) Error() string {
	return fmt.Sprintf("[PATCH /dcim/interface-templates/][%d] dcim_interface-templates_bulk_partial_update default  %+v", o._statusCode, o.Payload)
}

func (o *DcimInterfaceTemplatesBulkPartialUpdateDefault) String() string {
	return fmt.Sprintf("[PATCH /dcim/interface-templates/][%d] dcim_interface-templates_bulk_partial_update default  %+v", o._statusCode, o.Payload)
}

func (o *DcimInterfaceTemplatesBulkPartialUpdateDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *DcimInterfaceTemplatesBulkPartialUpdateDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
