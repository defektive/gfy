// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/defektive/gfy/scanner"

)

var cfgFile string
var Source string
var Destination string
// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gfy",
	Short: "Go file yourself",
	Long: `A simple program to scan your photos, remove dupes, and save them
	in a pretty directory structure`,
// Uncomment the following line if your bare application
// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sorting photos in "+ Source)
		fmt.Println("Moving photos to "+ Destination)

		files := scanner.ScanDir(Source, Destination)
		for _, f := range files {
			fmt.Println(f.SortedFullPath())
			os.MkdirAll(f.SortedPath(), 0777)
			_, err := os.Stat(f.SortedFullPath())
			if os.IsNotExist(err) {
				fmt.Println("Moving")
			  os.Rename(f.Path, f.SortedFullPath())
			} else {
				fmt.Println(err)
			}

		}
	},
}


// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gfy.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read photos from")
	RootCmd.Flags().StringVarP(&Destination, "destination", "d", "", "Destination directory to move photos to")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".gfy") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
