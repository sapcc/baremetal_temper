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
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Module module
//
// swagger:model Module
type Module struct {

	// Asset tag
	//
	// A unique tag used to identify this device
	// Max Length: 50
	AssetTag *string `json:"asset_tag,omitempty"`

	// Comments
	Comments string `json:"comments,omitempty"`

	// Created
	// Read Only: true
	// Format: date-time
	Created *strfmt.DateTime `json:"created,omitempty"`

	// Custom fields
	CustomFields interface{} `json:"custom_fields,omitempty"`

	// Description
	// Max Length: 200
	Description string `json:"description,omitempty"`

	// device
	// Required: true
	Device *NestedDevice `json:"device"`

	// Display
	// Read Only: true
	Display string `json:"display,omitempty"`

	// ID
	// Read Only: true
	ID int64 `json:"id,omitempty"`

	// Last updated
	// Read Only: true
	// Format: date-time
	LastUpdated *strfmt.DateTime `json:"last_updated,omitempty"`

	// module bay
	// Required: true
	ModuleBay *NestedModuleBay `json:"module_bay"`

	// module type
	// Required: true
	ModuleType *NestedModuleType `json:"module_type"`

	// Serial number
	// Max Length: 50
	Serial string `json:"serial,omitempty"`

	// status
	Status *ModuleStatus `json:"status,omitempty"`

	// tags
	Tags []*NestedTag `json:"tags,omitempty"`

	// Url
	// Read Only: true
	// Format: uri
	URL strfmt.URI `json:"url,omitempty"`
}

// Validate validates this module
func (m *Module) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAssetTag(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreated(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDevice(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastUpdated(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateModuleBay(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateModuleType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSerial(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTags(formats); err != nil {
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

func (m *Module) validateAssetTag(formats strfmt.Registry) error {
	if swag.IsZero(m.AssetTag) { // not required
		return nil
	}

	if err := validate.MaxLength("asset_tag", "body", *m.AssetTag, 50); err != nil {
		return err
	}

	return nil
}

func (m *Module) validateCreated(formats strfmt.Registry) error {
	if swag.IsZero(m.Created) { // not required
		return nil
	}

	if err := validate.FormatOf("created", "body", "date-time", m.Created.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Module) validateDescription(formats strfmt.Registry) error {
	if swag.IsZero(m.Description) { // not required
		return nil
	}

	if err := validate.MaxLength("description", "body", m.Description, 200); err != nil {
		return err
	}

	return nil
}

func (m *Module) validateDevice(formats strfmt.Registry) error {

	if err := validate.Required("device", "body", m.Device); err != nil {
		return err
	}

	if m.Device != nil {
		if err := m.Device.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("device")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("device")
			}
			return err
		}
	}

	return nil
}

func (m *Module) validateLastUpdated(formats strfmt.Registry) error {
	if swag.IsZero(m.LastUpdated) { // not required
		return nil
	}

	if err := validate.FormatOf("last_updated", "body", "date-time", m.LastUpdated.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Module) validateModuleBay(formats strfmt.Registry) error {

	if err := validate.Required("module_bay", "body", m.ModuleBay); err != nil {
		return err
	}

	if m.ModuleBay != nil {
		if err := m.ModuleBay.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("module_bay")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("module_bay")
			}
			return err
		}
	}

	return nil
}

func (m *Module) validateModuleType(formats strfmt.Registry) error {

	if err := validate.Required("module_type", "body", m.ModuleType); err != nil {
		return err
	}

	if m.ModuleType != nil {
		if err := m.ModuleType.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("module_type")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("module_type")
			}
			return err
		}
	}

	return nil
}

func (m *Module) validateSerial(formats strfmt.Registry) error {
	if swag.IsZero(m.Serial) { // not required
		return nil
	}

	if err := validate.MaxLength("serial", "body", m.Serial, 50); err != nil {
		return err
	}

	return nil
}

func (m *Module) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if m.Status != nil {
		if err := m.Status.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("status")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("status")
			}
			return err
		}
	}

	return nil
}

func (m *Module) validateTags(formats strfmt.Registry) error {
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
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("tags" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Module) validateURL(formats strfmt.Registry) error {
	if swag.IsZero(m.URL) { // not required
		return nil
	}

	if err := validate.FormatOf("url", "body", "uri", m.URL.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this module based on the context it is used
func (m *Module) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCreated(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDevice(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDisplay(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLastUpdated(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateModuleBay(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateModuleType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTags(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateURL(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Module) contextValidateCreated(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "created", "body", m.Created); err != nil {
		return err
	}

	return nil
}

func (m *Module) contextValidateDevice(ctx context.Context, formats strfmt.Registry) error {

	if m.Device != nil {
		if err := m.Device.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("device")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("device")
			}
			return err
		}
	}

	return nil
}

func (m *Module) contextValidateDisplay(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "display", "body", string(m.Display)); err != nil {
		return err
	}

	return nil
}

func (m *Module) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "id", "body", int64(m.ID)); err != nil {
		return err
	}

	return nil
}

func (m *Module) contextValidateLastUpdated(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "last_updated", "body", m.LastUpdated); err != nil {
		return err
	}

	return nil
}

func (m *Module) contextValidateModuleBay(ctx context.Context, formats strfmt.Registry) error {

	if m.ModuleBay != nil {
		if err := m.ModuleBay.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("module_bay")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("module_bay")
			}
			return err
		}
	}

	return nil
}

