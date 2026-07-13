package ats_view

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
	util "resuming/api/ats/util"
	ats_validator "resuming/api/ats/validator"
	systemconfig "resuming/system-config"
)

func ExtractResume() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to retrieve session data."})
			return
		}

		public_user_id, ok := retrieved_public_user_id.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to process session."})
			return
		}

		resume, err := c.FormFile("resume_file")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "File required."})
			return
		}

		//// using public session id as the name of the file
		//// in Linux, os.TempDir() points to the /tmp folder at root
		//// we shall setup cloud storage in the future
		temporary_file_path := filepath.Join(os.TempDir(), public_user_id)

		if err := c.SaveUploadedFile(resume, temporary_file_path); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file."})
			return
		}
		defer func() { _ = os.Remove(temporary_file_path) }()

		f, r, err := pdf.Open(temporary_file_path)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process resume."})
			return
		}
		defer func() { _ = f.Close() }()

		var sb strings.Builder
		for i := 1; i <= r.NumPage(); i++ {
			page := r.Page(i)
			if page.V.IsNull() {
				continue
			}
			text, err := page.GetPlainText(nil)
			if err != nil {
				continue
			}
			sb.WriteString(text)
		}

		resume_content := sb.String()

		c.Set("resume_content", resume_content)
		c.Next()
	}
}

func ParseResume() gin.HandlerFunc {
	return func(c *gin.Context) {
		resume_content, exists := c.Get("resume_content")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve resume data."})
			return
		}

		serialised_resume_content, err := json.Marshal(resume_content)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to serialise resume data."})
			return
		}

		resp, err := http.Post(
			systemconfig.AiModelsUri+"resume-sections-model-predict",
			"application/json",
			bytes.NewBuffer(serialised_resume_content),
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve resume data."})
			return
		}
		defer func() { _ = resp.Body.Close() }()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to read response body."})
			return
		}

		var resume_sections map[string][]string
		if err := json.Unmarshal(body, &resume_sections); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse response body."})
			return
		}

		c.Set("resume_sections", resume_sections)
		c.Next()
	}
}

func UserInputJobDesc() gin.HandlerFunc {
	return func(c *gin.Context) {
		job_desc := c.PostForm("job_desc")

		polished_job_desc, err := ats_validator.ValidateJobDesc(job_desc)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.Set("job_desc", polished_job_desc)
		c.Next()
	}
}

func WebScrapeJobDesc() gin.HandlerFunc {
	return func(c *gin.Context) {
		company := c.PostForm("company")
		job_title := c.PostForm("job_title")

		content := util.JobDescWebScrape(company, job_title)

		polished_job_desc, err := ats_validator.ValidateJobDesc(content)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.Set("job_desc", polished_job_desc)
		c.Next()
	}
}
