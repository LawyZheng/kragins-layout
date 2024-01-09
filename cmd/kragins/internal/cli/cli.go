package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/kardianos/service"
	"github.com/lawyzheng/daemon"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lawyzheng/kragins/cmd/kragins/internal/app/conf"
	"github.com/lawyzheng/kragins/pkg/buildinfo"
	"github.com/lawyzheng/kragins/pkg/process"
)

var (
	svcConfig = &service.Config{
		Name:             "KraginsDemo",              //服务显示名称
		DisplayName:      "Krangins Demo",            //服务名称
		Description:      "Demo Service For Kragins", //服务描述
		WorkingDirectory: process.ExecutablePath(),
		Arguments:        []string{"serve"}, //启动参数
	}

	cfgFile   = ""
	cfg       *conf.Config
	rootCmd   *cobra.Command
	daemonCtl *daemon.Controller

	afterRun func()
)

func Execute() error {
	defer func() {
		if afterRun != nil {
			afterRun()
		}
	}()

	return rootCmd.Execute()
}

func readConfig() (err error) {
	svcConfig.Arguments = append(svcConfig.Arguments, "--config", cfgFile)

	if !filepath.IsAbs(cfgFile) {
		cfgFile = filepath.Join(process.ExecutablePath(), cfgFile)
	}

	cfg, err = conf.Read(cfgFile)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %s", err)
	}
	return nil
}

func newDaemonCtl() error {
	svc, cleanup, err := wireApp(cfg)
	defer func() {
		afterRun = cleanup
	}()
	if err != nil {
		return fmt.Errorf("初始化应用失败：%s", err)
	}

	daemonCtl, err = daemon.NewController(svcConfig, svc)
	if err != nil {
		return fmt.Errorf("连接服务失败: %s", err)
	}
	return nil
}

func init() {
	// CLI for Serve
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Run demo on shell",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := readConfig(); err != nil {
				return err
			}
			if err := newDaemonCtl(); err != nil {
				return err
			}
			if err := daemonCtl.Run(); err != nil {
				return fmt.Errorf("运行服务失败: %s", err)
			}
			return nil
		},
	}

	// CLI for Service
	svcCmd := &cobra.Command{
		Use:   "service",
		Short: "Run demo as a daemon service",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := readConfig(); err != nil {
				return err
			}
			if err := newDaemonCtl(); err != nil {
				return err
			}
			return nil
		},
	}

	// CLI for Service start
	svcCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Run demo as service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := daemonCtl.Start(); err != nil {
				return fmt.Errorf("启动服务失败: %s", err)
			}
			fmt.Println("启动服务成功!")
			return nil
		},
	})

	// CLI for Service stop
	svcCmd.AddCommand(&cobra.Command{
		Use:   "stop",
		Short: "Stop demo service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := daemonCtl.Stop(); err != nil {
				return fmt.Errorf("停止服务失败: %s", err)
			}
			fmt.Println("停止服务成功!")
			return nil
		},
	})

	// CLI for Service restart
	svcCmd.AddCommand(&cobra.Command{
		Use:   "restart",
		Short: "Restart demo service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := daemonCtl.Restart(); err != nil {
				return fmt.Errorf("重启服务失败: %s", err)
			}
			fmt.Println("重启服务成功!")
			return nil
		},
	})

	// CLI for Service uninstall
	svcCmd.AddCommand(&cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall demo service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := daemonCtl.Uninstall(); err != nil {
				return fmt.Errorf("卸载服务失败: %s", err)
			}
			fmt.Println("卸载服务成功!")
			return nil
		},
	})

	// CLI for Service install
	svcCmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install demo as service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := daemonCtl.Install(); err != nil {
				return fmt.Errorf("安装服务失败: %s", err)
			}
			fmt.Println("安装服务成功!")
			return nil
		},
	})

	// CLI for Database
	dbCmd := &cobra.Command{
		Use:   "database",
		Short: "Interact with demo database",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := readConfig(); err != nil {
				return err
			}
			return nil
		},
	}

	// CLI for Database Migrate
	dbCmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Migrate in database struct",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("暂无数据库持久化要求")
		},
	})

	verCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			info := map[string]interface{}{
				"Version":   buildinfo.Version(),
				"Build":     buildinfo.Build(),
				"BuildTime": buildinfo.BuildTime(),
			}
			data, _ := json.MarshalIndent(info, "", "    ")
			fmt.Println(string(data))
		},
	}

	rootCmd = &cobra.Command{
		Use:   "kragins",
		Short: "kraings is the demo for kragins",
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file (default is <executable path>/config.yaml)")
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(svcCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(verCmd)
}
