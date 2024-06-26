package main

import (
	"goTest/config"
	"goTest/run"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("no .env files found")
	}
	// Создаем конфигурацию приложения
	conf := config.NewAppConf()
	// Создаем инстанс приложения
	App := run.NewApp(conf)

	exitCode := App.
		// Инициализация
		Bootstrap().
		// Запуск
		Run()

	os.Exit(exitCode)
}
