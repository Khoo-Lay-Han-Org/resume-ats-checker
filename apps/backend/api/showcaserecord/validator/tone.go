package showcaserecord_validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	typing "resuming/api/showcaserecord/typing"
	systemconfig "resuming/system-config"
)

func CheckTone[T any](request T) error {
	val := reflect.ValueOf(request)

	for _, field := range val.Fields() {
		field_value := field.Interface()
		translation_payload := map[string]string{
			"phrase": field_value.(string),
		}
		translation_json_data, _ := json.Marshal(translation_payload)
		resp1, err := http.Post(
			systemconfig.AiModelsUri+"translation_model_predict",
			"application/json",
			bytes.NewBuffer(translation_json_data),
		)
		if err != nil {
			return errors.New("failed to validate data")
		}
		defer func() { _ = resp1.Body.Close() }()

		var translation_response typing.TranslationModelResponse
		if err := json.NewDecoder(resp1.Body).Decode(&translation_response); err != nil {
			return errors.New("failed to validate data")
		}

		tone_payload := map[string]string{
			"phrase": translation_response.Prediction,
		}

		tone_json_data, _ := json.Marshal(tone_payload)

		resp2, err := http.Post(
			systemconfig.AiModelsUri+"tone-detection-model-predict",
			"application/json",
			bytes.NewBuffer(tone_json_data),
		)
		if err != nil {
			return errors.New("failed to validate data")
		}
		defer func() { _ = resp2.Body.Close() }()

		var tone_response typing.ToneDetectionModelResponse
		if err := json.NewDecoder(resp2.Body).Decode(&tone_response); err != nil {
			return errors.New("failed to validate data")
		}

		if tone_response.Prediction[0][0] >= 0.07 {
			return errors.New("data provided by user has a horrible tone")
		}
	}

	return nil
}
