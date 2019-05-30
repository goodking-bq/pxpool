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
	"pxpool/crawler"
	"pxpool/models"
	"pxpool/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "启动网络爬虫",
	Long:  `运行爬虫程序.`,
	Run: func(cmd *cobra.Command, args []string) {
		cManager := crawler.NewCrawl(logger, config, models.ProxyChan)
		cManager.Start()
		if config.Url == "" {
			storage.StartStorage(gCtx, &storager, models.ProxyChan)
		} else {
			for {
				select {
				case proxy := <-models.ProxyChan:
					go PostProxy(proxy, config.Url)
				case <-gCtx.Done():
					close(models.ProxyChan)
				}
			}
		}
		cManager.ExitSignal <- true
		defer gCancal()
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)
	crawlCmd.Flags().IntVarP(&config.Crawl.Ticker, "ticker", "t", 0, "扫描间隔时间")
	viper.BindPFlag("crawl.ticker", crawlCmd.Flags().Lookup("ticker"))
}
