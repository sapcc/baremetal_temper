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
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/netbox-community/go-netbox/v3/netbox/models"
)

// DcimPlatformsListReader is a Reader for the DcimPlatformsList structure.
type DcimPlatformsListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DcimPlatformsListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDcimPlatformsListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDcimPlatformsListDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDcimPlatformsListOK creates a DcimPlatformsListOK with default headers values
func NewDcimPlatformsListOK() *DcimPlatformsListOK {
	return &DcimPlatformsListOK{}
}

/*
DcimPlatformsListOK describes a response with status code 200, with default header values.

DcimPlatformsListOK dcim platforms list o k
*/
type DcimPlatformsListOK struct {
	Payload *DcimPlatformsListOKBody
}

// IsSuccess returns true when this dcim platforms list o k response has a 2xx status code
func (o *DcimPlatformsListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this dcim platforms list o k response has a 3xx status code
func (o *DcimPlatformsListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this dcim platforms list o k response has a 4xx status code
func (o *DcimPlatformsListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this dcim platforms list o k response has a 5xx status code
func (o *DcimPlatformsListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this dcim platforms list o k response a status code equal to that given
func (o *DcimPlatformsListOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the dcim platforms list o k response
func (o *DcimPlatformsListOK) Code() int {
	return 200
}

func (o *DcimPlatformsListOK) Error() string {
	return fmt.Sprintf("[GET /dcim/platforms/][%d] dcimPlatformsListOK  %+v", 200, o.Payload)
}

func (o *DcimPlatformsListOK) String() string {
	return fmt.Sprintf("[GET /dcim/platforms/][%d] dcimPlatformsListOK  %+v", 200, o.Payload)
}

func (o *DcimPlatformsListOK) GetPayload() *DcimPlatformsListOKBody {
	return o.Payload
}

func (o *DcimPlatformsListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(DcimPlatformsListOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDcimPlatformsListDefault creates a DcimPlatformsListDefault with default headers values
func NewDcimPlatformsListDefault(code int) *DcimPlatformsListDefault {
	return &DcimPlatformsListDefault{
		_statusCode: code,
	}
}

/*
DcimPlatformsListDefault describes a response with status code -1, with default header values.

DcimPlatformsListDefault dcim platforms list default
*/
type DcimPlatformsListDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this dcim platforms list default response has a 2xx status code
func (o *DcimPlatformsListDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this dcim platforms list default response has a 3xx status code
func (o *DcimPlatformsListDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this dcim platforms list default response has a 4xx status code
func (o *DcimPlatformsListDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this dcim platforms list default response has a 5xx status code
func (o *DcimPlatformsListDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this dcim platforms list default response a status code equal to that given
func (o *DcimPlatformsListDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the dcim platforms list default response
func (o *DcimPlatformsListDefault) Code() int {
	return o._statusCode
}

func (o *DcimPlatformsListDefault) Error() string {
	return fmt.Sprintf("[GET /dcim/platforms/][%d] dcim_platforms_list default  %+v", o._statusCode, o.Payload)
}

func (o *DcimPlatformsListDefault) String() string {
	return fmt.Sprintf("[GET /dcim/platforms/][%d] dcim_platforms_list default  %+v", o._statusCode, o.Payload)
}

func (o *DcimPlatformsListDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *DcimPlatformsListDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
DcimPlatformsListOKBody dcim platforms list o k body
swagger:model DcimPlatformsListOKBody
*/
type DcimPlatformsListOKBody struct {

	// count
	// Required: true
	Count *int64 `json:"count"`

	// next
	// Format: uri
	Next *strfmt.URI `json:"next,omitempty"`

	// previous
	// Format: uri
	Previous *strfmt.URI `json:"previous,omitempty"`

	// results
	// Required: true
	Results []*models.Platform `json:"results"`
}

// Validate validates this dcim platforms list o k body
func (o *DcimPlatformsListOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCount(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateNext(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validatePrevious(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateResults(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DcimPlatformsListOKBody) validateCount(formats strfmt.Registry) error {

	if err := validate.Required("dcimPlatformsListOK"+"."+"count", "body", o.Count); err != nil {
		return err
	}

	return nil
}

func (o *DcimPlatformsListOKBody) validateNext(formats strfmt.Registry) error {
	if swag.IsZero(o.Next) { // not required
		return nil
	}

	if err := validate.FormatOf("dcimPlatformsListOK"+"."+"next", "body", "uri", o.Next.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *DcimPlatformsListOKBody) validatePrevious(formats strfmt.Registry) error {
	if swag.IsZero(o.Previous) { // not required
		return nil
	}

	if err := validate.FormatOf("dcimPlatformsListOK"+"."+"previous", "body", "uri", o.Previous.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *DcimPlatformsListOKBody) validateResults(formats strfmt.Registry) error {

	if err := validate.Required("dcimPlatformsListOK"+"."+"results", "body", o.Results); err != nil {
		return err
	}

	for i := 0; i < len(o.Results); i++ {
		if swag.IsZero(o.Results[i]) { // not required
			continue
		}

		if o.Results[i] != nil {
			if err := o.Results[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("dcimPlatformsListOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("dcimPlatformsListOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this dcim platforms list o k body based on the context it is used
func (o *DcimPlatformsListOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateResults(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DcimPlatformsListOKBody) contextValidateResults(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Results); i++ {

		if o.Results[i] != nil {
			if err := o.Results[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("dcimPlatformsListOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("dcimPlatformsListOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *DcimPlatformsListOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DcimPlatformsListOKBody) UnmarshalBinary(b []byte) error {
	var res DcimPlatformsListOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
