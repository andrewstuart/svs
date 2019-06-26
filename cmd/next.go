package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "increment the version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		v := cur(os.Stdin)
		switch args[0] {
		case "major":
			v.Major++
			v.Minor = 0
			v.Patch = 0
		case "minor":
			v.Minor++
			v.Patch = 0
		case "patch":
			v.Patch++
		}
		fmt.Print(v)
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
