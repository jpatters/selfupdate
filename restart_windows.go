package selfupdate

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jpatters/selfupdate/internal/osext"
	"golang.org/x/sys/windows"
)

func restart(exiter func(error), executable string, requiresAdmin bool) error {
	var err error
	if executable == "" {
		executable, err = osext.Executable()
		if err != nil {
			return err
		}
	}

	if requiresAdmin {
		err = restartAsAdmin(executable)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		_, err = os.StartProcess(executable, os.Args, &os.ProcAttr{
			Dir:   wd,
			Env:   os.Environ(),
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Sys:   &syscall.SysProcAttr{},
		})
	}

	if exiter != nil {
		exiter(err)
	} else if err == nil {
		os.Exit(0)
	}
	return err
}

func restartAsAdmin(executable string) error {
	verb := "runas"
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(executable)
	cwdPtr, _ := syscall.UTF16PtrFromString(filepath.Dir(executable))
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	return windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
}
