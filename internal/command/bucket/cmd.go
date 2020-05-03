package bucket

import (
	"fmt"
	"os"

	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Command returns a cobra command for `bucket` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bucket",
		Short: "Display the buckets and their files",
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
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)

	table.SetAutoWrapText(false)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
	)

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
		if len(files) == 0 {
			table = append(table, []string{id, "EMPTY"})
		}
		for _, file := range files {
			table = append(table, []string{id, file})
		}
	}

	return table
}
