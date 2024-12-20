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
ID int
Username string
Points int
ClassID string
Password string
}
 
 
type AdminPageData struct {
Students []User
ClassID string
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
class_id TEXT NOT NULL
);
`)
if err != nil {
log.Fatal("Ошибка создания таблицы:", err)
}
 
// Задаем начальных пользователей
users := []User{
{Username: "admin", Password: "adminpass", ClassID: "admin"}, // Администратор
{Username: "student1", Password: "pass123", Points: 10, ClassID: "class1"},
{Username: "student2", Password: "pass123", Points: 20, ClassID: "class1"},
{Username: "student3", Password: "pass123", Points: 15, ClassID: "class2"},
{Username: "student4", Password: "pass123", Points: 25, ClassID: "class3"},
{Username: "student5", Password: "pass123", Points: 30, ClassID: "class4"},

{Username: "Максим", Password: "pass123", Points: 10, ClassID: "class5а"},
{Username: "Айрат", Password: "pass123", Points: 10, ClassID: "class5б"},
{Username: "Арни", Password: "pass123", Points: 10, ClassID: "class5в"},
{Username: "Иван Токарев", Password: "pass123", Points: 10, ClassID: "class6а"},
{Username: "Дмитрий", Password: "pass123", Points: 10, ClassID: "class6б"},
{Username: "Сергей", Password: "pass123", Points: 10, ClassID: "class6в"},
{Username: "Константин Йогурт", Password: "pass123", Points: 10, ClassID: "class7а"},
{Username: "Алмаз Табачков", Password: "pass123", Points: 10, ClassID: "class7б"},
{Username: "Олег Шерстяной", Password: "pass123", Points: 10, ClassID: "class7в"},
{Username: "Гусейн Гасанов", Password: "pass123", Points: 10, ClassID: "class8а"},
{Username: "Иван Золо", Password: "pass123", Points: 10, ClassID: "class8б"},
{Username: "Илон МаКС", Password: "pass123", Points: 10, ClassID: "class8в"},
{Username: "Владислав", Password: "pass123", Points: 10, ClassID: "class9а"},
{Username: "Александр", Password: "pass123", Points: 10, ClassID: "class9б"},
{Username: "Максим Барабанов", Password: "pass123", Points: 10, ClassID: "class9в"},
{Username: "Мария", Password: "pass123", Points: 10, ClassID: "class10-1"},
{Username: "Роман", Password: "pass123", Points: 10, ClassID: "class10-2"},
{Username: "Олег Мангал", Password: "pass123", Points: 10, ClassID: "class11-1"},
{Username: "Савелий", Password: "pass123", Points: 10, ClassID: "class11-2"},
}
 
 
for _, user := range users {
_, err := db.Exec("INSERT OR IGNORE INTO users (username, password, points, class_id) VALUES (?, ?, ?, ?)",
user.Username, user.Password, user.Points, user.ClassID)
if err != nil {
log.Printf("Ошибка добавления пользователя %s: %v\n", user.Username, err)
}
}
}
 
 
func loginHandler(w http.ResponseWriter, r *http.Request) {
if r.Method == http.MethodPost {
username := r.FormValue("username")
password := r.FormValue("password")
 
var dbPassword string
var classID string
err := db.QueryRow("SELECT password, class_id FROM users WHERE username = ?", username).Scan(&dbPassword, &classID)
if err != nil || password != dbPassword {
http.Error(w, "Неверные учётные данные", http.StatusUnauthorized)
return
}
 
if classID == "admin" { 
http.Redirect(w, r, "/admin?class_id=1", http.StatusSeeOther)
} else {
http.Redirect(w, r, "/home?username="+username, http.StatusSeeOther)
}
return
}
 
tmpl.ExecuteTemplate(w, "login.html", nil)
}
 
 
func homeHandler(w http.ResponseWriter, r *http.Request) {
username := r.URL.Query().Get("username")
 
var user User
err := db.QueryRow("SELECT id, username, points FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Points)
if err != nil {
http.Error(w, "Пользователь не найден", http.StatusInternalServerError)
return
}
 
tmpl.ExecuteTemplate(w, "home.html", user)
}
 
 
func adminHandler(w http.ResponseWriter, r *http.Request) {
classID := r.URL.Query().Get("class_id")
if classID == "" {
classID = "class1"
}
 
rows, err := db.Query("SELECT id, username, points FROM users WHERE class_id = ?", classID)
if err != nil {
http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
return
}
defer rows.Close()
 
students := []User{}
for rows.Next() {
var user User
if err := rows.Scan(&user.ID, &user.Username, &user.Points); err != nil {
http.Error(w, "Ошибка сканирования пользователей", http.StatusInternalServerError)
return
}
students = append(students, user)
}
 
data := AdminPageData{
Students: students,
ClassID: classID,
}
tmpl.ExecuteTemplate(w, "admin.html", data)
}
 
 
func updatePointsHandler(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Недопустимый метод запуска", http.StatusMethodNotAllowed)
return
}
 
id := r.FormValue("id")
change := r.FormValue("change")
 
changeValue, err := strconv.Atoi(change)
if err != nil {
http.Error(w, "Недопустимое значение баллов", http.StatusBadRequest)
return
}
 
var currentPoints int
err = db.QueryRow("SELECT points FROM users WHERE id = ?", id).Scan(&currentPoints)
if err != nil {
http.Error(w, "Пользователь не найден", http.StatusInternalServerError)
return
}
newPoints := currentPoints + changeValue
if newPoints < 0 {
newPoints = 0
}
 
_, err = db.Exec("UPDATE users SET points = ? WHERE id = ?", newPoints, id)
if err != nil {
http.Error(w, "Ошибка обновления баллов", http.StatusInternalServerError)
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