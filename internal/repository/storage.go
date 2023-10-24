package repository

import (
	"database/sql"
	"fmt"
	"user_service/internal/logging"
	"user_service/types"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetPersons(string, string, string, string, int, int) ([]*types.Person, error)
	AddPerson(*types.Person) error
	UpdatePerson(int64, *types.Person) error
	DeletePersonById(int64) error
}

type PostgresStore struct {
	db     *sql.DB
	logger logging.Logger
}

func NewPostgresStore(connStr string, logger logging.Logger) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db:     db,
		logger: logger,
	}, nil
}

func (s *PostgresStore) GetPersons(filterGender, filterAgeplus, filterAgeminus, filterNationality string, limit, offset int) ([]*types.Person, error) {
	str_query := generateQueryStringWithFilters(filterGender, filterAgeplus, filterAgeminus, filterNationality, limit, offset)
	rows, err := s.db.Query(str_query)
	if err != nil {
		s.logger.Infof("can't execute query GetPersons, error: %v", err)
		return nil, err
	}

	persons := []*types.Person{}
	for rows.Next() {
		person := new(types.Person)
		err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
		)

		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func (s *PostgresStore) AddPerson(person *types.Person) error {
	query := `insert into persons (name, surname, patronymic, age, gender, nationality)
	values ($1, $2, $3, $4, $5, $6)`
	if person.Name != "" && person.Surname != "" {
		_, err := s.db.Exec(
			query,
			person.Name,
			person.Surname,
			person.Patronymic,
			person.Age,
			person.Gender,
			person.Nationality)

		if err != nil {
			s.logger.Debugf("can't execute query AddPerson, error: %v", err)
			return err
		}
	} else {
		s.logger.Info("name or surname is empty")
	}

	return nil
}

func (s *PostgresStore) UpdatePerson(id int64, person *types.Person) error {
	query := `update persons set name = $1, 
	surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6
	where id = $7`
	if person.Name != "" && person.Surname != "" {
		_, err := s.db.Exec(
			query,
			person.Name,
			person.Surname,
			person.Patronymic,
			person.Age,
			person.Gender,
			person.Nationality,
			id)
		if err != nil {
			s.logger.Debugf("can't execute query DeletePersonById, error: %v", err)
			return err
		}
	} else {
		s.logger.Info("name or surname is empty")
	}

	return nil
}

func (s *PostgresStore) DeletePersonById(id int64) error {
	_, err := s.db.Query("delete from persons where id = $1", id)
	if err != nil {
		s.logger.Infof("can't execute query DeletePersonById, error: %v", err)
		return err
	}
	return nil
}

// Generates a sql query string with given parameters
func generateQueryStringWithFilters(filterGender, filterAgeplus, filterAgeminus, filterNationality string, limit, offset int) string {
	does_already_have_clause := false
	str_query := "select * from persons"
	if filterGender != "" {
		str_query = fmt.Sprintf("%s where gender = '%s'", str_query, filterGender)
		does_already_have_clause = true
	}
	if filterAgeplus != "" && !does_already_have_clause {
		str_query = fmt.Sprintf("%s where age > %s", str_query, filterAgeplus)
		does_already_have_clause = true
	} else if filterAgeplus != "" {
		str_query = fmt.Sprintf("%s and age > %s", str_query, filterAgeplus)
	}
	if filterAgeminus != "" && !does_already_have_clause {
		str_query = fmt.Sprintf("%s where age < %s", str_query, filterAgeminus)
		does_already_have_clause = true
	} else if filterAgeminus != "" {
		str_query = fmt.Sprintf("%s and age < %s", str_query, filterAgeminus)
	}
	if filterNationality != "" && !does_already_have_clause {
		str_query = fmt.Sprintf("%s where nationality = '%s'", str_query, filterNationality)
		does_already_have_clause = true
	} else if filterNationality != "" {
		str_query = fmt.Sprintf("%s and nationality = '%s'", str_query, filterNationality)
	}

	str_query = str_query + " order by id asc"

	if limit != -1 {
		str_query = str_query + fmt.Sprintf(" limit %d ", limit)
	}

	if offset != -1 {
		str_query = str_query + fmt.Sprintf(" offset %d ", offset)
	}
	return str_query
}
