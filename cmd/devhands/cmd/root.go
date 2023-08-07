/*
Copyright © 2023 Alexey Samoylov <alexey.samoylov@gmail.com>

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
	"bytes"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/metalagman/devhands/internal/devhands/config"
	"github.com/metalagman/devhands/pkg/logger"
	"io/fs"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg = config.Config{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "devhands",
	Run: func(cmd *cobra.Command, args []string) {
		logger.CheckErr(cmd.Help())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initDotEnv)
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogger)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Set high log verbosity")
	rootCmd.PersistentFlags().BoolP("pretty", "p", false, "Set pretty log formatting (instead of json)")
}

func initDotEnv() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		logger.CheckErr(fmt.Errorf(".env load: %w", err))
	}
}

func initConfig() {
	viper.SetConfigType("toml")
	var defaultConfig = []byte(`
[server]
listen=":8080"
timeout_read="5s"
timeout_write="60s"
timeout_idle="1m"
[log]
verbose=1
pretty=1
`)
	logger.CheckErr(viper.ReadConfig(bytes.NewBuffer(defaultConfig)))

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	logger.CheckErr(viper.BindPFlag("log.verbose", rootCmd.PersistentFlags().Lookup("verbose")))
	logger.CheckErr(viper.BindPFlag("log.pretty", rootCmd.PersistentFlags().Lookup("pretty")))

	logger.CheckErr(viper.Unmarshal(&cfg))
}

func initLogger() {
	logger.NewGlobal(cfg.Logger)
}
