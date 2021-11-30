package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const (
	USERNAME = "root"
	PASSWORD = ""
	NETWORK = "tcp"
	SERVER = "127.0.0.1"
	PORT = 3306
	DATABASE = "twitter"
)

type User struct {
	ID          string
	Username    string
	Password    string
	ProfileName string
	ProfileImg  string
}

type Post struct {
	ID          string
	Username    string
	ProfileName string
	ProfileImg  string
	Text        string
	Img         string
	Time        string
}

func GetDB() *sql.DB {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("Error while launching MySQL: ", err)
	}
	return db
}

func CreateUserTable() error {
	db := GetDB()
	defer db.Close()
	sql := `CREATE TABLE IF NOT EXISTS users(
		id INT(4) PRIMARY KEY AUTO_INCREMENT,
        username VARCHAR(64) NOT NULL UNIQUE,
        password VARCHAR(64) NOT NULL,
        profilename VARCHAR(64) NOT NULL UNIQUE,
        profileimg VARCHAR(128)
	); `

	if _, err := db.Exec(sql); err != nil {
		fmt.Println("Error while creating table: ", err)
		return err
	}
	fmt.Println("User Table created")
	return nil
}

func CreatePostTable() error {
	db := GetDB()
	defer db.Close()
	sql := `CREATE TABLE IF NOT EXISTS users(
		id INT(4) PRIMARY KEY AUTO_INCREMENT,
        username VARCHAR(64) NOT NULL UNIQUE,
        profilename VARCHAR(64) NOT NULL UNIQUE,
        profileimg VARCHAR(128),
        text VARCHAR(1024) NOT NULL,
        img VARCHAR(256),
        time VARCHAR(64) NOT NULL
	); `

	if _, err := db.Exec(sql); err != nil {
		fmt.Println("Error while creating table: ", err)
		return err
	}
	fmt.Println("Post Table created")
	return nil
}

func InsertUser(username, password, profilename, profileimg string) string{
	db := GetDB()
	defer db.Close()
	if _, err := QueryUser(username); err == nil {
		return "User already exists"
	} else {
		_, err := db.Exec("insert INTO users(username,password,profilename,profileimg) values(?,?,?,?)",username, password, profilename, profileimg)
		if err != nil{
			return "Error while creating user"
		}
		fmt.Println("\nUser created: ", username)
		return ""
	}
}

func QueryUser(username string) (*User, error) {
	db := GetDB()
	defer db.Close()
	user := new(User)
	row := db.QueryRow("select * from users where username=?", username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.ProfileName, &user.ProfileImg); err != nil {
		fmt.Printf("Error while mapping user: %v", err)
		return user, err
	}
	fmt.Println("\nUser found: ", *user)
	return user, nil
}

func main() {
	QueryUser("test")
}