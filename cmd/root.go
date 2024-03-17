package cmd

import (
	"github.com/spf13/cobra"
	"github.com/starjun/toes-cli/config"
)

var rootCmd = &cobra.Command{
	Use:   "toes",
	Short: "提效 toes-cli",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.CfgPath, "config",
		"f", "config.toes.yaml", "设置配置文件路径")
	rootCmd.PersistentFlags().StringVarP(&config.Dsn, "dsn",
		"", "", "设置链接数据库字符串，默认 无，PS:当未配置数据库连接，生成的模板基于一个 user demo 对象")
	rootCmd.PersistentFlags().StringVarP(&config.Style, "style",
		"s", "toes", "设置生成代码风格")
	rootCmd.PersistentFlags().StringVarP(&config.Table, "table",
		"t", "*", "设置需要生成的表名，PS：默认表示当前数据库下所有表，否则为指定的表名")
	rootCmd.PersistentFlags().StringVarP(&config.Path, "path",
		"", "", "设置项目Path")

	rootCmd.PersistentFlags().StringVarP(&config.RootPackage, "root_package",
		"", "toesRoot", "模块路径")
	rootCmd.PersistentFlags().StringVarP(&config.ToesRoot, "toes_root",
		"", "./", "toes项目根目录")
	rootCmd.PersistentFlags().BoolVarP(&config.CreateQueryConfig, "create_query_config",
		"", false, "是否创建综合查询配置文件，true：创建，false：不创建")

	config.LoadConfig(config.CfgPath)
}
