package main

import (
	"context"
	"order-processing/internal/store"
	"time"

	"go.uber.org/zap"
)

type application struct {
	store     *store.Storage
	logger    *zap.SugaredLogger
	processor processorConfig
	ctx       context.Context
}

type config struct {
	processor processorConfig
	db        dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type processorConfig struct {
	workers int
}