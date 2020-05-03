package add

import (
	"fmt"
	"path/filepath"

	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/spf13/cobra"
)

// Command returns a cobra command for `init` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Group files in a bucket",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bucketID, err := run(cmd, args)
			if err != nil {
				utils.PrintError(err)
				return
			}

			fmt.Println(bucketID)
		},
	}

	cmd.Flags().Duration("ttl", 0, "Remove a bucket after a period of time")

	return cmd
}

func run(cmd *cobra.Command, args []string) (string, error) {
	ttl, err := cmd.Flags().GetDuration("ttl")
	if err != nil {
		return "", err
	}

	var files []string
	for _, file := range args {
		filePath, err := filepath.Abs(file)
		if err != nil {
			return "", err
		}

		if !utils.IsPathExist(filePath) {
			return "", fmt.Errorf("path %s is not valid or do not exist", filePath)
		}
		files = append(files, filePath)
	}

	bucketID, err := hog.AddFiles(files)
	if err != nil {
		return "", err
	}

	err = bucketTTL(ttl, bucketID)

	return bucketID, nil
}
