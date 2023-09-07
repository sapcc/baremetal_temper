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

package ipam

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netbox-community/go-netbox/v3/netbox/models"
)

// IpamVlanGroupsBulkUpdateReader is a Reader for the IpamVlanGroupsBulkUpdate structure.
type IpamVlanGroupsBulkUpdateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *IpamVlanGroupsBulkUpdateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewIpamVlanGroupsBulkUpdateOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewIpamVlanGroupsBulkUpdateDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewIpamVlanGroupsBulkUpdateOK creates a IpamVlanGroupsBulkUpdateOK with default headers values
func NewIpamVlanGroupsBulkUpdateOK() *IpamVlanGroupsBulkUpdateOK {
	return &IpamVlanGroupsBulkUpdateOK{}
}

/*
IpamVlanGroupsBulkUpdateOK describes a response with status code 200, with default header values.

IpamVlanGroupsBulkUpdateOK ipam vlan groups bulk update o k
*/
type IpamVlanGroupsBulkUpdateOK struct {
	Payload *models.VLANGroup
}

// IsSuccess returns true when this ipam vlan groups bulk update o k response has a 2xx status code
func (o *IpamVlanGroupsBulkUpdateOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this ipam vlan groups bulk update o k response has a 3xx status code
func (o *IpamVlanGroupsBulkUpdateOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this ipam vlan groups bulk update o k response has a 4xx status code
func (o *IpamVlanGroupsBulkUpdateOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this ipam vlan groups bulk update o k response has a 5xx status code
func (o *IpamVlanGroupsBulkUpdateOK) IsServerError() bool {
	return false
}

// IsCode returns true when this ipam vlan groups bulk update o k response a status code equal to that given
func (o *IpamVlanGroupsBulkUpdateOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the ipam vlan groups bulk update o k response
func (o *IpamVlanGroupsBulkUpdateOK) Code() int {
	return 200
}

func (o *IpamVlanGroupsBulkUpdateOK) Error() string {
	return fmt.Sprintf("[PUT /ipam/vlan-groups/][%d] ipamVlanGroupsBulkUpdateOK  %+v", 200, o.Payload)
}

func (o *IpamVlanGroupsBulkUpdateOK) String() string {
	return fmt.Sprintf("[PUT /ipam/vlan-groups/][%d] ipamVlanGroupsBulkUpdateOK  %+v", 200, o.Payload)
}

func (o *IpamVlanGroupsBulkUpdateOK) GetPayload() *models.VLANGroup {
	return o.Payload
}

func (o *IpamVlanGroupsBulkUpdateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.VLANGroup)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewIpamVlanGroupsBulkUpdateDefault creates a IpamVlanGroupsBulkUpdateDefault with default headers values
func NewIpamVlanGroupsBulkUpdateDefault(code int) *IpamVlanGroupsBulkUpdateDefault {
	return &IpamVlanGroupsBulkUpdateDefault{
		_statusCode: code,
	}
}

/*
IpamVlanGroupsBulkUpdateDefault describes a response with status code -1, with default header values.

IpamVlanGroupsBulkUpdateDefault ipam vlan groups bulk update default
*/
type IpamVlanGroupsBulkUpdateDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this ipam vlan groups bulk update default response has a 2xx status code
func (o *IpamVlanGroupsBulkUpdateDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this ipam vlan groups bulk update default response has a 3xx status code
func (o *IpamVlanGroupsBulkUpdateDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this ipam vlan groups bulk update default response has a 4xx status code
func (o *IpamVlanGroupsBulkUpdateDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this ipam vlan groups bulk update default response has a 5xx status code
func (o *IpamVlanGroupsBulkUpdateDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this ipam vlan groups bulk update default response a status code equal to that given
func (o *IpamVlanGroupsBulkUpdateDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the ipam vlan groups bulk update default response
func (o *IpamVlanGroupsBulkUpdateDefault) Code() int {
	return o._statusCode
}

func (o *IpamVlanGroupsBulkUpdateDefault) Error() string {
	return fmt.Sprintf("[PUT /ipam/vlan-groups/][%d] ipam_vlan-groups_bulk_update default  %+v", o._statusCode, o.Payload)
}

func (o *IpamVlanGroupsBulkUpdateDefault) String() string {
	return fmt.Sprintf("[PUT /ipam/vlan-groups/][%d] ipam_vlan-groups_bulk_update default  %+v", o._statusCode, o.Payload)
}

func (o *IpamVlanGroupsBulkUpdateDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *IpamVlanGroupsBulkUpdateDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}