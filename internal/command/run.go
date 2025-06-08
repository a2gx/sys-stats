package command

import (
	"fmt"
	"github.com/a2gx/sys-stats/internal/config"
	"github.com/a2gx/sys-stats/internal/daemon"
	"github.com/spf13/cobra"
)

func NewRun() *cobra.Command {
	var background bool
	var configFile string
	var host string
	var port int

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the daemon to collect system statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Загружаем конфигурацию
			cfg, err := config.NewConfig(configFile)
			if err != nil {
				return err
			}

			// Устанавливаем параметры для gRPC сервера
			cfg.AddrGRPC = fmt.Sprintf("%s:%d", host, port)

			// Создаем менеджер демона с параметрами
			dm := daemon.NewDaemonManager(cfg)

			// Запускаем демон
			// + передаем параметры для перезапуска в фоновом режиме
			return dm.StartDaemon(daemon.StartDaemonFlags{
				Background: background,
				ConfigFile: configFile,
				Host:       host,
				Port:       port,
			})
		},
	}

	cmd.Flags().BoolVarP(&background, "background", "b", false, "Run the daemon in background mode")
	cmd.Flags().StringVar(&host, "host", "localhost", "Host fot gRPC server")
	cmd.Flags().IntVar(&port, "port", 50055, "Port fot gRPC server")
	cmd.Flags().StringVar(&configFile, "config", "./configs/config.yaml", "Path to the configuration file")

	return cmd
}
