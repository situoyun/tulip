/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"git.autops.xyz/autops/base/logs"
	"git.autops.xyz/autops/base/utils/pflagenv"

	"github.com/spf13/cobra"
)

var (
	cfg logs.LumberjackWrapperConfig

	SQLConn string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tulip",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tulip.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVar(&SQLConn, "sql_conn", "root:cmstop@tcp(mysql:3306)/sandy?charset=utf8&parseTime=True&loc=Local", "mysql conn")

	rootCmd.PersistentFlags().StringVar(&cfg.Path, "log", "/opt", "日志文件存储路径")
	rootCmd.PersistentFlags().IntVar(&cfg.MaxSize, "log_size", 50, "单个日志文件最大大小（MB）")
	rootCmd.PersistentFlags().IntVar(&cfg.MaxBackups, "log_backup", 10, "最多存储日志文件数目")
	rootCmd.PersistentFlags().IntVar(&cfg.MaxAge, "log_day", 7, "保留的历史日志文件天数")
	rootCmd.PersistentFlags().IntVar(&cfg.BufferSize, "log_buffer", 0, "日志写入缓冲区大小（byte）")
	rootCmd.PersistentFlags().BoolVar(&cfg.IsCompress, "log_compress", false, "日志文件是否压缩")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	if err := pflagenv.ParseSet("", rootCmd.PersistentFlags()); err != nil {
		fmt.Printf("pflagenv.ParseSet error:%v\n", err)
	}
}
