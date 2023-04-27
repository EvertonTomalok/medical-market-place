package entities

import "github.com/EvertonTomalok/marketplace-health/internal/domain/dto"

type Facility struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func (f *Facility) ToDTO() dto.FacilityDTO {
	return dto.FacilityDTO{
		ID:       f.ID,
		Name:     f.Name,
		IsActive: f.IsActive,
	}
}
