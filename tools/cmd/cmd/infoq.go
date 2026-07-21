package cmd

import (
	"fmt"
	"news4coder/internal/official"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var infoqDemoMode bool

var infoqCmd = &cobra.Command{
	Use:   "infoq",
	Short: "🎯 专注模式 - 获取 InfoQ 中文站热点内容",
	Long: `专注模式：直接从 InfoQ 中文站热点清单获取最新技术资讯。

这是官方信息源，使用专用抓取器直接获取原站热点内容，
无需搜索引擎中转，内容质量更高、更新更及时。`,
	Example: `  # 获取 InfoQ 热点内容
  news4coder infoq
  
  # 演示模式
  news4coder infoq --demo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取 InfoQ 官方源配置
		registry := official.GetRegistry()
		source, exists := registry.Get("infoq")
		if !exists {
			return fmt.Errorf("InfoQ 官方源未配置")
		}

		cyan := color.New(color.FgCyan).SprintFunc()
		magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()

		fmt.Printf("%s %s 专注模式 - 正在获取 %s 的热点内容...\n", magenta("🎯"), cyan("⟳"), source.Name)
		fmt.Println()

		if infoqDemoMode {
			// 演示模式
			results := generateDemoResults(source.Name, source.URL)
			displayOfficialResults(results, source.Name, source.URL)
			return nil
		}

		// 创建专用抓取器
		factory := official.NewFetcherFactory()
		fetcher, err := factory.Create(source)
		if err != nil {
			return fmt.Errorf("创建抓取器失败: %w", err)
		}

		// 执行抓取
		results, err := fetcher.Fetch()
		if err != nil {
			return fmt.Errorf("获取内容失败: %w", err)
		}

		// 显示结果
		displayOfficialResults(results, source.Name, source.URL)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoqCmd)
	infoqCmd.Flags().BoolVarP(&infoqDemoMode, "demo", "d", false, "演示模式（使用模拟数据）")
}
