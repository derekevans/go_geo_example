
package cmd

import (
	// "fmt"
	
	"geo/internal/load"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load data exported from MyJohnDeere.",
	Long: "Load data exported from MyJohnDeere into a PostGIS database.",
	Run: func(cmd *cobra.Command, args []string) {
		inDir, _ := cmd.Flags().GetString("dir")
		load.LoadMyJD(inDir)
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().String("dir", "", "Path to input directory")
	loadCmd.MarkFlagRequired("dir")
}
