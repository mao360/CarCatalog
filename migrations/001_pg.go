package migrations

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(up, down)
}

func up(tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS people (
			    "id" SERIAL PRIMARY KEY,
			    "name" VARCHAR(20),
			    "surname" VARCHAR(20),
			    "patronymic" VARCHAR(20));
				CREATE TABLE IF NOT EXISTS cars (
  				"car_id" SERIAL PRIMARY KEY, 
                "reg_num" VARCHAR(9),
                "mark" VARCHAR(20),
    			"model" VARCHAR(20),
    			"year" INTEGER,
    			"owner" SERIAL REFERENCES people(id));`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("migrationsUp: %w", err)
	}
	return nil
}

func down(tx *sql.Tx) error {
	query := `DROP TABLE cars;
				DROP TABLE people;`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("migrationsDown: %w", err)
	}
	return nil
}
