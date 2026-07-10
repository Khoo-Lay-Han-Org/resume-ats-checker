package ats_util

import "unicode"

func CheckIfPrintable(content string) (bool, []int) {
	unusual_index := []int{}

	for i, r := range content {
		if !unicode.IsPrint(r) {
			unusual_index = append(unusual_index, i+1)
		}
	}

	if len(unusual_index) > 0 {
		return false, unusual_index
	}

	return true, nil
}

func CountLine(resume_sections map[string][]string) int {
	line_count := 0

	for range resume_sections["summary"] {
		line_count += 1
	}
	for range resume_sections["skills"] {
		line_count += 1
	}
	for range resume_sections["personal_info"] {
		line_count += 1
	}
	for range resume_sections["experience"] {
		line_count += 1
	}
	for range resume_sections["education"] {
		line_count += 1
	}

	return line_count
}
