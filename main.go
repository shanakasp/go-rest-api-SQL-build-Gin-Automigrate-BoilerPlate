package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/restgoq")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Auto migrate the database schema
	autoMigrate()

	router := gin.Default()

	router.POST("/users", createUser) // POST /users
	// router.GET("/users/:id", getUser)   // GET /users/:id
	router.GET("/users", getUsers)     // GET /users
	// router.PUT("/users/:id", updateUser) // PUT /users/:id
	// router.DELETE("/users/:id", deleteUser) // DELETE /users/:id
	
	router.Run(":8080")
}

func autoMigrate() {
	// Create users table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS userss (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL DEFAULT NULL
	);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err.Error())
	}
}

func createUser(c *gin.Context) {

	var body struct {
	    Name string
	}
	c.Bind(&body)

	//prepare SQL statment
	stmt, err := db.Prepare("INSERT INTO userss(name) VALUES(?)")
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	// Get last inserted ID
	id, err := stmt.Exec(body.Name)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//return the created user as JSON
	user := User {
		ID : fmt.Sprintf("%d", id),
		Name : body.Name,
	}

	c.JSON(http.StatusCreated, user)
}

func getUsers(c *gin.Context) {
    c.Header("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, name FROM userss")
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
	    var user User
	    if err := rows.Scan(&user.ID, &user.Name); err != nil {
	        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	        return
	    }
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	    return
	}

	c.JSON(http.StatusOK, users)
}

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)

// 	stmt, err := db.Prepare("SELECT id, name FROM users WHERE id = ?")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer stmt.Close()

// 	var user User
// 	if err := stmt.QueryRow(params["id"]).Scan(&user.ID, &user.Name); err != nil {
// 		if err == sql.ErrNoRows {
// 			http.Error(w, "User not found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	json.NewEncoder(w).Encode(user)
// }

// func updateUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)

// 	stmt, err := db.Prepare("UPDATE users SET name = ? WHERE id = ?")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer stmt.Close()

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer r.Body.Close()

// 	var user User
// 	if err := json.Unmarshal(body, &user); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	_, err = stmt.Exec(user.Name, params["id"])
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprintf(w, "User with ID = %s was updated", params["id"])
// }

// func deleteUser(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)

// 	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(params["id"])
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprintf(w, "User with ID = %s was deleted", params["id"])
// }
