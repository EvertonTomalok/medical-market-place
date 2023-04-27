package helpers

import (
	"time"

	"github.com/EvertonTomalok/marketplace-health/internal/domain/dto"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
)

func GetRangeBetweenDates(shiftsTaken []entities.Shift) []dto.BetweenDates {
	dates := make(map[time.Time]dto.BetweenDates, 0)

	for _, s := range shiftsTaken {
		if value, ok := dates[s.Start]; ok {
			if s.End.After(value.End) {
				dates[s.Start] = dto.BetweenDates{
					Start: s.Start, End: s.End,
				}
			}
		} else {
			dates[s.Start] = dto.BetweenDates{
				Start: s.Start, End: s.End,
			}
		}
	}

	betweeDatesArray := make([]dto.BetweenDates, 0)
	for _, value := range dates {
		betweeDatesArray = append(betweeDatesArray, value)
	}

	return betweeDatesArray
}
