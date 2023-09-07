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

// DcimVirtualChassisReadReader is a Reader for the DcimVirtualChassisRead structure.
type DcimVirtualChassisReadReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DcimVirtualChassisReadReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDcimVirtualChassisReadOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDcimVirtualChassisReadDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDcimVirtualChassisReadOK creates a DcimVirtualChassisReadOK with default headers values
func NewDcimVirtualChassisReadOK() *DcimVirtualChassisReadOK {
	return &DcimVirtualChassisReadOK{}
}

/*
DcimVirtualChassisReadOK describes a response with status code 200, with default header values.

DcimVirtualChassisReadOK dcim virtual chassis read o k
*/
type DcimVirtualChassisReadOK struct {
	Payload *models.VirtualChassis
}

// IsSuccess returns true when this dcim virtual chassis read o k response has a 2xx status code
func (o *DcimVirtualChassisReadOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this dcim virtual chassis read o k response has a 3xx status code
func (o *DcimVirtualChassisReadOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this dcim virtual chassis read o k response has a 4xx status code
func (o *DcimVirtualChassisReadOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this dcim virtual chassis read o k response has a 5xx status code
func (o *DcimVirtualChassisReadOK) IsServerError() bool {
	return false
}

// IsCode returns true when this dcim virtual chassis read o k response a status code equal to that given
func (o *DcimVirtualChassisReadOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the dcim virtual chassis read o k response
func (o *DcimVirtualChassisReadOK) Code() int {
	return 200
}

func (o *DcimVirtualChassisReadOK) Error() string {
	return fmt.Sprintf("[GET /dcim/virtual-chassis/{id}/][%d] dcimVirtualChassisReadOK  %+v", 200, o.Payload)
}

func (o *DcimVirtualChassisReadOK) String() string {
	return fmt.Sprintf("[GET /dcim/virtual-chassis/{id}/][%d] dcimVirtualChassisReadOK  %+v", 200, o.Payload)
}

func (o *DcimVirtualChassisReadOK) GetPayload() *models.VirtualChassis {
	return o.Payload
}

func (o *DcimVirtualChassisReadOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.VirtualChassis)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDcimVirtualChassisReadDefault creates a DcimVirtualChassisReadDefault with default headers values
func NewDcimVirtualChassisReadDefault(code int) *DcimVirtualChassisReadDefault {
	return &DcimVirtualChassisReadDefault{
		_statusCode: code,
	}
}

/*
DcimVirtualChassisReadDefault describes a response with status code -1, with default header values.

DcimVirtualChassisReadDefault dcim virtual chassis read default
*/
type DcimVirtualChassisReadDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this dcim virtual chassis read default response has a 2xx status code
func (o *DcimVirtualChassisReadDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this dcim virtual chassis read default response has a 3xx status code
func (o *DcimVirtualChassisReadDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this dcim virtual chassis read default response has a 4xx status code
func (o *DcimVirtualChassisReadDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this dcim virtual chassis read default response has a 5xx status code
func (o *DcimVirtualChassisReadDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this dcim virtual chassis read default response a status code equal to that given
func (o *DcimVirtualChassisReadDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the dcim virtual chassis read default response
func (o *DcimVirtualChassisReadDefault) Code() int {
	return o._statusCode
}

func (o *DcimVirtualChassisReadDefault) Error() string {
	return fmt.Sprintf("[GET /dcim/virtual-chassis/{id}/][%d] dcim_virtual-chassis_read default  %+v", o._statusCode, o.Payload)
}

func (o *DcimVirtualChassisReadDefault) String() string {
	return fmt.Sprintf("[GET /dcim/virtual-chassis/{id}/][%d] dcim_virtual-chassis_read default  %+v", o._statusCode, o.Payload)
}

func (o *DcimVirtualChassisReadDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *DcimVirtualChassisReadDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
