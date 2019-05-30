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
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ticker int
)

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "启动网络爬虫",
	Long:  `运行爬虫程序.`,
	Run: func(cmd *cobra.Command, args []string) {
		cManager := crawler.NewDefaultCrawl(models.ProxyChan)
		if ticker > 0 {
			cManager.StartAndTicker(ticker)
		}
		if postURL == "" {
			storage.StartStorage(gCtx, &storager, models.ProxyChan)
		} else {
			for {
				select {
				case proxy := <-models.ProxyChan:
					go PostProxy(proxy, postURL)
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
	crawlCmd.PreRunE = crawlPreRunE
	scanCmd.Flags().StringVarP(&postURL, "post", "p", "", "结果提交地址")
	scanCmd.Flags().IntVarP(&ticker, "ticker", "t", 0, "扫描间隔时间")
	initStorager()
}

func crawlPreRunE(cmd *cobra.Command, args []string) error {
	if postURL == "" {
		postURL = viper.GetStringMapString("crawl")["post"]
	}
	if ticker == 0 {
		ts := viper.GetStringMapString("crawl")["ticker"]
		_t, err := strconv.ParseInt(ts, 10, 0)
		if err != nil {
			ticker = 0
		} else {
			ticker = int(_t)
		}
	}
	return nil
}
