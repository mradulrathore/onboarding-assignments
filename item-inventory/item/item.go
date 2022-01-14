package item

import (
	"fmt"

	"github.com/mradulrathore/onboarding-assignments/item-inventory/item/enum"
)

type Item struct {
	Name     string        `json:"name"`
	Price    float64       `json:"price"`
	Quantity int           `json:"quantity"`
	Type     enum.ItemType `json:"type"`
}

func New(name string, price float64, quantity int, typeItem string) (item Item, err error) {
	item.Name = name
	item.Price = price
	item.Quantity = quantity
	item.Type, err = enum.ItemTypeString(typeItem)
	if err != nil {
		return
	}

	err = item.validate()
	return
}

// func checkNegativeValue(value interface{}) (err error) {

// 	val, _ := value.(int)
// 	if val < 0 {
// 		err = NegativeQuantErr
// 	}
// 	return
// }

func (item Item) validate() (err error) {
	if item.Quantity < 0 {
		err = NegativeQuantErr
	}
	if item.Price < 0 {
		err = NegativePriceErr
	}
	return
	//return validation.ValidateStruct(&item, validation.Field(&item.Quantity, validation.By(checkNegativeValue)))
}

func (item Item) String() string {
	return fmt.Sprintf("[%s, %g, %d,%s,%g,%g]", item.Name, item.Price, item.Quantity, item.Type.String(), item.getTax(), item.getEffectivePrice())
}

func (item Item) getTax() (tax float64) {
	switch item.Type {
	case enum.Raw:
		//raw: 12.5% of the item cost
		tax = RAWItmTaxRate * item.Price
	case enum.Manufactured:
		// manufactured: 12.5% of the item cost + 2% of (item cost + 12.5% of the item cost)
		tax = ManufacturedItmTaxRate*item.Price + ManufacturedItmExtraTaxRate*(item.Price+ManufacturedItmTaxRate*item.Price)
	case enum.Imported:
		//imported: 10% import duty on item cost + a surcharge
		tax = ImportDuty * item.Price
	}

	return
}

func (item Item) getEffectivePrice() (effectivePrice float64) {
	surcharge := 0.0
	tax := item.getTax()

	switch item.Type {
	case enum.Raw:
		effectivePrice = item.Price + tax + surcharge
	case enum.Manufactured:
		effectivePrice = item.Price + tax + surcharge
	case enum.Imported:
		priceTemp := ImportDuty*item.Price + tax
		surcharge = item.importSurcharge(priceTemp)
		effectivePrice = priceTemp + surcharge
	}

	return
}

func (item Item) importSurcharge(price float64) float64 {
	if price <= ImportDutyLimit1 {
		return ImportDutyLimit1SurchargeAmt
	} else if price <= ImportDutyLimit2 {
		return ImportDutyLimit2SurchargeAmt
	} else {
		return price * ExceedeImportDutyLimit2SurchargeRate
	}
}