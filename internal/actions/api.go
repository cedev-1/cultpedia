package actions

import (
	"bufio"
	"cultpedia/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var apiData models.APIData

func RunAPIServer(serverPort ...string) {

	if err := loadData(); err != nil {
		log.Fatalf("Error loading data: %v", err)
	}

	http.HandleFunc("/api/questions", handleQuestions)
	http.HandleFunc("/api/geography/countries", handleCountries)
	http.HandleFunc("/api/geography/regions", handleRegions)
	http.HandleFunc("/api/geography/continents", handleContinents)
	http.HandleFunc("/api/geography/flags/", handleFlags)
	http.HandleFunc("/api/", handleRoot)
	http.HandleFunc("/", handleRoot)

	fmt.Printf("Cultpedia API server running on http://localhost:%s\n", serverPort[0])

	if err := http.ListenAndServe(":"+serverPort[0], nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func loadData() error {
	var err error

	if err := loadManifests(); err != nil {
		return fmt.Errorf("error loading manifests: %w", err)
	}

	apiData.Questions, err = loadQuestions()
	if err != nil {
		return fmt.Errorf("error loading questions: %w", err)
	}

	apiData.Countries, err = loadCountries()
	if err != nil {
		return fmt.Errorf("error loading countries: %w", err)
	}

	apiData.Regions, err = loadRegions()
	if err != nil {
		return fmt.Errorf("error loading regions: %w", err)
	}

	apiData.Continents, err = loadContinents()
	if err != nil {
		return fmt.Errorf("error loading continents: %w", err)
	}

	return nil
}

func loadManifests() error {
	geoFile, err := os.Open("datasets/geography/manifest.json")
	if err != nil {
		return err
	}
	defer geoFile.Close()

	var geoManifest struct {
		Version   string `json:"version"`
		UpdatedAt string `json:"updated_at"`
	}
	if err := json.NewDecoder(geoFile).Decode(&geoManifest); err != nil {
		return err
	}
	apiData.Manifests.Geography = models.ManifestInfo{
		Version:   geoManifest.Version,
		UpdatedAt: geoManifest.UpdatedAt,
	}

	qFile, err := os.Open("datasets/general-knowledge/manifest.json")
	if err != nil {
		return err
	}
	defer qFile.Close()

	var qManifest struct {
		Version   string `json:"version"`
		UpdatedAt string `json:"updated_at"`
	}
	if err := json.NewDecoder(qFile).Decode(&qManifest); err != nil {
		return err
	}
	apiData.Manifests.GeneralKnowledge = models.ManifestInfo{
		Version:   qManifest.Version,
		UpdatedAt: qManifest.UpdatedAt,
	}

	return nil
}

func loadQuestions() ([]models.Question, error) {
	file, err := os.Open("datasets/general-knowledge/questions.ndjson")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var questions []models.Question
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var q models.Question
		if err := json.Unmarshal(scanner.Bytes(), &q); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, scanner.Err()
}

func loadCountries() ([]models.Country, error) {
	file, err := os.Open("datasets/geography/countries.ndjson")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var countries []models.Country
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var c models.Country
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}

	return countries, scanner.Err()
}

func loadRegions() ([]models.Region, error) {
	file, err := os.Open("datasets/geography/regions.ndjson")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var regions []models.Region
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var r models.Region
		if err := json.Unmarshal(scanner.Bytes(), &r); err != nil {
			return nil, err
		}
		regions = append(regions, r)
	}

	return regions, scanner.Err()
}

func loadContinents() ([]models.Continent, error) {
	file, err := os.Open("datasets/geography/continents.ndjson")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var continents []models.Continent
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var c models.Continent
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			return nil, err
		}
		continents = append(continents, c)
	}

	return continents, scanner.Err()
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	response := models.APIRootResponse{
		Name:        "Cultpedia API",
		Version:     "1.0",
		Description: "API for Cultpedia questions and geography data",
		Datasets: map[string]interface{}{
			"general_knowledge": map[string]string{
				"version":    apiData.Manifests.GeneralKnowledge.Version,
				"updated_at": apiData.Manifests.GeneralKnowledge.UpdatedAt,
			},
			"geography": map[string]string{
				"version":    apiData.Manifests.Geography.Version,
				"updated_at": apiData.Manifests.Geography.UpdatedAt,
			},
		},
		Endpoints: []models.EndpointInfo{
			{
				Path:        "/api/questions",
				Method:      "GET",
				Description: "Get all questions",
			},
			{
				Path:        "/api/geography/countries",
				Method:      "GET",
				Description: "Get all countries",
			},
			{
				Path:        "/api/geography/regions",
				Method:      "GET",
				Description: "Get all regions",
			},
			{
				Path:        "/api/geography/continents",
				Method:      "GET",
				Description: "Get all continents",
			},
			{
				Path:        "/api/geography/flags/{code}",
				Method:      "GET",
				Description: "Get country flag SVG (use ISO Alpha2 code)",
			},
		},
		Stats: map[string]int{
			"questions":  len(apiData.Questions),
			"countries":  len(apiData.Countries),
			"regions":    len(apiData.Regions),
			"continents": len(apiData.Continents),
		},
	}
	
	json.NewEncoder(w).Encode(response)
}

func handleQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  apiData.Questions,
		"count": len(apiData.Questions),
	})
}

func handleCountries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  apiData.Countries,
		"count": len(apiData.Countries),
	})
}

func handleRegions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  apiData.Regions,
		"count": len(apiData.Regions),
	})
}

func handleContinents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  apiData.Continents,
		"count": len(apiData.Continents),
	})
}

func handleFlags(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/api/geography/flags/")
	code = strings.TrimSuffix(code, ".svg")

	if code == "" {
		http.Error(w, "Country code required", http.StatusBadRequest)
		return
	}

	flagPath := filepath.Join("datasets", "geography", "assets", "flags", "svg", code+".svg")

	if _, err := os.Stat(flagPath); os.IsNotExist(err) {
		http.Error(w, "Flag not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeFile(w, r, flagPath)
}
