package app

import (
	"context"
	"user-syncer/pkg/db"
	"user-syncer/pkg/domain"
	"user-syncer/pkg/provider"
	"user-syncer/pkg/syncer"

	"k8s.io/klog/v2"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	kubeOptions := domain.NewKubeOptions()
	httpProviderOptions := provider.NewOptions()
	rootCmd := &cobra.Command{
		Use:   "user-syncer",
		Short: "A syncer for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.TODO()

			dbConfig, err := db.NewConfigFromEnv()
			if err != nil {
				klog.Error(err)
				return err
			}
			klog.Infof("http Options %+v", httpProviderOptions)
			klog.Infof("kube Options %+v", kubeOptions)
			klog.Infof("db Options %+v", dbConfig)
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

			userProvider, err := provider.NewHttpProvider(httpProviderOptions)
			if err != nil {
				klog.Error(err)
				return err
			}

			orgProvider, err := provider.NewOrgProvider(httpProviderOptions)
			if err != nil {
				klog.Error(err)
				return err
			}

			ksProvider := provider.NewKSProvider(kubernetesClient, httpProviderOptions.Source)

			task := []*domain.Task{
				{
					Syncer:   ksSyncer,
					Provider: userProvider,
				},
				{
					Syncer:   dbSyncer,
					Provider: ksProvider,
				},
				{
					Syncer:   orgDBSyncer,
					Provider: orgProvider,
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

	return rootCmd
}
