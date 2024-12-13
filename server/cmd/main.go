package main

import (
	"flag"
	"log"

	_ "github.com/lib/pq"
	"github.com/lzimin05/IDZ/internal/api"
	"github.com/lzimin05/IDZ/internal/config"
	"github.com/lzimin05/IDZ/internal/provider"
	"github.com/lzimin05/IDZ/internal/usecase"
)

func main() {
	// Считываем аргументы командной строки
	configPath := flag.String("config-path", "./configs/IDZ.yaml", "путь к файлу конфигурации")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	prv := provider.NewProvider(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBname)
	use := usecase.NewUsecase(cfg.Usecase.DefaultMessage, prv)
	srv := api.NewServer(cfg.IP, cfg.Port, cfg.API.MaxMessageSize, use)

	srv.Run()
}
