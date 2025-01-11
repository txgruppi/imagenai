//go:generate go run github.com/abice/go-enum@v0.6.0 -f=$GOFILE --marshal
package api

// ENUM(RAW,DNG,JPEG=JPG)
type imageFormat string
