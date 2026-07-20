package ats_view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"resuming/ai"
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

		body, err := ai.SkillsKeywordPredict(resume_content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve request data."})
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

		body, err := ai.SkillsKeywordPredict(job_desc_str)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve request data."})
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

				score, err := ai.TextSimilarityPredict(skill_name1, skill_name2)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to calculate score."})
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
