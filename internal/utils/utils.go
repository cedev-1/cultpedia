package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cultpedia/internal/models"
)

const (
	ManifestFile             = "datasets/general-knowledge/manifest.json"
	QuestionsFile            = "datasets/general-knowledge/questions.ndjson"
	ThemesFile               = "datasets/general-knowledge/themes.ndjson"
	SubthemesFile            = "datasets/general-knowledge/subthemes.ndjson"
	TagsFile                 = "datasets/general-knowledge/tags.ndjson"
	NewQuestionFile          = "datasets/new-question.json"
	NewQuestionTrueFalseFile = "datasets/new-question-true-false.json"

	GeographyManifestFile  = "datasets/geography/manifest.json"
	CountriesFile          = "datasets/geography/countries.ndjson"
	ContinentsFile         = "datasets/geography/continents.ndjson"
	RegionsFile            = "datasets/geography/regions.ndjson"
	FlagsSVGDir            = "datasets/geography/assets/flags/svg"
)

func LoadQuestions() ([]models.Question, error) {
	data, err := os.ReadFile(QuestionsFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var questions []models.Question
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var q models.Question
		if err := json.Unmarshal([]byte(line), &q); err != nil {
			return nil, fmt.Errorf("json parsing error at line %d: %v", len(questions)+1, err)
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func SaveQuestion(q models.Question) error {
	minified, err := json.Marshal(q)
	if err != nil {
		return fmt.Errorf("minification error: %v", err)
	}
	ndjsonLine := string(minified) + "\n"

	if _, err := os.Stat(QuestionsFile); err == nil {
		f, err := os.Open(QuestionsFile)
		if err != nil {
			return fmt.Errorf("error opening file for check: %v", err)
		}
		defer func() { _ = f.Close() }()
		stat, err := f.Stat()
		if err != nil {
			return fmt.Errorf("error getting file stat: %v", err)
		}
		size := stat.Size()
		if size > 0 {
			buf := make([]byte, 1)
			_, err := f.ReadAt(buf, size-1)
			if err != nil {
				return fmt.Errorf("error reading file end: %v", err)
			}
			if buf[0] != '\n' {
				ndjsonLine = "\n" + ndjsonLine
			}
		}
	}

	f, err := os.OpenFile(QuestionsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err := f.WriteString(ndjsonLine); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func SlugExists(slug string) bool {
	questions, err := LoadQuestions()
	if err != nil {
		return false
	}
	for _, q := range questions {
		if q.Slug == slug {
			return true
		}
	}
	return false
}

func DetectModifiedTemplateFile() (filePath string, questionType string) {
	if isTemplateModified(NewQuestionFile, "default-question-slug") {
		return NewQuestionFile, "single_choice"
	}
	if isTemplateModified(NewQuestionTrueFalseFile, "default-true-false-question-slug") {
		return NewQuestionTrueFalseFile, "true_false"
	}
	return "", ""
}

func isTemplateModified(filePath, defaultSlug string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	var q models.Question
	if err := json.Unmarshal(data, &q); err != nil {
		return false
	}
	return q.Slug != defaultSlug
}

func PrintHelp() {
	helpText := `
Cultpedia - Question Dataset Management Tool

USAGE FOR CONTRIBUTORS:
  ./cultpedia                  Launch interactive UI (recommended for adding questions)

USAGE FOR MAINTAINERS:
  ./cultpedia [command]

COMMANDS:
  help                       Show this help message
  
  Questions Dataset:
  validate                      Validate the questions dataset for consistency and correctness
  check-duplicates              Check for duplicate questions in the dataset
  check-translations            Check for missing translations in the dataset
  add                           Add a new question to the dataset via interactive prompts
  sync-themes                   Synchronize themes and subthemes with the questions dataset
  bump-version                  Increment version and update manifest (automated in CI)
  
  Geography Dataset:
  validate-geography            Validate the geography dataset (countries, continents, regions)
  check-geography-duplicates    Check for duplicate entries in geography dataset
  check-geography-translations  Check for missing translations in geography dataset
  bump-geography-version        Increment geography version and update checksums (automated in CI)
  
  General:
  init [dataset-name]        Initialize a new Cultpedia dataset structure

CONTRIBUTION GUIDE:
  For questions: Fork → Edit template file → Use TUI to add → Create PR
  For code: Fork → Edit code → Run tests → Create PR

For more info, see CONTRIBUTING.md in the docs/ folder.
Or visit:
  https://docs.culturae.me/cultpedia/

Thank you for contributing to Cultpedia!
`
	fmt.Println(helpText)
}

func LoadCountries() ([]models.Country, error) {
	data, err := os.ReadFile(CountriesFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var countries []models.Country
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var c models.Country
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			return nil, fmt.Errorf("json parsing error at line %d: %v", len(countries)+1, err)
		}
		countries = append(countries, c)
	}
	return countries, nil
}

func LoadContinents() ([]models.Continent, error) {
	data, err := os.ReadFile(ContinentsFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var continents []models.Continent
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var c models.Continent
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			return nil, fmt.Errorf("json parsing error at line %d: %v", len(continents)+1, err)
		}
		continents = append(continents, c)
	}
	return continents, nil
}

func LoadRegions() ([]models.Region, error) {
	data, err := os.ReadFile(RegionsFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var regions []models.Region
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var r models.Region
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			return nil, fmt.Errorf("json parsing error at line %d: %v", len(regions)+1, err)
		}
		regions = append(regions, r)
	}
	return regions, nil
}
