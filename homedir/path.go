package homedir

import (
	"os/user"
	"path"
)

// Base returns the user's home directory
func Base() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

// FullPath joins variadic path inputs and returns
// a fully qualified filepath based on the user's home directory
func FullPath(paths ...string) (string, error) {
	dir, err := Base()
	if err != nil {
		return "", err
	}

	var fullPath = dir
	for _, p := range paths {
		fullPath = path.Join(fullPath, p)
	}

	return fullPath, nil
}
