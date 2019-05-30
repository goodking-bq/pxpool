package cmd

import (
	"pxpool/web"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "启动web api等",
	Long:  `启动网站，api等.`,
	Run: func(cmd *cobra.Command, args []string) {
		api := web.NewApp(config, storager, logger)
		// api.SetBind(bind)
		// api.SetPort(port)
		api.Run()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVarP(&config.Web.Bind, "bind", "b", "", "侦听的ip地址")
	viper.BindPFlag("web.bind", webCmd.Flags().Lookup("bind"))
	webCmd.Flags().IntVarP(&config.Web.Port, "port", "p", 3000, "侦听的端口")
	viper.BindPFlag("web.port", webCmd.Flags().Lookup("port"))
}
