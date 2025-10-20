package dto

type DTOHireMiner struct {
	MinerType string `json:"miner_type"`
}

func (v *DTOHireMiner) Validate() error {
	if v.MinerType == "" {
		return errEmptyType
	}
	return nil
}

type DTOREq struct {
	Pick    string
	Vent    string
	Trolley string
	AllBuy  bool
}

type DTOresponce struct {
}
