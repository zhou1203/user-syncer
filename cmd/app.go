package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"user-export/pkg/httpprovider"
	"user-export/pkg/ksgenerator"
)

func NewCommand() *cobra.Command {
	keOptions := ksgenerator.NewOptions()
	httpProviderOptions := httpprovider.NewOptions()
	rootCmd := &cobra.Command{
		Use:   "user-export",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			generator, err := ksgenerator.NewKSGenerator(keOptions)
			if err != nil {
				return err
			}
			err = generator.Generate(context.Background(), httpprovider.NewFakeProvider(httpProviderOptions))
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
