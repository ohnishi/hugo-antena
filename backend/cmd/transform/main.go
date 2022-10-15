package main

import (
	"fmt"
	"time"

	"github.com/ohnishi/antena/backend/common/command"
	"github.com/ohnishi/antena/backend/common/env"
	"github.com/spf13/cobra"
)

func newAmazonTransformCommand(use string) *cobra.Command {
	var (
		dates []string
		src   string
		dest  string
	)
	cmd := &cobra.Command{
		Use:   use,
		Short: fmt.Sprintf("Transform %s items", use),
		Args:  cobra.NoArgs,
		RunE: command.WithLoggingE(func(cmd *cobra.Command, args []string) error {
			return command.EachDate(dates, func(date time.Time) error {
				return transformAmazonItems(use, src, dest, date)
			})
		}),
	}
	command.SetDatesFlag(cmd.Flags(), &dates, "date for which the URL list file(s) is generated")
	_ = cmd.MarkFlagRequired("date")
	cmd.Flags().StringVar(&src, "src", env.DataDir("fanza/fetch"), "output path into which 5ch threads is written.")
	cmd.Flags().StringVar(&dest, "dest", env.DataDir("fanza/transform"), "output path into which 5ch threads is written.")

	return cmd
}

func main() {
	rootCmd := &cobra.Command{Use: "transform"}

	amazonCmd := &cobra.Command{Use: "amazon"}
	amazonCmd.AddCommand(
		newAmazonTransformCommand("book"),
	)

	rootCmd.AddCommand(
		amazonCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
