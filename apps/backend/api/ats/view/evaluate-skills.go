package ats_view

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	systemconfig "resuming/system-config"
)

func ResumeSkillsCheck() echo.HandlerFunc {
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

		req, err := http.NewRequest("POST", systemconfig.AiModelsUri+"skills-keyword-predict", bytes.NewBuffer(ai_request_content_bytes))
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

		retrieved_body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to read response body."})
		}

		var body []map[string]any
		err = json.Unmarshal(retrieved_body, &body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse JSON response."})
		}

		c.Set("resume_skills", body)
		return nil
	}
}

func JobDescSkillsCheck() echo.HandlerFunc {
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

		ai_request_content_bytes, err := json.Marshal(ai_request_content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to prepare request data."})
		}

		req, err := http.NewRequest("POST", systemconfig.AiModelsUri+"skills-keyword-predict", bytes.NewBuffer(ai_request_content_bytes))
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

		retrieved_body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to read response body."})
		}

		var body []map[string]any
		err = json.Unmarshal(retrieved_body, &body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse JSON response."})
		}

		c.Set("job_desc_skills", body)
		return nil
	}
}

func OverallSkillsCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_resume_skills := c.Get("resume_skills")
		if retrieved_resume_skills == nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve resume skills."})
		}

		resume_skills, ok := retrieved_resume_skills.([]map[string]any)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process resume skills."})
		}

		retrieved_job_desc_skills := c.Get("job_desc_skills")
		if retrieved_job_desc_skills == nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve job description skills."})
		}

		job_desc_skills, ok := retrieved_job_desc_skills.([]map[string]any)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process job description skills."})
		}

		total_requirements := len(job_desc_skills)
		total_matched := 0

		for _, item1 := range job_desc_skills {
			skill_name1, ok := item1["word"].(string)
			if !ok {
				continue
			}
			for _, item2 := range resume_skills {
				skill_name2, ok := item2["word"].(string)
				if !ok {
					continue
				}

				ai_similarity_check_request := map[string]string{
					"text1": skill_name1,
					"text2": skill_name2,
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

				if score > 0.7 {
					total_matched += 1
				}
			}
		}

		if total_matched > total_requirements {
			total_matched = total_requirements
		}

		var finalScore int
		if total_requirements > 0 {
			finalScore = (total_matched * 100) / total_requirements
		}

		c.Set("skills_score", finalScore)
		return nil
	}
}
