// Пакет config представляет инструменты для работы с конфигурациями сервера и агента
package config

import (
	"net"
	"os"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

// NetAddress строка вида [<хост>]:<порт> и реализует интерфейс pflag.Value
type NetAddress string

func (a NetAddress) String() string {
	return string(a)
}

func (a *NetAddress) Set(s string) error {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return err
	}
	if len(host) == 0 {
		// если не указан хост, то используем localhost по умолчанию
		s = "localhost:" + port
	}
	*a = NetAddress(s)
	return nil
}

func (a *NetAddress) Type() string {
	return "string"
}

// Config конфигурация агента
type Config struct {
	// адрес и порт сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `-a`
	RunAddress NetAddress `env:"RUN_ADDRESS"`
	// адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `-d`
	DarabaseURI string `env:"DATABASE_URI"`
	// адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	// уровень логирования: переменная окружения ОС `LOG_LEVEL` или флаг `-l`
	LogLevel string `env:"LOG_LEVEL"`
}

func parseFromCmd(c *Config) error {
	cmd := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	runAddress := NetAddress("localhost:8080")
	cmd.VarP(&runAddress, "run-address", "a", "адрес и порт запуска сервиса")
	accrualSystemAddress := cmd.StringP("accrual-address", "r", "http://accrual", "адрес системы расчёта начислений")
	databaseURI := cmd.StringP("database-uri", "d", "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable", "адрес подключения к базе данных")
	logLevel := cmd.StringP("log-level", "l", "info", "уровень логирования")

	if err := cmd.Parse(os.Args[1:]); err != nil {
		return err
	}

	*c = Config{
		RunAddress:           runAddress,
		AccrualSystemAddress: *accrualSystemAddress,
		DarabaseURI:          *databaseURI,
		LogLevel:             *logLevel,
	}
	return nil
}

func parseFromEnv(c *Config) error {
	if err := env.Parse(c); err != nil {
		return err
	}
	if len(c.RunAddress) != 0 {
		if err := c.RunAddress.Set(c.RunAddress.String()); err != nil {
			return err
		}
	}
	return nil
}

// Parse разбирает настройки из аргументов командной строки
// и переменных окружения. Переменные окружения имеют более высокий
// приоритет, чем аргументы.
func Parse(c *Config) error {
	if err := parseFromCmd(c); err != nil {
		return err
	}
	if err := parseFromEnv(c); err != nil {
		return err
	}
	return nil
}
