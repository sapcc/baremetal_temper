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

// IpamRirsCreateReader is a Reader for the IpamRirsCreate structure.
type IpamRirsCreateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *IpamRirsCreateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewIpamRirsCreateCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewIpamRirsCreateDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewIpamRirsCreateCreated creates a IpamRirsCreateCreated with default headers values
func NewIpamRirsCreateCreated() *IpamRirsCreateCreated {
	return &IpamRirsCreateCreated{}
}

/*
IpamRirsCreateCreated describes a response with status code 201, with default header values.

IpamRirsCreateCreated ipam rirs create created
*/
type IpamRirsCreateCreated struct {
	Payload *models.RIR
}

// IsSuccess returns true when this ipam rirs create created response has a 2xx status code
func (o *IpamRirsCreateCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this ipam rirs create created response has a 3xx status code
func (o *IpamRirsCreateCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this ipam rirs create created response has a 4xx status code
func (o *IpamRirsCreateCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this ipam rirs create created response has a 5xx status code
func (o *IpamRirsCreateCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this ipam rirs create created response a status code equal to that given
func (o *IpamRirsCreateCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the ipam rirs create created response
func (o *IpamRirsCreateCreated) Code() int {
	return 201
}

func (o *IpamRirsCreateCreated) Error() string {
	return fmt.Sprintf("[POST /ipam/rirs/][%d] ipamRirsCreateCreated  %+v", 201, o.Payload)
}

func (o *IpamRirsCreateCreated) String() string {
	return fmt.Sprintf("[POST /ipam/rirs/][%d] ipamRirsCreateCreated  %+v", 201, o.Payload)
}

func (o *IpamRirsCreateCreated) GetPayload() *models.RIR {
	return o.Payload
}

func (o *IpamRirsCreateCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RIR)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewIpamRirsCreateDefault creates a IpamRirsCreateDefault with default headers values
func NewIpamRirsCreateDefault(code int) *IpamRirsCreateDefault {
	return &IpamRirsCreateDefault{
		_statusCode: code,
	}
}

/*
IpamRirsCreateDefault describes a response with status code -1, with default header values.

IpamRirsCreateDefault ipam rirs create default
*/
type IpamRirsCreateDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this ipam rirs create default response has a 2xx status code
func (o *IpamRirsCreateDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this ipam rirs create default response has a 3xx status code
func (o *IpamRirsCreateDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this ipam rirs create default response has a 4xx status code
func (o *IpamRirsCreateDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this ipam rirs create default response has a 5xx status code
func (o *IpamRirsCreateDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this ipam rirs create default response a status code equal to that given
func (o *IpamRirsCreateDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the ipam rirs create default response
func (o *IpamRirsCreateDefault) Code() int {
	return o._statusCode
}

func (o *IpamRirsCreateDefault) Error() string {
	return fmt.Sprintf("[POST /ipam/rirs/][%d] ipam_rirs_create default  %+v", o._statusCode, o.Payload)
}

func (o *IpamRirsCreateDefault) String() string {
	return fmt.Sprintf("[POST /ipam/rirs/][%d] ipam_rirs_create default  %+v", o._statusCode, o.Payload)
}

func (o *IpamRirsCreateDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *IpamRirsCreateDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
