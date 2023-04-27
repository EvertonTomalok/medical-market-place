package dto

import (
	"time"
)

type ShiftDTO struct {
	ID         int64       `json:"id"`
	Start      time.Time   `json:"start"`
	End        time.Time   `json:"end"`
	Profession string      `json:"profession"`
	IsDeleted  bool        `json:"is_deleted"`
	FacilityId int64       `json:"facility_id,omitempty"`
	Facility   FacilityDTO `json:"facility,omitempty"`
	WorkerId   int64       `json:"worker_id,omitempty"`
}

type GroupedByDateShift map[string][]ShiftDTO

type BetweenDates struct {
	Start time.Time
	End   time.Time
}
