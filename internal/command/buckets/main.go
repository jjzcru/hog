package buckets

import (
	"fmt"
	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

// Command returns a cobra command for `init` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buckets",
		Short: "Return a list of all the buckets",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := run()
			if err != nil {
				utils.PrintError(err)
				return
			}
		},
	}

	return cmd
}

func run() error {

	hogPath, err := hog.GetPath()
	if err != nil {
		return err
	}

	h, err := getHog(hogPath)
	if err != nil {
		return err
	}

	data := transform(h)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Bucket", "Files"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()

	return nil
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

func transform(h hog.Hog) [][]string {
	var table [][]string
	for id, files := range h.Buckets {
		for _, file := range files {
			table = append(table, []string{id, file})
		}
	}

	return table
}
