// Copyright © 2016 Dylan Clendenin <dylan@betterdoctor.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "duncan",
	Short: "Duncan is a Docker deployment tool",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	checkExecutableVersion()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.duncan.yml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".duncan") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("could not find config file: %s\n With error %v", viper.ConfigFileUsed(), err)
		os.Exit(1)
	}
}

func checkExecutableVersion() {
	if Version == "" {
		fmt.Println("Version is only set at compile time, skipping self update")
		return
	}

	latest, found, err := selfupdate.DetectLatest("betterdoctor/duncan")
	if err != nil {
		fmt.Printf("Error occurred while detecting version: %s\n", err)
		return
	}

	v := semver.MustParse(Version[1:])
	if !found || latest.Version.LTE(v) {
		fmt.Printf("Current version is the latest: %s\n", v)
		return
	}

	fmt.Println("Your duncan version is out of date!")
	fmt.Printf("Do you want to update to %s (yes/no): ", latest.Version)
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if strings.TrimSpace(input) != "yes" {
		fmt.Println("Skipping update")
		return
	}

	exe, err := os.Executable()
	if err != nil {
		fmt.Println("Could not locate executable path")
		return
	}
	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		fmt.Printf("Error occurred while updating binary: %s\n", err)
		return
	}
	fmt.Printf("Successfully updated to version: %s\n", latest.Version)
}
