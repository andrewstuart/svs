package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/blang/semver/v4"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func sortStrings(strs []string) []string {
	type Version struct {
		semver semver.Version
		orig   string
	}

	var versions []Version

	for _, str := range strs {
		v, err := semver.ParseTolerant(str)
		if err != nil {
			// Ignore if it is not a semver.
			continue
		}

		versions = append(
			versions,
			Version{
				semver: v,
				orig:   str,
			},
		)
	}

	// Sort versions.
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].semver.LT(versions[j].semver)
	})

	// Prepare result slice with original strings sorted by semver.
	var res []string
	for _, v := range versions {
		res = append(res, v.orig)
	}

	return res
}

var cfgFile string

func cur(rc io.Reader) (string, error) {
	var vs []string
	br := bufio.NewReader(os.Stdin)

	var err error
	var bs []byte
	for err == nil {
		bs, _, err = br.ReadLine()
		for _, f := range strings.Fields(string(bs)) {
			vs = append(vs, f)
		}
	}
	if err == io.EOF {
		err = nil
	}

	strs := sortStrings(vs)
	return strs[len(strs)-1], err
}

func fromGit() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", errors.Wrap(err, "git open failed")
	}

	t, err := repo.Tags()
	if err != nil {
		return "", errors.Wrap(err, "git tags failed")
	}
	var tags []string
	t.ForEach(func(r *plumbing.Reference) error {
		tags = append(tags, path.Base(r.Name().String()))
		return nil
	})
	sorted := sortStrings(tags)
	return sorted[len(sorted)-1], nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "svs",
	Short: "semver sort",
	Run: func(c *cobra.Command, args []string) {
		if len(args) > 0 && args[0] == "-" {
			v, err := cur(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(v)
			return
		}
		t, err := fromGit()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(t)
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
