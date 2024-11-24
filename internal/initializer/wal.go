package initializer

import (
	"errors"

	"github.com/courteo/key-value/internal/database/filesystem"
	"github.com/courteo/key-value/internal/database/storage/wal"
	"github.com/courteo/key-value/internal/domain/config"
	"github.com/courteo/key-value/pkg/common"
	"go.uber.org/zap"
)

func CreateWAL(cfg *config.WAL, logger *zap.Logger) (*wal.WAL, error) {
	flushingBatchSize := defaultFlushingBatchSize
	flushingBatchTimeout := defaultFlushingBatchTimeout
	maxSegmentSize := defaultMaxSegmentSize
	dataDirectory := defaultWALDataDirectory

	if cfg.FlushingBatchLength != 0 {
		flushingBatchSize = cfg.FlushingBatchLength
	}

	if cfg.FlushingBatchTimeout != 0 {
		flushingBatchTimeout = cfg.FlushingBatchTimeout
	}

	if cfg.MaxSegmentSize != "" {
		size, err := common.ParseSize(cfg.MaxSegmentSize)
		if err != nil {
			return nil, errors.New("max segment size is incorrect")
		}

		maxSegmentSize = size
	}

	if cfg.DataDirectory != "" {
		// TODO: need to create a directory,
		// if it is missing
		dataDirectory = cfg.DataDirectory
	}

	segmentsDirectory := filesystem.NewSegmentsDirectory(dataDirectory)
	reader, err := wal.NewLogsReader(segmentsDirectory)
	if err != nil {
		return nil, err
	}

	segment := filesystem.NewSegment(dataDirectory, maxSegmentSize)
	writer, err := wal.NewLogsWriter(segment, logger)
	if err != nil {
		return nil, err
	}

	return wal.New(writer, reader, flushingBatchTimeout, flushingBatchSize)
}
