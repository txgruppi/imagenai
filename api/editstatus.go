//go:generate go run github.com/abice/go-enum@v0.6.0 -f=$GOFILE --marshal
package api

// ENUM(
// Unknown="",
// Pending,
// InProgress=In Progress,
// Failed,
// Completed,
// )
type editStatus string
