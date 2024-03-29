// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/sealstore/sealbuild/pkg"
	"github.com/sealstore/sealbuild/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string
var tmpFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sealbuild",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		config := &utils.Config{}
		utils.LoadIni(config, cfgFile)
		utils.VarsConfig = config
		pkg.Build(config.AppEnable, tmpFile)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&cfgFile, "conf", "f", "config.ini", "sealbuild config.ini file location")
	rootCmd.Flags().StringVarP(&tmpFile, "template", "t", "", "config template file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
