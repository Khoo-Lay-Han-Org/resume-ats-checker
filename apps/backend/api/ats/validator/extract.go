package ats_validator

import "github.com/bobch27/valtra-go"

func ValidatePDFFile() {}

func ValidateJobDesc(job_desc string) (string, error) {
	v := valtra.NewCollector()

	polished_job_desc := valtra.Val(job_desc, "Job description").
		Transform(valtra.TrimSpace()).
		Validate(valtra.Required[string]("Job description is required")).
		Collect(v)

	if !v.IsValid() {
		return job_desc, v.Errors()[0]
	}

	return polished_job_desc, nil
}
