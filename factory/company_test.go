package factory_test

import (
	"coalFactory/equipment"
	"coalFactory/factory"
	"context"
	"testing"
	"time"
)

func TestPassiveIncome(t *testing.T) {

	t.Run("passive", func(t *testing.T) {

		eq := equipment.NewEquipment()

		var incomes []int64

		comp := factory.NewCompany(context.Background(), eq)
		go comp.PassiveIncome()
		for i := 0; i < 3; i++ {
			select {
			case coal := <-comp.Income:
				if coal == 0 {
					t.Errorf("expected to receive values")
				}
				incomes = append(incomes, int64(coal))
			case <-time.After(2 * time.Second):
				t.Errorf("income is taking too long")

			}
		}

		for _, v := range incomes {
			if v != 1 {
				t.Errorf("value exceeds base income")
			}
		}

	})

}
