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

{Username: "Абросов Андрей", Password: "ajhpun17", Points: 10, ClassID: "class7А"},
{Username: "Гимальдинов Артём", Password: "7rcle8ng", Points: 10, ClassID: "class7А"},
{Username: "Дейнеко Макар", Password: "dw86iqnh", Points: 10, ClassID: "class7А"},
{Username: "Егоров Артём", Password: "5v6csijh", Points: 10, ClassID: "class7А"},
{Username: "Завалишин Дмитрий", Password: "i77wd744", Points: 10, ClassID: "class7А"},
{Username: "Завалишина Марина", Password: "puzzluf0", Points: 10, ClassID: "class7А"},
{Username: "Захаров Тимофей", Password: "s9h7n7pv", Points: 10, ClassID: "class7А"},
{Username: "Крашенинников Макар", Password: "t2aodytr", Points: 10, ClassID: "class7А"},
{Username: "Курчик Анна", Password: "aepazin6", Points: 10, ClassID: "class7А"},
{Username: "Лексин Егор", Password: "7xc3cavw", Points: 10, ClassID: "class7А"},
{Username: "Матрасова Милена", Password: "24y4414j", Points: 10, ClassID: "class7А"},
{Username: "Мухлынина Виктория", Password: "gfcoyqw5", Points: 10, ClassID: "class7А"},
{Username: "Николаев Лев", Password: "vk1s69wb", Points: 10, ClassID: "class7А"},
{Username: "Пастухов Дмитрий", Password: "grhx4cdq", Points: 10, ClassID: "class7А"},
{Username: "Перегримова Дарья", Password: "0qcv87np", Points: 10, ClassID: "class7А"},
{Username: "Пичугов Егор", Password: "ht0yhljg", Points: 10, ClassID: "class7А"},
{Username: "Ральников Кирилл", Password: "mmwgnlb6", Points: 10, ClassID: "class7А"},
{Username: "Ремизов Валерий", Password: "lwsd4udc", Points: 10, ClassID: "class7А"},
{Username: "Ремизова Анастасия", Password: "59itzvus", Points: 10, ClassID: "class7А"},
{Username: "Савчук Игорь", Password: "g6oxufhb", Points: 10, ClassID: "class7А"},
{Username: "Самильянов Максим", Password: "9ntv1q7h", Points: 10, ClassID: "class7А"},
{Username: "Ситникова Татьяна", Password: "n36qd19t", Points: 10, ClassID: "class7А"},
{Username: "Таран Григорий", Password: "7ha9042z", Points: 10, ClassID: "class7А"},

{Username: "Голонастикова Владислава", Password: "keqc6r56", Points: 10, ClassID: "class7Б"},
{Username: "Даулетшина Милана", Password: "5ytv9z5r", Points: 10, ClassID: "class7Б"},
{Username: "Дэль Даниил", Password: "q59h5ro5", Points: 10, ClassID: "class7Б"},
{Username: "Еремеева Екатерина", Password: "6jzf5nkc", Points: 10, ClassID: "class7Б"},
{Username: "Занина Виктория", Password: "7r8bfn66", Points: 10, ClassID: "class7Б"},
{Username: "Кокорин Дмитрий", Password: "339rwt5g", Points: 10, ClassID: "class7Б"},
{Username: "Кокорин Константин", Password: "woj8lu8z", Points: 10, ClassID: "class7Б"},
{Username: "Колосова София", Password: "q4l9bzy5", Points: 10, ClassID: "class7Б"},
{Username: "Лоцманов Сергей", Password: "cn9ao1wh", Points: 10, ClassID: "class7Б"},
{Username: "Миланич Ярослав", Password: "wiuq5amp", Points: 10, ClassID: "class7Б"},
{Username: "Миниахметова Дарья", Password: "lmv650gn", Points: 10, ClassID: "class7Б"},
{Username: "Мирзаянова Малика", Password: "e56cswkn", Points: 10, ClassID: "class7Б"},
{Username: "Осипенко Константин", Password: "vg5eyrou", Points: 10, ClassID: "class7Б"},
{Username: "Пешков Гаврил", Password: "m4g3zz93", Points: 10, ClassID: "class7Б"},
{Username: "Рябых Варвара", Password: "qqnfvow7", Points: 10, ClassID: "class7Б"},
{Username: "Соколов Илья", Password: "o0n7bz92", Points: 10, ClassID: "class7Б"},
{Username: "Стафеев Владимир", Password: "wu9oz12w", Points: 10, ClassID: "class7Б"},
{Username: "Уколов Игорь", Password: "pmhf96un", Points: 10, ClassID: "class7Б"},
{Username: "Ханипов Владислав", Password: "zcq6nx7l", Points: 10, ClassID: "class7Б"},
{Username: "Хейфиц Анна", Password: "78tbfxuk", Points: 10, ClassID: "class7Б"},

{Username: "Виноградов Дмитрий", Password: "unhzuuu8", Points: 10, ClassID: "class7В"},
{Username: "Владимирова Марина", Password: "yz9kzkbe", Points: 10, ClassID: "class7В"},
{Username: "Гаврилова Ева", Password: "r0xq7lbw", Points: 10, ClassID: "class7В"},
{Username: "Гулиева Элина", Password: "98jim6rh", Points: 10, ClassID: "class7В"},
{Username: "Казьмин Макар", Password: "6m3ijzfy", Points: 10, ClassID: "class7В"},
{Username: "Калугина Полина", Password: "f2fyeyd0", Points: 10, ClassID: "class7В"},
{Username: "Калужинская София", Password: "iceilui1", Points: 10, ClassID: "class7В"},
{Username: "Каюмов Малик", Password: "wi8xjt0q", Points: 10, ClassID: "class7В"},
{Username: "Кожевников Михаил", Password: "l5tytlvw", Points: 10, ClassID: "class7В"},
{Username: "Кружалова Антонина", Password: "p5mkfzcn", Points: 10, ClassID: "class7В"},
{Username: "Минлыгараев Роман", Password: "55jr2s74", Points: 10, ClassID: "class7В"},
{Username: "Пантелеев Семён", Password: "ay12rydm", Points: 10, ClassID: "class7В"},
{Username: "Петров Егор", Password: "2ztil8rx", Points: 10, ClassID: "class7В"},
{Username: "Пильщиков Сергей", Password: "818hdm0m", Points: 10, ClassID: "class7В"},
{Username: "Савин Максим", Password: "ag46jxeh", Points: 10, ClassID: "class7В"},
{Username: "Сарафанова Вероника", Password: "1fp894vb", Points: 10, ClassID: "class7В"},
{Username: "Сибирцева Валерия", Password: "i3j387ab", Points: 10, ClassID: "class7В"},
{Username: "Стародубцев Александр", Password: "927a1wcz", Points: 10, ClassID: "class7В"},
{Username: "Стародубцев Дмитрий", Password: "v67xrvvo", Points: 10, ClassID: "class7В"},
{Username: "Стародубцева Мария", Password: "15ap99km", Points: 10, ClassID: "class7В"},
{Username: "Ступак Матвей", Password: "oh3603bg", Points: 10, ClassID: "class7В"},
{Username: "Табаков Виктор", Password: "qbbv0ufu", Points: 10, ClassID: "class7В"},
{Username: "Токарев Матвей", Password: "s4gn8ak5", Points: 10, ClassID: "class7В"},
{Username: "Чебыкина Арина", Password: "9kpz2808", Points: 10, ClassID: "class7В"},

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