package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	dbURL          = "postgresql://postgres:NNA2s*123@localhost:5432/requests?sslmode=disable"
	backupDir      = "C:\\Backup"
	backupFileName = "backup.sql"
)

func main() {
	certFile := "./sertificate/__enostr_ru.full.crt"
	keyFile := "./sertificate/__enostr_ru.key"
	router := gin.Default()

	// Инициализация хранилища сессий
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.GET("/", showHomePage)
	router.GET("/add-request", showAddRequestPage)
	router.GET("/add-techRequest", showAddTechRequestPage)
	router.POST("/add-techRequest", addTechRequestPage)
	router.POST("/add-request", validateRequestData, addRequest)
	router.GET("/view-requests", viewRequests)
	router.POST("/add-technic", addTechnic)
	router.GET("/tech-requests", viewTechRequests)
	router.GET("/tech-accounting", showAccountingPage)
	router.POST("/mark-as-completed", markRequestAsCompleted)
	router.POST("/tech-done", doneTech)
	router.GET("/tasks", showTasksPage)
	router.GET("/send-question", showQuestionPage)
	router.POST("/login", handleLogin)
	router.POST("/loginAddRequest", handleLoginAddRequest)
	router.POST("/loginAddTechRequest", handleLoginAddTechRequest)

	// Middleware для проверки авторизации
	router.Use(checkAuthMiddleware())

	//router.Run(":80")
	err := router.RunTLS(":443", certFile, keyFile)
	if err != nil {
		log.Fatal("Failed to start HTTPS server: ", err)
	}
	c := cron.New()

	// Расписание выполнения резервного копирования (каждые 2 дня в 00:00)
	c.AddFunc("0 0 */2 * *", func() {
		err := performBackup()
		if err != nil {
			log.Println("Ошибка при выполнении резервного копирования:", err)
		}
	})

	c.Start()

	// Запуск бесконечного цикла для ожидания выполнения задач
	select {}
}

func addTechnic(c *gin.Context) {
	model := c.PostForm("model")
	serial_number := c.PostForm("serial-number")
	status := c.PostForm("status")
	name := c.PostForm("name")
	// Подключение к базе данных
	db, err := sql.Open("postgres", "postgresql://postgres:NNA2s*123@localhost:5432/requests?sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		fmt.Printf("error: %v", err)
		return
	}
	defer db.Close()
	// Вставка данных в таблицу requests
	insertQuery := "INSERT INTO equipment (model, serial_number, status, name) VALUES ($1, $2, $3, $4)"
	_, err = db.Exec(insertQuery, model, serial_number, status, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add request"})
		return
	}
	c.Redirect(http.StatusSeeOther, "tech-accounting")
}

func showAccountingPage(c *gin.Context) {
	session := sessions.Default(c)
	flashMessages := session.Flashes()
	session.Save()
	if auth := session.Get("authenticated3"); auth != nil && auth.(bool) {
		c.HTML(http.StatusOK, "tech_accounting.html", gin.H{
			"Authenticated": true,
			"FlashMessages": flashMessages,
		})
		return
	}
	c.HTML(http.StatusOK, "tech_accounting.html", gin.H{
		"Authenticated": false,
		"FlashMessages": flashMessages,
	})
	return
}

func doneTech(c *gin.Context) {
	techRequestID, _ := strconv.Atoi(c.PostForm("techRequestID"))

	// Update the request status in the database to completed
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE techrequests SET complete = $1 WHERE id = $2")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(true, techRequestID)
	if err != nil {
		log.Fatal(err)
	}

	// Set the response status as completed
	c.Redirect(http.StatusSeeOther, "/tech-requests")
}

func handleLoginAddTechRequest(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Проверка логина и пароля
	if username == "tech" && password == "P@Stech" {
		// Авторизация прошла успешно, устанавливаем флаг авторизации в сессии
		session.Set("authenticated3", true)
		session.Save()
		c.Redirect(http.StatusSeeOther, "/tech-requests")

	} else {
		session.AddFlash("Неверный логин или пароль")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/tech-requests")
	}
}

func viewTechRequests(c *gin.Context) {
	session := sessions.Default(c)
	flashMessages := session.Flashes()
	session.Save()

	if auth := session.Get("authenticated3"); auth != nil && auth.(bool) {
		category := c.Query("category")

		techrequests, _ := getRequestsFromtechDB(category)

		for i := range techrequests {
			// Пометка цветом в соответствии с условиями
			if techrequests[i].Employed {
				techrequests[i].Color = "green" // Зеленый
			} else if daysUntilDeadline(techrequests[i].Date) <= 3 {
				techrequests[i].Color = "yellow" // Желтый
			} else if techrequests[i].Priority == "Срочно" {
				techrequests[i].Color = "red" // Красный
			} else {
				techrequests[i].Color = "none" // Нет пометки
			}
		}

		c.HTML(http.StatusOK, "technical_requests.html", gin.H{
			"Authenticated": true,
			"Requests":      techrequests,
			"FlashMessages": flashMessages,
		})
		return
	}
	c.HTML(http.StatusOK, "technical_requests.html", gin.H{
		"Authenticated": false,
		"FlashMessages": flashMessages,
	})
	return
}

