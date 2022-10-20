package app

import (
	"user-generator/pkg/db"
	"user-generator/pkg/domain"
	"user-generator/pkg/provider"

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

			return nil
		},
	}
	fs := rootCmd.Flags()
	fs.AddFlagSet(kubeOptions.Flags())
	fs.AddFlagSet(httpProviderOptions.Flags())
	fs.AddFlagSet(dbConfig.Flags())

	return rootCmd
}
