package update

import (
	"fmt"
	"path/filepath"

	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/spf13/cobra"
)

// Command returns a cobra command for `update` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the files in a bucket by its id",
		Args:  cobra.ExactArgs(2),
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

	bucketID, err := getBucketID(h, id)
	if err != nil {
		return err
	}

	for i, file := range files {
		filePath, err := filepath.Abs(file)
		if err != nil {
			return err
		}

		if !utils.IsPathExist(filePath) {
			return fmt.Errorf("path '%s' is not valid or do not exist", filePath)
		}

		files[i] = filePath
	}

	h.Buckets[bucketID] = files

	return hog.Save(h)
}

func getBucketID(h hog.Hog, id string) (string, error) {
	var ids []string
	response := ""

	for k := range h.Buckets {
		isSubstring, err := utils.IsSubstring(id, k)
		if err != nil {
			return response, err
		}

		if isSubstring {
			ids = append(ids, k)
		}
	}

	ids = utils.RemoveDuplicate(ids)

	if len(ids) == 0 {
		return response, fmt.Errorf("no bucket was found with id '%s'", id)
	}

	if len(ids) > 1 {
		return response, fmt.Errorf("more than one bucket exists with the id '%s'", id)
	}

	response = ids[0]

	return response, nil
}
