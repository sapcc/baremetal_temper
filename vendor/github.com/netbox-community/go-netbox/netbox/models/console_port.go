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

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ConsolePort console port
//
// swagger:model ConsolePort
type ConsolePort struct {

	// cable
	Cable *NestedCable `json:"cable,omitempty"`

	// Connected endpoint
	//
	//
	// Return the appropriate serializer for the type of connected object.
	//
	// Read Only: true
	ConnectedEndpoint map[string]string `json:"connected_endpoint,omitempty"`

	// Connected endpoint type
	// Read Only: true
	ConnectedEndpointType string `json:"connected_endpoint_type,omitempty"`

	// connection status
	ConnectionStatus *ConsolePortConnectionStatus `json:"connection_status,omitempty"`

	// Description
	// Max Length: 200
	Description string `json:"description,omitempty"`

	// device
	// Required: true
	Device *NestedDevice `json:"device"`

	// ID
	// Read Only: true
	ID int64 `json:"id,omitempty"`

	// Label
	//
	// Physical label
	// Max Length: 64
	Label string `json:"label,omitempty"`

	// Name
	// Required: true
	// Max Length: 64
	// Min Length: 1
	Name *string `json:"name"`

	// tags
	Tags []*NestedTag `json:"tags,omitempty"`

	// type
	Type *ConsolePortType `json:"type,omitempty"`

	// Url
	// Read Only: true
	// Format: uri
	URL strfmt.URI `json:"url,omitempty"`
}

// Validate validates this console port
func (m *ConsolePort) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCable(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateConnectionStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDevice(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLabel(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTags(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateURL(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ConsolePort) validateCable(formats strfmt.Registry) error {

	if swag.IsZero(m.Cable) { // not required
		return nil
	}

	if m.Cable != nil {
		if err := m.Cable.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cable")
			}
			return err
		}
	}

	return nil
}

func (m *ConsolePort) validateConnectionStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.ConnectionStatus) { // not required
		return nil
	}

	if m.ConnectionStatus != nil {
		if err := m.ConnectionStatus.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("connection_status")
			}
			return err
		}
	}

	return nil
}

func (m *ConsolePort) validateDescription(formats strfmt.Registry) error {

	if swag.IsZero(m.Description) { // not required
		return nil
	}

	if err := validate.MaxLength("description", "body", string(m.Description), 200); err != nil {
		return err
	}

	return nil
}

func (m *ConsolePort) validateDevice(formats strfmt.Registry) error {

	if err := validate.Required("device", "body", m.Device); err != nil {
		return err
	}

	if m.Device != nil {
		if err := m.Device.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("device")
			}
			return err
		}
	}

	return nil
}

func (m *ConsolePort) validateLabel(formats strfmt.Registry) error {

	if swag.IsZero(m.Label) { // not required
		return nil
	}

	if err := validate.MaxLength("label", "body", string(m.Label), 64); err != nil {
		return err
	}

	return nil
}

func (m *ConsolePort) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	if err := validate.MinLength("name", "body", string(*m.Name), 1); err != nil {
		return err
	}

	if err := validate.MaxLength("name", "body", string(*m.Name), 64); err != nil {
		return err
	}

	return nil
}

