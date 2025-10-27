package equipment_test

import (
	"coalFactory/equipment"
	"testing"
)

func TestAllBuy(t *testing.T) {

	eq := equipment.NewEquipmet()

	eq1 := equipment.NewEquipmet()

	_, err := eq.Buy("pick")
	if err != nil {
		t.Errorf("pick не была куплена %v", err)
	}
	_, err = eq.Buy("vent")
	if err != nil {
		t.Errorf("vent не была куплена %v", err)
	}
	_, err = eq.Buy("trolley")
	if err != nil {
		t.Errorf("trolley не была куплена %v", err)
	}
	_, err = eq.Buy("dragon")
	if err == nil {
		t.Errorf("dragon не была куплена %v", err)
	}

	if eq1.AllBuyed() {
		t.Errorf("оборудование не было куплено")
	}

	if !eq.AllBuyed() {
		t.Errorf("оборудование не было куплено")
	}

}
