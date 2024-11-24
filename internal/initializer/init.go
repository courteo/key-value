package initializer

import (
	"errors"
	"time"

	in_memory "github.com/courteo/key-value/internal/database/storage/engine/in-memory"
	"github.com/courteo/key-value/internal/database/storage/replication"
	"github.com/courteo/key-value/internal/database/storage/wal"
	"github.com/courteo/key-value/internal/domain/config"
	logger2 "github.com/courteo/key-value/pkg/logger"
	"github.com/courteo/key-value/pkg/tcp"
	"go.uber.org/zap"
)

const (
	masterType = "master"
	slaveType  = "slave"
)

const (
	defaultReplicationSyncInterval = time.Second
	defaultMaxReplicasNumber       = 5
)

const (
	defaultFlushingBatchSize    = 100
	defaultFlushingBatchTimeout = time.Millisecond * 10
	defaultMaxSegmentSize       = 10 << 20
	defaultWALDataDirectory     = "./data/spider/wal"
)

const defaultAddress = ":3000"

type Initializer struct {
	wal    *wal.WAL
	engine *in_memory.Engine
	server *tcp.Server
	slave  *replication.Slave
	master *replication.Master
	logger *zap.Logger
}

func New(cfg *config.Config) (*Initializer, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	logger, err := logger2.New(cfg.Logging)
	if err != nil {
		return nil, err
	}

	address := getAddress(cfg.Network)

	server, err := tcp.NewServer(address, logger)
	if err != nil {
		return nil, err
	}

	engine, err := createEngine(cfg.Engine, logger)
	if err != nil {
		return nil, err
	}

	wal, err := CreateWAL(cfg.WAL, logger)
	if err != nil {
		return nil, err
	}

	replica, err := CreateReplica(*cfg.Replication, *cfg.WAL, logger)
	if err != nil {
		return nil, err
	}

	initializer := &Initializer{
		wal:    wal,
		engine: engine,
		server: server,
		logger: logger,
	}

	switch v := replica.(type) {
	case *replication.Slave:
		initializer.slave = v
	case *replication.Master:
		initializer.master = v
	}

	return initializer, nil
}

func getAddress(cfg *config.Network) string {
	if cfg == nil || cfg.Address == "" {
		return defaultAddress
	}

	return cfg.Address
}
