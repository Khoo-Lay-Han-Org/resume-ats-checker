package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func callPython(modelName string, input any) ([]byte, error) {
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	cmd := exec.Command(".venv/bin/python3", "-m", "ai.predict", modelName)
	cmd.Stdin = bytes.NewReader(inputJSON)

	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	result, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("python process failed: %w\nstderr: %s", err, stderr.String())
	}

	return result, nil
}

func JobTypeModelPredict(text string) (string, error) {
	input := map[string]string{"text": text}
	result, err := callPython("job-type-model-predict", input)
	if err != nil {
		return "", err
	}

	var response string
	if err := json.Unmarshal(result, &response); err != nil {
		return "", err
	}

	return response, nil
}

func ResumeSectionModelPredict(text string) (map[string][]string, error) {
	input := map[string]string{"text": text}
	result, err := callPython("resume-sections-model-predict", input)
	if err != nil {
		return nil, err
	}

	var response map[string][]string
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response, nil
}

type ToneDetectionResponse struct {
	Prediction [][]float64 `json:"prediction"`
}

func ToneDetectionModelPredict(phrase string) (*ToneDetectionResponse, error) {
	input := map[string]string{"phrase": phrase}
	result, err := callPython("tone-detection-model-predict", input)
	if err != nil {
		return nil, err
	}

	var response ToneDetectionResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type TranslationResponse struct {
	Prediction string `json:"prediction"`
}

func TranslationModelPredict(text string, tgtLang string) (*TranslationResponse, error) {
	if tgtLang == "" {
		tgtLang = "eng_Latn"
	}
	input := map[string]string{"text": text, "tgt_lang": tgtLang}
	result, err := callPython("translation-model-predict", input)
	if err != nil {
		return nil, err
	}

	var response TranslationResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func SkillsKeywordPredict(text string) ([]map[string]any, error) {
	input := map[string]string{"text": text}
	result, err := callPython("skills-keyword-predict", input)
	if err != nil {
		return nil, err
	}

	var response []map[string]any
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func TextSimilarityPredict(text1, text2 string) (float64, error) {
	input := map[string]string{"text1": text1, "text2": text2}
	result, err := callPython("text-similarity-predict", input)
	if err != nil {
		return 0, err
	}

	var response float64
	if err := json.Unmarshal(result, &response); err != nil {
		return 0, err
	}

	return response, nil
}
