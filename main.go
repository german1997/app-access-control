package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {

	DBConnection()

	r := gin.Default()

	// define the routes
	r.GET("/user/:id", GetTeacherHandler)

	err := r.Run(checkPort())

	if err != nil {
		log.Fatalf("impossible to start server: %s", err)
	}
}

func DBConnection() {

	cfg := mysql.Config{
		User:   "root",
		Passwd: "EA5a2eaCedbBFD4ddeGCHHDdbcHe2HDD",
		Net:    "tcp",
		Addr:   "monorail.proxy.rlwy.net:21442",
		DBName: "railway",
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		fmt.Println("Db fallo al conectarse: ", err)
	}
}

func GetTeacherHandler(c *gin.Context) {

	userid := c.Param("id")

	cc, err := Validate(userid)

	if err == 1 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	if err == 2 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "impossible to retrieve teacher"})
		return
	}

	c.JSON(http.StatusOK, cc)

}

func Validate(userid string) (*User, int) {
	var user User

	query := "SELECT * FROM user where rut = '" + userid + "'"

	row := db.QueryRow(query)

	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Rut); err != nil {
		if err == sql.ErrNoRows {
			return nil, 1
		}
		return nil, 2
	}

	return &user, 3
}

type Response struct {
	Message string `json:"message"`
}

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Rut       string `json:"rut"`
}

func checkPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}
