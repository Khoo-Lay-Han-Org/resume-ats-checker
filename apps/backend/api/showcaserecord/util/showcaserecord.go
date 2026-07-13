package showcaserecord_util

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/lib/pq"
	valkey "github.com/valkey-io/valkey-go"
	"gorm.io/datatypes"
	typing "resuming/api/showcaserecord/typing"
	validator "resuming/api/showcaserecord/validator"
	"resuming/database"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func InsertShowCaseRecordData(request any, public_user_id string) error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			user, dbErr := database.FindUserByPublicId(public_user_id)
			if dbErr != nil {
				return dbErr
			}
			if syncErr := database.SyncIndividualShowCaseRecordDataSessionStore(public_user_id, user); syncErr != nil {
				return syncErr
			}
			data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

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
			typed_new_value := field_value.(pq.StringArray)
			typed_old_value := existing_values.(pq.StringArray)
			deserialised_data[field_name] = append(typed_old_value, typed_new_value...)
		case "job_experience", "education", "certificate", "project":
			typed_new_value := field_value.(datatypes.JSON)
			typed_old_value := existing_values.(datatypes.JSON)
			deserialised_data[field_name] = append(typed_old_value, typed_new_value...)
		}
	}

	serialised_showcaserecord_data, err := json.Marshal(deserialised_data)
	if err != nil {
		return err
	}

	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
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
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			user, dbErr := database.FindUserByPublicId(public_user_id)
			if dbErr != nil {
				return dbErr
			}
			if syncErr := database.SyncIndividualShowCaseRecordDataSessionStore(public_user_id, user); syncErr != nil {
				return syncErr
			}
			data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

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
			field_slice := deserialised_data[field_name].(pq.StringArray)
			field_slice[index] = typed_new_value
			deserialised_data[field_name] = field_slice
		case "job_experience", "education", "certificate", "project":
			typed_new_value := field_value.(datatypes.JSON)
			field_slice := deserialised_data[field_name].([]any)
			field_slice[index] = typed_new_value
			deserialised_data[field_name] = field_slice
		}
	}

	serialised_showcaserecord_data, err := json.Marshal(deserialised_data)
	if err != nil {
		return err
	}

	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
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
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			user, dbErr := database.FindUserByPublicId(public_user_id)
			if dbErr != nil {
				return dbErr
			}
			if syncErr := database.SyncIndividualShowCaseRecordDataSessionStore(public_user_id, user); syncErr != nil {
				return syncErr
			}
			data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

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

	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":showcaserecord_data").Value(string(serialised_showcaserecord_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return err
	}

	return nil
}

func ValidateData[T typing.AddressSection | typing.CertificateSection | typing.EducationSection | typing.EmailSection | typing.JobExperienceSection | typing.LanguageSection | typing.NameSection | typing.PhoneNumberSection | typing.ProjectSection | typing.SkillSection | typing.SocialMediaSection](request any) (T, error) {
	switch request.(type) {
	case typing.NameSection:
		typed_request, ok := request.(typing.NameSection)
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
	case typing.EmailSection:
		typed_request, ok := request.(typing.EmailSection)
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
	case typing.PhoneNumberSection:
		typed_request, ok := request.(typing.PhoneNumberSection)
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
	case typing.AddressSection:
		typed_request, ok := request.(typing.AddressSection)
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
	case typing.SocialMediaSection:
		typed_request, ok := request.(typing.SocialMediaSection)
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
	case typing.JobExperienceSection:
		typed_request, ok := request.(typing.JobExperienceSection)
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
	case typing.EducationSection:
		typed_request, ok := request.(typing.EducationSection)
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
	case typing.SkillSection:
		typed_request, ok := request.(typing.SkillSection)
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
	case typing.CertificateSection:
		typed_request, ok := request.(typing.CertificateSection)
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
	case typing.LanguageSection:
		typed_request, ok := request.(typing.LanguageSection)
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
	case typing.ProjectSection:
		typed_request, ok := request.(typing.ProjectSection)
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
