package app

import (
	"context"
	"github.com/spf13/cobra"
	"user-generator/pkg/httpprovider"
	"user-generator/pkg/ksgenerator"
)

func NewCommand() *cobra.Command {
	keOptions := ksgenerator.NewOptions()
	httpProviderOptions := httpprovider.NewOptions()
	rootCmd := &cobra.Command{
		Use:   "user-generator",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			generator, err := ksgenerator.NewKSGenerator(keOptions)
			if err != nil {
				return err
			}
			provider, err := httpprovider.NewHttpProvider(httpProviderOptions)
			if err != nil {
				return err
			}
			err = generator.Generate(context.Background(), provider)
			if err != nil {
				return err
			}
			return nil
		},
	}
	fs := rootCmd.Flags()
	fs.AddFlagSet(keOptions.Flags())
	fs.AddFlagSet(httpProviderOptions.Flags())

	return rootCmd
}
