package usecases

import (
	"context"
	"time"

	"github.com/EvertonTomalok/marketplace-health/internal/app/helpers"
	"github.com/EvertonTomalok/marketplace-health/internal/domain/custom_errors"
	"github.com/EvertonTomalok/marketplace-health/internal/domain/dto"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/sirupsen/logrus"
)

type chanReturn struct {
	data interface{}
	err  error
}

type multiplexingShiftsAndDocsReturn struct {
	documentsIDs []int64
	shifts       []entities.Shift
	errors       []error
}

type multiplexingFacilitiesActiveAndDocsReturn struct {
	documentsIDs []int64
	facilities   []entities.Facility
	errors       []error
}

type ShiftUsecases struct {
	logger *logrus.Logger
}

func (su ShiftUsecases) FindAvailableShiftsFasterVersion(
	ctx context.Context,
	profession entities.Profession,
	workerID int64,
	facilityID *int64,
	startTime *time.Time,
	endTime *time.Time,
	offset int64,
	limit int64,
) ([]dto.ShiftDTO, []error) {
	_, err := su.checkWorker(ctx, workerID)
	if err != nil {
		return []dto.ShiftDTO{}, []error{err}
	}

	var infoReturn multiplexingFacilitiesActiveAndDocsReturn = su.getWorkerDocumentsAndFacilities(ctx, workerID)
	if len(infoReturn.errors) > 0 {
		return []dto.ShiftDTO{}, infoReturn.errors
	}
	facilities := infoReturn.facilities
	workerDocumentsIds := infoReturn.documentsIDs

	if len(facilities) == 0 {
		su.logger.Debug("No facilities active. Return empty shifts...")
		return []dto.ShiftDTO{}, []error{}
	}

	allFacilitiesIds := []int64{}
	for _, f := range facilities {
		allFacilitiesIds = append(allFacilitiesIds, f.ID)
	}

	facilitiesRequirementsAggregated, err := FacilityRequirementsRepository.FindRequirementsByFacilitiesId(ctx, allFacilitiesIds)
	if err != nil {
		return []dto.ShiftDTO{}, []error{}
	}
	facilitiesWorkerHasDocumentsRequired := []int64{}
	for facilityId, requirements := range facilitiesRequirementsAggregated {
		if helpers.WorkerHasAllDocumentsRequired(workerDocumentsIds, requirements.DocumentsId) {
			facilitiesWorkerHasDocumentsRequired = append(facilitiesWorkerHasDocumentsRequired, facilityId)
		}
	}

	shiftsAlreadyTaken, err := ShiftRepository.FindShifts(
		ctx,
		profession,
		startTime,
		endTime,
		&workerID,
		0,
		0, // get all shifts from user
		nil,
	)
	if err != nil {
		return []dto.ShiftDTO{}, []error{err}
	}
	shiftsDatesTaken := helpers.GetRangeBetweenDates(shiftsAlreadyTaken)

	shifts, err := ShiftRepository.FindShiftsByFacilities(
		ctx,
		facilitiesWorkerHasDocumentsRequired,
		profession,
		startTime,
		endTime,
		offset,
		limit,
		&shiftsDatesTaken,
	)
	if err != nil {
		return []dto.ShiftDTO{}, []error{err}
	}

	shiftsDTO := make([]dto.ShiftDTO, 0)
	for _, s := range shifts {
		shiftsDTO = append(shiftsDTO, s.ToDTO())
	}

	return shiftsDTO, []error{}
}

