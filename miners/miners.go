package miners

import (
	"github.com/google/uuid"
)

type (
	MinerType string
	Coal      int64
	Energy    int64
)

type MinerInfo struct {
	ID        uuid.UUID
	MinerType MinerType
	CoalPower Coal
	Energy    int64
	Cost      int64
}
