package start

import (
	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/server"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/spf13/cobra"
)

// Command returns a cobra command for `init` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a server",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().IntP("port", "p", 1618, "Port where the server is going to run")
	cmd.Flags().BoolP("query", "q", false, "Enables graphql playground endpoint 🎮")
	cmd.Flags().BoolP("auth", "a", false, "Enables authorization for endpoints")
	cmd.Flags().StringP("token", "t", "", "Set a specific token for authorization")
	cmd.Flags().BoolP("detached", "d", false, "Run in detached mode and return the PID")

	return cmd
}

func run(cmd *cobra.Command) error {
	isDetached, err := cmd.Flags().GetBool("detached")
	if err != nil {
		return err
	}

	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}

	isQueryEnabled, err := cmd.Flags().GetBool("query")
	if err != nil {
		return err
	}

	isAuthEnable, err := cmd.Flags().GetBool("auth")
	if err != nil {
		return err
	}

	token, err := cmd.Flags().GetString("token")
	if err != nil {
		return err
	}

	if isAuthEnable {
		if len(token) == 0 {
			token = utils.GetToken()
		}
	} else {
		token = ""
	}

	if isDetached {
		return detached(token)
	}

	hogPath, err := hog.GetPath()
	if err != nil {
		return nil
	}

	return server.Start(port, hogPath, isQueryEnabled, token)
}
