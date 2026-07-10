package showcaserecord_typing

type NameSection struct {
	Name string `json:"name" binding:"required"`
	// Nil if index not provided
	Index *int `json:"index"`
}

type EmailSection struct {
	Email string `json:"mail" binding:"required"`
	Index *int   `json:"index"`
}

type PhoneNumberSection struct {
	PhoneNumber string `json:"phonenumber" binding:"required"`
	Index       *int   `json:"index"`
}

type AddressSection struct {
	Address string `json:"address" binding:"required"`
	Index   *int   `json:"index"`
}

type SocialMediaSection struct {
	SocialMedia string `json:"socialmedia" binding:"required"`
	Index       *int   `json:"index"`
}

type JobExperienceSection struct {
	CompanyName    string `json:"companyname" binding:"required"`
	JobTitle       string `json:"jobtitle" binding:"required"`
	JobDescription string `json:"jobdescription" binding:"required"`
	Skill          string `json:"skill" binding:"required"`
	StartDate      string `json:"startdate" binding:"required"`
	EndDate        string `json:"enddate" binding:"required"`
	Index          *int   `json:"index"`
}

type EducationSection struct {
	InstitutionName string `json:"institutionname" binding:"required"`
	DegreeType      string `json:"degreetype" binding:"required"`
	Score           string `json:"score" binding:"required"`
	Skill           string `json:"skill" binding:"required"`
	StartDate       string `json:"startdate" binding:"required"`
	EndDate         string `json:"enddate" binding:"required"`
	Index           *int   `json:"index"`
}

type SkillSection struct {
	Skill string `json:"skill" binding:"required"`
	Index *int   `json:"index"`
}

type CertificateSection struct {
	Certificate string `json:"certificate" binding:"required"`
	Index       *int   `json:"index"`
}

type LanguageSection struct {
	Language string `json:"language" binding:"required"`
	Index    *int   `json:"index"`
}

type ProjectSection struct {
	ProjectTitle       string `json:"projecttitle" binding:"required"`
	ProjectDescription string `json:"projectdescription" binding:"required"`
	StartDate          string `json:"startdate" binding:"required"`
	EndDate            string `json:"enddate" binding:"required"`
	Index              *int   `json:"index"`
}

type ToneDetectionModelResponse struct {
	Prediction [][]float64 `json:"prediction"`
}

type TranslationModelResponse struct {
	Prediction string `json:"prediction"`
}

type SpecificPortoflioDataRequest struct {
	SectionTitle string `json:"sectiontitle" binding:"required"`
	Index        string `json:"index" binding:"required"`
}
