package ats_view

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ledongthuc/pdf"
	util "resuming/api/ats/util"
	ats_validator "resuming/api/ats/validator"
	systemconfig "resuming/system-config"
)

func ExtractResume() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve session data."})
		}

		public_user_id, ok := retrieved_public_user_id.(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to process session."})
		}

		resume, err := c.FormFile("resume_file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "File required."})
		}

		temporary_file_path := filepath.Join(os.TempDir(), public_user_id)

		src, err := resume.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to open file."})
		}
		defer src.Close()

		dst, err := os.Create(temporary_file_path)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to save file."})
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to save file."})
		}
		defer func() { _ = os.Remove(temporary_file_path) }()

		f, r, err := pdf.Open(temporary_file_path)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process resume."})
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
		return nil
	}
}

func ParseResume() echo.HandlerFunc {
	return func(c echo.Context) error {
		resume_content := c.Get("resume_content")
		if resume_content == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve resume data."})
		}

		serialised_resume_content, err := json.Marshal(resume_content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to serialise resume data."})
		}

		resp, err := http.Post(
			systemconfig.AiModelsUri+"resume-sections-model-predict",
			"application/json",
			bytes.NewBuffer(serialised_resume_content),
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve resume data."})
		}
		defer func() { _ = resp.Body.Close() }()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to read response body."})
		}

		var resume_sections map[string][]string
		if err := json.Unmarshal(body, &resume_sections); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse response body."})
		}

		c.Set("resume_sections", resume_sections)
		return nil
	}
}

func UserInputJobDesc() echo.HandlerFunc {
	return func(c echo.Context) error {
		job_desc := c.FormValue("job_desc")

		polished_job_desc, err := ats_validator.ValidateJobDesc(job_desc)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		c.Set("job_desc", polished_job_desc)
		return nil
	}
}

func WebScrapeJobDesc() echo.HandlerFunc {
	return func(c echo.Context) error {
		company := c.FormValue("company")
		job_title := c.FormValue("job_title")

		content := util.JobDescWebScrape(company, job_title)

		polished_job_desc, err := ats_validator.ValidateJobDesc(content)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		c.Set("job_desc", polished_job_desc)
		return nil
	}
}
