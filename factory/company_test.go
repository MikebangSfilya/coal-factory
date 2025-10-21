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
		eq := equipment.NewEquipmet()

		var incomes []int

		comp := factory.NewCompany(context.Background(), eq)
		go comp.PassiveIncome()
		for i := 0; i < 3; i++ {
			select {
			case coal := <-comp.Income:
				if coal == 0 {
					t.Errorf("ожидалось получение значений")
				}
				incomes = append(incomes, coal)
			case <-time.After(2 * time.Second):
				t.Errorf("доход не идет слишком долго")

			}
		}

		for _, v := range incomes {
			if v != 1 {
				t.Errorf("значение выше базового дохода")
			}
		}

	})

}
