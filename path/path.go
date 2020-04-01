package path

import (
	"os/user"
	"path"
)

// FullPath returns a fully qualified filepath
// based on the user's home directory
func FullPath(filename string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	home := usr.HomeDir
	return path.Join(home, filename), nil
}
