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
	"github.com/defektive/gfy/scanner"
	"github.com/spf13/cobra"
	"os"
)

var Source string
var Destination string

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import photos from SOURCE to DESTINATION",
	Long: `Photos will be moved to a path based on date taken,
	eg: 2015/05/13/GFY_[hash].jpg

	Photos not moved are **potential** duplicates.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sorting photos in " + Source)
		fmt.Println("Moving photos to " + Destination)

		files := scanner.ScanDir(Source)
		for _, f := range files {
			os.MkdirAll(f.SortedPath(Destination), 0777)
			_, err := os.Stat(f.SortedFullPath(Destination))
			if os.IsNotExist(err) {
				fmt.Printf("Moving %s >> %s\n", f.Path, f.SortedFullPath(Destination))
				os.Rename(f.Path, f.SortedFullPath(Destination))
			} else if err == nil {
				fmt.Printf("Duplicate %s >> %s\n", f.Path, f.SortedFullPath(Destination))
			} else {
				fmt.Printf("Error[%s] %s >> %s\n", err, f.Path, f.SortedFullPath(Destination))
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	importCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read photos from")
	importCmd.Flags().StringVarP(&Destination, "destination", "d", "", "Destination directory to move photos to")

}
