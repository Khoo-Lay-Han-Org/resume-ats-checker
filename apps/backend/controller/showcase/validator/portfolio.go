package showcase_validator

import (
	dto "resuming/controller/showcase/dto"
	"resuming/shared/validation"
)

func ValidateNamePortfolioData(request dto.NameSection) (dto.NameSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.NameSection{}, err
	}
	return cleaned.(dto.NameSection), nil
}

func ValidateEmailPortfolioData(request dto.EmailSection) (dto.EmailSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.EmailSection{}, err
	}
	return cleaned.(dto.EmailSection), nil
}

func ValidatePhoneNumberPortfolioData(request dto.PhoneNumberSection) (dto.PhoneNumberSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.PhoneNumberSection{}, err
	}
	return cleaned.(dto.PhoneNumberSection), nil
}

func ValidateAddressPortfolioData(request dto.AddressSection) (dto.AddressSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.AddressSection{}, err
	}
	return cleaned.(dto.AddressSection), nil
}

func ValidateSocialMediaPortfolioData(request dto.SocialMediaSection) (dto.SocialMediaSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.SocialMediaSection{}, err
	}
	return cleaned.(dto.SocialMediaSection), nil
}

func ValidateJobExperiencePortfolioData(request dto.JobExperienceSection) (dto.JobExperienceSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.JobExperienceSection{}, err
	}
	return cleaned.(dto.JobExperienceSection), nil
}

func ValidateEducationPortfolioData(request dto.EducationSection) (dto.EducationSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.EducationSection{}, err
	}
	return cleaned.(dto.EducationSection), nil
}

func ValidateSkillPortfolioData(request dto.SkillSection) (dto.SkillSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.SkillSection{}, err
	}
	return cleaned.(dto.SkillSection), nil
}

func ValidateLanguagePortfolioData(request dto.LanguageSection) (dto.LanguageSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.LanguageSection{}, err
	}
	return cleaned.(dto.LanguageSection), nil
}

func ValidateCertificatePortfolioData(request dto.CertificateSection) (dto.CertificateSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.CertificateSection{}, err
	}
	return cleaned.(dto.CertificateSection), nil
}

func ValidateProjectPortfolioData(request dto.ProjectSection) (dto.ProjectSection, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.ProjectSection{}, err
	}
	return cleaned.(dto.ProjectSection), nil
}
