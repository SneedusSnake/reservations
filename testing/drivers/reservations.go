package drivers

type Reservations interface{
	UserRequestsSubjectsList()
	UserSeesSubjects(subject ...string)
}
