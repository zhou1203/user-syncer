package app

import (
	"context"
	"github.com/spf13/cobra"
	"user-generator/pkg"
	"user-generator/pkg/db"
	"user-generator/pkg/provider"
	"user-generator/pkg/syncer"
)

func NewCommand() *cobra.Command {
	kubeOptions := pkg.NewKubeOptions()
	httpProviderOptions := provider.NewOptions()
	rootCmd := &cobra.Command{
		Use:   "user-syncer",
		Short: "A syncer for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dbConfig, err := db.NewConfigFromEnv()
			if err != nil {
				return err
			}

			database, err := db.Connect(dbConfig, nil)
			if err != nil {
				return err
			}

			client, err := pkg.NewKubernetesClient(kubeOptions)
			if err != nil {
				return err
			}

			ksGenerator := syncer.NewKSSyncer(client)

			httpProvider, err := provider.NewHttpProvider(httpProviderOptions)
			if err != nil {
				return err
			}
			err = ksGenerator.Sync(context.Background(), httpProvider)
			if err != nil {
				return err
			}

			err = syncer.NewDBSyncer(database).Sync(context.Background(), provider.NewKSProvider(client))
			if err != nil {
				return err
			}
			return nil
		},
	}
	fs := rootCmd.Flags()
	fs.AddFlagSet(kubeOptions.Flags())
	fs.AddFlagSet(httpProviderOptions.Flags())

	return rootCmd
}
