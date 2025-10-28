package dto_in

import "fmt"

var (
	ErrEmptyMinerType      = fmt.Errorf("empty MinerType")
	ErrUnknownCommandMiner = fmt.Errorf("unknown MinerType")

	ErrUnknownCommandItem = fmt.Errorf("unknown ItemType")
	ErrEmptyItemType      = fmt.Errorf("empty ItemType")
)
