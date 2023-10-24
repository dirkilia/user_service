package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"user_service/internal/logging"
	"user_service/internal/repository"
	"user_service/types"

	"github.com/gorilla/mux"
)

func NewAPIServer(listenAddr string, store repository.Storage, logger logging.Logger) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
		logger:     logger,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/persons", makeHTTPHandleFunc(s.handlePersons, s.logger))

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handlePersons(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetPersons(w, r)
	case "POST":
		return s.handleAddPerson(w, r)
	case "DELETE":
		return s.handleDeletePersonById(w, r)
	case "PATCH":
		return s.handleUpdatePerson(w, r)
	}

	return nil
}

func (s *APIServer) handleGetPersons(w http.ResponseWriter, r *http.Request) error {
	var err error
	// Code below get query parameters for limit, offset, gender, age, nationality from url
	strLimit := r.URL.Query().Get("limit")
	limit := -1
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < 1 {
			http.Error(w, "limit query parameter is no valid number", http.StatusBadRequest)
			return err
		}
	}

	strOffset := r.URL.Query().Get("offset")
	offset := -1
	if strOffset != "" {
		offset, err = strconv.Atoi(strOffset)
		if err != nil || offset < -1 {
			http.Error(w, "offset query parameter is no valid number", http.StatusBadRequest)
			return err
		}
	}

	filter := r.URL.Query().Get("gender")
	filterGender := ""
	if filter != "" {
		err = createStringFilter(w, "gender", filter, &filterGender)
		if err != nil {
			return err
		}
	}

	// ageplus parameter means that we search rows more than a given parameter value
	filter = r.URL.Query().Get("ageplus")
	filterAgeplus := ""
	if filter != "" {
		err = createStringFilter(w, "ageplus", filter, &filterAgeplus)
		if err != nil {
			return err
		}
	}

	// ageminus parameter means that we search rows less than a given parameter value
	filter = r.URL.Query().Get("ageminus")
	filterAgeminus := ""
	if filter != "" {
		err = createStringFilter(w, "ageminus", filter, &filterAgeminus)
		if err != nil {
			return err
		}
	}

	filter = r.URL.Query().Get("nationality")
	filterNationality := ""
	if filter != "" {
		err = createStringFilter(w, "nationality", filter, &filterNationality)
		if err != nil {
			return err
		}
	}

	persons, err := s.store.GetPersons(filterGender, filterAgeplus, filterAgeminus, filterNationality, limit, offset)
	if err != nil {
		s.logger.Debugf("can't execute GetPersons, error: %v", err)
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(persons)
	if err != nil {
		s.logger.Debugf("can't Encode persons, error: %v", err)
		return err
	}

	return nil
}

func (s *APIServer) handleAddPerson(w http.ResponseWriter, r *http.Request) error {
	req := new(types.Person)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Debugf("can't decode to &req, error: %v", err)
		return err
	}
	err := s.fetchExternalApiData(req)
	if err != nil {
		s.logger.Debugf("can't execute fetchExternalApiData, error: %v", err)
		return err
	}

	person := types.NewPerson(req.Name, req.Surname, req.Patronymic, req.Gender, req.Nationality, req.Age)

	if err := s.store.AddPerson(person); err != nil {
		s.logger.Debugf("can't execute AddPerson, error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (s *APIServer) handleUpdatePerson(w http.ResponseWriter, r *http.Request) error {
	req := new(types.Person)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Debugf("can't decode to &req, error: %v", err)
		return err
	}

	if err := s.store.UpdatePerson(req.ID, req); err != nil {
		s.logger.Debugf("can't execute UpdatePerson, error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (s *APIServer) handleDeletePersonById(w http.ResponseWriter, r *http.Request) error {
	var id struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		s.logger.Debugf("can't decode to &id, error: %v", err)
		return err
	}

	err := s.store.DeletePersonById(id.ID)
	if err != nil {
		s.logger.Debugf("can't execute DeletePersonById, error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}
