package main

type PriceCategory string

const (
	PriceCategoryFree    PriceCategory = "FREE"
	PriceCategoryRegular PriceCategory = "PAID"

	StudentPrice int = 50
	RegularPrice int = 100
	OnlinePrice  int = 100
)

type Price struct {
	Student  int
	Regular  int
	Online   int
	Category PriceCategory
}

func GetPriceFromShowType(showType ShowType) Price {
	switch showType {
	case ShowTypeClashOfTitans:
		return Price{
			Student:  StudentPrice,
			Regular:  RegularPrice,
			Online:   OnlinePrice,
			Category: PriceCategoryRegular,
		}
	case ShowTypeDuoLab:
		return Price{
			Student:  StudentPrice,
			Regular:  RegularPrice,
			Online:   OnlinePrice,
			Category: PriceCategoryRegular,
		}
	case ShowTypeStoryNight:
		return Price{
			Student:  StudentPrice,
			Regular:  RegularPrice,
			Online:   OnlinePrice,
			Category: PriceCategoryRegular,
		}
	case Project:
		return Price{
			Student:  0,
			Regular:  0,
			Online:   0,
			Category: PriceCategoryFree,
		}
	case ShowTypeCProject:
		return Price{
			Student:  0,
			Regular:  0,
			Online:   0,
			Category: PriceCategoryFree,
		}
	case ShowTypeJam:
		return Price{
			Student:  0,
			Regular:  0,
			Online:   0,
			Category: PriceCategoryFree,
		}
	case ShowTypeRegular:
		return Price{
			Student:  50,
			Regular:  100,
			Online:   100,
			Category: PriceCategoryRegular,
		}
	default:
		return Price{
			Student:  50,
			Regular:  100,
			Online:   100,
			Category: PriceCategoryRegular,
		}
	}
}