func (m *ConsolePort) validateTags(formats strfmt.Registry) error {

	if swag.IsZero(m.Tags) { // not required
		return nil
	}

	for i := 0; i < len(m.Tags); i++ {
		if swag.IsZero(m.Tags[i]) { // not required
			continue
		}

		if m.Tags[i] != nil {
			if err := m.Tags[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("tags" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ConsolePort) validateType(formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	if m.Type != nil {
		if err := m.Type.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("type")
			}
			return err
		}
	}

	return nil
}

func (m *ConsolePort) validateURL(formats strfmt.Registry) error {

	if swag.IsZero(m.URL) { // not required
		return nil
	}

	if err := validate.FormatOf("url", "body", "uri", m.URL.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ConsolePort) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConsolePort) UnmarshalBinary(b []byte) error {
	var res ConsolePort
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ConsolePortConnectionStatus Connection status
//
// swagger:model ConsolePortConnectionStatus
type ConsolePortConnectionStatus struct {

	// label
	// Required: true
	// Enum: [Not Connected Connected]
	Label *string `json:"label"`

	// value
	// Required: true
	// Enum: [false true]
	Value *bool `json:"value"`
}

// Validate validates this console port connection status
func (m *ConsolePortConnectionStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLabel(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValue(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var consolePortConnectionStatusTypeLabelPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["Not Connected","Connected"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		consolePortConnectionStatusTypeLabelPropEnum = append(consolePortConnectionStatusTypeLabelPropEnum, v)
	}
}

const (

	// ConsolePortConnectionStatusLabelNotConnected captures enum value "Not Connected"
	ConsolePortConnectionStatusLabelNotConnected string = "Not Connected"

	// ConsolePortConnectionStatusLabelConnected captures enum value "Connected"
	ConsolePortConnectionStatusLabelConnected string = "Connected"
)

// prop value enum
func (m *ConsolePortConnectionStatus) validateLabelEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, consolePortConnectionStatusTypeLabelPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *ConsolePortConnectionStatus) validateLabel(formats strfmt.Registry) error {

	if err := validate.Required("connection_status"+"."+"label", "body", m.Label); err != nil {
		return err
	}

	// value enum
	if err := m.validateLabelEnum("connection_status"+"."+"label", "body", *m.Label); err != nil {
		return err
	}

	return nil
}

var consolePortConnectionStatusTypeValuePropEnum []interface{}

func init() {
	var res []bool
	if err := json.Unmarshal([]byte(`[false,true]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		consolePortConnectionStatusTypeValuePropEnum = append(consolePortConnectionStatusTypeValuePropEnum, v)
	}
}

// prop value enum
func (m *ConsolePortConnectionStatus) validateValueEnum(path, location string, value bool) error {
	if err := validate.EnumCase(path, location, value, consolePortConnectionStatusTypeValuePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *ConsolePortConnectionStatus) validateValue(formats strfmt.Registry) error {

	if err := validate.Required("connection_status"+"."+"value", "body", m.Value); err != nil {
		return err
	}

	// value enum
	if err := m.validateValueEnum("connection_status"+"."+"value", "body", *m.Value); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ConsolePortConnectionStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConsolePortConnectionStatus) UnmarshalBinary(b []byte) error {
	var res ConsolePortConnectionStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ConsolePortType Type
//
// swagger:model ConsolePortType
type ConsolePortType struct {

	// label
	// Required: true
	// Enum: [DE-9 DB-25 RJ-11 RJ-12 RJ-45 USB Type A USB Type B USB Type C USB Mini A USB Mini B USB Micro A USB Micro B Other]
	Label *string `json:"label"`

	// value
	// Required: true
	// Enum: [de-9 db-25 rj-11 rj-12 rj-45 usb-a usb-b usb-c usb-mini-a usb-mini-b usb-micro-a usb-micro-b other]
	Value *string `json:"value"`
}

// Validate validates this console port type
func (m *ConsolePortType) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLabel(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValue(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var consolePortTypeTypeLabelPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["DE-9","DB-25","RJ-11","RJ-12","RJ-45","USB Type A","USB Type B","USB Type C","USB Mini A","USB Mini B","USB Micro A","USB Micro B","Other"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		consolePortTypeTypeLabelPropEnum = append(consolePortTypeTypeLabelPropEnum, v)
	}
}

const (

	// ConsolePortTypeLabelDE9 captures enum value "DE-9"
	ConsolePortTypeLabelDE9 string = "DE-9"

	// ConsolePortTypeLabelDB25 captures enum value "DB-25"
	ConsolePortTypeLabelDB25 string = "DB-25"

	// ConsolePortTypeLabelRJ11 captures enum value "RJ-11"
	ConsolePortTypeLabelRJ11 string = "RJ-11"

	// ConsolePortTypeLabelRJ12 captures enum value "RJ-12"
	ConsolePortTypeLabelRJ12 string = "RJ-12"

	// ConsolePortTypeLabelRJ45 captures enum value "RJ-45"
	ConsolePortTypeLabelRJ45 string = "RJ-45"

	// ConsolePortTypeLabelUSBTypeA captures enum value "USB Type A"
	ConsolePortTypeLabelUSBTypeA string = "USB Type A"

	// ConsolePortTypeLabelUSBTypeB captures enum value "USB Type B"
	ConsolePortTypeLabelUSBTypeB string = "USB Type B"

	// ConsolePortTypeLabelUSBTypeC captures enum value "USB Type C"
	ConsolePortTypeLabelUSBTypeC string = "USB Type C"

	// ConsolePortTypeLabelUSBMiniA captures enum value "USB Mini A"
	ConsolePortTypeLabelUSBMiniA string = "USB Mini A"

	// ConsolePortTypeLabelUSBMiniB captures enum value "USB Mini B"
	ConsolePortTypeLabelUSBMiniB string = "USB Mini B"

	// ConsolePortTypeLabelUSBMicroA captures enum value "USB Micro A"
	ConsolePortTypeLabelUSBMicroA string = "USB Micro A"

	// ConsolePortTypeLabelUSBMicroB captures enum value "USB Micro B"
	ConsolePortTypeLabelUSBMicroB string = "USB Micro B"

	// ConsolePortTypeLabelOther captures enum value "Other"
	ConsolePortTypeLabelOther string = "Other"
)

// prop value enum
func (m *ConsolePortType) validateLabelEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, consolePortTypeTypeLabelPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *ConsolePortType) validateLabel(formats strfmt.Registry) error {

	if err := validate.Required("type"+"."+"label", "body", m.Label); err != nil {
		return err
	}

	// value enum
	if err := m.validateLabelEnum("type"+"."+"label", "body", *m.Label); err != nil {
		return err
	}

	return nil
}

var consolePortTypeTypeValuePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["de-9","db-25","rj-11","rj-12","rj-45","usb-a","usb-b","usb-c","usb-mini-a","usb-mini-b","usb-micro-a","usb-micro-b","other"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		consolePortTypeTypeValuePropEnum = append(consolePortTypeTypeValuePropEnum, v)
	}
}

const (

	// ConsolePortTypeValueDe9 captures enum value "de-9"
	ConsolePortTypeValueDe9 string = "de-9"

	// ConsolePortTypeValueDb25 captures enum value "db-25"
	ConsolePortTypeValueDb25 string = "db-25"

	// ConsolePortTypeValueRj11 captures enum value "rj-11"
	ConsolePortTypeValueRj11 string = "rj-11"

	// ConsolePortTypeValueRj12 captures enum value "rj-12"
	ConsolePortTypeValueRj12 string = "rj-12"

	// ConsolePortTypeValueRj45 captures enum value "rj-45"
	ConsolePortTypeValueRj45 string = "rj-45"

	// ConsolePortTypeValueUsba captures enum value "usb-a"
	ConsolePortTypeValueUsba string = "usb-a"

	// ConsolePortTypeValueUsbb captures enum value "usb-b"
	ConsolePortTypeValueUsbb string = "usb-b"

	// ConsolePortTypeValueUsbc captures enum value "usb-c"
	ConsolePortTypeValueUsbc string = "usb-c"

	// ConsolePortTypeValueUsbMinia captures enum value "usb-mini-a"
	ConsolePortTypeValueUsbMinia string = "usb-mini-a"

	// ConsolePortTypeValueUsbMinib captures enum value "usb-mini-b"
	ConsolePortTypeValueUsbMinib string = "usb-mini-b"

	// ConsolePortTypeValueUsbMicroa captures enum value "usb-micro-a"
	ConsolePortTypeValueUsbMicroa string = "usb-micro-a"

	// ConsolePortTypeValueUsbMicrob captures enum value "usb-micro-b"
	ConsolePortTypeValueUsbMicrob string = "usb-micro-b"

	// ConsolePortTypeValueOther captures enum value "other"
	ConsolePortTypeValueOther string = "other"
)

// prop value enum
func (m *ConsolePortType) validateValueEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, consolePortTypeTypeValuePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *ConsolePortType) validateValue(formats strfmt.Registry) error {

	if err := validate.Required("type"+"."+"value", "body", m.Value); err != nil {
		return err
	}

	// value enum
	if err := m.validateValueEnum("type"+"."+"value", "body", *m.Value); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ConsolePortType) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ConsolePortType) UnmarshalBinary(b []byte) error {
	var res ConsolePortType
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
