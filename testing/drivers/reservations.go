package drivers

type Reservations interface{
	AdminAddsSubject(subject string)
	AdminAddsTagsToSubject(subject string, tags ...string)
	
	UserRequestsSubjectsList()
	UserRequestsSubjectTags(subject string)

	UserSeesSubjects(subject ...string)
	UserSeesSubjectTags(tags ...string)
}
