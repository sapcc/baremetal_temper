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

// DcimRegionsBulkUpdateReader is a Reader for the DcimRegionsBulkUpdate structure.
type DcimRegionsBulkUpdateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DcimRegionsBulkUpdateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDcimRegionsBulkUpdateOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDcimRegionsBulkUpdateDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDcimRegionsBulkUpdateOK creates a DcimRegionsBulkUpdateOK with default headers values
func NewDcimRegionsBulkUpdateOK() *DcimRegionsBulkUpdateOK {
	return &DcimRegionsBulkUpdateOK{}
}

/*
DcimRegionsBulkUpdateOK describes a response with status code 200, with default header values.

DcimRegionsBulkUpdateOK dcim regions bulk update o k
*/
type DcimRegionsBulkUpdateOK struct {
	Payload *models.Region
}

// IsSuccess returns true when this dcim regions bulk update o k response has a 2xx status code
func (o *DcimRegionsBulkUpdateOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this dcim regions bulk update o k response has a 3xx status code
func (o *DcimRegionsBulkUpdateOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this dcim regions bulk update o k response has a 4xx status code
func (o *DcimRegionsBulkUpdateOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this dcim regions bulk update o k response has a 5xx status code
func (o *DcimRegionsBulkUpdateOK) IsServerError() bool {
	return false
}

// IsCode returns true when this dcim regions bulk update o k response a status code equal to that given
func (o *DcimRegionsBulkUpdateOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the dcim regions bulk update o k response
func (o *DcimRegionsBulkUpdateOK) Code() int {
	return 200
}

func (o *DcimRegionsBulkUpdateOK) Error() string {
	return fmt.Sprintf("[PUT /dcim/regions/][%d] dcimRegionsBulkUpdateOK  %+v", 200, o.Payload)
}

func (o *DcimRegionsBulkUpdateOK) String() string {
	return fmt.Sprintf("[PUT /dcim/regions/][%d] dcimRegionsBulkUpdateOK  %+v", 200, o.Payload)
}

func (o *DcimRegionsBulkUpdateOK) GetPayload() *models.Region {
	return o.Payload
}

func (o *DcimRegionsBulkUpdateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Region)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDcimRegionsBulkUpdateDefault creates a DcimRegionsBulkUpdateDefault with default headers values
func NewDcimRegionsBulkUpdateDefault(code int) *DcimRegionsBulkUpdateDefault {
	return &DcimRegionsBulkUpdateDefault{
		_statusCode: code,
	}
}

/*
DcimRegionsBulkUpdateDefault describes a response with status code -1, with default header values.

DcimRegionsBulkUpdateDefault dcim regions bulk update default
*/
type DcimRegionsBulkUpdateDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this dcim regions bulk update default response has a 2xx status code
func (o *DcimRegionsBulkUpdateDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this dcim regions bulk update default response has a 3xx status code
func (o *DcimRegionsBulkUpdateDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this dcim regions bulk update default response has a 4xx status code
func (o *DcimRegionsBulkUpdateDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this dcim regions bulk update default response has a 5xx status code
func (o *DcimRegionsBulkUpdateDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this dcim regions bulk update default response a status code equal to that given
func (o *DcimRegionsBulkUpdateDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the dcim regions bulk update default response
func (o *DcimRegionsBulkUpdateDefault) Code() int {
	return o._statusCode
}

func (o *DcimRegionsBulkUpdateDefault) Error() string {
	return fmt.Sprintf("[PUT /dcim/regions/][%d] dcim_regions_bulk_update default  %+v", o._statusCode, o.Payload)
}

func (o *DcimRegionsBulkUpdateDefault) String() string {
	return fmt.Sprintf("[PUT /dcim/regions/][%d] dcim_regions_bulk_update default  %+v", o._statusCode, o.Payload)
}

func (o *DcimRegionsBulkUpdateDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *DcimRegionsBulkUpdateDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
