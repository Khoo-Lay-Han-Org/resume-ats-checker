package ats_view

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	systemconfig "resuming/system-config"
)

func ResumeSkillsCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_resume_content, exists := c.Get("resume_content")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve resume sections."})
			return
		}

		resume_content, ok := retrieved_resume_content.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process resume content."})
			return
		}

		ai_request_content := map[string]string{
			"text": resume_content,
		}

		ai_request_content_bytes, err := json.Marshal(ai_request_content)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to prepare request data."})
			return
		}

		req, err := http.NewRequest("POST", systemconfig.AiModelsUri+"skills-keyword-predict", bytes.NewBuffer(ai_request_content_bytes))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to prepare request data."})
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve request data."})
			return
		}
		defer func() { _ = resp.Body.Close() }()

		retrieved_body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to read response body."})
			return
		}

		var body []map[string]any
		err = json.Unmarshal(retrieved_body, &body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse JSON response."})
			return
		}

		c.Set("resume_skills", body)
		c.Next()
	}
}

func JobDescSkillsCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		job_desc, exists := c.Get("job_desc")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve job description."})
			return
		}

		job_desc_str, ok := job_desc.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process job description."})
			return
		}

		ai_request_content := map[string]string{
			"text": job_desc_str,
		}

		ai_request_content_bytes, err := json.Marshal(ai_request_content)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to prepare request data."})
			return
		}

		req, err := http.NewRequest("POST", systemconfig.AiModelsUri+"skills-keyword-predict", bytes.NewBuffer(ai_request_content_bytes))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to prepare request data."})
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve request data."})
			return
		}
		defer func() { _ = resp.Body.Close() }()

		retrieved_body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to read response body."})
			return
		}

		var body []map[string]any
		err = json.Unmarshal(retrieved_body, &body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse JSON response."})
			return
		}

		c.Set("job_desc_skills", body)
		c.Next()
	}
}

func OverallSkillsCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_resume_skills, exists := c.Get("resume_skills")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve resume skills."})
			return
		}

		resume_skills, ok := retrieved_resume_skills.([]map[string]any)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process resume skills."})
			return
		}

		retrieved_job_desc_skills, exists := c.Get("job_desc_skills")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve job description skills."})
			return
		}

		job_desc_skills, ok := retrieved_job_desc_skills.([]map[string]any)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process job description skills."})
			return
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
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to calculate score."})
					return
				}

				resp, err := http.Post(
					systemconfig.AiModelsUri+"text-similarity-predict",
					"application/json",
					bytes.NewBuffer(json_ai_similarity_check_request),
				)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to calculate score."})
					return
				}
				defer func() { _ = resp.Body.Close() }()

				var score int
				if err := json.NewDecoder(resp.Body).Decode(&score); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to decode response."})
					return
				}

				if score > 70 {
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
		c.Next()
	}
}
