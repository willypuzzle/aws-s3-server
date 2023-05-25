package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type Bucket struct {
	Id   int
	Name string
}

type Object struct {
	Id          int
	BucketId    int
	Key         string
	Data        []byte
	ContentType string
}

var DbUser = ""
var DbPass = ""
var DbHost = ""
var DbPort = ""
var DbName = ""

const DbOption string = "?charset=utf8"

func tableExists(db *sql.DB, databaseName string, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s')", databaseName, tableName)

	err := db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func getEnvironmentVariable(variableName string) string {
	env, err := os.LookupEnv(variableName)
	if err == false {
		fmt.Println(fmt.Sprintf("EnvironmentVariable '%s' is not set", variableName))
	}
	return env
}

func getEnvironmentHost() string {
	return getEnvironmentVariable("DB_HOST")
}

func getEnvironmentUsername() string {
	return getEnvironmentVariable("DB_USER")
}

func getEnvironmentPort() string {
	return getEnvironmentVariable("DB_PORT")
}

func getEnvironmentPassword() string {
	return getEnvironmentVariable("DB_PASSWORD")
}

func getEnvironmentDatabaseName() string {
	return getEnvironmentVariable("DB_NAME")
}

func initEnvVariable() {
	if DbUser == "" || DbPass == "" || DbHost == "" || DbName == "" || DbPort == "" {
		DbUser = getEnvironmentUsername()
		DbPass = getEnvironmentPassword()
		DbHost = getEnvironmentHost()
		DbName = getEnvironmentDatabaseName()
		DbPort = getEnvironmentPort()
	}
}

func openDatabase() *sql.DB {
	initEnvVariable()
	dsn := DbUser + ":" + DbPass + "@" + DbHost + DbPort + "/" + DbName + DbOption
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	bucketExists, _ := tableExists(db, DbName, "Buckets")
	objectExists, _ := tableExists(db, DbName, "Objects")

	if bucketExists == false || objectExists == false {
		err = migrate(db)
	}

	return db
}

func closeDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {

	}
}

func migrate(db *sql.DB) error {
	migration := `
		CREATE TABLE IF NOT EXISTS Buckets (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS Objects (
			id INT PRIMARY KEY AUTO_INCREMENT,
			bucket_id INT NOT NULL,
			key VARCHAR(255) NOT NULL,
			object_data BLOB,
		    content_type VARCHAR(255) NOT NULL
			FOREIGN KEY (bucket_id) REFERENCES Buckets(id)
		);
	`
	_, err := db.Exec(migration)
	return err
}

func insertBucket(bucket *Bucket) error {
	db := openDatabase()
	res, err := db.Exec("INSERT INTO Buckets (name) VALUES (?)", bucket.Name)
	if err != nil {
		return err
	}
	closeDatabase(db)
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	bucket.Id = int(id)
	return nil
}

func insertObject(object *Object) error {
	db := openDatabase()
	stmt, err := db.Prepare("INSERT INTO Objects (bucket_id, key, object_data, content_type) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)
	res, err := stmt.Exec(object.BucketId, object.Key, object.Data, object.ContentType)
	if err != nil {
		return err
	}
	closeDatabase(db)
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	object.Id = int(id)
	return nil
}

func selectBuckets() ([]Bucket, error) {
	db := openDatabase()
	rows, err := db.Query("SELECT id, name FROM Buckets")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var buckets []Bucket
	for rows.Next() {
		var bucket Bucket
		err := rows.Scan(&bucket.Id, &bucket.Name)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return buckets, nil
}
