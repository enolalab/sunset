// Package model defines the application data models.
package model

// User represents a user in the system.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// FullName returns the display name for the user.
func (u *User) FullName() string {
	return u.Name
}

// isValid checks whether the user data is valid.
func (u *User) isValid() bool {
	return u.Name != "" && u.Email != ""
}
