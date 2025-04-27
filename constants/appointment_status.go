package constants

type status struct {
	SCHEDULED string
	COMPLETED string
	CANCELLED string
}

var AppointmentStatus = status{
	SCHEDULED: "scheduled",
	COMPLETED: "completed",
	CANCELLED: "cancelled",
}
