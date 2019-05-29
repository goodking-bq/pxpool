package cmd

import (
	"os"
	"pxpool/models"
	"pxpool/scanner"
	"pxpool/storage"
	"runtime"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	maxConcurrency int64
	cidrFile       string
	cidrOne        string
	ports          string
)

func ExecuteScan(cmd *cobra.Command, args []string) {

	scanner := scanner.NewScanner(maxConcurrency, logger)
	if cidrFile != "" {
		file, err := os.Open(cidrFile)
		if err != nil {
			logger.Errorf("文件 %s 不存在", cidrFile)
			return
		}
		defer file.Close()
		go scanner.FromFile(file, models.ProxyChan)
	} else if cidrOne != "" {
		go scanner.FromCidr([]byte(cidrOne), models.ProxyChan)
	} else {
		logger.Errorln("未给扫描目标")
		return
	}
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		var lastTick int64
		for i := int64(0); i < maxConcurrency; i++ {
			scanner.DoChan <- true
		}
		for {
			//for scanner.Doing.Get() < maxConcurrency {
			select {
			case <-ticker.C:
				logger.Infof("ip数：%d,正在执行: %d,已完成: %d\n", scanner.IPcount.Get(), scanner.Doing.Get(), scanner.ScanCount.Get())
				if scanner.ScanCount.Get() != lastTick {
					lastTick = scanner.ScanCount.Get()
				} else {
					logger.Fatalln("扫描停止。")
				}
			case address := <-scanner.Chan:
				<-scanner.DoChan
				scanner.Wg.Add(1)
				go scanner.ScanIP(address, models.ProxyChan)
				scanner.Doing.Inc(1)

			case <-gCtx.Done():
				close(models.ProxyChan)
				break
			}
			//}
		}
	}()
	storage.StartStorage(gCtx, &storager, models.ProxyChan)
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
	scanCmd.PreRunE = preRunE
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&cidrFile, "file", "f", "", "从文件加载目标，一行一个。")
	scanCmd.Flags().StringVar(&ports, "ports", "80", "需要扫描的端口")
	scanCmd.Flags().StringVar(&cidrOne, "cidr", "", "给定cidr作为扫描目标")
	scanCmd.Flags().Int64VarP(&maxConcurrency, "maxconcurrency", "C", 0, "扫描并发数量")
}

func preRunE(cmd *cobra.Command, args []string) error {
	if cidrFile == "" {
		cidrFile = viper.GetStringMapString("scanner")["file"]
	}
	if cidrOne == "" {
		cidrOne = viper.GetStringMapString("scanner")["cidr"]
	}
	if ports == "" {
		ports = viper.GetStringMapString("scanner")["ports"]
	}
	if maxConcurrency == 0 {
		u, err := strconv.ParseUint(viper.GetStringMapString("scanner")["maxconcurrency"], 10, 64)
		if err != nil {
			u = 100
		}
		maxConcurrency = int64(u)
		if maxConcurrency == 0 {
			maxConcurrency = int64(runtime.NumGoroutine())
		}
	}
	return nil
}
