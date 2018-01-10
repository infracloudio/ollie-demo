// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Attributes attributes
// swagger:model Attributes
type Attributes struct {

	// command
	Command string `json:"command,omitempty"`

	// direction
	Direction uint16 `json:"direction,omitempty"`

	// duration
	Duration uint16 `json:"duration,omitempty"`

	// speed
	Speed uint8 `json:"speed,omitempty"`
}

// Validate validates this attributes
func (m *Attributes) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Attributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Attributes) UnmarshalBinary(b []byte) error {
	var res Attributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}