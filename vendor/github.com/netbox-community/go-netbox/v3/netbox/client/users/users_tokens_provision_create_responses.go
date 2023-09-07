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

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// UsersTokensProvisionCreateReader is a Reader for the UsersTokensProvisionCreate structure.
type UsersTokensProvisionCreateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UsersTokensProvisionCreateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewUsersTokensProvisionCreateCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewUsersTokensProvisionCreateDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUsersTokensProvisionCreateCreated creates a UsersTokensProvisionCreateCreated with default headers values
func NewUsersTokensProvisionCreateCreated() *UsersTokensProvisionCreateCreated {
	return &UsersTokensProvisionCreateCreated{}
}

/*
UsersTokensProvisionCreateCreated describes a response with status code 201, with default header values.

UsersTokensProvisionCreateCreated users tokens provision create created
*/
type UsersTokensProvisionCreateCreated struct {
}

// IsSuccess returns true when this users tokens provision create created response has a 2xx status code
func (o *UsersTokensProvisionCreateCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this users tokens provision create created response has a 3xx status code
func (o *UsersTokensProvisionCreateCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this users tokens provision create created response has a 4xx status code
func (o *UsersTokensProvisionCreateCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this users tokens provision create created response has a 5xx status code
func (o *UsersTokensProvisionCreateCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this users tokens provision create created response a status code equal to that given
func (o *UsersTokensProvisionCreateCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the users tokens provision create created response
func (o *UsersTokensProvisionCreateCreated) Code() int {
	return 201
}

func (o *UsersTokensProvisionCreateCreated) Error() string {
	return fmt.Sprintf("[POST /users/tokens/provision/][%d] usersTokensProvisionCreateCreated ", 201)
}

func (o *UsersTokensProvisionCreateCreated) String() string {
	return fmt.Sprintf("[POST /users/tokens/provision/][%d] usersTokensProvisionCreateCreated ", 201)
}

func (o *UsersTokensProvisionCreateCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUsersTokensProvisionCreateDefault creates a UsersTokensProvisionCreateDefault with default headers values
func NewUsersTokensProvisionCreateDefault(code int) *UsersTokensProvisionCreateDefault {
	return &UsersTokensProvisionCreateDefault{
		_statusCode: code,
	}
}

/*
UsersTokensProvisionCreateDefault describes a response with status code -1, with default header values.

UsersTokensProvisionCreateDefault users tokens provision create default
*/
type UsersTokensProvisionCreateDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this users tokens provision create default response has a 2xx status code
func (o *UsersTokensProvisionCreateDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this users tokens provision create default response has a 3xx status code
func (o *UsersTokensProvisionCreateDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this users tokens provision create default response has a 4xx status code
func (o *UsersTokensProvisionCreateDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this users tokens provision create default response has a 5xx status code
func (o *UsersTokensProvisionCreateDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this users tokens provision create default response a status code equal to that given
func (o *UsersTokensProvisionCreateDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the users tokens provision create default response
func (o *UsersTokensProvisionCreateDefault) Code() int {
	return o._statusCode
}

func (o *UsersTokensProvisionCreateDefault) Error() string {
	return fmt.Sprintf("[POST /users/tokens/provision/][%d] users_tokens_provision_create default  %+v", o._statusCode, o.Payload)
}

func (o *UsersTokensProvisionCreateDefault) String() string {
	return fmt.Sprintf("[POST /users/tokens/provision/][%d] users_tokens_provision_create default  %+v", o._statusCode, o.Payload)
}

func (o *UsersTokensProvisionCreateDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *UsersTokensProvisionCreateDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
