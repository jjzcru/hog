package set

import (
	"errors"
	"fmt"
	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// Command returns a cobra command for `set` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Edit hog configuration values",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
	}

	cmd.AddCommand(
		protocol(),
		domain(),
		port(),
	)

	return cmd
}

func protocol() *cobra.Command {
	return &cobra.Command{
		Use:   "protocol",
		Short: "Set protocol value",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			proto := args[0]
			if proto != "http" && proto != "https" {
				utils.PrintError(errors.New("the only valid protocols are http or https"))
				return
			}

			h, err := hog.Get()
			if err != nil {
				utils.PrintError(err)
				return
			}
			h.Protocol = proto

			err = hog.Save(h)
			if err != nil {
				utils.PrintError(err)
				return
			}
		},
	}
}

func domain() *cobra.Command {
	return &cobra.Command{
		Use:   "domain",
		Short: "Set domain value",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			h, err := hog.Get()
			if err != nil {
				utils.PrintError(err)
				return
			}
			h.Domain = args[0]
			err = hog.Save(h)
			if err != nil {
				utils.PrintError(err)
				return
			}
		},
	}
}

func port() *cobra.Command {
	return &cobra.Command{
		Use:   "port",
		Short: "Set port value",
		Args:  validatePort(),
		Run: func(cmd *cobra.Command, args []string) {
			h, err := hog.Get()
			if err != nil {
				utils.PrintError(err)
				return
			}

			port, err := strconv.Atoi(args[0])
			if err != nil {
				utils.PrintError(err)
				return
			}

			h.Port = port
			err = hog.Save(h)
			if err != nil {
				utils.PrintError(err)
				return
			}
		},
	}
}

func validatePort() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg, received %d", len(args))
		}

		_, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("'%s' is not an valid port number", args[0])
		}

		return nil
	}
}
