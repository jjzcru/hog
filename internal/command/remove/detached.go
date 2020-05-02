package remove

import (
	"github.com/jjzcru/hog/pkg/utils"
	"os"
	"os/exec"
)

// Detached runs ox in detached mode
func Detached() error {
	command := utils.RemoveDetachedFlag(os.Args)
	cmd := exec.Command(command[0], command[1:]...)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}
