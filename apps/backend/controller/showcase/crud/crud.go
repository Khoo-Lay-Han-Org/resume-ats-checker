package showcase_crud

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	dto "resuming/controller/showcase/dto"
	validator "resuming/controller/showcase/validator"
	"resuming/service"
	shared_find "resuming/shared/find"
	"resuming/systemconfig"
)

func InsertShowCaseRecordData(request any, public_user_id string) error {
	data, err := shared_find.GetShowcaseRecordData(public_user_id)
	if err != nil {
		return err
	}
	ctx := context.Background()

	json_data := []byte(data)
	var deserialised_data map[string]any
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return err
	}

	field_extract := reflect.TypeOf(request)
	value_extract := reflect.ValueOf(request)

	for i := 0; i < field_extract.NumField(); i++ {
		field_name := field_extract.Field(i).Name
		field_value := value_extract.Field(i).Interface()

		existing_values := deserialised_data[field_name]
		switch field_name {
		case "name", "email", "phone_number", "address", "social_media", "skill", "language":
			typed_new_value := field_value.([]string)
			typed_old_value := existing_values.([]string)
			deserialised_data[field_name] = append(typed_old_value, typed_new_value...)
		case "job_experience", "education", "certificate", "project":
			typed_new_value := field_value.([]byte)
			typed_old_value := existing_values.([]byte)
			deserialised_data[field_name] = append(typed_old_value, typed_new_value...)
		}
	}

	serialised_showcaserecord_data, err := json.Marshal(deserialised_data)
	if err != nil {
		return err
	}

	err = service.Valkey.Do(
		ctx,
		service.Valkey.B().Set().
			Key(public_user_id+":showcaserecord_data").Value(string(serialised_showcaserecord_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return err
	}

	return nil
}

func EditShowCaseRecordData[T any](request T, index int, public_user_id string) error {
	data, err := shared_find.GetShowcaseRecordData(public_user_id)
	if err != nil {
		return err
	}
	ctx := context.Background()

	json_data := []byte(data)
	var deserialised_data map[string]any
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return err
	}

	field_extract := reflect.TypeOf(request)
	value_extract := reflect.ValueOf(request)

	for i := 0; i < field_extract.NumField(); i++ {
		field_name := field_extract.Field(i).Name
		field_value := value_extract.Field(i).Interface()

		switch field_name {
		case "name", "email", "phone_number", "address", "social_media", "skill", "language":
			typed_new_value := field_value.(string)
			field_slice := deserialised_data[field_name].([]string)
			field_slice[index] = typed_new_value
			deserialised_data[field_name] = field_slice
		case "job_experience", "education", "certificate", "project":
			typed_new_value := field_value.([]byte)
			field_slice := deserialised_data[field_name].([]any)
			field_slice[index] = typed_new_value
			deserialised_data[field_name] = field_slice
		}
	}

	serialised_showcaserecord_data, err := json.Marshal(deserialised_data)
	if err != nil {
		return err
	}

	err = service.Valkey.Do(
		ctx,
		service.Valkey.B().Set().
			Key(public_user_id+":showcaserecord_data").Value(string(serialised_showcaserecord_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return err
	}

	return nil
}

func DeleteShowCaseRecordData(field_name string, index int, public_user_id string) error {
	data, err := shared_find.GetShowcaseRecordData(public_user_id)
	if err != nil {
		return err
	}
	ctx := context.Background()

	json_data := []byte(data)
	var deserialised_data map[string]any
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return err
	}

	field_slice := deserialised_data[field_name].([]any)
	field_slice = append(field_slice[:index], field_slice[index+1:]...)
	deserialised_data[field_name] = field_slice

	serialised_showcaserecord_data, err := json.Marshal(deserialised_data)
	if err != nil {
		return err
	}

	err = service.Valkey.Do(
		ctx,
		service.Valkey.B().Set().
			Key(public_user_id+":showcaserecord_data").Value(string(serialised_showcaserecord_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return err
	}

	return nil
}

func ValidateData[T dto.AddressSection | dto.CertificateSection | dto.EducationSection | dto.EmailSection | dto.JobExperienceSection | dto.LanguageSection | dto.NameSection | dto.PhoneNumberSection | dto.ProjectSection | dto.SkillSection | dto.SocialMediaSection](request any) (T, error) {
	switch request.(type) {
	case dto.NameSection:
		typed_request, ok := request.(dto.NameSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidateNamePortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	case dto.EmailSection:
		typed_request, ok := request.(dto.EmailSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidateEmailPortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	case dto.PhoneNumberSection:
		typed_request, ok := request.(dto.PhoneNumberSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidatePhoneNumberPortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	case dto.AddressSection:
		typed_request, ok := request.(dto.AddressSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidateAddressPortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	case dto.SocialMediaSection:
		typed_request, ok := request.(dto.SocialMediaSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidateSocialMediaPortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	case dto.JobExperienceSection:
		typed_request, ok := request.(dto.JobExperienceSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidateJobExperiencePortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	case dto.EducationSection:
		typed_request, ok := request.(dto.EducationSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidateEducationPortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	case dto.SkillSection:
		typed_request, ok := request.(dto.SkillSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidateSkillPortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	case dto.CertificateSection:
		typed_request, ok := request.(dto.CertificateSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidateCertificatePortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	case dto.LanguageSection:
		typed_request, ok := request.(dto.LanguageSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidateLanguagePortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	case dto.ProjectSection:
		typed_request, ok := request.(dto.ProjectSection)
		if !ok {
			var zero T
			return zero, errors.New("failed to process data")
		}

		validated_request, err := validator.ValidateProjectPortfolioData(typed_request)
		if err != nil {
			var zero T
			return zero, err
		}

		return any(validated_request).(T), nil
	}

	var zero T
	return zero, errors.New("failed to process data")
}
