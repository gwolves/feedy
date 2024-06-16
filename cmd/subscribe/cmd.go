package subscribe

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/gwolves/feedy/internal/app"
)

func NewCommand() *cobra.Command {
	var (
		channelID string
		groupID   string
		url       string
		name      string
	)

	cmd := cobra.Command{
		Use:   "subscribe",
		Short: "subscribe RSS/Atom feed. You can assign a name for notification bot.",
		Run: func(cmd *cobra.Command, args []string) {
			u := app.MustInitUsecase()

			ctx := context.Background()
			err := u.Subscribe(ctx, channelID, groupID, url, name)
			if err != nil {
				log.Println("subscribe error", err)
			}
		},
	}

	cmd.Flags().StringVar(&channelID, "channel", "", "channel id to subscribe feed")
	cmd.Flags().StringVar(&groupID, "group", "", "group id to subscribe feed")
	cmd.Flags().StringVar(&url, "url", "", "url of target feed")
	cmd.Flags().StringVar(&name, "name", "", "alias for subscription")
	cmd.MarkFlagRequired("channel")
	cmd.MarkFlagRequired("group")
	cmd.MarkFlagRequired("url")

	return &cmd
}
