package types

type Person struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int64  `json:"age"`
	Gender      string `json:"sex"`
	Nationality string `json:"nationality"`
}

func NewPerson(name, surname, patronymic, gender, nationality string, age int64) *Person {
	return &Person{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality}
}
