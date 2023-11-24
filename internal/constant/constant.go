package constant

const AppName = "medbuddy"

const (
	Practitioner string = "practitioner"
	Patient      string = "patient"
	CounterKey   string = "counter"
)

const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

var Roles = map[string]int{
	Practitioner: 1,
	Patient:      2,
}

const (
	UsersCollection         = "users"
	PatientsCollection      = "patients"
	PractitionersCollection = "practitioners"
	MedicineCollection      = "medicines"
	MedicationCollection    = "medications"
	DosageCollection        = "dosages"
)

const (
	DosageTaken    = "taken"
	DosageSkipped  = "skipped"
	DosageNotTaken = "not taken"
)
