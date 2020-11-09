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

// FullPath returns a fully qualified filepath
// based on the user's home directory
func FullPath(filename string) (string, error) {
	dir, err := Base()
	if err != nil {
		return "", err
	}

	return path.Join(dir, filename), nil
}
