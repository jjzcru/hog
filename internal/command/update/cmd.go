package update

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
		Use:   "update",
		Short: "Update the files in a bucket by its id",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			files := args[1:]
			err := run(id, files)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	return cmd
}

func run(id string, files []string) error {
	h, err := hog.Get()
	if err != nil {
		return err
	}

	if _, ok := h.Buckets[id]; !ok {
		return fmt.Errorf("bucket with id '%s' do not exist", id)
	}

	for i, file := range files {
		filePath, err := filepath.Abs(file)
		if err != nil {
			return err
		}

		if !utils.IsPathExist(filePath) {
			return fmt.Errorf("path %s is not valid or do not exist", filePath)
		}

		files[i] = filePath
	}

	h.Buckets[id] = files

	return hog.Save(h)
}
