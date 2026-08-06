package showcase_dto

type NameSection struct {
	Name string `json:"name" binding:"required" mod:"trim,ucase" validate:"required,min=3,max=50"`
	// Nil if index not provided
	Index *int `json:"index"`
}

type EmailSection struct {
	Email string `json:"mail" binding:"required" mod:"trim,lcase" validate:"required,min=4,max=30,email,emailmx"`
	Index *int   `json:"index"`
}

type PhoneNumberSection struct {
	PhoneNumber string `json:"phonenumber" binding:"required" mod:"trim" validate:"required,min=3,max=17"`
	Index       *int   `json:"index"`
}

type AddressSection struct {
	Address string `json:"address" binding:"required" mod:"trim" validate:"required,min=3,max=300"`
	Index   *int   `json:"index"`
}

type SocialMediaSection struct {
	SocialMedia string `json:"socialmedia" binding:"required" mod:"trim" validate:"required,min=3,max=100"`
	Index       *int   `json:"index"`
}

type JobExperienceSection struct {
	CompanyName    string `json:"companyname" binding:"required" mod:"trim,ucase" validate:"required,min=3,max=100"`
	JobTitle       string `json:"jobtitle" binding:"required" mod:"trim,ucase" validate:"required,min=3,max=100"`
	JobDescription string `json:"jobdescription" binding:"required" mod:"trim" validate:"required,min=3,max=1000"`
	Skill          string `json:"skill" binding:"required"`
	StartDate      string `json:"startdate" binding:"required" mod:"trim" validate:"required,min=3,max=50"`
	EndDate        string `json:"enddate" binding:"required" mod:"trim" validate:"required,min=3,max=50"`
	Index          *int   `json:"index"`
}

type EducationSection struct {
	InstitutionName string `json:"institutionname" binding:"required" mod:"trim,ucase" validate:"required,min=3,max=100"`
	DegreeType      string `json:"degreetype" binding:"required" mod:"trim,ucase" validate:"required,min=3,max=100"`
	Score           string `json:"score" binding:"required" mod:"trim" validate:"min=3,max=20"`
	Skill           string `json:"skill" binding:"required"`
	StartDate       string `json:"startdate" binding:"required" mod:"trim" validate:"required,min=3,max=50"`
	EndDate         string `json:"enddate" binding:"required" mod:"trim" validate:"required,min=3,max=50"`
	Index           *int   `json:"index"`
}

type SkillSection struct {
	Skill string `json:"skill" binding:"required" mod:"trim,ucase" validate:"required,min=3,max=30"`
	Index *int   `json:"index"`
}

type CertificateSection struct {
	Certificate string `json:"certificate" binding:"required" mod:"trim,ucase" validate:"required,min=3,max=30"`
	Index       *int   `json:"index"`
}

type LanguageSection struct {
	Language string `json:"language" binding:"required" mod:"trim,ucase" validate:"required,min=3,max=30"`
	Index    *int   `json:"index"`
}

type ProjectSection struct {
	ProjectTitle       string `json:"projecttitle" binding:"required" mod:"trim" validate:"required,min=3,max=30"`
	ProjectDescription string `json:"projectdescription" binding:"required" mod:"trim" validate:"required,min=3,max=1000"`
	StartDate          string `json:"startdate" binding:"required" mod:"trim" validate:"required,min=3,max=50"`
	EndDate            string `json:"enddate" binding:"required" mod:"trim" validate:"required,min=3,max=50"`
	Index              *int   `json:"index"`
}

type SpecificPortoflioDataRequest struct {
	SectionTitle string `json:"sectiontitle" binding:"required" mod:"trim,lcase" validate:"required,min=3,max=300"`
	Index        string `json:"index" binding:"required"`
}
