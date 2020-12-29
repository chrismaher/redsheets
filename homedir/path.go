package homedir

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

// Root returns the user's home directory
func Home() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

// FullPath joins variadic path inputs and returns
// a fully qualified filepath based on the user's home directory
func AbsPath(paths ...string) (string, error) {
	path := filepath.Join(paths...)

	// return if the path is already absolute
	if filepath.IsAbs(path) {
		return path, nil
	}

	home, err := Home()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, path), nil
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateIfNotExists(path string, contents []byte) error {
	if !Exists(path) {
		if err := ioutil.WriteFile(path, contents, 0644); err != nil {
			return err
		}
	}
	return nil
}
