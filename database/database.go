package database

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgconn/stmtcache"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jadevelopmentgrp/Ticket-Dashboard/config"
	database "github.com/jadevelopmentgrp/Tickets-Database"
	"github.com/sirupsen/logrus"
)

var Client *database.Database

func ConnectToDatabase() {
	config, err := pgxpool.ParseConfig(config.Conf.Database.Uri)
	if err != nil {
		panic(err)
	}

	config.ConnConfig.LogLevel = pgx.LogLevelWarn
	config.ConnConfig.Logger = logrusadapter.NewLogger(logrus.New())

	config.MinConns = 1
	config.MaxConns = 3

	config.ConnConfig.BuildStatementCache = func(conn *pgconn.PgConn) stmtcache.Cache {
		return stmtcache.New(conn, stmtcache.ModeDescribe, 512)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	Client = database.NewDatabase(pool)
}
