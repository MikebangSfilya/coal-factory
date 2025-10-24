package dto

import (
	"coalFactory/equipment"
	"coalFactory/factory/statistic"
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
		return ErrEmptyMinerType
	}

}

type DTORBuyItem struct {
	ItemType string `json:"item_type"`
}

func (v *DTORBuyItem) Validate() error {
	if v.ItemType == "" {
		return ErrEmptyItemType
	}
	switch v.ItemType {
	case equipment.PickType, equipment.TrolleyType, equipment.VentType:
		return nil
	default:
		return ErrEmptyItemType
	}

}

type DTORespItem struct {
	Status string `json:"status"`
	Item   string `json:"Item"`
}

func NewResp(itemType string) DTORespItem {
	return DTORespItem{
		Status: "purchased",
		Item:   itemType,
	}
}

type DTOStats struct {
	Balance      int64
	TotalBalance int64
	TotalTime    string
}

func DtoStatsNew(companyStats statistic.CompanyStats) DTOStats {
	return DTOStats{
		Balance:      companyStats.GetBalance(),
		TotalBalance: companyStats.GetTotalBalance(),
		TotalTime:    companyStats.TimeCompleted(),
	}
}
