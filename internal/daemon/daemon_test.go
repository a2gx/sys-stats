package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/a2gx/sys-stats/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProcessManager реализация ProcessManager для тестирования
type MockProcessManager struct {
	mock.Mock
}

func (m *MockProcessManager) StartProcess(cmd *exec.Cmd) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func (m *MockProcessManager) FindProcess(pid int) (Process, error) {
	args := m.Called(pid)
	return args.Get(0).(Process), args.Error(1)
}

// MockProcess реализация Process для тестирования
type MockProcess struct {
	mock.Mock
}

func (m *MockProcess) Signal(sig os.Signal) error {
	args := m.Called(sig)
	return args.Error(0)
}

func TestNewDaemonManager(t *testing.T) {
	cfg := &config.Config{}
	customPidPath := "/var/run/custom.pid"
	customLogPath := "/var/log/custom.log"
	mockPM := new(MockProcessManager)

	t.Run("DefaultValues", func(t *testing.T) {
		dm := NewDaemonManager(cfg)

		assert.NotNil(t, dm, "DaemonManager должен быть создан")
		assert.Equal(t, cfg, dm.cfg, "Конфигурация должна быть сохранена")
		assert.Equal(t, "/tmp/daemon.pid", dm.pidFilePath, "PID-файл должен иметь значение по умолчанию")
		assert.Equal(t, "/tmp/daemon.log", dm.logFilePath, "Log-файл должен иметь значение по умолчанию")
		_, ok := dm.processManager.(*SystemProcessManager)
		assert.True(t, ok, "ProcessManager должен быть SystemProcessManager по умолчанию")
	})

	t.Run("WithOptions", func(t *testing.T) {
		dm := NewDaemonManager(
			cfg,
			WithPidFilePath(customPidPath),
			WithLogFilePath(customLogPath),
			WithProcessManager(mockPM),
		)

		assert.NotNil(t, dm, "DaemonManager должен быть создан")
		assert.Equal(t, cfg, dm.cfg, "Конфигурация должна быть сохранена")
		assert.Equal(t, customPidPath, dm.pidFilePath, "PID-файл должен быть установлен через опцию")
		assert.Equal(t, customLogPath, dm.logFilePath, "Log-файл должен быть установлен через опцию")
		assert.Same(t, mockPM, dm.processManager, "ProcessManager должен быть установлен через опцию")
	})

	t.Run("WithPidFilePath", func(t *testing.T) {
		dm := &DaemonManager{pidFilePath: "/tmp/daemon.pid"}
		customPath := "/var/run/custom.pid"

		option := WithPidFilePath(customPath)
		option(dm)

		assert.Equal(t, customPath, dm.pidFilePath, "Опция WithPidFilePath должна изменить путь к PID-файлу")
	})

	t.Run("WithLogFilePath", func(t *testing.T) {
		dm := &DaemonManager{logFilePath: "/tmp/daemon.log"}
		customPath := "/var/log/custom.log"

		option := WithLogFilePath(customPath)
		option(dm)

		assert.Equal(t, customPath, dm.logFilePath, "Опция WithLogFilePath должна изменить путь к лог-файлу")
	})

	t.Run("WithProcessManager", func(t *testing.T) {
		dm := &DaemonManager{processManager: &SystemProcessManager{}}
		mockPM := new(MockProcessManager)

		option := WithProcessManager(mockPM)
		option(dm)

		assert.Same(t, mockPM, dm.processManager, "Опция WithProcessManager должна заменить ProcessManager")
	})
}

func TestStartDaemon(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "daemon-test-")
	if err != nil {
		t.Fatalf("Не удалось создать временную директорию: %v", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	pidPath := filepath.Join(tmpDir, "test.pid")
	logPath := filepath.Join(tmpDir, "test.log")

	// Базовые флаги для тестов
	baseFlags := StartDaemonFlags{
		Background: true,
		ConfigFile: "config.yaml",
		Host:       "localhost",
		Port:       8080,
	}

	t.Run("BackgroundMode", func(t *testing.T) {
		mockPM := new(MockProcessManager)

		// Настраиваем ожидания для мока
		mockCmd := &exec.Cmd{}
		mockCmd.Process = &os.Process{Pid: 12345}
		mockPM.On("StartProcess", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			// Копируем PID в переданную команду, чтобы имитировать успешный запуск
			cmd := args.Get(0).(*exec.Cmd)
			cmd.Process = mockCmd.Process
		})

		// Создаем менеджер демона с временными путями к файлам
		dm := NewDaemonManager(
			&config.Config{},
			WithPidFilePath(pidPath),
			WithLogFilePath(logPath),
			WithProcessManager(mockPM),
		)

		// Проверяем запуск демона в фоновом режиме
		err := dm.StartDaemon(baseFlags)

		// Проверяем результат
		assert.NoError(t, err, "Запуск демона должен быть успешным")
		mockPM.AssertExpectations(t)

		// Проверяем, что PID файл был создан
		pidBytes, err := os.ReadFile(pidPath)
		assert.NoError(t, err, "PID-файл должен существовать")
		assert.Equal(t, "12345", string(pidBytes), "PID должен быть записан в файл")

		// Проверяем, что лог-файл был создан
		_, err = os.Stat(logPath)
		assert.NoError(t, err, "Лог-файл должен быть создан")
	})

	t.Run("DaemonAlreadyStarted", func(t *testing.T) {
		// Создаем PID файл для имитации уже запущенного демона
		err := os.WriteFile(pidPath, []byte("12345"), 0644)
		assert.NoError(t, err, "Должны создать тестовый PID-файл")

		// Создаем менеджер демона
		dm := NewDaemonManager(
			&config.Config{},
			WithPidFilePath(pidPath),
			WithLogFilePath(logPath),
		)

		// Пытаемся запустить демон
		err = dm.StartDaemon(baseFlags)

		// Должна быть ошибка, т.к. демон "уже запущен"
		assert.Error(t, err, "Запуск должен вернуть ошибку, если демон уже запущен")
		assert.Contains(t, err.Error(), "already running", "Ошибка должна указывать, что демон уже запущен")
	})

	t.Run("ErrorStartedDaemon", func(t *testing.T) {
		// Настраиваем мок с ошибкой запуска
		mockPM := new(MockProcessManager)
		mockPM.On("StartProcess", mock.Anything).Return(fmt.Errorf("ошибка запуска"))

		// Создаем менеджер демона
		dm := NewDaemonManager(
			&config.Config{},
			WithPidFilePath(pidPath),
			WithLogFilePath(logPath),
			WithProcessManager(mockPM),
		)

		// Пытаемся запустить демон
		err := dm.StartDaemon(baseFlags)

		// Проверяем результат
		assert.Error(t, err, "Должна быть ошибка при неудачном запуске процесса")
		assert.Contains(t, err.Error(), "daemon already running", "Сообщение должно содержать информацию об ошибке запуска")
	})
}
