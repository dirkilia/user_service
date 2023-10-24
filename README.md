# user_service

Service gets name, surname, patronymic(optional), adds age, gender and nationality via external apis to the PostgreSQL database
>https://api.agify.io/?name=Dmitriy - for age
>https://api.genderize.io/?name=Dmitriy - for gender
>https://api.nationalize.io/?name=Dmitriy - for nationality

## API REST Methods
1. handleGetPersons to get list of persons \(GET\)
2. handleAddPerson to create a new person row in database \(POST\)
3. handleUpdatePerson to update an existing person row in database by id \(PATCH\)
4. handleDeletePersonById to delete an existing person row in database by id \(DELETE\)

## Other

Database structure created with migrations, use
```shell
make migrationsup
```

To run server use 
```shell
make run
```

Code is covered with info- and debug- logs using [logrus](https://github.com/sirupsen/logrus)