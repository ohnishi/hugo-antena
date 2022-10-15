package main

import (
	"fmt"
	"time"

	"github.com/ohnishi/antena/backend/common/command"
	"github.com/ohnishi/antena/backend/common/env"
	"github.com/spf13/cobra"
)

func newAmazonPublishCommand(use string) *cobra.Command {
	var (
		dates []string
		src   string
		dest  string
	)
	cmd := &cobra.Command{
		Use:   use,
		Short: fmt.Sprintf("Publish %s items", use),
		Args:  cobra.NoArgs,
		RunE: command.WithLoggingE(func(cmd *cobra.Command, args []string) error {
			return command.EachDate(dates, func(date time.Time) error {
				return publishAmazonItems(use, src, dest, date)
			})
		}),
	}
	command.SetDatesFlag(cmd.Flags(), &dates, "date for which the URL list file(s) is generated")
	_ = cmd.MarkFlagRequired("date")
	cmd.Flags().StringVar(&src, "src", env.DataDir("itbook/transform"), "output path into which 5ch threads is written.")
	cmd.Flags().StringVar(&dest, "dest", "./itbook/content/posts", "output path into which 5ch threads is written.")

	return cmd
}

func main() {
	rootCmd := &cobra.Command{Use: "publish"}

	amazonCmd := &cobra.Command{Use: "amazon"}
	amazonCmd.AddCommand(
		newAmazonPublishCommand("book"),
	)

	rootCmd.AddCommand(
		amazonCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
