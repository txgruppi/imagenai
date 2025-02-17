// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package api

import (
	"errors"
	"fmt"
)

const (
	// EditStatusUnknown is a editStatus of type Unknown.
	EditStatusUnknown editStatus = ""
	// EditStatusPending is a editStatus of type Pending.
	EditStatusPending editStatus = "Pending"
	// EditStatusInProgress is a editStatus of type InProgress.
	EditStatusInProgress editStatus = "In Progress"
	// EditStatusFailed is a editStatus of type Failed.
	EditStatusFailed editStatus = "Failed"
	// EditStatusCompleted is a editStatus of type Completed.
	EditStatusCompleted editStatus = "Completed"
)

var ErrInvalideditStatus = errors.New("not a valid editStatus")

// String implements the Stringer interface.
func (x editStatus) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x editStatus) IsValid() bool {
	_, err := ParseeditStatus(string(x))
	return err == nil
}

var _editStatusValue = map[string]editStatus{
	"":            EditStatusUnknown,
	"Pending":     EditStatusPending,
	"In Progress": EditStatusInProgress,
	"Failed":      EditStatusFailed,
	"Completed":   EditStatusCompleted,
}

// ParseeditStatus attempts to convert a string to a editStatus.
func ParseeditStatus(name string) (editStatus, error) {
	if x, ok := _editStatusValue[name]; ok {
		return x, nil
	}
	return editStatus(""), fmt.Errorf("%s is %w", name, ErrInvalideditStatus)
}

// MarshalText implements the text marshaller method.
func (x editStatus) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *editStatus) UnmarshalText(text []byte) error {
	tmp, err := ParseeditStatus(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
