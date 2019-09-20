package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func Execute() error {
	var rootCmd = &cobra.Command{
		Use:   "gopenapi",
		Short: "An OpenAPI utility for Go",
	}
	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "The generator utility",
	}

	var format string
	var output string
	var generateSpecCmd = &cobra.Command{
		Use:   "spec [optional path]",
		Short: "The spec generator utility",
		Long:  "The spec generator utility can GenerateSpec specifications from source code",

		Run: func(cmd *cobra.Command, args []string) {
			if err := GenerateSpec(format, output, args); err != nil {
				println(err)
				os.Exit(1)
			}
		},
	}
	generateSpecCmd.Flags().StringVarP(&format, "format", "f", "json", "The format of the output. May be json or yaml")
	generateSpecCmd.Flags().StringVarP(&output, "output", "o", "-", "Where the output should be directed. May be '-' (stdout) or a path to a file")

	generateCmd.AddCommand(generateSpecCmd)
	rootCmd.AddCommand(generateCmd)

	return rootCmd.Execute()
}