func daysUntilDeadline(date time.Time) int {
	now := time.Now()
	duration := date.Sub(now)
	return int(duration.Hours() / 24)
}

func addTechRequestPage(c *gin.Context) {
	category := c.PostForm("technic")
	characteristic := c.PostForm("characteristic")
	employmentDate := c.PostForm("employmentDate")
	deadline := c.PostForm("deadline")
	date := time.Time{}
	if employmentDate != "" {
		date, _ = time.Parse("02.01.2006", employmentDate)
	} else if deadline != "" {
		date, _ = time.Parse("02.01.2006", deadline)
	}
	isEmployed := c.PostForm("employment") == "on"
	employee := c.PostForm("employee")
	priority := c.PostForm("priority")
	object := ""
	complete := false
	if c.PostForm("location") == "office" {
		object = "Офис"
	} else {
		object = c.PostForm("objectName")
	}
	// Подключение к базе данных
	db, err := sql.Open("postgres", "postgresql://postgres:NNA2s*123@localhost:5432/requests?sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		fmt.Printf("error: %v", err)
		return
	}
	defer db.Close()

	// Вставка данных в таблицу techrequests
	insertQuery := "INSERT INTO techrequests (category, characteristic, employee, date, priority, employed, object, complete) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err = db.Exec(insertQuery, category, characteristic, employee, date, priority, isEmployed, object, complete)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add request"})
		fmt.Printf("error: %v", err)
		return
	}

	// Добавление флэш-сообщения о успешной отправке
	session := sessions.Default(c)
	session.AddFlash("Отправка успешно проведена")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/")
}

func showAddTechRequestPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_technical_request.html", gin.H{})
}

func handleLoginAddRequest(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Проверка логина и пароля
	if username == "master" && password == "P@Smaster" {
		// Авторизация прошла успешно, устанавливаем флаг авторизации в сессии
		session.Set("authenticated2", true)
		session.Save()
		c.Redirect(http.StatusSeeOther, "/add-request")

	} else {
		session.AddFlash("Неверный логин или пароль")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/add-request")
	}
}
func showQuestionPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_question.html", nil)
}

// Middleware для проверки авторизации
func checkAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// Проверяем, авторизован ли пользователь
		if auth := session.Get("authenticated"); auth != nil && auth.(bool) {
			c.Next()
			return
		}

		// Если не авторизован, перенаправляем на страницу входа
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
	}
}

func handleLogin(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Проверка логина и пароля
	if username == "kadr" && password == "P@Skadr" {
		// Авторизация прошла успешно, устанавливаем флаг авторизации в сессии
		session.Set("authenticated", true)
		session.Save()
		c.Redirect(http.StatusSeeOther, "/view-requests")

	} else {
		session.AddFlash("Неверные данные для входа в просмотр заявок")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/view-requests")
	}
}

func viewRequests(c *gin.Context) {
	session := sessions.Default(c)
	flashMessages := session.Flashes()
	session.Save()
	// Проверяем, авторизован ли пользователь
	if auth := session.Get("authenticated"); auth != nil && auth.(bool) {
		category := c.Query("category")

		// Получите данные заявок из базы данных или другого источника в соответствии с выбранной категорией
		// В данном примере создан простой список заявок
		requests, _ := getRequestsFromDB(category)
		c.HTML(http.StatusOK, "view_requests.html", gin.H{
			"Authenticated": true,
			"Requests":      requests,
			"FlashMessages": flashMessages,
		})
		return
	}
	c.HTML(http.StatusOK, "view_requests.html", gin.H{
		"Authenticated": false,
		"FlashMessages": flashMessages,
	})
	return
}

func showHomePage(c *gin.Context) {
	session := sessions.Default(c)
	flashMessages := session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"FlashMessages": flashMessages,
	})
}

func showTasksPage(c *gin.Context) {
	c.HTML(http.StatusOK, "tasks.html", nil)
}

func showAddRequestPage(c *gin.Context) {
	session := sessions.Default(c)
	flashMessages := session.Flashes()
	session.Save()
	if auth := session.Get("authenticated2"); auth != nil && auth.(bool) {
		c.HTML(http.StatusOK, "add_request.html", gin.H{
			"FlashMessages": flashMessages,
			"Authenticated": true,
		})
		return
	}
	c.HTML(http.StatusOK, "add_request.html", gin.H{
		"Authenticated": false,
		"FlashMessages": flashMessages,
	})
	return

}

func addRequest(c *gin.Context) {
	object := c.PostForm("object")
	department := c.PostForm("department")
	brigade := c.PostForm("brigade")
	employees := c.PostForm("employees")

	// Подключение к базе данных
	db, err := sql.Open("postgres", "postgresql://postgres:NNA2s*123@localhost:5432/requests?sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		fmt.Printf("error: %v", err)
		return
	}
	defer db.Close()

	completed := false
	// Вставка данных в таблицу requests
	insertQuery := "INSERT INTO requests (object, department, brigade, employees, completed) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.Exec(insertQuery, object, department, brigade, employees, completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add request"})
		return
	}

	// Добавление флэш-сообщения о успешной отправке
	session := sessions.Default(c)
	session.AddFlash("Отправка успешно проведена")
	session.Save()

	c.Redirect(http.StatusSeeOther, "/")
}

