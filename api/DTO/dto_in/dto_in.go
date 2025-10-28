package dto_in

import (
	"coalFactory/equipment"
	"coalFactory/miners"
)

type DTOHireMiner struct {
	MinerType string `json:"miner_type"`
}

func (v *DTOHireMiner) Validate() error {
	if v.MinerType == "" {
		return ErrEmptyMinerType
	}
	switch v.MinerType {
	case miners.MinerTypeLittle, miners.MinerTypeNormal, miners.MinerTypePowerful:
		return nil
	default:
		return ErrUnknownCommandMiner
	}

}

type DTOBuyItem struct {
	ItemType string `json:"item_type"`
}

func (v *DTOBuyItem) Validate() error {
	if v.ItemType == "" {
		return ErrEmptyItemType
	}
	switch v.ItemType {
	case equipment.PickType, equipment.TrolleyType, equipment.VentType:
		return nil
	default:
		return ErrUnknownCommandItem
	}

}