func (m *Module) contextValidateModuleType(ctx context.Context, formats strfmt.Registry) error {

	if m.ModuleType != nil {
		if err := m.ModuleType.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("module_type")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("module_type")
			}
			return err
		}
	}

	return nil
}

func (m *Module) contextValidateStatus(ctx context.Context, formats strfmt.Registry) error {

	if m.Status != nil {
		if err := m.Status.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("status")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("status")
			}
			return err
		}
	}

	return nil
}

func (m *Module) contextValidateTags(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Tags); i++ {

		if m.Tags[i] != nil {
			if err := m.Tags[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("tags" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("tags" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Module) contextValidateURL(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "url", "body", strfmt.URI(m.URL)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Module) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Module) UnmarshalBinary(b []byte) error {
	var res Module
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ModuleStatus Status
//
// swagger:model ModuleStatus
type ModuleStatus struct {

	// label
	// Required: true
	// Enum: [Offline Active Planned Staged Failed Decommissioning]
	Label *string `json:"label"`

	// value
	// Required: true
	// Enum: [offline active planned staged failed decommissioning]
	Value *string `json:"value"`
}

// Validate validates this module status
func (m *ModuleStatus) Validate(formats strfmt.Registry) error {
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

var moduleStatusTypeLabelPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["Offline","Active","Planned","Staged","Failed","Decommissioning"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		moduleStatusTypeLabelPropEnum = append(moduleStatusTypeLabelPropEnum, v)
	}
}

const (

	// ModuleStatusLabelOffline captures enum value "Offline"
	ModuleStatusLabelOffline string = "Offline"

	// ModuleStatusLabelActive captures enum value "Active"
	ModuleStatusLabelActive string = "Active"

	// ModuleStatusLabelPlanned captures enum value "Planned"
	ModuleStatusLabelPlanned string = "Planned"

	// ModuleStatusLabelStaged captures enum value "Staged"
	ModuleStatusLabelStaged string = "Staged"

	// ModuleStatusLabelFailed captures enum value "Failed"
	ModuleStatusLabelFailed string = "Failed"

	// ModuleStatusLabelDecommissioning captures enum value "Decommissioning"
	ModuleStatusLabelDecommissioning string = "Decommissioning"
)

// prop value enum
func (m *ModuleStatus) validateLabelEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, moduleStatusTypeLabelPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *ModuleStatus) validateLabel(formats strfmt.Registry) error {

	if err := validate.Required("status"+"."+"label", "body", m.Label); err != nil {
		return err
	}

	// value enum
	if err := m.validateLabelEnum("status"+"."+"label", "body", *m.Label); err != nil {
		return err
	}

	return nil
}

var moduleStatusTypeValuePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["offline","active","planned","staged","failed","decommissioning"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		moduleStatusTypeValuePropEnum = append(moduleStatusTypeValuePropEnum, v)
	}
}

const (

	// ModuleStatusValueOffline captures enum value "offline"
	ModuleStatusValueOffline string = "offline"

	// ModuleStatusValueActive captures enum value "active"
	ModuleStatusValueActive string = "active"

	// ModuleStatusValuePlanned captures enum value "planned"
	ModuleStatusValuePlanned string = "planned"

	// ModuleStatusValueStaged captures enum value "staged"
	ModuleStatusValueStaged string = "staged"

	// ModuleStatusValueFailed captures enum value "failed"
	ModuleStatusValueFailed string = "failed"

	// ModuleStatusValueDecommissioning captures enum value "decommissioning"
	ModuleStatusValueDecommissioning string = "decommissioning"
)

// prop value enum
func (m *ModuleStatus) validateValueEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, moduleStatusTypeValuePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *ModuleStatus) validateValue(formats strfmt.Registry) error {

	if err := validate.Required("status"+"."+"value", "body", m.Value); err != nil {
		return err
	}

	// value enum
	if err := m.validateValueEnum("status"+"."+"value", "body", *m.Value); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this module status based on context it is used
func (m *ModuleStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ModuleStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModuleStatus) UnmarshalBinary(b []byte) error {
	var res ModuleStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
