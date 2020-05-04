package add

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
)

// Command returns a cobra command for `add` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Group files in a bucket",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			response, err := run(cmd, args)
			if err != nil {
				utils.PrintError(err)
				return
			}

			if len(response) > 0 {
				fmt.Println(response)
			}
		},
	}

	cmd.Flags().Duration("ttl", 0, "Remove a bucket after a period of time")
	cmd.Flags().BoolP("url", "u", false, "Return a share url as response")
	cmd.Flags().BoolP("qr", "q", false, "Return a qr code with the url as response")

	return cmd
}

func run(cmd *cobra.Command, args []string) (string, error) {
	ttl, err := cmd.Flags().GetDuration("ttl")
	if err != nil {
		return "", err
	}

	isUrl, err := cmd.Flags().GetBool("url")
	if err != nil {
		return "", err
	}

	isQr, err := cmd.Flags().GetBool("qr")
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

	// Evaluate TTL
	if ttl > 0 {
		rmvCmd := fmt.Sprintf("remove %s --ttl %s", bucketID, ttl.String())
		cmd := exec.Command("hog", strings.Split(rmvCmd, " ")...)

		err := cmd.Start()
		if err != nil {
			return "", err
		}

		err = cmd.Process.Release()
		if err != nil {
			return "", err
		}
	}

	var h hog.Hog
	var url string

	if isUrl || isQr {
		h, err = hog.Get()
		if err != nil {
			return "", err
		}

		url = fmt.Sprintf("%s://%s:%d/download/%s", h.Protocol, h.Domain, h.Port, bucketID)
		bucketID = ""
	}

	if isQr {
		qrterminal.Generate(url, qrterminal.M, os.Stdout)
	}

	if isUrl {
		fmt.Println(url)
	}

	return bucketID, nil
}
