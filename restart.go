package selfupdate

import (
	"os"

	"github.com/jpatters/selfupdate/internal/osext"
)

func restart(exiter func(error), executable string, requiresAdmin bool, inheritArgs bool) error {
	var err error
	if executable == "" {
		executable, err = osext.Executable()
		if err != nil {
			return err
		}
	}

	args := []string{}
	if inheritArgs {
		args = os.Args[1:]
	}

	if requiresAdmin {
		err = restartAsAdmin(executable, args)
	} else {
		err = restartAsUser(executable, args)
	}

	if exiter != nil {
		exiter(err)
	} else if err == nil {
		os.Exit(0)
	}
	return err
}
