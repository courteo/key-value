package initializer

import (
	"errors"

	"github.com/courteo/key-value/internal/database/storage/replication"
	"github.com/courteo/key-value/internal/domain/config"
	"github.com/courteo/key-value/pkg/common"
	"github.com/courteo/key-value/pkg/tcp"
	"go.uber.org/zap"
)

func CreateReplica(replicaCfg config.Replication, walCfg config.WAL, logger *zap.Logger) (interface{}, error) {
	supportedTypes := map[string]struct{}{
		masterType: {},
		slaveType:  {},
	}

	if _, found := supportedTypes[replicaCfg.ReplicaType]; !found {
		return nil, errors.New("replica type is incorrect")
	}

	if replicaCfg.MasterAddress == "" {
		return nil, errors.New("master address is incorrect")
	}

	maxMessageSize := defaultMaxSegmentSize
	masterAddress := replicaCfg.MasterAddress
	syncInterval := defaultReplicationSyncInterval
	walDirectory := defaultWALDataDirectory

	if replicaCfg.SyncInterval != 0 {
		syncInterval = replicaCfg.SyncInterval
	}

	if walCfg.DataDirectory != "" {
		walDirectory = walCfg.DataDirectory
	}

	if walCfg.MaxSegmentSize != "" {
		size, _ := common.ParseSize(walCfg.MaxSegmentSize)
		maxMessageSize = size
	}

	idleTimeout := syncInterval * 3
	if replicaCfg.ReplicaType == masterType {
		maxReplicasNumber := defaultMaxReplicasNumber
		if replicaCfg.MaxReplicasNumber != 0 {
			maxReplicasNumber = replicaCfg.MaxReplicasNumber
		}

		var options []tcp.ServerOption
		options = append(options, tcp.WithServerIdleTimeout(idleTimeout))
		options = append(options, tcp.WithServerBufferSize(uint(maxMessageSize)))
		options = append(options, tcp.WithServerMaxConnectionsNumber(uint(maxReplicasNumber)))
		server, err := tcp.NewServer(masterAddress, logger, options...)
		if err != nil {
			return nil, err
		}

		return replication.NewMaster(server, walDirectory, logger)
	} else {
		var options []tcp.ClientOption
		options = append(options, tcp.WithClientIdleTimeout(idleTimeout))
		options = append(options, tcp.WithClientBufferSize(uint(maxMessageSize)))
		client, err := tcp.NewClient(masterAddress, options...)
		if err != nil {
			return nil, err
		}

		return replication.NewSlave(client, walDirectory, syncInterval, logger)
	}
}
