//go:generate go run github.com/abice/go-enum@v0.6.0 -f=$GOFILE --marshal
package api

// ENUM(
// Crop=crop,
// HDRMerge=hdr_merge,
// PerspectiveCorrection=perspective_correction,
// PortraitCrop=portrait_crop,
// Straighten=straighten,
// SubjectMask=subject_mask,
// SmoothSkin=smooth_skin,
// )
type tool string
