package share

import (
	"fmt"
	"github.com/jjzcru/hog/pkg/hog"
	"github.com/jjzcru/hog/pkg/utils"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
	"os"
)

// Command returns a cobra command for `update` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "share",
		Short: "Share a bucket by its id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			err := run(cmd, id)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().BoolP("qr", "q", false, "Return a qr code with the url as response")
	cmd.Flags().String("protocol", "", "Overwrite the protocol in the url")
	cmd.Flags().String("domain", "", "Overwrite the domain in the url")
	cmd.Flags().IntP("port", "p", 0, "Overwrite the port in the url")

	return cmd
}

func run(cmd *cobra.Command, id string) error {
	isQr, err := cmd.Flags().GetBool("qr")
	if err != nil {
		return err
	}

	protocol, err := cmd.Flags().GetString("protocol")
	if err != nil {
		return err
	}

	domain, err := cmd.Flags().GetString("domain")
	if err != nil {
		return err
	}

	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}

	h, err := hog.Get()
	if err != nil {
		return err
	}

	if len(protocol) > 0 {
		if protocol != "http" && protocol != "https" {
			return fmt.Errorf("protocol '%s' is invalid, only http or https", protocol)
		}
		h.Protocol = protocol
	}

	if len(domain) > 0 {
		h.Domain = domain
	}

	if port > 0 {
		h.Port = port
	}

	bucketID, err := getBucketID(h, id)
	if err != nil {
		return err
	}

	url := hog.Url(h, bucketID)
	if isQr {
		qrterminal.Generate(url, qrterminal.M, os.Stdout)
		return nil
	}

	fmt.Println(url)

	return nil
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
