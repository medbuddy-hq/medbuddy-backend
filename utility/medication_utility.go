package utility

import (
	"fmt"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/errors"
	"medbuddy-backend/internal/model"
	"time"
)

func GetDosages(startDate time.Time, medic *model.MedicationRequest) ([]model.Dosage, errors.InternalError) {
	if medic.DailyDosage != len(medic.DosageTimes) {
		return nil, errors.BadRequestError("'dailyDosage' should match the length of dosage times")
	}

	now := time.Now()
	currentDay := time.Date(now.Year(), now.Month(), now.Day(), 00, 00, 00, 00, time.Local)
	if startDate.Before(currentDay) {
		return nil, errors.BadRequestError("invalid startDate")
	}

	times, err := parseTime(medic.DosageTimes)
	if err != nil {
		return nil, errors.BadRequestError(err.Error())
	}

	var dosages []model.Dosage
	var counter int
	var previousTimeRef = startDate
	for i := 0; i < medic.TotalNumberOfDosage; i++ {
		currentDosageTime := times[counter]
		reminderTime := time.Date(previousTimeRef.Year(), previousTimeRef.Month(), previousTimeRef.Day(), currentDosageTime.Hour(),
			currentDosageTime.Minute(), currentDosageTime.Second(), currentDosageTime.Nanosecond(), time.Local)
		if i > 0 {
			// check if this current reminder time is before the previous saved time
			// if it is, then add a day to the current reminder time to move to the next day
			if reminderTime.Before(dosages[i-1].ReminderTime) {
				reminderTime = reminderTime.Add(24 * time.Hour)
			}
		}

		previousTimeRef = reminderTime
		dose := model.Dosage{
			ReminderTime: reminderTime,
			Status:       constant.DosageNotTaken,
			IsActive:     true,
		}
		dosages = append(dosages, dose)

		counter++
		if counter >= len(times) {
			counter = 0
		}
	}

	if len(dosages) > 0 {
		if dosages[0].ReminderTime.Before(time.Now()) {
			return nil, errors.BadRequestError("invalid reminder time for your first dosage")
		}
	}

	return dosages, nil
}

func parseTime(times []string) ([]time.Time, error) {
	var result []time.Time
	if len(times) <= 0 {
		return nil, fmt.Errorf("no dosage time defined")
	}

	for _, t := range times {
		v, err := time.Parse(time.TimeOnly, t)
		if err != nil {
			return nil, fmt.Errorf("invalid dosage time: %v", t)
		}

		result = append(result, v)
	}

	return result, nil
}

func MedicationToMedicationResponse(medic *model.Medication) model.MedicationResponse {
	return model.MedicationResponse{
		ID:                  medic.ID,
		Name:                medic.Name,
		StartDate:           medic.StartDate,
		EndDate:             medic.EndDate,
		DosageQuantity:      medic.DosageQuantity,
		DailyDosage:         medic.DailyDosage,
		DosagesTaken:        medic.DosagesTaken,
		TotalNumberOfDosage: medic.TotalNumberOfDosage,
		Treatment:           medic.Treatment,
		Comment:             medic.Comment,
		CreatedAt:           medic.CreatedAt,
		UpdatedAt:           medic.UpdatedAt,
		MedicineID:          medic.MedicineID,
		PatientID:           medic.PatientID,
	}
}