func (su ShiftUsecases) FindAvailableShiftsSlowerVersion(
	ctx context.Context,
	profession entities.Profession,
	workerID int64,
	facilityID *int64,
	startTime *time.Time,
	endTime *time.Time,
	offset int64,
	limit int64,
) ([]dto.ShiftDTO, []error) {
	_, err := su.checkWorker(ctx, workerID)
	if err != nil {
		return []dto.ShiftDTO{}, []error{err}
	}

	infoReturn := su.getShiftsAndDocuments(
		ctx,
		profession,
		workerID,
		startTime,
		endTime,
		offset,
		limit,
	)
	if len(infoReturn.errors) > 0 {
		return []dto.ShiftDTO{}, infoReturn.errors
	}

	shifts := infoReturn.shifts
	workerDocuments := infoReturn.documentsIDs

	facilitiesIDS := make([]int64, 0)
	for _, shift := range shifts {
		facilitiesIDS = append(facilitiesIDS, shift.Facility.ID)
	}
	facilitiesRequirements, err := FacilityRequirementsRepository.FindRequirementsByFacilitiesId(ctx, facilitiesIDS)
	if err != nil {
		return []dto.ShiftDTO{}, []error{err}
	}
	var shiftsReturn []dto.ShiftDTO = make([]dto.ShiftDTO, 0)
	for _, s := range shifts {
		req, exists := facilitiesRequirements[s.Facility.ID]
		if exists && helpers.WorkerHasAllDocumentsRequired(workerDocuments, req.DocumentsId) {
			shiftsReturn = append(shiftsReturn, s.ToDTO())
		}
		shiftsReturn = append(shiftsReturn, s.ToDTO())
	}

	return shiftsReturn, make([]error, 0)
}

func (su ShiftUsecases) shiftConflict(shift entities.Shift, shiftsTaken []entities.Shift) bool {
	for _, shiftTaken := range shiftsTaken {
		shiftTakenConflictWithShiftStart := shiftTaken.Start.Before(shift.Start) && shiftTaken.End.After(shift.Start)
		shiftTakenConflictWithShiftEnd := shiftTaken.Start.Before(shift.End) && shiftTaken.End.After(shift.End)
		shiftWrapshiftTaken := shift.Start.Before(shiftTaken.Start) && shift.End.After(shiftTaken.End)

		if shiftTakenConflictWithShiftStart || shiftTakenConflictWithShiftEnd || shiftWrapshiftTaken {
			return true
		}
	}
	return false
}

func (su ShiftUsecases) checkWorker(ctx context.Context, workerId int64) (entities.Worker, error) {
	worker, er := WorkerRepository.FindById(ctx, workerId)

	if er != nil {
		return worker, er
	}
	if worker.ID == 0 {
		return worker, custom_errors.WorkerNotFound{Id: workerId}
	}
	if !worker.IsActive {
		return worker, custom_errors.WorkerInactive{Id: workerId}
	}

	return worker, nil
}

func (su ShiftUsecases) getWorkerDocumentsAndFacilities(
	ctx context.Context,
	workerId int64,
) multiplexingFacilitiesActiveAndDocsReturn {

	info := multiplexingFacilitiesActiveAndDocsReturn{}

	workerDocumentsChan := make(chan chanReturn)
	facilitiesWhitOpenShiftChan := make(chan chanReturn)

	go func() {
		documentsIds, err := DocumentWorkerRepository.FindDocumentsIds(ctx, workerId)

		workerDocumentsChan <- chanReturn{
			data: documentsIds,
			err:  err,
		}
	}()
	go func() {
		facilities, err := FacilityRepository.FindActive(ctx)
		facilitiesWhitOpenShiftChan <- chanReturn{
			data: facilities,
			err:  err,
		}
	}()

	for i := 0; i < 2; i++ {
		// Await both of these values
		// simultaneously, putting each one in array when arrives and its not nil.
		select {
		case docs := <-workerDocumentsChan:
			info.documentsIDs = docs.data.([]int64)
			if docs.err != nil {
				info.errors = append(info.errors, docs.err)
			}
		case facilities := <-facilitiesWhitOpenShiftChan:
			info.facilities = facilities.data.([]entities.Facility)
			if facilities.err != nil {
				info.errors = append(info.errors, facilities.err)
			}
		}
	}

	return info
}

