// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Request request
// swagger:model Request
type Request struct {

	// dialog state
	DialogState string `json:"dialogState,omitempty"`

	// intent
	Intent *Intent `json:"intent,omitempty"`

	// locale
	Locale string `json:"locale,omitempty"`

	// request Id
	RequestID string `json:"requestId,omitempty"`

	// timestamp
	Timestamp string `json:"timestamp,omitempty"`

	// type
	Type string `json:"type,omitempty"`
}

// Validate validates this request
func (m *Request) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateIntent(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Request) validateIntent(formats strfmt.Registry) error {

	if swag.IsZero(m.Intent) { // not required
		return nil
	}

	if m.Intent != nil {

		if err := m.Intent.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("intent")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Request) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Request) UnmarshalBinary(b []byte) error {
	var res Request
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
