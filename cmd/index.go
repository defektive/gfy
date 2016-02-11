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
	"encoding/json"
	"fmt"
	"github.com/defektive/gfy/scanner"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type SimplePhoto struct {
	Path      string
	ThumbPath string
	Date      time.Time
	Hash      string
}

var Path string

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Generate an index of your picture data",
	Long:  `Scans all your pics, saves the date, hash, and file path into an index`,
	Run: func(cmd *cobra.Command, args []string) {
		files := scanner.ScanDir(Path)
		db := scanner.OpenDb(Path)
		photos := []*SimplePhoto{}
		defer db.Close()
		for _, f := range files {
			fmt.Printf("%s\t%s\t%s\n", f.Date, f.Hash(), f.Path)
			f.Thumbnail(Path)
			// db.Add(f.Path, f.Date, f.Hash())
			simple := &SimplePhoto{
				Path:      f.Path,
				ThumbPath: f.Thumbnail(Path),
				Date:      f.Datetime(),
				Hash:      f.Hash(),
			}
			photos = append(photos, simple)
		}

		b, err := json.Marshal(photos)
		if err != nil {
			fmt.Println("error:", err)
		}
		os.Stdout.Write(b)
	},
}

func init() {
	RootCmd.AddCommand(indexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	indexCmd.Flags().StringVarP(&Path, "path", "p", "", "Path to your photos")
}
