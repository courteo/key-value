package config

import (
	"time"
)

type Config struct {
	Engine      *Engine      `yaml:"engine"`
	WAL         *WAL         `yaml:"wal"`
	Replication *Replication `yaml:"replication"`
	Network     *Network     `yaml:"network"`
	Logging     *Logging     `yaml:"logging"`
}

type Engine struct {
	Type             string `yaml:"type"`
	PartitionsNumber uint   `yaml:"partitions_number"`
}

type WAL struct {
	FlushingBatchLength  int           `yaml:"flushing_batch_length"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       string        `yaml:"max_segment_size"`
	DataDirectory        string        `yaml:"data_directory"`
}

type Replication struct {
	ReplicaType       string        `yaml:"replica_type"`
	MasterAddress     string        `yaml:"master_address"`
	SyncInterval      time.Duration `yaml:"sync_interval"`
	MaxReplicasNumber int           `yaml:"max_replicas_number"`
}

type Network struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type Logging struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}
