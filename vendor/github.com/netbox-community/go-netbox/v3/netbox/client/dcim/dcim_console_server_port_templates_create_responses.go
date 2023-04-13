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

// DcimConsoleServerPortTemplatesCreateReader is a Reader for the DcimConsoleServerPortTemplatesCreate structure.
type DcimConsoleServerPortTemplatesCreateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DcimConsoleServerPortTemplatesCreateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewDcimConsoleServerPortTemplatesCreateCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDcimConsoleServerPortTemplatesCreateDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDcimConsoleServerPortTemplatesCreateCreated creates a DcimConsoleServerPortTemplatesCreateCreated with default headers values
func NewDcimConsoleServerPortTemplatesCreateCreated() *DcimConsoleServerPortTemplatesCreateCreated {
	return &DcimConsoleServerPortTemplatesCreateCreated{}
}

/*
DcimConsoleServerPortTemplatesCreateCreated describes a response with status code 201, with default header values.

DcimConsoleServerPortTemplatesCreateCreated dcim console server port templates create created
*/
type DcimConsoleServerPortTemplatesCreateCreated struct {
	Payload *models.ConsoleServerPortTemplate
}

// IsSuccess returns true when this dcim console server port templates create created response has a 2xx status code
func (o *DcimConsoleServerPortTemplatesCreateCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this dcim console server port templates create created response has a 3xx status code
func (o *DcimConsoleServerPortTemplatesCreateCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this dcim console server port templates create created response has a 4xx status code
func (o *DcimConsoleServerPortTemplatesCreateCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this dcim console server port templates create created response has a 5xx status code
func (o *DcimConsoleServerPortTemplatesCreateCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this dcim console server port templates create created response a status code equal to that given
func (o *DcimConsoleServerPortTemplatesCreateCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the dcim console server port templates create created response
func (o *DcimConsoleServerPortTemplatesCreateCreated) Code() int {
	return 201
}

func (o *DcimConsoleServerPortTemplatesCreateCreated) Error() string {
	return fmt.Sprintf("[POST /dcim/console-server-port-templates/][%d] dcimConsoleServerPortTemplatesCreateCreated  %+v", 201, o.Payload)
}

func (o *DcimConsoleServerPortTemplatesCreateCreated) String() string {
	return fmt.Sprintf("[POST /dcim/console-server-port-templates/][%d] dcimConsoleServerPortTemplatesCreateCreated  %+v", 201, o.Payload)
}

func (o *DcimConsoleServerPortTemplatesCreateCreated) GetPayload() *models.ConsoleServerPortTemplate {
	return o.Payload
}

func (o *DcimConsoleServerPortTemplatesCreateCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ConsoleServerPortTemplate)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDcimConsoleServerPortTemplatesCreateDefault creates a DcimConsoleServerPortTemplatesCreateDefault with default headers values
func NewDcimConsoleServerPortTemplatesCreateDefault(code int) *DcimConsoleServerPortTemplatesCreateDefault {
	return &DcimConsoleServerPortTemplatesCreateDefault{
		_statusCode: code,
	}
}

/*
DcimConsoleServerPortTemplatesCreateDefault describes a response with status code -1, with default header values.

DcimConsoleServerPortTemplatesCreateDefault dcim console server port templates create default
*/
type DcimConsoleServerPortTemplatesCreateDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this dcim console server port templates create default response has a 2xx status code
func (o *DcimConsoleServerPortTemplatesCreateDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this dcim console server port templates create default response has a 3xx status code
func (o *DcimConsoleServerPortTemplatesCreateDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this dcim console server port templates create default response has a 4xx status code
func (o *DcimConsoleServerPortTemplatesCreateDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this dcim console server port templates create default response has a 5xx status code
func (o *DcimConsoleServerPortTemplatesCreateDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this dcim console server port templates create default response a status code equal to that given
func (o *DcimConsoleServerPortTemplatesCreateDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the dcim console server port templates create default response
func (o *DcimConsoleServerPortTemplatesCreateDefault) Code() int {
	return o._statusCode
}

func (o *DcimConsoleServerPortTemplatesCreateDefault) Error() string {
	return fmt.Sprintf("[POST /dcim/console-server-port-templates/][%d] dcim_console-server-port-templates_create default  %+v", o._statusCode, o.Payload)
}

func (o *DcimConsoleServerPortTemplatesCreateDefault) String() string {
	return fmt.Sprintf("[POST /dcim/console-server-port-templates/][%d] dcim_console-server-port-templates_create default  %+v", o._statusCode, o.Payload)
}

func (o *DcimConsoleServerPortTemplatesCreateDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *DcimConsoleServerPortTemplatesCreateDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
