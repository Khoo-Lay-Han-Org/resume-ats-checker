package showcaserecord_validator

import (
	"errors"
	"reflect"

	"resuming/ai"
)

func CheckTone[T any](request T) error {
	val := reflect.ValueOf(request)

	for _, field := range val.Fields() {
		field_value := field.Interface()

		translation_response, err := ai.TranslationModelPredict(field_value.(string), "eng_Latn")
		if err != nil {
			return errors.New("failed to validate data")
		}

		tone_response, err := ai.ToneDetectionModelPredict(translation_response.Prediction)
		if err != nil {
			return errors.New("failed to validate data")
		}

		if tone_response.Prediction[0][0] >= 0.07 {
			return errors.New("data provided by user has a horrible tone")
		}
	}

	return nil
}
