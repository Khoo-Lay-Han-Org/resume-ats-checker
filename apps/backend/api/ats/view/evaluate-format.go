package ats_view

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	typing "resuming/api/ats/typing"
	util "resuming/api/ats/util"
)

func SectionExistenceCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_resume_sections := c.Get("resume_sections")
		if retrieved_resume_sections == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to extract parsed resume content."})
		}

		resume_sections := retrieved_resume_sections.(map[string][]string)

		score := 0
		section_errors := map[string]string{}

		if len(resume_sections["summary"]) != 0 {
			data_num := 0
			for i := range resume_sections["summary"] {
				i = i + 1
				if i >= 3 {
					break
				}
				if i > 1 {
					score += 2
					continue
				}
				score += 20
				data_num = i
			}
			if data_num < 1 {
				section_errors["summary"] = "Too little amount of summary."
			} else {
				section_errors["summary"] = ""
			}
		}
		if resume_sections["personal_info"] != nil {
			data_num := 0
			for i := range resume_sections["personal_info"] {
				i = i + 1
				if i >= 6 {
					break
				}
				if i > 2 {
					score += 2
					continue
				}
				score += 10
				data_num = i
			}
			if data_num < 1 {
				section_errors["personal_info"] = "Too little amount of personal info."
			} else {
				section_errors["personal_info"] = ""
			}
		}
		if resume_sections["skills"] != nil {
			data_num := 0
			for i := range resume_sections["skills"] {
				i = i + 1
				if i >= 10 {
					break
				}
				if i > 1 {
					score += 2
					continue
				}
				score += 20
				data_num = i
			}
			if data_num < 1 {
				section_errors["skills"] = "Too little amount of skills."
			} else {
				section_errors["skills"] = ""
			}
		}
		if resume_sections["experience"] != nil {
			data_num := 0
			for i := range resume_sections["experience"] {
				i = i + 1
				if i >= 10 {
					break
				}
				if i > 2 {
					score += 2
					continue
				}
				score += 10
				data_num = i
			}
			if data_num < 1 {
				section_errors["experience"] = "Too little amount of experience."
			} else {
				section_errors["experience"] = ""
			}
		}
		if resume_sections["education"] != nil {
			data_num := 0
			for i := range resume_sections["education"] {
				i = i + 1
				if i >= 5 {
					break
				}
				if i > 2 {
					score += 2
					continue
				}
				score += 10
				data_num = i
			}
			if data_num < 1 {
				section_errors["education"] = "Too little amount of education."
			} else {
				section_errors["education"] = ""
			}
		}

		if resume_sections["certificates"] != nil {
			for i := range resume_sections["certificates"] {
				i = i + 1
				if i >= 4 {
					break
				}
				if i > 2 {
					score += 1
					continue
				}
				score += 3
			}
		}
		if resume_sections["objective"] != nil {
			for i := range resume_sections["objective"] {
				i = i + 1
				if i >= 4 {
					break
				}
				if i > 2 {
					score += 1
					continue
				}
				score += 3
			}
		}

		if score > 100 {
			score = 100
		}

		c.Set("section_existence_score", score)
		return nil
	}
}

func FormattingCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_resume_sections := c.Get("resume_sections")
		if retrieved_resume_sections == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to extract parsed resume content."})
		}

		resume_sections, ok := retrieved_resume_sections.(map[string][]string)
		if !ok {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to process resume sections."})
		}

		line_count := util.CountLine(resume_sections)
		error_count := 0
		format_error := map[string][]typing.FormatCheckErrorStruct{
			"summary":       {},
			"skills":        {},
			"personal_info": {},
			"experience":    {},
			"education":     {},
		}

		// summary
		for _, item := range resume_sections["summary"] {
			if strings.Contains(item, "  ") {
				idx := strings.Index(item, "  ")

				message := "Double spaces detected."
				start := idx
				end := idx + 2

				format_error["summary"] = append(format_error["summary"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if found, position := util.CheckIfPrintable(item); !found {
				for _, item := range position {
					message := "Unusual text found."
					start := item
					end := item + 1

					format_error["summary"] = append(format_error["summary"], typing.FormatCheckErrorStruct{
						Message: message,
						Start:   start,
						End:     end,
					})

					error_count += 1
				}
			}
			if strings.Count(item, "  ") > 2 {
				idx := strings.Index(item, "  ")

				message := "Likely table content found."
				start := idx
				end := idx + 1

				format_error["summary"] = append(format_error["summary"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if strings.Contains(item, "\t") {
				idx := strings.Index(item, "\t")

				message := "Likely table content found."
				start := idx
				end := idx + 1

				format_error["summary"] = append(format_error["summary"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if len(item) < 50 || len(item) > 500 {
				message := ""
				if len(item) < 50 {
					message = "Sentence too short."
				} else {
					message = "Sentence too long."
				}
				start := 0
				end := len(item)

				format_error["summary"] = append(format_error["summary"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}

		}

		// personal information
		for _, item := range resume_sections["personal_info"] {
			if strings.Contains(item, "  ") {
				idx := strings.Index(item, "  ")

				message := "Double spaces detected."
				start := idx
				end := idx + 2

				format_error["personal_info"] = append(format_error["personal_info"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if found, position := util.CheckIfPrintable(item); !found {
				for _, item := range position {
					message := "Unusual text found."
					start := item
					end := item + 1

					format_error["personal_info"] = append(format_error["personal_info"], typing.FormatCheckErrorStruct{
						Message: message,
						Start:   start,
						End:     end,
					})

					error_count += 1
				}
			}
			if strings.Count(item, "  ") > 2 {
				idx := strings.Index(item, "  ")

				message := "Likely table content found."
				start := idx
				end := idx + 1

				format_error["personal_info"] = append(format_error["personal_info"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if strings.Contains(item, "\t") {
				idx := strings.Index(item, "\t")

				message := "Likely table content found."
				start := idx
				end := idx + 1

				format_error["personal_info"] = append(format_error["personal_info"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if len(item) < 50 || len(item) > 500 {
				message := ""
				if len(item) < 50 {
					message = "Sentence too short."
				} else {
					message = "Sentence too long."
				}
				start := 0
				end := len(item)

				format_error["personal_info"] = append(format_error["personal_info"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}

		}

		// skills
		for _, item := range resume_sections["skills"] {
			if strings.Contains(item, "  ") {
				idx := strings.Index(item, "  ")

				message := "Double spaces detected."
				start := idx
				end := idx + 2

				format_error["skills"] = append(format_error["skills"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if found, position := util.CheckIfPrintable(item); !found {
				for _, item := range position {
					message := "Unusual text found."
					start := item
					end := item + 1

					format_error["skills"] = append(format_error["skills"], typing.FormatCheckErrorStruct{
						Message: message,
						Start:   start,
						End:     end,
					})

					error_count += 1
				}
			}
			if strings.Count(item, "  ") > 2 {
				idx := strings.Index(item, "  ")

				message := "Likely table content found."
				start := idx
				end := idx + 1

				format_error["skills"] = append(format_error["skills"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if strings.Contains(item, "\t") {
				idx := strings.Index(item, "\t")

				message := "Likely table content found."
				start := idx
				end := idx + 1

				format_error["skills"] = append(format_error["skills"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if len(item) < 50 || len(item) > 500 {
				message := ""
				if len(item) < 50 {
					message = "Sentence too short."
				} else {
					message = "Sentence too long."
				}
				start := 0
				end := len(item)

				format_error["skills"] = append(format_error["skills"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
		}

		// experience
		for _, item := range resume_sections["experience"] {
			if strings.Contains(item, "  ") {
				idx := strings.Index(item, "  ")

				message := "Double spaces detected."
				start := idx
				end := idx + 2

				format_error["experience"] = append(format_error["experience"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if found, position := util.CheckIfPrintable(item); !found {
				for _, item := range position {
					message := "Unusual text found."
					start := item
					end := item + 1

					format_error["experience"] = append(format_error["experience"], typing.FormatCheckErrorStruct{
						Message: message,
						Start:   start,
						End:     end,
					})

					error_count += 1
				}
			}
			if strings.Count(item, "  ") > 2 {
				idx := strings.Index(item, "  ")

				message := "Likely table content found."
				start := idx
				end := idx + 1

				format_error["experience"] = append(format_error["experience"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if strings.Contains(item, "\t") {
				idx := strings.Index(item, "\t")

				message := "Likely table content found."
				start := idx
				end := idx + 1

				format_error["experience"] = append(format_error["experience"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if len(item) < 50 || len(item) > 500 {
				message := ""
				if len(item) < 50 {
					message = "Sentence too short."
				} else {
					message = "Sentence too long."
				}
				start := 0
				end := len(item)

				format_error["experience"] = append(format_error["experience"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}

		}

		// education
		for _, item := range resume_sections["education"] {
			if strings.Contains(item, "  ") {
				idx := strings.Index(item, "  ")

				message := "Double spaces detected."
				start := idx
				end := idx + 2

				format_error["education"] = append(format_error["education"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if found, position := util.CheckIfPrintable(item); !found {
				for _, item := range position {
					message := "Unusual text found."
					start := item
					end := item + 1

					format_error["education"] = append(format_error["education"], typing.FormatCheckErrorStruct{
						Message: message,
						Start:   start,
						End:     end,
					})

					error_count += 1
				}
			}
			if strings.Count(item, "  ") > 2 {
				idx := strings.Index(item, "  ")

				message := "Likely table content found."
				start := idx
				end := idx + 1

				format_error["education"] = append(format_error["education"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if strings.Contains(item, "\t") {
				idx := strings.Index(item, "\t")

				message := "Likely table content found."
				start := idx
				end := idx + 1

				format_error["education"] = append(format_error["education"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}
			if len(item) < 50 || len(item) > 500 {
				message := ""
				if len(item) < 50 {
					message = "Sentence too short."
				} else {
					message = "Sentence too long."
				}
				start := 0
				end := len(item)

				format_error["education"] = append(format_error["education"], typing.FormatCheckErrorStruct{
					Message: message,
					Start:   start,
					End:     end,
				})

				error_count += 1
			}

		}

		var total_score int
		if line_count > 0 {
			total_score = (error_count * 100) / line_count
		}

		c.Set("format_score", total_score)
		return nil
	}
}
