package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/jadevelopmentgrp/Ticket-Worker/i18n"
	archiverclient "github.com/jadevelopmentgrp/Tickets-Archiver-Client"
	app "github.com/jadevelopmentgrp/Tickets-Dashboard/app/http"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/app/http/endpoints/api/ticket/livechat"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/config"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/database"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/redis"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/rpc"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/rpc/cache"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/utils"
	"github.com/jadevelopmentgrp/Tickets-Utilities/chatrelay"
	"github.com/jadevelopmentgrp/Tickets-Utilities/model"
	"github.com/jadevelopmentgrp/Tickets-Utilities/observability"
	"github.com/jadevelopmentgrp/Tickets-Utilities/premium"
	"github.com/jadevelopmentgrp/Tickets-Utilities/secureproxy"
	"github.com/rxdn/gdl/rest/request"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	startPprof()

	cfg, err := config.LoadConfig()
	utils.Must(err)
	config.Conf = cfg

	var logger *zap.Logger
	if config.Conf.JsonLogs {
		loggerConfig := zap.NewProductionConfig()
		loggerConfig.Level.SetLevel(config.Conf.LogLevel)

		logger, err = loggerConfig.Build(
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
			zap.WrapCore(observability.ZapSentryAdapter(observability.EnvironmentProduction)),
		)
	} else {
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.Level.SetLevel(config.Conf.LogLevel)
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		logger, err = loggerConfig.Build(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	}

	if err != nil {
		panic(fmt.Errorf("failed to initialise zap logger: %w", err))
	}

	logger.Info("Connecting to database")
	database.ConnectToDatabase()

	logger.Info("Connecting to cache")
	cache.Instance = cache.NewCache()

	logger.Info("Initialising microservice clients")
	utils.ArchiverClient = archiverclient.NewArchiverClient(archiverclient.NewProxyRetriever(config.Conf.Bot.ObjectStore), []byte(config.Conf.Bot.AesKey))
	utils.SecureProxyClient = secureproxy.NewSecureProxy(config.Conf.SecureProxyUrl)

	utils.LoadEmoji()

	i18n.Init()

	if config.Conf.Bot.ProxyUrl != "" {
		request.RegisterHook(utils.ProxyHook)
	}

	logger.Info("Connecting to Redis")
	redis.Client = redis.NewRedisClient()

	socketManager := livechat.NewSocketManager()
	go socketManager.Run()

	go ListenChat(redis.Client, socketManager)

	if !config.Conf.Debug {
		rpc.PremiumClient = premium.NewPremiumLookupClient(
			redis.Client.Client,
			cache.Instance.PgCache,
			database.Client,
		)
	} else {
		c := premium.NewMockLookupClient(premium.Whitelabel, model.EntitlementSourcePatreon)
		rpc.PremiumClient = &c
	}

	logger.Info("Starting server")
	app.StartServer(logger, socketManager)
}

func ListenChat(client *redis.RedisClient, sm *livechat.SocketManager) {
	ch := make(chan chatrelay.MessageData)
	go chatrelay.Listen(client.Client, ch)

	for event := range ch {
		sm.BroadcastMessage(event)
	}
}

func startPprof() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/{action}", pprof.Index)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	go func() {
		http.ListenAndServe(":6060", mux)
	}()
}
