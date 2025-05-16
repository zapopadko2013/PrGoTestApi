package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// response структуры
type agifyResponse struct {
	Age int `json:"age"`
}

type genderizeResponse struct {
	Gender string `json:"gender"`
}

type nationalizeResponse struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

// EnrichPerson обогащает данные по имени через внешние API
func Enrich(name string) (int, string, string, error) {
	log.Printf("[DEBUG] Enrich: start enrichment for name=%s\n", name)

	age, err := getAge(name)
	if err != nil {
		log.Printf("[ERROR] getAge failed for name=%s: %v\n", name, err)
		return 0, "", "", err
	}
	log.Printf("[INFO] getAge: name=%s, age=%d\n", name, age)

	gender, err := getGender(name)
	if err != nil {
		log.Printf("[ERROR] getGender failed for name=%s: %v\n", name, err)
		return 0, "", "", err
	}
	log.Printf("[INFO] getGender: name=%s, gender=%s\n", name, gender)

	nationality, err := getNationality(name)
	if err != nil {
		log.Printf("[ERROR] getNationality failed for name=%s: %v\n", name, err)
		return 0, "", "", err
	}
	log.Printf("[INFO] getNationality: name=%s, nationality=%s\n", name, nationality)

	log.Printf("[DEBUG] Enrich: enrichment complete for name=%s\n", name)
	return age, gender, nationality, nil
}

func getAge(name string) (int, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	log.Printf("[DEBUG] getAge: fetching URL %s\n", url)
	var result agifyResponse
	err := fetchJSON(url, &result)
	if err != nil {
		log.Printf("[ERROR] getAge: failed fetching %s: %v\n", url, err)
	}
	return result.Age, err
}

func getGender(name string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	log.Printf("[DEBUG] getGender: fetching URL %s\n", url)
	var result genderizeResponse
	err := fetchJSON(url, &result)
	if err != nil {
		log.Printf("[ERROR] getGender: failed fetching %s: %v\n", url, err)
	}
	return result.Gender, err
}

func getNationality(name string) (string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	log.Printf("[DEBUG] getNationality: fetching URL %s\n", url)
	var result nationalizeResponse
	err := fetchJSON(url, &result)
	if err != nil {
		log.Printf("[ERROR] getNationality: failed fetching %s: %v\n", url, err)
		return "", err
	}
	if len(result.Country) > 0 {
		log.Printf("[INFO] getNationality: country detected %s with probability %f\n", result.Country[0].CountryID, result.Country[0].Probability)
		return result.Country[0].CountryID, nil
	}
	log.Printf("[INFO] getNationality: no country detected, returning 'unknown'\n")
	return "unknown", nil
}

func fetchJSON(url string, target interface{}) error {
	log.Printf("[DEBUG] fetchJSON: start fetching %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[ERROR] fetchJSON: http.Get failed for %s: %v\n", url, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("API request failed: %s", resp.Status)
		log.Printf("[ERROR] fetchJSON: %s\n", errMsg)
		return fmt.Errorf(errMsg)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(target); err != nil {
		log.Printf("[ERROR] fetchJSON: JSON decode failed for %s: %v\n", url, err)
		return err
	}

	log.Printf("[DEBUG] fetchJSON: successfully fetched and decoded %s\n", url)
	return nil
}