func getRequestsFromDB(category string) ([]map[string]interface{}, error) {
	db, err := sql.Open("postgres", "postgresql://postgres:NNA2s*123@localhost:5432/requests?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()
	// Построение SQL-запроса в зависимости от выбранной категории
	var query string
	var queryParams []interface{}
	if category == "all" {
		query = "SELECT * FROM requests ORDER BY completed ASC"
	} else {
		query = "SELECT * FROM requests WHERE object = $1 ORDER BY completed ASC"
		queryParams = append(queryParams, category)
	}

	rows, err := db.Query(query, queryParams...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %v", err)
	}

	var requests []map[string]interface{}
	for rows.Next() {
		// Создание слайса с указателями на значения для каждого столбца
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		// Сканирование строки результата запроса в указатели значений
		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Создание карты значений для каждой заявки
		request := make(map[string]interface{})
		for i, column := range columns {
			if column == "employees" {
				// Разделение имен сотрудников по запятой
				employees := strings.Split(values[i].(string), ",")
				request[column] = employees
			} else if column == "completed" {
				request[column] = values[i]
			} else {
				request[column] = values[i]
			}
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}

	return requests, nil
}
func getRequestsFromtechDB(category string) ([]TechRequest, error) {
	db, err := sql.Open("postgres", "postgresql://postgres:NNA2s*123@localhost:5432/requests?sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var techRequests []TechRequest
	query := ""
	queryParams := []interface{}{}

	if category == "all" {
		query = "SELECT * FROM techrequests ORDER BY complete, date ASC"
	} else if category == "Трудоустройство" {
		query = "SELECT * FROM techrequests WHERE employed = $1 ORDER BY complete, date ASC"
		queryParams = append(queryParams, true)
	} else if category == "Высокий приоритет" {
		query = "SELECT * FROM techrequests WHERE priority = $1 ORDER BY complete, date ASC"
		queryParams = append(queryParams, "Срочно")
	} else if category == "Офис" {
		query = "SELECT * FROM techrequests WHERE object = $1 ORDER BY complete, date ASC"
		queryParams = append(queryParams, category)
	} else if category == "На объекте" {
		query = "SELECT * FROM techrequests WHERE object != $1 ORDER BY complete, date ASC"
		queryParams = append(queryParams, "Офис")
	} else {
		query = "SELECT * FROM techrequests WHERE category = $1 ORDER BY complete, date ASC"
		queryParams = append(queryParams, category)
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var techRequest TechRequest
		err := rows.Scan(&techRequest.ID, &techRequest.Category, &techRequest.Characteristic, &techRequest.Employee, &techRequest.Date, &techRequest.Priority, &techRequest.Employed, &techRequest.Object, &techRequest.Complete)
		if err != nil {
			return nil, err
		}
		techRequests = append(techRequests, techRequest)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return techRequests, nil
}

type TechRequest struct {
	ID             int
	Category       string
	Characteristic string
	Employee       string
	Date           time.Time
	Priority       string
	Employed       bool
	Object         string
	Complete       bool
	Color          string
}

var data = ""

func validateRequestData(c *gin.Context) {
	object := c.PostForm("object")
	department := c.PostForm("department")
	brigade := c.PostForm("brigade")
	employees := c.PostForm("employees")

	if object == "" || department == "" || brigade == "" || employees == "" {
		// Если хотя бы одно поле пустое, отправляем ошибку
		session := sessions.Default(c)
		session.AddFlash("Заполните все поля")
		session.Save()
		c.Redirect(http.StatusSeeOther, "/add-request")
		c.Abort()
		return
	}

	// Дополнительные проверки и валидации данных, если необходимо

	// Если данные прошли валидацию, продолжаем выполнение следующих обработчиков
	c.Next()
}
func markRequestAsCompleted(c *gin.Context) {
	// Get the request ID from the form
	requestID, _ := strconv.Atoi(c.PostForm("requestID"))

	// Update the request status in the database to completed
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE requests SET completed = $1 WHERE id = $2")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(true, requestID)
	if err != nil {
		log.Fatal(err)
	}

	// Set the response status as completed
	c.Redirect(http.StatusSeeOther, "/view-requests")
}
func performBackup() error {

	// Генерация имени файла резервной копии на основе текущей даты и времени
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupFile := fmt.Sprintf(backupFileName, backupDir, timestamp)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Выполнение команды pg_dump для создания резервной копии
	cmd := exec.Command("pg_dump", "-U", "postgres", "-d", "requests", "-f", backupFile)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", "NNA2s*123"))
	err = cmd.Run()
	if err != nil {
		return err
	}

	log.Println("Резервная копия успешно создана:", backupFile)
	return nil
}
