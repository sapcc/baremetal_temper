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

// IpamServiceTemplatesReadReader is a Reader for the IpamServiceTemplatesRead structure.
type IpamServiceTemplatesReadReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *IpamServiceTemplatesReadReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewIpamServiceTemplatesReadOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewIpamServiceTemplatesReadDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewIpamServiceTemplatesReadOK creates a IpamServiceTemplatesReadOK with default headers values
func NewIpamServiceTemplatesReadOK() *IpamServiceTemplatesReadOK {
	return &IpamServiceTemplatesReadOK{}
}

/*
IpamServiceTemplatesReadOK describes a response with status code 200, with default header values.

IpamServiceTemplatesReadOK ipam service templates read o k
*/
type IpamServiceTemplatesReadOK struct {
	Payload *models.ServiceTemplate
}

// IsSuccess returns true when this ipam service templates read o k response has a 2xx status code
func (o *IpamServiceTemplatesReadOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this ipam service templates read o k response has a 3xx status code
func (o *IpamServiceTemplatesReadOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this ipam service templates read o k response has a 4xx status code
func (o *IpamServiceTemplatesReadOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this ipam service templates read o k response has a 5xx status code
func (o *IpamServiceTemplatesReadOK) IsServerError() bool {
	return false
}

// IsCode returns true when this ipam service templates read o k response a status code equal to that given
func (o *IpamServiceTemplatesReadOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the ipam service templates read o k response
func (o *IpamServiceTemplatesReadOK) Code() int {
	return 200
}

func (o *IpamServiceTemplatesReadOK) Error() string {
	return fmt.Sprintf("[GET /ipam/service-templates/{id}/][%d] ipamServiceTemplatesReadOK  %+v", 200, o.Payload)
}

func (o *IpamServiceTemplatesReadOK) String() string {
	return fmt.Sprintf("[GET /ipam/service-templates/{id}/][%d] ipamServiceTemplatesReadOK  %+v", 200, o.Payload)
}

func (o *IpamServiceTemplatesReadOK) GetPayload() *models.ServiceTemplate {
	return o.Payload
}

func (o *IpamServiceTemplatesReadOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ServiceTemplate)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewIpamServiceTemplatesReadDefault creates a IpamServiceTemplatesReadDefault with default headers values
func NewIpamServiceTemplatesReadDefault(code int) *IpamServiceTemplatesReadDefault {
	return &IpamServiceTemplatesReadDefault{
		_statusCode: code,
	}
}

/*
IpamServiceTemplatesReadDefault describes a response with status code -1, with default header values.

IpamServiceTemplatesReadDefault ipam service templates read default
*/
type IpamServiceTemplatesReadDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this ipam service templates read default response has a 2xx status code
func (o *IpamServiceTemplatesReadDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this ipam service templates read default response has a 3xx status code
func (o *IpamServiceTemplatesReadDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this ipam service templates read default response has a 4xx status code
func (o *IpamServiceTemplatesReadDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this ipam service templates read default response has a 5xx status code
func (o *IpamServiceTemplatesReadDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this ipam service templates read default response a status code equal to that given
func (o *IpamServiceTemplatesReadDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the ipam service templates read default response
func (o *IpamServiceTemplatesReadDefault) Code() int {
	return o._statusCode
}

func (o *IpamServiceTemplatesReadDefault) Error() string {
	return fmt.Sprintf("[GET /ipam/service-templates/{id}/][%d] ipam_service-templates_read default  %+v", o._statusCode, o.Payload)
}

func (o *IpamServiceTemplatesReadDefault) String() string {
	return fmt.Sprintf("[GET /ipam/service-templates/{id}/][%d] ipam_service-templates_read default  %+v", o._statusCode, o.Payload)
}

func (o *IpamServiceTemplatesReadDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *IpamServiceTemplatesReadDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
