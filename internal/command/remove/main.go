package remove

import (
"fmt"
	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/spf13/cobra"
	"time"
)

// Command returns a cobra command for `init` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove file reference by id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd, args[0])
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().Duration("ttl", 0, "Remove a bucket after a period of time")
	cmd.Flags().String("deadline", "", "Remove a bucket a particular time")
	cmd.Flags().BoolP("detached", "d", false, "")
	return cmd
}

func run(cmd *cobra.Command, id string) error {
	isDetached, err := cmd.Flags().GetBool("detached")
	if err != nil {
		return err
	}

	ttl, err := cmd.Flags().GetDuration("ttl")
	if err != nil {
		return err
	}

	deadline, err := cmd.Flags().GetString("deadline")
	if err != nil {
		return err
	}

	hogPath, err := hog.GetPath()
	if err != nil {
		return err
	}

	h, err := getHog(hogPath)
	if err != nil {
		return err
	}

	if len(deadline) > 0 {
		deadlineTime, err := utils.GetTimeFromString(deadline)
		if err != nil {
			return err
		}

		now := time.Now()
		if deadlineTime.Before(now) {
			return fmt.Errorf("deadline can't be before of current time")
		}
	}

	err = validate(h, id)
	if err != nil {
		return err
	}

	if isDetached {
		return Detached()
	}

	delayCmd(ttl, deadline)

	return remove(hogPath, h, id)
}

func getHog(hogPath string) (hog.Hog, error) {
	var h hog.Hog

	if !utils.IsPathExist(hogPath) {
		return h, fmt.Errorf("hog path '%s' do not exist", hogPath)
	}

	h, err := hog.FromPath(hogPath)
	if err != nil {
		return h, err
	}

	return h, nil

}

func validate(h hog.Hog, id string) error{
	if _, ok := h.Buckets[id]; !ok {
		return fmt.Errorf("bucket with id '%s' do not exist", id)
	}
	return nil
}

func remove(hogPath string, h hog.Hog, id string) error{
	delete(h.Buckets, id)
	return hog.Save(hogPath, h)
}

func delayCmd(ttl time.Duration, deadline string) {
	var deadlineDuration time.Duration
	var ttlDuration time.Duration
	var sleepDuration time.Duration

	if len(deadline) > 0 {
		startTime, _ := utils.GetTimeFromString(deadline)
		now := time.Now()
		diff := startTime.Sub(now)

		deadlineDuration = diff
	}

	if ttl > 0 {
		ttlDuration = ttl
	}

	if deadlineDuration > 0 && ttlDuration > 0 {
		if deadlineDuration < ttlDuration {
			sleepDuration = deadlineDuration
		} else {
			sleepDuration = ttlDuration
		}
	} else if deadlineDuration > 0 {
		sleepDuration = deadlineDuration
	} else if ttlDuration > 0 {
		sleepDuration = ttlDuration
	}

	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
	}
}
