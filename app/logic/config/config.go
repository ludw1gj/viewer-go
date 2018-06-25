// Package config contains miscellaneous configuration values and methods to get/set them.
package config

var usersDirectory string

// SetUsersDirectory sets the user's directory.
func SetUsersDirectory(dir string) {
	usersDirectory = dir
}

// GetUsersDirectory gets the user's directory.
func GetUsersDirectory() string {
	return usersDirectory
}
