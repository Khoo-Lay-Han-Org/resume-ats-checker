package ats_view

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OverallScore() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_section_existence_score, exists := c.Get("section_existence_score")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Score not found."})
			return
		}

		retrieved_format_score, exists := c.Get("format_score")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Score not found."})
			return
		}

		retrieved_job_type_score, exists := c.Get("job_type_score")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Score not found."})
			return
		}

		retrieved_skills_score, exists := c.Get("skills_score")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Score not found."})
			return
		}

		section_existence_score, ok := retrieved_section_existence_score.(int)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process section existence score."})
			return
		}

		format_score, ok := retrieved_format_score.(int)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process format score."})
			return
		}

		job_type_score, ok := retrieved_job_type_score.(int)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process job type score."})
			return
		}

		skills_score, ok := retrieved_skills_score.(int)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process skills score."})
			return
		}

		overall_score := (section_existence_score + format_score + job_type_score + skills_score) / 4

		c.Set("response_data", overall_score)
	}
}
