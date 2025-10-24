package dto_in

import "fmt"

var (
	ErrEmptyMinerType     = fmt.Errorf("empty MinerType")
	ErrUnknowCommandMiner = fmt.Errorf("unknow MinerType")

	ErrUnknowCommandItem = fmt.Errorf("unknow ItemType")
	ErrEmptyItemType     = fmt.Errorf("empty ItemType")
)
