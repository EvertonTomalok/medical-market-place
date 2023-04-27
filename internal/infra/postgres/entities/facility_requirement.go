package entities

type FacilityRequirement struct {
	ID         int64 `json:"id"`
	FacilityId int64 `json:"facility_id"`
	DocumentId int64 `json:"document_id"`
}

type FacilityRequirementAggregated struct {
	FacilityId  int64   `json:"facility_id"`
	DocumentsId []int64 `json:"documents_id"`
}
