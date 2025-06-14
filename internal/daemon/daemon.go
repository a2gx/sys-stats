package daemon

import (
	"github.com/a2gx/sys-stats/internal/config"
	"os"
	"os/exec"
)

// DaemonManager управляет жизненным циклом демона
type DaemonManager struct {
	cfg            *config.Config
	pidFilePath    string
	logFilePath    string
	processManager ProcessManager
}

// ProcessManager операции с процессами для возможности тестирования
type ProcessManager interface {
	StartProcess(cmd *exec.Cmd) error
	FindProcess(pid int) (Process, error)
}

// Process абстракция над системным процессом
type Process interface {
	Signal(sig os.Signal) error
}

// SystemProcessManager реальная имплементация ProcessManager
type SystemProcessManager struct{}

func (spm *SystemProcessManager) StartProcess(cmd *exec.Cmd) error {
	return cmd.Start()
}

func (spm *SystemProcessManager) FindProcess(pid int) (Process, error) {
	return os.FindProcess(pid)
}

// NewDaemonManager создает новый менеджер демона
func NewDaemonManager(cfg *config.Config, options ...Option) *DaemonManager {
	dm := &DaemonManager{
		cfg:            cfg,
		pidFilePath:    "/tmp/daemon.pid",
		logFilePath:    "/tmp/daemon.log",
		processManager: &SystemProcessManager{},
	}

	// Применяем опции конфигурации
	for _, option := range options {
		option(dm)
	}

	return dm
}

// Option представляет функциональную опцию для конфигурации DaemonManager
type Option func(*DaemonManager)

// WithPidFilePath устанавливает путь к PID-файлу
func WithPidFilePath(path string) Option {
	return func(dm *DaemonManager) {
		dm.pidFilePath = path
	}
}

// WithLogFilePath устанавливает путь к лог-файлу
func WithLogFilePath(path string) Option {
	return func(dm *DaemonManager) {
		dm.logFilePath = path
	}
}

// WithProcessManager устанавливает менеджер процессов (для тестирования)
func WithProcessManager(pm ProcessManager) Option {
	return func(dm *DaemonManager) {
		dm.processManager = pm
	}
}
