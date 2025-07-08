package analyzer

func AnalyzeFiles(files []string) []Finding {
	var results []Finding
	for _, file := range files {
		results = append(results, AnalyzeCleanCode(file)...)
		results = append(results, analyzePerformance(file)...)
		results = append(results, analyzeSecurity(file)...)
	}
	return results
}
