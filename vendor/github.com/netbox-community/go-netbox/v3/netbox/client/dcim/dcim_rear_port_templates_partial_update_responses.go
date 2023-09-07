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

// DcimRearPortTemplatesPartialUpdateReader is a Reader for the DcimRearPortTemplatesPartialUpdate structure.
type DcimRearPortTemplatesPartialUpdateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DcimRearPortTemplatesPartialUpdateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDcimRearPortTemplatesPartialUpdateOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDcimRearPortTemplatesPartialUpdateDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDcimRearPortTemplatesPartialUpdateOK creates a DcimRearPortTemplatesPartialUpdateOK with default headers values
func NewDcimRearPortTemplatesPartialUpdateOK() *DcimRearPortTemplatesPartialUpdateOK {
	return &DcimRearPortTemplatesPartialUpdateOK{}
}

/*
DcimRearPortTemplatesPartialUpdateOK describes a response with status code 200, with default header values.

DcimRearPortTemplatesPartialUpdateOK dcim rear port templates partial update o k
*/
type DcimRearPortTemplatesPartialUpdateOK struct {
	Payload *models.RearPortTemplate
}

// IsSuccess returns true when this dcim rear port templates partial update o k response has a 2xx status code
func (o *DcimRearPortTemplatesPartialUpdateOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this dcim rear port templates partial update o k response has a 3xx status code
func (o *DcimRearPortTemplatesPartialUpdateOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this dcim rear port templates partial update o k response has a 4xx status code
func (o *DcimRearPortTemplatesPartialUpdateOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this dcim rear port templates partial update o k response has a 5xx status code
func (o *DcimRearPortTemplatesPartialUpdateOK) IsServerError() bool {
	return false
}

// IsCode returns true when this dcim rear port templates partial update o k response a status code equal to that given
func (o *DcimRearPortTemplatesPartialUpdateOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the dcim rear port templates partial update o k response
func (o *DcimRearPortTemplatesPartialUpdateOK) Code() int {
	return 200
}

func (o *DcimRearPortTemplatesPartialUpdateOK) Error() string {
	return fmt.Sprintf("[PATCH /dcim/rear-port-templates/{id}/][%d] dcimRearPortTemplatesPartialUpdateOK  %+v", 200, o.Payload)
}

func (o *DcimRearPortTemplatesPartialUpdateOK) String() string {
	return fmt.Sprintf("[PATCH /dcim/rear-port-templates/{id}/][%d] dcimRearPortTemplatesPartialUpdateOK  %+v", 200, o.Payload)
}

func (o *DcimRearPortTemplatesPartialUpdateOK) GetPayload() *models.RearPortTemplate {
	return o.Payload
}

func (o *DcimRearPortTemplatesPartialUpdateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RearPortTemplate)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDcimRearPortTemplatesPartialUpdateDefault creates a DcimRearPortTemplatesPartialUpdateDefault with default headers values
func NewDcimRearPortTemplatesPartialUpdateDefault(code int) *DcimRearPortTemplatesPartialUpdateDefault {
	return &DcimRearPortTemplatesPartialUpdateDefault{
		_statusCode: code,
	}
}

/*
DcimRearPortTemplatesPartialUpdateDefault describes a response with status code -1, with default header values.

DcimRearPortTemplatesPartialUpdateDefault dcim rear port templates partial update default
*/
type DcimRearPortTemplatesPartialUpdateDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this dcim rear port templates partial update default response has a 2xx status code
func (o *DcimRearPortTemplatesPartialUpdateDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this dcim rear port templates partial update default response has a 3xx status code
func (o *DcimRearPortTemplatesPartialUpdateDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this dcim rear port templates partial update default response has a 4xx status code
func (o *DcimRearPortTemplatesPartialUpdateDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this dcim rear port templates partial update default response has a 5xx status code
func (o *DcimRearPortTemplatesPartialUpdateDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this dcim rear port templates partial update default response a status code equal to that given
func (o *DcimRearPortTemplatesPartialUpdateDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the dcim rear port templates partial update default response
func (o *DcimRearPortTemplatesPartialUpdateDefault) Code() int {
	return o._statusCode
}

func (o *DcimRearPortTemplatesPartialUpdateDefault) Error() string {
	return fmt.Sprintf("[PATCH /dcim/rear-port-templates/{id}/][%d] dcim_rear-port-templates_partial_update default  %+v", o._statusCode, o.Payload)
}

func (o *DcimRearPortTemplatesPartialUpdateDefault) String() string {
	return fmt.Sprintf("[PATCH /dcim/rear-port-templates/{id}/][%d] dcim_rear-port-templates_partial_update default  %+v", o._statusCode, o.Payload)
}

func (o *DcimRearPortTemplatesPartialUpdateDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *DcimRearPortTemplatesPartialUpdateDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
