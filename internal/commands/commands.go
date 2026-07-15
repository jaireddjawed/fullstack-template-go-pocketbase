// Package commands registers custom CLI commands on the PocketBase root
// command (which already provides `serve`, `migrate`, `superuser`, ...).
package commands

import (
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"

	"github.com/jaireddjawed/fullstack-template-golang/internal/seed"
)

// Register adds all custom commands to the app.
func Register(pb *pocketbase.PocketBase) {
	pb.RootCmd.AddCommand(&cobra.Command{
		Use:   "seed",
		Short: "Seed the database with development data (idempotent)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// The app is bootstrapped by Execute(), but migrations only run
			// on `serve`/`migrate up`, so apply any pending ones first.
			if err := pb.RunAllMigrations(); err != nil {
				return err
			}
			return seed.Run(pb)
		},
	})
}
