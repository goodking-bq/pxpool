package cmd

import (
	"pxpool/web"
	"strconv"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	bind   string
	port   int
	secret string
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "启动web api等",
	Long:  `启动网站，api等.`,
	Run: func(cmd *cobra.Command, args []string) {
		api := web.NewApp(storager, secret, logger)
		// api.SetBind(bind)
		// api.SetPort(port)
		api.Run()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.PreRunE = webPreRunE
	webCmd.Flags().StringVarP(&bind, "bind", "b", "", "侦听的ip地址")
	webCmd.Flags().IntVarP(&port, "port", "p", 0, "侦听的端口")
	webCmd.Flags().StringVarP(&secret, "secret", "s", "", "secret验证")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crawlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crawlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func webPreRunE(cmd *cobra.Command, args []string) error {
	if bind == "" {
		bind = viper.GetStringMapString("web")["bind"]
	}
	if port == 0 {
		portStr := viper.GetStringMapString("web")["port"]
		if portStr != "" {
			_port, err := strconv.ParseInt(portStr, 10, 16)
			if err != nil {
				return err
			}
			port = int(_port)
		} else {
			port = 3000
		}
	}
	if secret == "" {
		secret = viper.GetStringMapString("web")["secret"]
	}
	return nil
}
