package ats_view

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	systemconfig "resuming/system-config"
)

// Relevance
func ResumeJobTypeCheck() gin.HandlerFunc {
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

		req, err := http.NewRequest("POST", systemconfig.AiModelsUri+"job_type_model_predict", bytes.NewBuffer(ai_request_content_bytes))
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

		retrieved_body, _ := io.ReadAll(resp.Body)
		body := string(retrieved_body)

		c.Set("resume_job_type", body)
		c.Next()
	}
}

func JobDescJobTypeCheck() gin.HandlerFunc {
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

		json_ai_request_content, err := json.Marshal(ai_request_content)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
			return
		}

		resp, err := http.Post(
			systemconfig.AiModelsUri+"job-type-model-predict",
			"application/json",
			bytes.NewBuffer(json_ai_request_content),
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve request data."})
			return
		}
		defer func() { _ = resp.Body.Close() }()

		var response string
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to decode response."})
			return
		}

		c.Set("job_desc_job_type", response)
		c.Next()
	}
}

func JobTypeRelevanceCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		job_desc_job_type, exists := c.Get("job_desc_job_type")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve data."})
			return
		}

		resume_job_type, exists := c.Get("resume_job_type")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve data."})
			return
		}

		job_desc_job_type_str, ok := job_desc_job_type.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process job description job type."})
			return
		}

		resume_job_type_str, ok := resume_job_type.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process resume job type."})
			return
		}

		ai_similarity_check_request := map[string]string{
			"text1": job_desc_job_type_str,
			"text2": resume_job_type_str,
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

		c.Set("job_type_score", score)
		c.Next()
	}
}
