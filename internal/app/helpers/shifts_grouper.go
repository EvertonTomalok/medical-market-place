package helpers

import (
	"github.com/EvertonTomalok/marketplace-health/internal/domain/dto"
)

func GroupShiftByDate(shifts []dto.ShiftDTO) dto.GroupedByDateShift {
	var grouped dto.GroupedByDateShift = make(dto.GroupedByDateShift)

	for _, shift := range shifts {
		date := shift.Start.Format("2006-01-02")

		if values, ok := grouped[date]; ok {
			grouped[date] = append(values, shift)
		} else {
			grouped[date] = []dto.ShiftDTO{shift}
		}
	}

	return grouped
}
