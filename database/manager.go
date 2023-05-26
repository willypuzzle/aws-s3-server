package database

import (
	"aws-s3-server/types"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
)

type Database struct {
	DbUser string
	DbPass string
	DbHost string
	DbPort string
	DbName string
}

func Builder() *Database {
	return &Database{
		DbUser: getEnvironmentUsername(),
		DbHost: getEnvironmentHost(),
		DbPass: getEnvironmentPassword(),
		DbName: getEnvironmentDatabaseName(),
		DbPort: getEnvironmentPort(),
	}
}

const DbOption string = "?charset=utf8&parseTime=true"

func (d *Database) tableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s')", d.DbName, tableName)

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

func (d *Database) openDatabase() *sql.DB {
	dsn := d.DbUser + ":" + d.DbPass + "@tcp(" + d.DbHost + ":" + d.DbPort + ")/" + d.DbName + DbOption
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	bucketExists, _ := d.tableExists(db, "Buckets")
	objectExists, _ := d.tableExists(db, "Objects")

	if bucketExists == false || objectExists == false {
		migrate(db)
	}

	return db
}

func closeDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		fmt.Println("Database closing failed")
	}
}

func migrate(db *sql.DB) {
	migration1 := `
		CREATE TABLE IF NOT EXISTS Buckets (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL UNIQUE
		);`
	_, err1 := db.Exec(migration1)
	if err1 != nil {
		fmt.Println("Migration of Buckets failed.")
		os.Exit(1)
	}

	migration2 := `CREATE TABLE IF NOT EXISTS Objects (
			id INT PRIMARY KEY AUTO_INCREMENT,
			bucket_id INT NOT NULL,
			key_path VARCHAR(255) NOT NULL UNIQUE,
			object_data BLOB,
		    content_type VARCHAR(255) NOT NULL,
    		uuid VARCHAR(255) NOT NULL UNIQUE,
            file_size INT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (bucket_id) REFERENCES Buckets(id)
		);
	`
	_, err2 := db.Exec(migration2)
	if err2 != nil {
		fmt.Println("Migration of Objects failed.")
		os.Exit(1)
	}
}

func (d *Database) BucketExists(name string) (int, error) {
	db := d.openDatabase()
	rows, err := db.Query("SELECT id FROM Buckets WHERE name = (?)", name)
	if err != nil {
		return 0, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Unable to close rows")
		}
	}(rows)

	var id int
	if rows.Next() == true {
		rows.Scan(&id)
	} else {
		return 0, nil
	}

	return id, nil
}

func (d *Database) InsertBucket(bucket *types.Bucket) error {
	db := d.openDatabase()
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

func (d *Database) InsertOrUpdateObject(object *types.Object) error {
	object.Size = len(object.Data)
	db := d.openDatabase()
	stmt, err := db.Prepare(
		"REPLACE INTO Objects (bucket_id, key_path, object_data, content_type, uuid, file_size) VALUES (?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)
	res, err := stmt.Exec(object.BucketId, object.Key, object.Data, object.ContentType, object.Uuid, object.Size)
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

func (d *Database) selectBuckets() ([]types.Bucket, error) {
	db := d.openDatabase()
	rows, err := db.Query("SELECT id, name FROM Buckets")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var buckets []types.Bucket
	for rows.Next() {
		var bucket types.Bucket
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

func (d *Database) SelectContents(bucketName string, prefix string) ([]types.Content, error) {
	db := d.openDatabase()
	rows, err := db.Query(`SELECT os.key_path, os.file_size, os.updated_at, os.uuid
                                 FROM Buckets AS bs
                                 INNER JOIN Objects AS os ON (bs.id = os.bucket_id)
                                 WHERE bs.name = ? AND os.key_path LIKE ?`, bucketName, prefix+"%")
	if err != nil {
		return nil, err
	}

	var contents []types.Content
	for rows.Next() {
		var content types.Content
		err := rows.Scan(&content.Key, &content.Size, &content.LastModified, &content.ETag)
		if err != nil {
			return nil, err
		}
		content.Key = strings.TrimPrefix(content.Key, prefix)
		contents = append(contents, content)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	return contents, nil
}
