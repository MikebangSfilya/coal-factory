package equipment_test

import (
	"coalFactory/config"
	"coalFactory/equipment"
	"coalFactory/factory"
	"context"
	"testing"
)

//Чтобы протестировать покупку нам нужно
//Создать оборудование
//Иметь деньги, для этого нужна компания (TODO переделать логику)
//Вызвать метод
//Сравнить с ожиданием

func TestBuyPick(t *testing.T) {
	cfg := config.Load()
	eq := equipment.NewEquipmet()
	company := factory.NewCompany(context.Background(), eq)

	company.SetBalance(10000)
	startBalance := company.GetBalance()

	//Покупка кирки
	res, err := company.Buy("pick")
	if err != nil {
		t.Errorf("Неожидонная ошибка %v", err)
	}

	if res == nil {
		t.Fatal("должен был вернуться объект оборудования")
	}

	//Проверка на покупку кирки
	if !res.Pick.IsBuyed {
		t.Errorf("Кирка должна была быть куплена")
	}

	endBalance := company.GetBalance()
	expectedBalance := startBalance - cfg.PickCost

	if endBalance != expectedBalance {
		t.Errorf("Неверный балананс, ожидалось %d, получилось %d", expectedBalance, endBalance)
	}

}

func TestBuyAll(t *testing.T) {
	cfg := config.Configurate{
		PickCost:    5000,
		VentCost:    15000,
		TrolleyCost: 50000,
	}
	equipment.Init(&cfg)
	tests := []struct {
		name         string
		startBalance int64
		item         string
		errorNeed    bool
		errorType    error
		wantBalance  int64
		wantBought   bool
	}{
		{
			name:         "all ok Pick",
			startBalance: 10_000,
			item:         "pick",
			errorNeed:    false,
			errorType:    nil,
			wantBalance:  10_000 - int64(equipment.PickCost),
			wantBought:   true,
		},
		{
			name:         "all ok Vent",
			startBalance: 20_000,
			item:         "vent",
			errorNeed:    false,
			errorType:    nil,
			wantBalance:  20_000 - int64(equipment.VentCost),
			wantBought:   true,
		},
		{
			name:         "all ok Trolley",
			startBalance: 55_000,
			item:         "trolley",
			errorNeed:    false,
			errorType:    nil,
			wantBalance:  55_000 - int64(equipment.TrolleyCost),
			wantBought:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			eq := equipment.NewEquipmet()
			comp := factory.NewCompany(context.Background(), eq)

			comp.SetBalance(int(tc.startBalance))

			res, err := comp.Buy(tc.item)
			if tc.errorNeed {
				if err == nil {
					t.Errorf("ожидалась ошибка, ошибики нет")
				}
				if tc.errorType != nil && err != tc.errorType {
					t.Errorf("ожидалась ошибка %v, получили %v", tc.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Неожиданная ошибка %v", err)
				}
			}

			gotBalance := comp.GetBalance()

			if gotBalance != int(tc.wantBalance) {
				t.Errorf("баланс неверный, ожидалось %d, получилось %d", tc.wantBalance, gotBalance)
			}

			if tc.wantBought {
				if res == nil {
					t.Errorf("ожидалось покупка оборудования, но ничего не было куплено")
				} else {
					switch tc.item {
					case "pick":
						if !res.Pick.IsBuyed {
							t.Errorf("кирка должна была быть куплена")
						}
					case "vent":
						if !res.Vent.IsBuyed {
							t.Errorf("вентиляция должна была быть куплена")
						}
					case "trolley":
						if !res.Trolley.IsBuyed {
							t.Errorf("вагонетка должна была быть куплена")
						}
					}
				}
			} else {
				if res != nil {
					t.Errorf("Покупка не ожидалась вообще")
				}
			}

		})
	}

}

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
