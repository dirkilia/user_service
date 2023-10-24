package api

import (
	"user_service/internal/logging"
	"user_service/internal/repository"
)

type APIServer struct {
	listenAddr string
	store      repository.Storage
	logger     logging.Logger
}

type Age struct {
	Age int64
}

type Gender struct {
	Gender string
}

type Nationality struct {
	Country []struct {
		Country_id  string
		Probability float64
	}
}
