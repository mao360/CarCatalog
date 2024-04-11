package repo

import (
	"database/sql"
	"fmt"
	"github.com/mao360/CarCatalog"
	"strconv"
	"strings"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetAll(queryParamsMap map[string]string, offset int) ([]CarCatalog.Car, error) {
	fmt.Println("POSTGRES GET ALL STARTED")
	filterStr := ""
	for key := range queryParamsMap {
		if queryParamsMap[key] != "" {
			if key == "year" {
				year, err := strconv.Atoi(queryParamsMap["year"])
				if err != nil {
					return nil, fmt.Errorf("postgres.go, GetAll, strconv.Atoi: %w", err)
				}
				filterStr = fmt.Sprintf(`%s %s=%d AND`, filterStr, key, year)
			} else {
				filterStr = fmt.Sprintf(`%s %s='%s' AND`, filterStr, key, queryParamsMap[key])
			}
		}
	}

	filterStr = strings.TrimSuffix(filterStr, "AND")
	if filterStr != "" {
		filterStr = fmt.Sprintf(`WHERE %s`, filterStr)
	}

	query := fmt.Sprintf(
		`SELECT name, surname, patronymic, reg_num, mark, model, year
				FROM people JOIN cars ON id = owner
				%s
				LIMIT 5 OFFSET %s;`, filterStr, strconv.Itoa(offset))

	cars := make([]CarCatalog.Car, 0)

	rows, err := r.db.Query(query)
	defer rows.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			return cars, nil
		}
		return nil, fmt.Errorf("postgres.go, GetAll, db.Query: %w", err)
	}

	for rows.Next() {
		var res CarCatalog.Car
		err = rows.Scan(&res.Owner.Name, &res.Owner.Surname, &res.Owner.Patronymic, &res.RegNum, &res.Mark, &res.Model, &res.Year)
		if err != nil {
			return nil, fmt.Errorf("postgres.go, GetAll, rows.Scan: %w", err)
		}
		cars = append(cars, res)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("postgres.go, GetAll, rows.Err(): %w", err)
	}
	return cars, nil
}

func (r *Repo) DeleteByID(regNum string) error {
	query := fmt.Sprintf(
		`DELETE FROM cars
				WHERE reg_num = '%s';`, regNum)
	res, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) ChangeByID(car CarCatalog.Car) error {
	updateColumns := ""
	switch {
	case car.Mark != "":
		updateColumns = fmt.Sprintf(`%s mark = '%s',`, updateColumns, car.Mark)
		fallthrough
	case car.Model != "":
		updateColumns = fmt.Sprintf(`%s model = '%s',`, updateColumns, car.Model)
		fallthrough
	case car.Year != 0:
		updateColumns = fmt.Sprintf(`%s year = %s,`, updateColumns, strconv.Itoa(car.Year))
	}
	updateColumns = strings.TrimSuffix(updateColumns, ",")

	query := fmt.Sprintf(
		`UPDATE cars
				SET %s
				WHERE reg_num = '%s';`, updateColumns, car.RegNum)
	res, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("postgres.go, ChangeByID, db.Exec: %w", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres.go, ChangeByID, RowsAffected: %w", err)
	}

	return nil
}

func (r *Repo) AddNew(cars []CarCatalog.Car) error {
	queryInPeople :=
		`INSERT INTO people (name, surname, patronymic)
		VALUES`

	queryInCars :=
		`INSERT INTO cars (reg_num, mark, model, year, owner)
		VALUES`

	for _, val := range cars {
		queryInPeople = fmt.Sprintf(`%s ('%s', '%s', '%s'),`, queryInPeople, val.Owner.Name, val.Owner.Surname, val.Owner.Patronymic)
		queryInCars = fmt.Sprintf(`%s ('%s', '%s', '%s', %s, (
									SELECT id
									FROM people
									WHERE name = '%s' AND surname = '%s' AND patronymic = '%s')),`,
			queryInCars, val.RegNum, val.Mark, val.Model, strconv.Itoa(val.Year), val.Owner.Name, val.Owner.Surname, val.Owner.Patronymic)
	}

	queryInPeople = strings.TrimSuffix(queryInPeople, ",")
	queryInCars = strings.TrimSuffix(queryInCars, ",")

	_, err := r.db.Exec(queryInPeople)
	if err != nil {
		return fmt.Errorf("postgres.go, AddNew, 1st db.Exec: %w", err)
	}

	_, err = r.db.Exec(queryInCars)
	if err != nil {
		return fmt.Errorf("postgres.go, AddNew, 2nd db.Exec: %w", err)
	}
	return nil
}
