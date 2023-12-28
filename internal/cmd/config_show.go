/*
Copyright © 2023 Curtis Jewell <bulbistry@curtisjewell.name>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configShowCmd represents the configShow command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows the current configuration",
	Long:  `Shows the current configuration`,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return initConfig(true)
	},
	Run: func(cmd *cobra.Command, args []string) {
		settings := viper.AllKeys()

		settingsSimple := make(map[string]string, len(settings))

		for _, key := range settings {
			envName := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
			settingsSimple[envName] = fmt.Sprint(viper.Get(key))
		}

		keys := make(sort.StringSlice, 0, len(settingsSimple))
		for k := range settingsSimple {
			keys = append(keys, k)
		}
		keys.Sort()

		for _, k := range keys {
			fmt.Println(k, "value is:", settingsSimple[k])
		}
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
}
