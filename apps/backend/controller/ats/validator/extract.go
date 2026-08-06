package ats_validator

import (
	"errors"
	"strings"
)

func ValidatePDFFile() {}

func ValidateJobDesc(job_desc string) (string, error) {
	polished_job_desc := strings.TrimSpace(job_desc)
	if polished_job_desc == "" {
		return job_desc, errors.New("Job description is required")
	}
	return polished_job_desc, nil
}
