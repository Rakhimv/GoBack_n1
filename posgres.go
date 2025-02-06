package main


import (
	"database/sql"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)


type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}


var db *sql.DB;

func main() {
	connStr := "host=localhost port=5432 user=postgres password=gravitifolz dbname=go sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к бд: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Ошибка при пинге бд: ", err)
	}

	router := gin.Default()


	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})



	router.GET("/users", getUsers)
	router.POST("/users", createUsers)
	router.DELETE("/users/:id", deleteUser)

	if err := router.Run(":6060"); err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
	}
}

func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()







	var users []User
	for rows.Next(){
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}


func createUsers(c *gin.Context) {
	var user User


	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err := db.QueryRow("INSERT INTO users(name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}


func deleteUser(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: ": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}