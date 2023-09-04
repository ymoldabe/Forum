package configs

import (
	"encoding/json"
	"os"
)

// Config представляет структуру конфигурации приложения.
type Config struct {
	Port    string
	Migrate string
	DB      struct {
		Dsn    string
		Driver string
	}
}

// NewConfig создает и возвращает новую конфигурацию приложения, загружая ее из JSON-файла.
func NewConfig() (Config, error) {
	// Открываем JSON-файл с конфигурацией.
	configFile, err := os.Open("./configs/config.json")
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()

	// Декодируем JSON-файл в структуру Config.
	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
