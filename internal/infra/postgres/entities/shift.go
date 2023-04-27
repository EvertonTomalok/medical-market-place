package entities

import (
	"database/sql"
	"time"

	"github.com/EvertonTomalok/marketplace-health/internal/domain/dto"
)

type Shift struct {
	ID         int64         `json:"id"`
	Start      time.Time     `json:"start"`
	End        time.Time     `json:"end"`
	Profession Profession    `json:"profession"`
	IsDeleted  bool          `json:"is_deleted"`
	Facility   Facility      `json:"facility"`
	WorkerId   sql.NullInt64 `json:"worker_id,omitempty"`
}

func (s Shift) ToDTO() dto.ShiftDTO {
	shiftDTO := dto.ShiftDTO{
		ID:         s.ID,
		Start:      s.Start,
		End:        s.End,
		Profession: s.Profession.ToString(),
		Facility:   dto.FacilityDTO(s.Facility),
		IsDeleted:  s.IsDeleted,
	}

	if s.WorkerId.Valid {
		shiftDTO.WorkerId = s.WorkerId.Int64
	}

	return shiftDTO
}

type ShiftAndFacilityRequirements struct {
	Shift
	FacilityRequirements []int64
}
