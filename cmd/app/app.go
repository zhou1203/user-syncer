package app

import (
	"context"
	"k8s.io/klog/v2"
	"user-syncer/pkg/db"
	"user-syncer/pkg/domain"
	"user-syncer/pkg/provider"
	"user-syncer/pkg/syncer"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	kubeOptions := domain.NewKubeOptions()
	httpProviderOptions := provider.NewOptions()
	dbConfig := &db.Config{}
	rootCmd := &cobra.Command{
		Use:   "user-syncer",
		Short: "A syncer for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.TODO()
			client, err := domain.NewKubernetesClient(kubeOptions)
			if err != nil {
				klog.Error(err)
				return err
			}
			dbConnect, err := db.Connect(dbConfig, nil)
			if err != nil {
				klog.Error(err)
				return err
			}

			ksSyncer := syncer.NewKSSyncer(client)
			dbSyncer := syncer.NewDBSyncer(dbConnect)

			fakeProvider := provider.NewFakeProvider(&provider.Options{Source: "test-source"})
			err = domain.NewSyncerOperator(ksSyncer).Sync(ctx, fakeProvider)
			if err != nil {
				klog.Error(err)
				return err
			}

			err = domain.NewSyncerOperator(dbSyncer).Sync(ctx, provider.NewKSProvider(client, "test-source"))
			if err != nil {
				klog.Error(err)
				return err
			}

			return nil
		},
	}
	fs := rootCmd.Flags()
	fs.AddFlagSet(kubeOptions.Flags())
	fs.AddFlagSet(httpProviderOptions.Flags())
	fs.AddFlagSet(dbConfig.Flags())

	return rootCmd
}
