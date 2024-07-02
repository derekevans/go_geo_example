
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "geo",
	Short: "This application is an example of using Go with geospatial data.",
	Long: "This application is an example of using Go with geospatial data.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
