package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAmazonFetchCommand(use string) *cobra.Command {
	var dest string
	cmd := &cobra.Command{
		Use:   use,
		Short: fmt.Sprintf("Fetch %s data", use),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := fetchAmazon(use, dest, 3)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&dest, "dest", "~/Desktop", "dir to save spotify json")

	return cmd
}

func main() {
	rootCmd := &cobra.Command{Use: "fetch"}

	amazonCmd := &cobra.Command{Use: "amazon"}
	amazonCmd.AddCommand(
		newAmazonFetchCommand("book"),
	)

	rootCmd.AddCommand(
		amazonCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
