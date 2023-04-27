package entities

type DocumentWorker struct {
	ID         int64 `json:"id"`
	WorkerId   int64 `json:"worker_id"`
	DocumentId int64 `json:"document_id"`
}
