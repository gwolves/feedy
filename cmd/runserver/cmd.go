package runserver

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/gwolves/feedy/internal/app"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "runserver",
		Short: "run api server",
		Run: func(cmd *cobra.Command, args []string) {
			server := app.MustInitHTTPServer()
			if err := server.Serve(); err != nil {
				log.Println("error", err)
			}
		},
	}
}
