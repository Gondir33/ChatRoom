package run

import (
	"context"
	"fmt"
	"goTest/config"
	"goTest/internal/infrastructure/component"
	"goTest/internal/infrastructure/router"
	"goTest/internal/infrastructure/server"
	"goTest/internal/modules"
	"goTest/internal/storages"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Application - интерфейс приложения
type Application interface {
	Runner
	Bootstraper
}

// Runner - интерфейс запуска приложения
type Runner interface {
	Run() int
}

// Bootstraper - интерфейс инициализации приложения
type Bootstraper interface {
	Bootstrap(options ...interface{}) Runner
}

type App struct {
	conf   config.AppConf
	logger *zap.Logger

	srv      server.Server
	Sig      chan os.Signal
	Storages *storages.Storages
	Servises *modules.Services
}

func NewApp(conf config.AppConf) *App {
	return &App{conf: conf, Sig: make(chan os.Signal, 1)}
}

func (a *App) Run() int {
	// на русском
	// создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	errGroup, ctx := errgroup.WithContext(ctx)

	// запускаем горутину для graceful shutdown
	// при получении сигнала SIGINT
	// вызываем cancel для контекста
	errGroup.Go(func() error {
		sigInt := <-a.Sig
		a.logger.Info("signal interrupt recieved", zap.Stringer("os_signal", sigInt))
		cancel()
		return nil
	})

	errGroup.Go(func() error {
		err := a.srv.Serve(ctx)
		if err != nil && err != http.ErrServerClosed {
			a.logger.Error("app: server error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return 1
	}
	return 0
}

func (a *App) Bootstrap(options ...interface{}) Runner {
	// инициализация логгера
	logger, _ := zap.NewProduction()
	a.logger = logger
	components := component.NewComponents(a.conf, a.logger)

	// инициализация базы данных sql и его адаптера
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		a.conf.DB.User, a.conf.DB.Password, a.conf.DB.Host, a.conf.DB.Port, a.conf.DB.Name)
	pool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		a.logger.Fatal("error init db", zap.Error(err))
	}

	// инициализация хранилищ
	newStorages := storages.NewStorages(pool)

	a.Storages = newStorages
	// инициализация сервисов
	services := modules.NewServices(newStorages, components)
	a.Servises = services
	controllers := modules.NewControllers(services, components)
	// инициализация роутера
	r := router.NewRouter(controllers, components)
	// конфигурация сервера
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.conf.Server.Port),
		Handler: r,
	}

	a.srv = server.NewHttpServer(a.conf.Server, srv, a.logger)
	return a
}
