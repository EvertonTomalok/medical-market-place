package dto

type WorkerDTO struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	IsActive   bool   `json:"is_active"`
	Profession string `json:"profession"`
}
