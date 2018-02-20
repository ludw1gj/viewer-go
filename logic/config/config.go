package config

var usersDirectory string

func SetUsersDirectory(dir string) {
	usersDirectory = dir
}

func GetUsersDirectory() string {
	return usersDirectory
}
