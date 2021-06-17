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
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/netbox-community/go-netbox/netbox/models"
)

// NewIpamPrefixesCreateParams creates a new IpamPrefixesCreateParams object
// with the default values initialized.
func NewIpamPrefixesCreateParams() *IpamPrefixesCreateParams {
	var ()
	return &IpamPrefixesCreateParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewIpamPrefixesCreateParamsWithTimeout creates a new IpamPrefixesCreateParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewIpamPrefixesCreateParamsWithTimeout(timeout time.Duration) *IpamPrefixesCreateParams {
	var ()
	return &IpamPrefixesCreateParams{

		timeout: timeout,
	}
}

// NewIpamPrefixesCreateParamsWithContext creates a new IpamPrefixesCreateParams object
// with the default values initialized, and the ability to set a context for a request
func NewIpamPrefixesCreateParamsWithContext(ctx context.Context) *IpamPrefixesCreateParams {
	var ()
	return &IpamPrefixesCreateParams{

		Context: ctx,
	}
}

// NewIpamPrefixesCreateParamsWithHTTPClient creates a new IpamPrefixesCreateParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewIpamPrefixesCreateParamsWithHTTPClient(client *http.Client) *IpamPrefixesCreateParams {
	var ()
	return &IpamPrefixesCreateParams{
		HTTPClient: client,
	}
}

/*IpamPrefixesCreateParams contains all the parameters to send to the API endpoint
for the ipam prefixes create operation typically these are written to a http.Request
*/
type IpamPrefixesCreateParams struct {

	/*Data*/
	Data *models.WritablePrefix

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the ipam prefixes create params
func (o *IpamPrefixesCreateParams) WithTimeout(timeout time.Duration) *IpamPrefixesCreateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the ipam prefixes create params
func (o *IpamPrefixesCreateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the ipam prefixes create params
func (o *IpamPrefixesCreateParams) WithContext(ctx context.Context) *IpamPrefixesCreateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the ipam prefixes create params
func (o *IpamPrefixesCreateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the ipam prefixes create params
func (o *IpamPrefixesCreateParams) WithHTTPClient(client *http.Client) *IpamPrefixesCreateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the ipam prefixes create params
func (o *IpamPrefixesCreateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithData adds the data to the ipam prefixes create params
func (o *IpamPrefixesCreateParams) WithData(data *models.WritablePrefix) *IpamPrefixesCreateParams {
	o.SetData(data)
	return o
}

// SetData adds the data to the ipam prefixes create params
func (o *IpamPrefixesCreateParams) SetData(data *models.WritablePrefix) {
	o.Data = data
}

// WriteToRequest writes these params to a swagger request
func (o *IpamPrefixesCreateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Data != nil {
		if err := r.SetBodyParam(o.Data); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
