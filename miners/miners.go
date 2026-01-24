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
	ID        uuid.UUID `json:"id"`
	MinerType MinerType `json:"minerType"`
	CoalPower Coal      `json:"coalPower"`
	Energy    Energy    `json:"energy"`
	Cost      Coal      `json:"cost"`
}
