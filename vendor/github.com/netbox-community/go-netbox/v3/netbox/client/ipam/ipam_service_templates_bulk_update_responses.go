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

// IpamServiceTemplatesBulkUpdateReader is a Reader for the IpamServiceTemplatesBulkUpdate structure.
type IpamServiceTemplatesBulkUpdateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *IpamServiceTemplatesBulkUpdateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewIpamServiceTemplatesBulkUpdateOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewIpamServiceTemplatesBulkUpdateDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewIpamServiceTemplatesBulkUpdateOK creates a IpamServiceTemplatesBulkUpdateOK with default headers values
func NewIpamServiceTemplatesBulkUpdateOK() *IpamServiceTemplatesBulkUpdateOK {
	return &IpamServiceTemplatesBulkUpdateOK{}
}

/*
IpamServiceTemplatesBulkUpdateOK describes a response with status code 200, with default header values.

IpamServiceTemplatesBulkUpdateOK ipam service templates bulk update o k
*/
type IpamServiceTemplatesBulkUpdateOK struct {
	Payload *models.ServiceTemplate
}

// IsSuccess returns true when this ipam service templates bulk update o k response has a 2xx status code
func (o *IpamServiceTemplatesBulkUpdateOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this ipam service templates bulk update o k response has a 3xx status code
func (o *IpamServiceTemplatesBulkUpdateOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this ipam service templates bulk update o k response has a 4xx status code
func (o *IpamServiceTemplatesBulkUpdateOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this ipam service templates bulk update o k response has a 5xx status code
func (o *IpamServiceTemplatesBulkUpdateOK) IsServerError() bool {
	return false
}

// IsCode returns true when this ipam service templates bulk update o k response a status code equal to that given
func (o *IpamServiceTemplatesBulkUpdateOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the ipam service templates bulk update o k response
func (o *IpamServiceTemplatesBulkUpdateOK) Code() int {
	return 200
}

func (o *IpamServiceTemplatesBulkUpdateOK) Error() string {
	return fmt.Sprintf("[PUT /ipam/service-templates/][%d] ipamServiceTemplatesBulkUpdateOK  %+v", 200, o.Payload)
}

func (o *IpamServiceTemplatesBulkUpdateOK) String() string {
	return fmt.Sprintf("[PUT /ipam/service-templates/][%d] ipamServiceTemplatesBulkUpdateOK  %+v", 200, o.Payload)
}

func (o *IpamServiceTemplatesBulkUpdateOK) GetPayload() *models.ServiceTemplate {
	return o.Payload
}

func (o *IpamServiceTemplatesBulkUpdateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ServiceTemplate)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewIpamServiceTemplatesBulkUpdateDefault creates a IpamServiceTemplatesBulkUpdateDefault with default headers values
func NewIpamServiceTemplatesBulkUpdateDefault(code int) *IpamServiceTemplatesBulkUpdateDefault {
	return &IpamServiceTemplatesBulkUpdateDefault{
		_statusCode: code,
	}
}

/*
IpamServiceTemplatesBulkUpdateDefault describes a response with status code -1, with default header values.

IpamServiceTemplatesBulkUpdateDefault ipam service templates bulk update default
*/
type IpamServiceTemplatesBulkUpdateDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this ipam service templates bulk update default response has a 2xx status code
func (o *IpamServiceTemplatesBulkUpdateDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this ipam service templates bulk update default response has a 3xx status code
func (o *IpamServiceTemplatesBulkUpdateDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this ipam service templates bulk update default response has a 4xx status code
func (o *IpamServiceTemplatesBulkUpdateDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this ipam service templates bulk update default response has a 5xx status code
func (o *IpamServiceTemplatesBulkUpdateDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this ipam service templates bulk update default response a status code equal to that given
func (o *IpamServiceTemplatesBulkUpdateDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the ipam service templates bulk update default response
func (o *IpamServiceTemplatesBulkUpdateDefault) Code() int {
	return o._statusCode
}

func (o *IpamServiceTemplatesBulkUpdateDefault) Error() string {
	return fmt.Sprintf("[PUT /ipam/service-templates/][%d] ipam_service-templates_bulk_update default  %+v", o._statusCode, o.Payload)
}

func (o *IpamServiceTemplatesBulkUpdateDefault) String() string {
	return fmt.Sprintf("[PUT /ipam/service-templates/][%d] ipam_service-templates_bulk_update default  %+v", o._statusCode, o.Payload)
}

func (o *IpamServiceTemplatesBulkUpdateDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *IpamServiceTemplatesBulkUpdateDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
