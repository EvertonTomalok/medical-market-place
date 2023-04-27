package entities

import "github.com/EvertonTomalok/marketplace-health/internal/domain/dto"

type Worker struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	IsActive   bool       `json:"is_active"`
	Profession Profession `json:"profession"`
}

func (w *Worker) ToDTO() dto.WorkerDTO {
	return dto.WorkerDTO{
		ID:         w.ID,
		Name:       w.Name,
		IsActive:   w.IsActive,
		Profession: w.Profession.ToString(),
	}
}