func (su ShiftUsecases) getShiftsAndDocuments(
	ctx context.Context,
	profession entities.Profession,
	workerID int64,
	startTime *time.Time,
	endTime *time.Time,
	offset int64,
	limit int64,
) multiplexingShiftsAndDocsReturn {
	info := multiplexingShiftsAndDocsReturn{}

	shiftsChan := make(chan chanReturn)
	documentsIdsChan := make(chan chanReturn)

	go func() {
		shiftsAlreadyTaken, err := ShiftRepository.FindShifts(
			ctx,
			profession,
			startTime,
			endTime,
			&workerID,
			0,
			0, // get all shifts from user
			nil,
		)
		if err != nil {
			shiftsChan <- chanReturn{
				data: []entities.Shift{},
				err:  err,
			}
			return
		}

		shiftsDatesTaken := helpers.GetRangeBetweenDates(shiftsAlreadyTaken)

		shiftsAndRequirements, err := ShiftRepository.FindShifts(
			ctx,
			profession,
			startTime,
			endTime,
			nil,
			offset,
			limit,
			&shiftsDatesTaken,
		)
		shiftsChan <- chanReturn{
			data: shiftsAndRequirements,
			err:  err,
		}
	}()
	go func() {
		workerDocuments, err := DocumentWorkerRepository.FindDocumentsIds(ctx, workerID)
		documentsIdsChan <- chanReturn{
			data: workerDocuments,
			err:  err,
		}
	}()

	for i := 0; i < 2; i++ {
		// Await both of these values
		// simultaneously, putting each one in array when arrives and its not nil.
		select {
		case docs := <-documentsIdsChan:
			info.documentsIDs = docs.data.([]int64)
			if docs.err != nil {
				info.errors = append(info.errors, docs.err)
			}
		case shift := <-shiftsChan:
			info.shifts = shift.data.([]entities.Shift)
			if shift.err != nil {
				info.errors = append(info.errors, shift.err)
			}
		}
	}

	return info
}

// type multiplexingWorkerAndFacilityReturn struct {
// 	worker   entities.Worker
// 	facility entities.Facility
// 	errors   []error
// }

// func (su ShiftUsecases) checkFacility(ctx context.Context, facilityId int64) (entities.Facility, error) {
// 	facility, er := FacilityRepository.FindById(ctx, facilityId)

// 	if er != nil {
// 		return facility, er
// 	}
// 	if facility.ID == 0 {
// 		return facility, custom_errors.FacilityNotFound{Id: facilityId}
// 	}
// 	if !facility.IsActive {
// 		return facility, custom_errors.FacilityInactive{Id: facilityId}
// 	}

// 	return facility, nil
// }

// func (su ShiftUsecases) checkWorkerAndFacility(ctx context.Context, workerId, facilityId int64) multiplexingWorkerAndFacilityReturn {
// 	info := multiplexingWorkerAndFacilityReturn{}
// 	workerChan := make(chan chanReturn)
// 	facilityChan := make(chan chanReturn)

// 	go func() {
// 		worker, err := su.checkWorker(ctx, workerId)
// 		workerChan <- chanReturn{
// 			data: worker,
// 			err:  err,
// 		}
// 	}()
// 	go func() {
// 		facility, err := su.checkFacility(ctx, facilityId)
// 		facilityChan <- chanReturn{
// 			data: facility,
// 			err:  err,
// 		}
// 	}()

// 	for i := 0; i < 2; i++ {
// 		// Await both of these values
// 		// simultaneously, putting each one in array when arrives and its not nil.
// 		select {
// 		case workerReturn := <-workerChan:
// 			info.worker = workerReturn.data.(entities.Worker)
// 			if workerReturn.err != nil {
// 				info.errors = append(info.errors, workerReturn.err)
// 			}
// 		case facilityReturn := <-facilityChan:
// 			info.facility = facilityReturn.data.(entities.Facility)
// 			if facilityReturn.err != nil {
// 				info.errors = append(info.errors, facilityReturn.err)
// 			}
// 		}
// 	}

// 	return info
// }
