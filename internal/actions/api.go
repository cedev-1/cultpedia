package actions

import (
	"cultpedia/internal/models"
	"cultpedia/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var apiData models.APIData

const defaultPort = "8080"

func RunAPIServer(serverPort ...string) {
	port := defaultPort
	if len(serverPort) > 0 && serverPort[0] != "" {
		port = serverPort[0]
	}

	if err := loadData(); err != nil {
		log.Fatalf("Error loading data: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/questions", handleQuestions)
	mux.HandleFunc("/api/geography/countries", handleCountries)
	mux.HandleFunc("/api/geography/regions", handleRegions)
	mux.HandleFunc("/api/geography/continents", handleContinents)
	mux.HandleFunc("/api/geography/flags/", handleFlags)
	mux.HandleFunc("/api/", handleRoot)
	mux.HandleFunc("/", handleRoot)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Cultpedia API server running on http://localhost:%s\n", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func loadData() error {
	var err error

	if err := loadManifests(); err != nil {
		return fmt.Errorf("error loading manifests: %w", err)
	}

	apiData.Questions, err = utils.LoadQuestions()
	if err != nil {
		return fmt.Errorf("error loading questions: %w", err)
	}

	apiData.Countries, err = utils.LoadCountries()
	if err != nil {
		return fmt.Errorf("error loading countries: %w", err)
	}

	apiData.Regions, err = utils.LoadRegions()
	if err != nil {
		return fmt.Errorf("error loading regions: %w", err)
	}

	apiData.Continents, err = utils.LoadContinents()
	if err != nil {
		return fmt.Errorf("error loading continents: %w", err)
	}

	return nil
}

func loadManifests() error {
	geoFile, err := os.Open(utils.GeographyManifestFile)
	if err != nil {
		return err
	}
	defer func() { _ = geoFile.Close() }()

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

	qFile, err := os.Open(utils.ManifestFile)
	if err != nil {
		return err
	}
	defer func() { _ = qFile.Close() }()

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

	_ = json.NewEncoder(w).Encode(response)
}

func handleQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  apiData.Questions,
		"count": len(apiData.Questions),
	})
}

func handleCountries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  apiData.Countries,
		"count": len(apiData.Countries),
	})
}

func handleRegions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  apiData.Regions,
		"count": len(apiData.Regions),
	})
}

func handleContinents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  apiData.Continents,
		"count": len(apiData.Continents),
	})
}

func handleFlags(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/api/geography/flags/")
	code = strings.TrimSuffix(code, ".svg")
	code = strings.ToLower(code)

	if code == "" {
		http.Error(w, "Country code required", http.StatusBadRequest)
		return
	}

	if len(code) != 2 || !isAlphaOnly(code) {
		http.Error(w, "Invalid country code format", http.StatusBadRequest)
		return
	}

	flagPath := filepath.Join(utils.FlagsSVGDir, code+".svg")

	if _, err := os.Stat(flagPath); os.IsNotExist(err) {
		http.Error(w, "Flag not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeFile(w, r, flagPath)
}

func isAlphaOnly(s string) bool {
	for _, c := range s {
		if c < 'a' || c > 'z' {
			return false
		}
	}
	return true
}
