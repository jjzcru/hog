package command

import (
	"github.com/jjzcru/hog/internal/command/get"
	"github.com/jjzcru/hog/internal/command/set"
	"github.com/jjzcru/hog/internal/command/share"
	"os"

	"github.com/jjzcru/hog/internal/command/add"
	"github.com/jjzcru/hog/internal/command/bucket"
	"github.com/jjzcru/hog/internal/command/remove"
	"github.com/jjzcru/hog/internal/command/start"
	"github.com/jjzcru/hog/internal/command/update"
	"github.com/jjzcru/hog/internal/command/version"
	"github.com/spf13/cobra"
)

func Execute() error {
	var rootCmd = &cobra.Command{
		Use:   "hog",
		Short: "Share your file system 🐗",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
	}

	rootCmd.AddCommand(
		add.Command(),
		remove.Command(),
		update.Command(),
		start.Command(),
		bucket.Command(),
		version.Command(),
		set.Command(),
		get.Command(),
		share.Command(),
	)

	return rootCmd.Execute()
}
