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
	// ImageFormatRAW is a imageFormat of type RAW.
	ImageFormatRAW imageFormat = "RAW"
	// ImageFormatDNG is a imageFormat of type DNG.
	ImageFormatDNG imageFormat = "DNG"
	// ImageFormatJPEG is a imageFormat of type JPEG.
	ImageFormatJPEG imageFormat = "JPG"
)

var ErrInvalidimageFormat = errors.New("not a valid imageFormat")

// String implements the Stringer interface.
func (x imageFormat) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x imageFormat) IsValid() bool {
	_, err := ParseimageFormat(string(x))
	return err == nil
}

var _imageFormatValue = map[string]imageFormat{
	"RAW": ImageFormatRAW,
	"DNG": ImageFormatDNG,
	"JPG": ImageFormatJPEG,
}

// ParseimageFormat attempts to convert a string to a imageFormat.
func ParseimageFormat(name string) (imageFormat, error) {
	if x, ok := _imageFormatValue[name]; ok {
		return x, nil
	}
	return imageFormat(""), fmt.Errorf("%s is %w", name, ErrInvalidimageFormat)
}

// MarshalText implements the text marshaller method.
func (x imageFormat) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *imageFormat) UnmarshalText(text []byte) error {
	tmp, err := ParseimageFormat(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
