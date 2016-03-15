package lti

const (
	MessageTypeBasicLaunchRequest = "basic-lti-launch-request"

	ParamLTIMessageType = "lti_message_type"
	ParamRoles = "roles"

	// roles that may be passed as part of an LTI launch
	// NOTE: just contains the small number of LIS roles used by Canvas, as
	//       of 081679e.  many more can be found in the LTI spec:
	//       http://www.imsglobal.org/specs/ltiv1p1p1/implementation-guide#toc-19
	RoleContentDeveloper = "ContentDeveloper"
	RoleInstructor = "Instructor"
	RoleLISInstroleAdministrator = "urn:lti:instrole:ims/lis/Administrator"
	RoleLISInstroleInstructor = "urn:lti:instrole:ims/lis/Instructor"
	RoleLISRoleContentDeveloper = "urn:lti:role:ims/lis/ContentDeveloper"
	RoleLISRoleInstructor = "urn:lti:role:ims/lis/Instructor"
	RoleLISRoleLearner = "urn:lti:role:ims/lis/Learner"
	RoleLISRoleNonCreditLearner = "urn:lti:role:ims/lis/Learner/NonCreditLearner"
	RoleLISRoleTeachingAssistant = "urn:lti:role:ims/lis/TeachingAssistant"
	RoleLISSysroleNone = "urn:lti:sysrole:ims/lis/None"
	RoleLISSysroleSysAdmin = "urn:lti:sysrole:ims/lis/SysAdmin"
	RoleLISSysroleUser = "urn:lti:sysrole:ims/lis/User"
	RoleLearner = "Learner"
    RoleLISInstroleObserver = "urn:lti:instrole:ims/lis/Observer"
    RoleLISInstroleStudent = "urn:lti:instrole:ims/lis/Student"
)

var staffRoles = []string{
	RoleContentDeveloper,
	RoleInstructor,
	RoleLISInstroleAdministrator,
	RoleLISInstroleInstructor,
	RoleLISRoleContentDeveloper,
	RoleLISRoleInstructor,
	RoleLISRoleTeachingAssistant,
	RoleLISSysroleSysAdmin,
}
func StaffRoles() []string {
	return staffRoles
}
