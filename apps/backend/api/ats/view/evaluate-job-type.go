package ats_view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"resuming/ai"
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

		body, err := ai.JobTypeModelPredict(resume_content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve request data."})
		}

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

		response, err := ai.JobTypeModelPredict(job_desc_str)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve request data."})
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

		score, err := ai.TextSimilarityPredict(job_desc_job_type_str, resume_job_type_str)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to calculate score."})
		}

		c.Set("job_type_score", int(score*100))
		return nil
	}
}
