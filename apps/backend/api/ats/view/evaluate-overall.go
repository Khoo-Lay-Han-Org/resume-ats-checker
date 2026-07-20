package ats_view

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func OverallScore() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_section_existence_score := c.Get("section_existence_score")
		if retrieved_section_existence_score == nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Score not found."})
		}

		retrieved_format_score := c.Get("format_score")
		if retrieved_format_score == nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Score not found."})
		}

		retrieved_job_type_score := c.Get("job_type_score")
		if retrieved_job_type_score == nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Score not found."})
		}

		retrieved_skills_score := c.Get("skills_score")
		if retrieved_skills_score == nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Score not found."})
		}

		section_existence_score, ok := retrieved_section_existence_score.(int)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process section existence score."})
		}

		format_score, ok := retrieved_format_score.(int)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process format score."})
		}

		job_type_score, ok := retrieved_job_type_score.(int)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process job type score."})
		}

		skills_score, ok := retrieved_skills_score.(int)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process skills score."})
		}

		overall_score := (section_existence_score + format_score + job_type_score + skills_score) / 4

		c.Set("response_data", overall_score)
		return nil
	}
}
