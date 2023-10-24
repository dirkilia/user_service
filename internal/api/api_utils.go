package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"user_service/internal/logging"
	"user_service/types"
)

var (
	AGIFY_URL       = "https://api.agify.io/?name="
	GENDERIZE_URL   = "https://api.genderize.io/?name="
	NATIONALIZE_URL = "https://api.nationalize.io/?name="
)

type apiFunc func(http.ResponseWriter, *http.Request) error

// makeHTTPHandleFunc logging errors that we can get from handlers
func makeHTTPHandleFunc(f apiFunc, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			logger.Info(err)
		}
	}
}

// Writes valid filter to the filterName
func createStringFilter(w http.ResponseWriter, filter_type string, filter string, filterName *string) error {
	var err error
	*filterName, err = validateAndReturnFilter(filter_type, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

// Validate filter parameter
func validateAndReturnFilter(filter_type string, filter string) (string, error) {
	switch filter_type {
	case "gender":
		digits, _ := regexp.MatchString(`\d`, filter)
		in_range, _ := regexp.MatchString(`male|female`, filter)
		if digits || !in_range {
			return "", errors.New("gender query parameter is not valid")
		}
		return filter, nil
	case "ageplus":
		nodigits, _ := regexp.MatchString(`\D`, filter)
		if nodigits {
			return "", errors.New("ageplus query parameter is not valid")
		}
		return filter, nil
	case "ageminus":
		nodigits, _ := regexp.MatchString(`\D`, filter)
		if nodigits {
			return "", errors.New("ageminus query parameter is not valid")
		}
		return filter, nil
	case "nationality":
		digits, _ := regexp.MatchString(`\d`, filter)
		if digits {
			return "", errors.New("nationality query parameter is not valid")
		}
		return filter, nil
	default:
		return "", errors.New("query parameter is undefined")
	}
}

func getAge(name string) (int64, error) {
	var age Age
	body, err := getInfoFromExternalApi(name, AGIFY_URL)
	if err != nil {
		return 0, errors.New("cant get body age")
	}

	err = json.Unmarshal(body, &age)
	if err != nil {
		return 0, errors.New("cant unmarshal to age")
	}

	if age.Age == 0 {
		return 0, errors.New("can't determine age, check if name is valid")
	}

	return age.Age, nil
}

func getGender(name string) (string, error) {
	var gender Gender
	body, err := getInfoFromExternalApi(name, GENDERIZE_URL)
	if err != nil {
		return "", errors.New("cant get body gender")
	}
	err = json.Unmarshal(body, &gender)

	if err != nil {
		return "", errors.New("cant unmarshal to gender")
	}
	if gender.Gender == "" {
		return "", errors.New("can't determine gender, check if name is valid")
	}
	return gender.Gender, nil
}

func getNationality(name string) (string, error) {
	var nationality Nationality
	body, err := getInfoFromExternalApi(name, NATIONALIZE_URL)
	if err != nil {
		log.Fatalf("cant get body nationality %s", err)
	}
	err = json.Unmarshal(body, &nationality)

	if err != nil {
		log.Fatalf("cant unmarshal to nationality %s", err)
	}

	if len(nationality.Country) == 0 {
		return "", errors.New("can't determine nationality, check if name is valid")
	}

	sort.Slice(nationality.Country, func(p, q int) bool {
		return nationality.Country[p].Probability < nationality.Country[q].Probability
	})

	most_probable := nationality.Country[len(nationality.Country)-1]

	return most_probable.Country_id, nil
}

// Returns []byte from given API url
func getInfoFromExternalApi(name string, url string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", url, name))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Fetches data from external API to the struct Person
func (s *APIServer) fetchExternalApiData(p *types.Person) error {
	age, err := getAge(p.Name)
	if err != nil {
		s.logger.Debugf("can't get age, error: %v", err)
		return err
	}
	p.Age = age
	gender, err := getGender(p.Name)
	if err != nil {
		s.logger.Debugf("can't get gender, error: %v", err)
		return err
	}
	p.Gender = gender
	nationality, err := getNationality(p.Name)
	if err != nil {
		s.logger.Debugf("can't get nationality, error: %v", err)
		return err
	}
	p.Nationality = nationality
	return nil
}
