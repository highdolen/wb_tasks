package config

// Config - описывает конфигурацию приложения
type Config struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	MetricsPath string `yaml:"metrics_path"`
	HealthPath  string `yaml:"health_path"`
	GCPercent   int    `yaml:"gc_percent"`
}

// Address - возвращает адрес сервера в формате host:port
func (c Config) Address() string {
	return c.Host + ":" + c.Port
}
