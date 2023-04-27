package custom_errors

import "fmt"

type WorkerNotFound struct {
	Id int64
}

func (wnf WorkerNotFound) Error() string {
	return fmt.Sprintf("worker with id %d not found", wnf.Id)
}

type WorkerInactive struct {
	Id int64
}

func (wnf WorkerInactive) Error() string {
	return fmt.Sprintf("worker with id %d is inactive", wnf.Id)
}

type FacilityNotFound struct {
	Id int64
}

func (fnf FacilityNotFound) Error() string {
	return fmt.Sprintf("facility with id %d not found", fnf.Id)
}

type FacilityInactive struct {
	Id int64
}

func (fnf FacilityInactive) Error() string {
	return fmt.Sprintf("facility with id %d is inactive", fnf.Id)
}
