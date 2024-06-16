package publish

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/gwolves/feedy/internal/app"
)

func NewCommand() *cobra.Command {
	var (
		id  int64
		all bool
	)

	cmd := cobra.Command{
		Use:   "publish",
		Short: "publish feed for subscription",
		Run: func(cmd *cobra.Command, args []string) {
			u := app.MustInitUsecase()

			ctx := context.Background()

			var err error
			if all {
				err = u.PublishAllFeeds(ctx)
			} else {
				err = u.PublishFeed(ctx, id)
			}

			if err != nil {
				log.Println("publish error", err)
			}
		},
	}

	cmd.Flags().Int64Var(&id, "id", 0, "feed id")
	cmd.Flags().BoolVar(&all, "all", false, "publish all feed")
	cmd.MarkFlagsOneRequired("id", "all")

	return &cmd
}
