package showcaserecord_validator

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/bobch27/valtra-go"
	typing "resuming/api/showcaserecord/typing"
)

func ValidateNamePortfolioData(request typing.NameSection) (typing.NameSection, error) {
	v := valtra.NewCollector()

	name := typing.NameSection{
		Name: valtra.Val(request.Name, "Name").
			Transform(valtra.TrimSpace(), valtra.Uppercase()).
			Validate(
				valtra.Required[string]("Name is required."),
				valtra.MinLengthString(3, "Name must be at least 3 characters"),
				valtra.MaxLengthString(50, "Name must be at most 50 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.NameSection{}, v.Errors()[0]
	}

	return name, nil
}

func ValidateEmailPortfolioData(request typing.EmailSection) (typing.EmailSection, error) {
	v := valtra.NewCollector()

	email := typing.EmailSection{
		Email: valtra.Val(request.Email, "Email").
			Transform(valtra.TrimSpace(), valtra.Lowercase()).
			Validate(
				valtra.Required[string]("Email is required."),
				valtra.MinLengthString(4, "Email must be at least 4 characters"),
				valtra.MaxLengthString(30, "Email must be at most 30 characters"),
				valtra.Email("Email must be in correct email format"),
				func(v valtra.Value[string]) error {
					if !validateEmailMX(v.Value()) {
						return fmt.Errorf("Email domain must have valid MX or A records")
					}
					return nil
				},
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.EmailSection{}, v.Errors()[0]
	}

	return email, nil
}

func ValidatePhoneNumberPortfolioData(request typing.PhoneNumberSection) (typing.PhoneNumberSection, error) {
	v := valtra.NewCollector()

	phone_number := typing.PhoneNumberSection{
		PhoneNumber: valtra.Val(request.PhoneNumber, "Phone number").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Phone number is required."),
				valtra.MinLengthString(3, "Phone number must be at least 3 characters"),
				valtra.MaxLengthString(17, "Phone number must be at most 17 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.PhoneNumberSection{}, v.Errors()[0]
	}

	return phone_number, nil
}

func ValidateAddressPortfolioData(request typing.AddressSection) (typing.AddressSection, error) {
	v := valtra.NewCollector()

	address := typing.AddressSection{
		Address: valtra.Val(request.Address, "Address").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Address is required."),
				valtra.MinLengthString(3, "Address must be at least 3 characters"),
				valtra.MaxLengthString(300, "Address must be at most 300 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.AddressSection{}, v.Errors()[0]
	}

	return address, nil
}

func ValidateSocialMediaPortfolioData(request typing.SocialMediaSection) (typing.SocialMediaSection, error) {
	v := valtra.NewCollector()

	social_media := typing.SocialMediaSection{
		SocialMedia: valtra.Val(request.SocialMedia, "Social media").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Social media is required."),
				valtra.MinLengthString(3, "Social media must be at least 3 characters"),
				valtra.MaxLengthString(100, "Social media must be at most 100 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.SocialMediaSection{}, v.Errors()[0]
	}

	return social_media, nil
}

func ValidateJobExperiencePortfolioData(request typing.JobExperienceSection) (typing.JobExperienceSection, error) {
	v := valtra.NewCollector()

	job_experience := typing.JobExperienceSection{
		CompanyName: valtra.Val(request.CompanyName, "Company name").
			Transform(valtra.TrimSpace(), valtra.Uppercase()).
			Validate(
				valtra.Required[string]("Company name is required."),
				valtra.MinLengthString(3, "Company name must be at least 3 characters"),
				valtra.MaxLengthString(100, "Company name must be at most 100 characters"),
			).
			Collect(v),

		JobTitle: valtra.Val(request.JobTitle, "Job title").
			Transform(valtra.TrimSpace(), valtra.Uppercase()).
			Validate(
				valtra.Required[string]("Job title is required."),
				valtra.MinLengthString(3, "Job title must be at least 3 characters"),
				valtra.MaxLengthString(100, "Job title must be at most 100 characters"),
			).
			Collect(v),

		JobDescription: valtra.Val(request.JobDescription, "Job description").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Job description is required."),
				valtra.MinLengthString(3, "Job description must be at least 3 characters"),
				valtra.MaxLengthString(1000, "Job description must be at most 1000 characters"),
			).
			Collect(v),

		StartDate: valtra.Val(request.StartDate, "Start date").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Start date is required."),
				valtra.MinLengthString(3, "Start date must be at least 3 characters"),
				valtra.MaxLengthString(50, "Start date must be at most 50 characters"),
			).
			Collect(v),

		EndDate: valtra.Val(request.EndDate, "End date").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("End date is required."),
				valtra.MinLengthString(3, "End date must be at least 3 characters"),
				valtra.MaxLengthString(50, "End date must be at most 50 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.JobExperienceSection{}, v.Errors()[0]
	}

	return job_experience, nil
}

func ValidateEducationPortfolioData(request typing.EducationSection) (typing.EducationSection, error) {
	v := valtra.NewCollector()

	education := typing.EducationSection{
		InstitutionName: valtra.Val(request.InstitutionName, "Institution name").
			Transform(valtra.TrimSpace(), valtra.Uppercase()).
			Validate(
				valtra.Required[string]("Institution name is required."),
				valtra.MinLengthString(3, "Institution name must be at least 3 characters"),
				valtra.MaxLengthString(100, "Institution name must be at most 100 characters"),
			).
			Collect(v),

		DegreeType: valtra.Val(request.DegreeType, "Degree type").
			Transform(valtra.TrimSpace(), valtra.Uppercase()).
			Validate(
				valtra.Required[string]("Degree type is required."),
				valtra.MinLengthString(3, "Degree type must be at least 3 characters"),
				valtra.MaxLengthString(100, "Degree type must be at most 100 characters"),
			).
			Collect(v),

		Score: valtra.Val(request.Score, "Score").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.MinLengthString(3, "Score must be at least 3 characters"),
				valtra.MaxLengthString(20, "Score must be at most 20 characters"),
			).
			Collect(v),

		StartDate: valtra.Val(request.StartDate, "Start date").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Start date is required."),
				valtra.MinLengthString(3, "Start date must be at least 3 characters"),
				valtra.MaxLengthString(50, "Start date must be at most 50 characters"),
			).
			Collect(v),

		EndDate: valtra.Val(request.EndDate, "End date").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("End date is required."),
				valtra.MinLengthString(3, "End date must be at least 3 characters"),
				valtra.MaxLengthString(50, "End date must be at most 50 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.EducationSection{}, v.Errors()[0]
	}

	return education, nil
}

func ValidateSkillPortfolioData(request typing.SkillSection) (typing.SkillSection, error) {
	v := valtra.NewCollector()

	skill := typing.SkillSection{
		Skill: valtra.Val(request.Skill, "Skill").
			Transform(valtra.TrimSpace(), valtra.Uppercase()).
			Validate(
				valtra.Required[string]("Skill is required."),
				valtra.MinLengthString(3, "Skill must be at least 3 characters"),
				valtra.MaxLengthString(30, "Skill must be at most 30 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.SkillSection{}, v.Errors()[0]
	}

	return skill, nil
}

func ValidateLanguagePortfolioData(request typing.LanguageSection) (typing.LanguageSection, error) {
	v := valtra.NewCollector()

	language := typing.LanguageSection{
		Language: valtra.Val(request.Language, "Language").
			Transform(valtra.TrimSpace(), valtra.Uppercase()).
			Validate(
				valtra.Required[string]("Language is required."),
				valtra.MinLengthString(3, "Language must be at least 3 characters"),
				valtra.MaxLengthString(30, "Language must be at most 30 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.LanguageSection{}, v.Errors()[0]
	}

	return language, nil
}

func ValidateCertificatePortfolioData(request typing.CertificateSection) (typing.CertificateSection, error) {
	v := valtra.NewCollector()

	certificate := typing.CertificateSection{
		Certificate: valtra.Val(request.Certificate, "Certificate").
			Transform(valtra.TrimSpace(), valtra.Uppercase()).
			Validate(
				valtra.Required[string]("Certificate is required."),
				valtra.MinLengthString(3, "Certificate must be at least 3 characters"),
				valtra.MaxLengthString(30, "Certificate must be at most 30 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.CertificateSection{}, v.Errors()[0]
	}

	return certificate, nil
}

func ValidateProjectPortfolioData(request typing.ProjectSection) (typing.ProjectSection, error) {
	v := valtra.NewCollector()

	project := typing.ProjectSection{
		ProjectTitle: valtra.Val(request.ProjectTitle, "Project title").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Project title is required."),
				valtra.MinLengthString(3, "Project title must be at least 3 characters"),
				valtra.MaxLengthString(30, "Project title must be at most 30 characters"),
			).
			Collect(v),

		ProjectDescription: valtra.Val(request.ProjectDescription, "Project description").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Project description is required."),
				valtra.MinLengthString(3, "Project description must be at least 3 characters"),
				valtra.MaxLengthString(1000, "Project description must be at most 1000 characters"),
			).
			Collect(v),

		StartDate: valtra.Val(request.StartDate, "Start date").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Start date is required."),
				valtra.MinLengthString(3, "Start date must be at least 3 characters"),
				valtra.MaxLengthString(50, "Start date must be at most 50 characters"),
			).
			Collect(v),

		EndDate: valtra.Val(request.EndDate, "End date").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("End date is required."),
				valtra.MinLengthString(3, "End date must be at least 3 characters"),
				valtra.MaxLengthString(50, "End date must be at most 50 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return typing.ProjectSection{}, v.Errors()[0]
	}

	return project, nil
}

func validateEmailMX(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]

	resolver := net.Resolver{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mx, err := resolver.LookupMX(ctx, domain)
	if err == nil && len(mx) > 0 {
		return true
	}

	ips, err := resolver.LookupIPAddr(ctx, domain)
	if err == nil && len(ips) > 0 {
		return true
	}

	return false
}
