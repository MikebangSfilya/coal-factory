package dto_out

import (
	"coalFactory/factory/statistic"
)

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
	Balance             int64
	TotalBalance        int64
	TotalTime           string
	LittleMinersHired   int
	NormalMinersHired   int
	PowerfulMinersHired int
}

func DtoStatsNew(companyStats statistic.CompanyStats) DTOStats {
	return DTOStats{
		Balance:             companyStats.GetBalance(),
		TotalBalance:        companyStats.GetTotalBalance(),
		TotalTime:           companyStats.TimeCompleted(),
		LittleMinersHired:   companyStats.LittleMinersHired,
		NormalMinersHired:   companyStats.NormalMinersHired,
		PowerfulMinersHired: companyStats.PowerfulMinersHired,
	}
}
