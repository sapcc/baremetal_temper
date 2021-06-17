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
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewDcimPowerPortTemplatesDeleteParams creates a new DcimPowerPortTemplatesDeleteParams object
// with the default values initialized.
func NewDcimPowerPortTemplatesDeleteParams() *DcimPowerPortTemplatesDeleteParams {
	var ()
	return &DcimPowerPortTemplatesDeleteParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDcimPowerPortTemplatesDeleteParamsWithTimeout creates a new DcimPowerPortTemplatesDeleteParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDcimPowerPortTemplatesDeleteParamsWithTimeout(timeout time.Duration) *DcimPowerPortTemplatesDeleteParams {
	var ()
	return &DcimPowerPortTemplatesDeleteParams{

		timeout: timeout,
	}
}

// NewDcimPowerPortTemplatesDeleteParamsWithContext creates a new DcimPowerPortTemplatesDeleteParams object
// with the default values initialized, and the ability to set a context for a request
func NewDcimPowerPortTemplatesDeleteParamsWithContext(ctx context.Context) *DcimPowerPortTemplatesDeleteParams {
	var ()
	return &DcimPowerPortTemplatesDeleteParams{

		Context: ctx,
	}
}

// NewDcimPowerPortTemplatesDeleteParamsWithHTTPClient creates a new DcimPowerPortTemplatesDeleteParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDcimPowerPortTemplatesDeleteParamsWithHTTPClient(client *http.Client) *DcimPowerPortTemplatesDeleteParams {
	var ()
	return &DcimPowerPortTemplatesDeleteParams{
		HTTPClient: client,
	}
}

/*DcimPowerPortTemplatesDeleteParams contains all the parameters to send to the API endpoint
for the dcim power port templates delete operation typically these are written to a http.Request
*/
type DcimPowerPortTemplatesDeleteParams struct {

	/*ID
	  A unique integer value identifying this power port template.

	*/
	ID int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the dcim power port templates delete params
func (o *DcimPowerPortTemplatesDeleteParams) WithTimeout(timeout time.Duration) *DcimPowerPortTemplatesDeleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the dcim power port templates delete params
func (o *DcimPowerPortTemplatesDeleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the dcim power port templates delete params
func (o *DcimPowerPortTemplatesDeleteParams) WithContext(ctx context.Context) *DcimPowerPortTemplatesDeleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the dcim power port templates delete params
func (o *DcimPowerPortTemplatesDeleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the dcim power port templates delete params
func (o *DcimPowerPortTemplatesDeleteParams) WithHTTPClient(client *http.Client) *DcimPowerPortTemplatesDeleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the dcim power port templates delete params
func (o *DcimPowerPortTemplatesDeleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the dcim power port templates delete params
func (o *DcimPowerPortTemplatesDeleteParams) WithID(id int64) *DcimPowerPortTemplatesDeleteParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the dcim power port templates delete params
func (o *DcimPowerPortTemplatesDeleteParams) SetID(id int64) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DcimPowerPortTemplatesDeleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt64(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
