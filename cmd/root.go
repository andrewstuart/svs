package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/blang/semver"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func cur(rc io.Reader) semver.Version {
	var vs semver.Versions
	br := bufio.NewReader(os.Stdin)

	var err error
	var bs []byte
	for err == nil {
		bs, _, err = br.ReadLine()
		for _, f := range strings.Fields(string(bs)) {
			sv, err := semver.Parse(strings.TrimPrefix(f, "v"))
			if err == nil {
				vs = append(vs, sv)
			}
		}
	}

	sort.Sort(vs)
	return vs[len(vs)-1]
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "svs",
	Short: "semver sort",
	Run: func(c *cobra.Command, args []string) {
		fmt.Print(cur(os.Stdin))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".svs" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".svs")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
