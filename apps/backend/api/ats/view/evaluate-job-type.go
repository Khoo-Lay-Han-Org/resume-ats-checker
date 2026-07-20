package ats_view

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	systemconfig "resuming/system-config"
)

// Relevance
func ResumeJobTypeCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_resume_content := c.Get("resume_content")
		if retrieved_resume_content == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve resume sections."})
		}

		resume_content, ok := retrieved_resume_content.(string)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process resume content."})
		}

		ai_request_content := map[string]string{
			"text": resume_content,
		}

		ai_request_content_bytes, err := json.Marshal(ai_request_content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to prepare request data."})
		}

		req, err := http.NewRequest("POST", systemconfig.AiModelsUri+"job-type-model-predict", bytes.NewBuffer(ai_request_content_bytes))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to prepare request data."})
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve request data."})
		}
		defer func() { _ = resp.Body.Close() }()

		retrieved_body, _ := io.ReadAll(resp.Body)
		body := string(retrieved_body)

		c.Set("resume_job_type", body)
		return nil
	}
}

func JobDescJobTypeCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		job_desc := c.Get("job_desc")
		if job_desc == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve job description."})
		}

		job_desc_str, ok := job_desc.(string)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process job description."})
		}

		ai_request_content := map[string]string{
			"text": job_desc_str,
		}

		json_ai_request_content, err := json.Marshal(ai_request_content)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		resp, err := http.Post(
			systemconfig.AiModelsUri+"job-type-model-predict",
			"application/json",
			bytes.NewBuffer(json_ai_request_content),
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve request data."})
		}
		defer func() { _ = resp.Body.Close() }()

		var response string
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to decode response."})
		}

		c.Set("job_desc_job_type", response)
		return nil
	}
}

func JobTypeRelevanceCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		job_desc_job_type := c.Get("job_desc_job_type")
		if job_desc_job_type == nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve data."})
		}

		resume_job_type := c.Get("resume_job_type")
		if resume_job_type == nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve data."})
		}

		job_desc_job_type_str, ok := job_desc_job_type.(string)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process job description job type."})
		}

		resume_job_type_str, ok := resume_job_type.(string)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process resume job type."})
		}

		ai_similarity_check_request := map[string]string{
			"text1": job_desc_job_type_str,
			"text2": resume_job_type_str,
		}

		json_ai_similarity_check_request, err := json.Marshal(ai_similarity_check_request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to calculate score."})
		}

		resp, err := http.Post(
			systemconfig.AiModelsUri+"text-similarity-predict",
			"application/json",
			bytes.NewBuffer(json_ai_similarity_check_request),
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to calculate score."})
		}
		defer func() { _ = resp.Body.Close() }()

		var score float64
		if err := json.NewDecoder(resp.Body).Decode(&score); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to decode response."})
		}

		c.Set("job_type_score", int(score*100))
		return nil
	}
}
