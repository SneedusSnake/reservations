package drivers

type Reservations interface{
	AdminAddsSubject(subject string)
	
	UserRequestsSubjectsList()
	UserSeesSubjects(subject ...string)
}
