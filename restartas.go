//go:build !windows
// +build !windows

package selfupdate

import (
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

func restartAsAdmin(executable string, args []string) error {
	return restartAsUser(executable, args)
}

func restartAsUser(executable string, _ []string) error {
	filename := filepath.Base(executable)

	inFile, err := os.Open(executable)
	if err != nil {
		return err
	}
	defer inFile.Close()

	u, err := user.Current()
	if err != nil {
		return err
	}

	downloadDir := filepath.Join(u.HomeDir, "Downloads")
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		return err
	}

	outFile, err := os.Create(filepath.Join(downloadDir, filename))
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		return err
	}

	return exec.Command("open", outFile.Name()).Start()
}
