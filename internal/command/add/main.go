package add

import (
	"fmt"
	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/spf13/cobra"
	"path/filepath"
)

// Command returns a cobra command for `init` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Make file/s accessible to the service",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			groupID, err := run(args)
			if err != nil {
				utils.PrintError(err)
				return
			}

			fmt.Println(groupID)
		},
	}

	return cmd
}

func run(args []string) (string, error) {
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

	return hog.AddFiles(files)
}
