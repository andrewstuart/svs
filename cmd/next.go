package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"
)

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "increment the version",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ver, err := fromGit()
		if len(args) > 1 && args[1] == "-" {
			ver, err = cur(os.Stdin)
		}
		if err != nil {
			log.Fatal(err)
		}

		vPref := strings.HasPrefix(ver, "v")
		v, _ := semver.ParseTolerant(ver)

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
		if vPref {
			fmt.Println("v" + v.String())
			return
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
