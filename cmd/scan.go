package cmd

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"pxpool/models"
	"pxpool/scanner"
	"pxpool/storage"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PostProxy(proxy *models.Proxy, urlStr string) {
	resp, err := http.PostForm(urlStr,
		url.Values{"host": {proxy.Host}, "port": {proxy.Host}, "category": {proxy.Category}})

	if err != nil {
		// handle error
		logger.WithFields(logrus.Fields{"proxy": proxy.URL()}).Errorf("代理提交失败: %s", err)
		return
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		logger.WithFields(logrus.Fields{"proxy": proxy.URL()}).Errorf("代理提交返回失败: %s", err)
		return
	}

	logger.WithFields(logrus.Fields{"proxy": proxy.URL()}).Errorln("代理提交成功！")
}

func ExecuteScan(cmd *cobra.Command, args []string) {
	scanner := scanner.NewScanner(config, logger)
	go scanner.MakeAddress()
	go scanner.Scan(gCtx)
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

	// exit := make(chan os.Signal)
	// signal.Notify(exit, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
	// signal := <-exit
	// cancal()
	// log.Println(signal)

}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "启动网络爬虫",
	Long:  `运行爬虫程序.`,
	Run:   ExecuteScan,
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&config.Scanner.File, "file", "f", "", "从文件加载目标，一行一个。")
	viper.BindPFlag("scanner.file", scanCmd.Flags().Lookup("file"))
	scanCmd.Flags().StringVar(&config.Scanner.PortString, "ports", "", "需要扫描的端口")
	viper.BindPFlag("scanner.portstring", scanCmd.Flags().Lookup("ports"))
	scanCmd.Flags().StringVar(&config.Scanner.Cidr, "cidr", "", "给定cidr作为扫描目标")
	viper.BindPFlag("scanner.cidr", scanCmd.Flags().Lookup("cidr"))
	scanCmd.Flags().IntVarP(&config.Scanner.MaxConcurrency, "maxconcurrency", "C", runtime.NumGoroutine(), "扫描并发数量")
	viper.BindPFlag("scanner.maxconcurrency", scanCmd.Flags().Lookup("maxconcurrency"))
}
