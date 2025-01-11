//go:generate go run github.com/abice/go-enum@v0.6.0 -f=$GOFILE --marshal
package api

// ENUM(
// None="",
// Wedding=WEDDING,
// RealEstate=REAL_ESTATE,
// School=SCHOOL,
// Sports=SPORTS,
// Events=EVENTS,
// Portraits=PORTRAITS,
// NoType=NO_TYPE,
// Other=OTHER,
// LandscapeNature=LANDSCAPE_NATURE,
// FamilyNewborn=FAMILY_NEWBORN,
// Boudoir=BOUDOIR,
// )
type photographyType string
