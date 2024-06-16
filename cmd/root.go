package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gwolves/feedy/cmd/publish"
	"github.com/gwolves/feedy/cmd/runserver"
	"github.com/gwolves/feedy/cmd/subscribe"
)

func MustExecute() {
	rootCmd := newCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "feedy",
		Short: "feedy is rss/atom feed integration for channel talk",
		Long:  "feedy is rss/atom feed integration for channel talk",
	}

	cmd.AddCommand(runserver.NewCommand())
	cmd.AddCommand(subscribe.NewCommand())
	cmd.AddCommand(publish.NewCommand())

	return &cmd
}
