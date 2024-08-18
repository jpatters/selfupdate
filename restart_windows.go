package selfupdate

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

func restartAsUser(executable string, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = os.StartProcess(executable, args, &os.ProcAttr{
		Dir:   wd,
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Sys:   &syscall.SysProcAttr{},
	})
	return err
}

func restartAsAdmin(executable string, args []string) error {
	verb := "runas"
	arg := strings.Join(args, " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(executable)
	cwdPtr, _ := syscall.UTF16PtrFromString(filepath.Dir(executable))
	argPtr, _ := syscall.UTF16PtrFromString(arg)

	var showCmd int32 = 1 //SW_NORMAL

	return windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
}
