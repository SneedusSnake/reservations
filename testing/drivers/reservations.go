package drivers

type Reservations interface{
	AdminAddsSubject(subject string)
	AdminAddsTagsToSubject(subject string, tags ...string)
	
	UserRequestsSubjectsList()
	UserRequestsSubjectTags(subject string)
	UserRequestsReservationsList()
	UserRequestsReservationForSubject(user string, subject string, minutes int)
	UserRequestsReservationRemoval(user string, subject string)

	UserSeesSubjects(subject ...string)
	UserSeesSubjectTags(tags ...string)
	UserAcquiredReservationForSubject(user string, subject string, until string)
	UserSeesReservations(reservations ...string)
	SubjectHasAlreadyBeenReservedBy(user string, until string)

	ClockSet(time string)
}
