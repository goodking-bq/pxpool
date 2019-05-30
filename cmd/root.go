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
	"context"
	"log"
	"os"
	"pxpool/models"
	"pxpool/storage"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	storager      storage.Storager
	logger        = logrus.New()
	gCtx, gCancal = context.WithCancel(context.Background())
	config        = models.DefaultConfig()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pxpool",
	Short: "pxpool 代理池全功能",
	Long: `pxpool 代理池全功能支持 网站爬虫 端口扫描 web api
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(args)
		cmd.Println(viper.GetStringMapString("bolt")["datapath"])

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Errorln(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(config.UnmarshalViper, initLog)
	cobra.OnInitialize(initStorager)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件 (默认是 config.yaml)")
	rootCmd.PersistentFlags().StringVar(&(config.StorageType), "storagetype", "bolt", "数据库存储类型 (默认是 bolt)")
	viper.BindPFlag("storagetype", rootCmd.PersistentFlags().Lookup("storagetype"))
	rootCmd.PersistentFlags().StringVar(&config.Log.File, "log", "", "日志文件路径 (默认是 os.stdout)")
	viper.BindPFlag("log.file", rootCmd.PersistentFlags().Lookup("log"))
	rootCmd.PersistentFlags().StringVar(&config.Log.Level, "loglevel", "", "日志级别 (默认是 info)")
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("loglevel"))
	rootCmd.PersistentFlags().StringVarP(&config.Url, "url", "u", "", "扫描或爬虫结果提交地址，默认保存到本地数据库。")
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	rootCmd.PersistentFlags().StringVarP(&config.Secret, "secret", "s", "", "secret验证")
	viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		_, err := os.Stat(cfgFile)
		if os.IsNotExist(err) {
			log.Fatalf("配置文件( %s )不存在", cfgFile)
		}
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			logger.Errorln(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pxpool")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Infoln("使用配置文件:", viper.ConfigFileUsed())
	}
	//cidrFile = viper.GetStringMapString("scanner")["file"]
}

func initStorager() {
	// if config.StorageType == "" {
	// 	storageType = viper.GetString("storagetype")
	// }
	switch strings.ToLower(config.StorageType) {
	case "bolt":
		bolt := storage.GetBoltStorage(viper.GetStringMapString("bolt")["datapath"])
		storager = bolt
		logger.Info("当前存储为 bolt")
		break
	}
}

func initLog() {
	switch strings.ToLower(config.Log.Level) {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
		break
	case "info":
		logger.SetLevel(logrus.InfoLevel)
		break
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
		break
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
		break
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
		break
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
		break
	}
	if config.Log.File != "" {
		pathMap := lfshook.PathMap{
			logrus.DebugLevel: config.Log.File,
			//logrus.InfoLevel:  logFile,
			//logrus.ErrorLevel: "./error.log",
		}
		logger.Hooks.Add(lfshook.NewHook(
			pathMap,
			&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"},
		))
	}
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	logger.Infoln("日志工作正常")
}
