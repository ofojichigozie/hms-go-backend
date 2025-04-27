package constants

type role struct {
	ADMIN        string
	DOCTOR       string
	RECEPTIONIST string
}

var Roles = role{
	ADMIN:        "admin",
	DOCTOR:       "doctor",
	RECEPTIONIST: "receptionist",
}
