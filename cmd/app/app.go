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
			fakeOptions := &provider.Options{Source: "test-source"}

			kubernetesClient, err := domain.NewKubernetesClient(kubeOptions)
			if err != nil {
				klog.Error(err)
				return err
			}
			database, err := db.Connect(dbConfig, nil)
			if err != nil {
				klog.Error(err)
				return err
			}

			ksSyncer := syncer.NewKSSyncer(kubernetesClient)
			dbSyncer := syncer.NewDBSyncer(database)
			orgDBSyncer := syncer.NewOrgDBSyncer(database)

			fakeUserProvider := provider.NewFakeProvider(fakeOptions)
			fakeOrgProvider := provider.NewFakeOrgProvider()
			ksProvider := provider.NewKSProvider(kubernetesClient, fakeOptions.Source)

			task := []*domain.Task{
				{
					Syncer:   ksSyncer,
					Provider: fakeUserProvider,
				},
				{
					Syncer:   dbSyncer,
					Provider: ksProvider,
				},
				{
					Syncer:   orgDBSyncer,
					Provider: fakeOrgProvider,
				},
			}

			err = domain.NewSyncerOperator(task...).Sync(ctx)
			if err != nil {
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
