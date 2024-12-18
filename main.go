package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var tmpl *template.Template

// User структура
type User struct {
	ID       int
	Username string
	Points   int
	ClassID  int
	Password string
}

// AdminPageData структура для передачи данных на admin.html
type AdminPageData struct {
	Students []User
	ClassID  int
}

// Иниц бд
func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			points INTEGER DEFAULT 0,
			class_id INTEGER NOT NULL
		);
	`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Задаем начальных пользователей
	users := []User{
		{Username: "admin", Password: "adminpass", ClassID: 0}, // Администратор
		{Username: "student1", Password: "pass123", Points: 10, ClassID: 1},
		{Username: "student2", Password: "pass123", Points: 20, ClassID: 1},
		{Username: "student3", Password: "pass123", Points: 15, ClassID: 2},
		{Username: "student4", Password: "pass123", Points: 25, ClassID: 3},
		{Username: "student5", Password: "pass123", Points: 30, ClassID: 4},
	}

	// Вставляем пользователей
	for _, user := range users {
		_, err := db.Exec("INSERT OR IGNORE INTO users (username, password, points, class_id) VALUES (?, ?, ?, ?)",
			user.Username, user.Password, user.Points, user.ClassID)
		if err != nil {
			log.Printf("Failed to add user %s: %v\n", user.Username, err)
		}
	}
}

// Страница входа
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var dbPassword string
		var classID int
		err := db.QueryRow("SELECT password, class_id FROM users WHERE username = ?", username).Scan(&dbPassword, &classID)
		if err != nil || password != dbPassword {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if classID == 0 { // Если admin
			http.Redirect(w, r, "/admin?class_id=1", http.StatusSeeOther)
		} else { // Если обычный пользователь
			http.Redirect(w, r, "/home?username="+username, http.StatusSeeOther)
		}
		return
	}

	tmpl.ExecuteTemplate(w, "login.html", nil)
}

// Главная страница для пользователей
func homeHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	var user User
	err := db.QueryRow("SELECT id, username, points FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Points)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "home.html", user)
}

// Панель администратора
func adminHandler(w http.ResponseWriter, r *http.Request) {
	classIDStr := r.URL.Query().Get("class_id")
	classID, err := strconv.Atoi(classIDStr)
	if err != nil || classID < 1 || classID > 4 {
		classID = 1
	}

	rows, err := db.Query("SELECT id, username, points FROM users WHERE class_id = ?", classID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	students := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Points); err != nil {
			http.Error(w, "Error scanning users", http.StatusInternalServerError)
			return
		}
		students = append(students, user)
	}

	data := AdminPageData{
		Students: students,
		ClassID:  classID,
	}
	tmpl.ExecuteTemplate(w, "admin.html", data)
}

// Обновление баллов
func updatePointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	change := r.FormValue("change")

	changeValue, err := strconv.Atoi(change)
	if err != nil {
		http.Error(w, "Invalid points value", http.StatusBadRequest)
		return
	}

	var currentPoints int
	err = db.QueryRow("SELECT points FROM users WHERE id = ?", id).Scan(&currentPoints)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Проверка, чтобы баллы не ушли в минус
	newPoints := currentPoints + changeValue
	if newPoints < 0 {
		newPoints = 0
	}

	_, err = db.Exec("UPDATE users SET points = ? WHERE id = ?", newPoints, id)
	if err != nil {
		http.Error(w, "Failed to update points", http.StatusInternalServerError)
		return
	}

	// Отправляем обновлённое значение баллов клиенту
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(newPoints)))
}

// Главная функция
func main() {
	initDB()
	defer db.Close()

	tmpl = template.Must(template.ParseFiles("templates/login.html", "templates/home.html", "templates/admin.html"))

	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/update-points", updatePointsHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
