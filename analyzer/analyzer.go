package analyzer

import "github.com/schollz/progressbar/v3"

func AnalyzeFiles(files []string) []Finding {
	var results []Finding
	bar := progressbar.Default(int64(len(files)))
	for _, file := range files {
		results = append(results, AnalyzeCleanCode(file)...)
		results = append(results, analyzePerformance(file)...)
		results = append(results, analyzeSecurity(file)...)
		bar.Add(1)
	}
	return results
}
