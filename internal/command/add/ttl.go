package add

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func bucketTTL(duration time.Duration, id string) error {
	if duration > 0 {
		rmvCmd := fmt.Sprintf("remove %s --ttl %s", id, duration.String())
		cmd := exec.Command("hog", strings.Split(rmvCmd, " ")...)

		err := cmd.Start()
		if err != nil {
			return err
		}

		err = cmd.Process.Release()
		if err != nil {
			return err
		}
	}
	return nil
}
