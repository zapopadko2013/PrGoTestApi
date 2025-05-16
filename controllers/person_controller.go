package controllers

import (
	"log"
	"net/http"
	"strconv"

	"PrGoRestApi/models"
	"PrGoRestApi/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePerson godoc
// @Summary Создание нового человека
// @Description Создает новую запись о человеке, с автоматическим обогащением данных (возраст, пол, национальность) по имени
// @Tags persons
// @Accept json
// @Produce json
// @Param person body models.PersonInput true "Информация о человеке"
// @Success 200 {object} models.Person
// @Failure 400 {object} map[string]interface{} "Ошибка в данных запроса"
// @Failure 500 {object} map[string]interface{} "Ошибка сервера"
// @Router /api/persons [post]
func CreatePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[INFO] CreatePerson: Начало обработки запроса")

		var input models.PersonInput
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Printf("[ERROR] CreatePerson: Ошибка биндинга JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("[DEBUG] CreatePerson: Входные данные: %+v\n", input)

		age, gender, nationality, err := services.Enrich(input.Name)
		if err != nil {
			log.Printf("[ERROR] CreatePerson: Ошибка обогащения данных: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обогащения данных"})
			return
		}

		person := models.Person{
			Name:        input.Name,
			Surname:     input.Surname,
			Patronymic:  input.Patronymic,
			Age:         &age,
			Gender:      &gender,
			Nationality: &nationality,
		}

		if err := db.Create(&person).Error; err != nil {
			log.Printf("[ERROR] CreatePerson: Ошибка сохранения в БД: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения в БД"})
			return
		}

		log.Printf("[INFO] CreatePerson: Человек успешно создан с ID %d\n", person.ID)
		c.JSON(http.StatusOK, person)
	}
}

// GetPersons godoc
// @Summary Получить всех людей
// @Description Возвращает список всех людей с возможностью фильтрации по имени, пагинации через limit и offset
// @Tags persons
// @Produce json
// @Param name query string false "Фильтр по имени"
// @Param limit query int false "Лимит"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.Person
// @Failure 500 {object} map[string]interface{} "Ошибка сервера"
// @Router /api/persons [get]
func GetPersons(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("[INFO] GetPersons: Начало обработки запроса")

		var people []models.Person
		query := db

		// Динамическая фильтрация по полям
		if name := c.Query("name"); name != "" {
			log.Printf("[DEBUG] Фильтр: name ILIKE %s", name)
			query = query.Where("name ILIKE ?", "%"+name+"%")
		}
		if surname := c.Query("surname"); surname != "" {
			log.Printf("[DEBUG] Фильтр: surname ILIKE %s", surname)
			query = query.Where("surname ILIKE ?", "%"+surname+"%")
		}
		if patronymic := c.Query("patronymic"); patronymic != "" {
			log.Printf("[DEBUG] Фильтр: patronymic ILIKE %s", patronymic)
			query = query.Where("patronymic ILIKE ?", "%"+patronymic+"%")
		}
		if ageStr := c.Query("age"); ageStr != "" {
			age, err := strconv.Atoi(ageStr)
			if err != nil {
				log.Printf("[WARN] Некорректный age: %s", ageStr)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный возраст"})
				return
			}
			log.Printf("[DEBUG] Фильтр: age = %d", age)
			query = query.Where("age = ?", age)
		}
		if gender := c.Query("gender"); gender != "" {
			log.Printf("[DEBUG] Фильтр: gender = %s", gender)
			query = query.Where("gender = ?", gender)
		}
		if nationality := c.Query("nationality"); nationality != "" {
			log.Printf("[DEBUG] Фильтр: nationality = %s", nationality)
			query = query.Where("nationality = ?", nationality)
		}

		// Обработка offset
		if offsetStr := c.Query("offset"); offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err != nil {
				log.Printf("[WARN] Некорректный offset: %s", offsetStr)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный offset"})
				return
			}
			query = query.Offset(offset)
		}

		// Обработка limit (если передан)
		if limitStr := c.Query("limit"); limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				log.Printf("[WARN] Некорректный limit: %s", limitStr)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный limit"})
				return
			}
			query = query.Limit(limit)
		}

		// Выполнение запроса
		if err := query.Find(&people).Error; err != nil {
			log.Printf("[ERROR] Ошибка получения данных: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных"})
			return
		}

		log.Printf("[INFO] Успешно: найдено %d человек", len(people))
		c.JSON(http.StatusOK, people)
	}
}

// GetPersonByID godoc
// @Summary Получить человека по ID
// @Description Возвращает данные человека по его уникальному идентификатору
// @Tags persons
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} models.Person
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Router /api/persons/{id} [get]
func GetPersonByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("[INFO] GetPersonByID: Запрос человека с ID %s\n", id)

		var person models.Person
		if err := db.First(&person, id).Error; err != nil {
			log.Printf("[ERROR] GetPersonByID: Пользователь не найден: %v\n", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}

		log.Printf("[INFO] GetPersonByID: Пользователь найден: %+v\n", person)
		c.JSON(http.StatusOK, person)
	}
}

// UpdatePerson godoc
// @Summary Обновить данные человека
// @Description Обновляет данные человека по ID с автоматическим обогащением данных, если имя изменилось
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param person body models.PersonInput true "Данные"
// @Success 200 {object} models.Person
// @Failure 400 {object} map[string]interface{} "Ошибка в данных запроса"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Failure 500 {object} map[string]interface{} "Ошибка сервера"
// @Router /api/persons/{id} [put]
func UpdatePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("[INFO] UpdatePerson: Обновление пользователя с ID %s\n", id)

		var person models.Person
		if err := db.First(&person, id).Error; err != nil {
			log.Printf("[ERROR] UpdatePerson: Пользователь не найден: %v\n", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}

		var input models.PersonInput
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Printf("[ERROR] UpdatePerson: Ошибка биндинга JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("[DEBUG] UpdatePerson: Входные данные: %+v\n", input)

		// Обогащаем данные, если имя изменилось или поля пусты
		if person.Name != input.Name || person.Age == nil || person.Gender == nil || person.Nationality == nil {
			age, gender, nationality, err := services.Enrich(input.Name)
			if err != nil {
				log.Printf("[ERROR] UpdatePerson: Ошибка обогащения данных: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обогащения данных"})
				return
			}

			person.Age = &age
			person.Gender = &gender
			person.Nationality = &nationality
		}

		person.Name = input.Name
		person.Surname = input.Surname
		person.Patronymic = input.Patronymic

		if err := db.Save(&person).Error; err != nil {
			log.Printf("[ERROR] UpdatePerson: Ошибка сохранения в БД: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения в БД"})
			return
		}

		log.Printf("[INFO] UpdatePerson: Пользователь с ID %s успешно обновлен\n", id)
		c.JSON(http.StatusOK, person)
	}
}

// DeletePerson godoc
// @Summary Удалить человека по ID
// @Description Удаляет человека из базы данных по ID
// @Tags persons
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} map[string]interface{} "Удалено успешно"
// @Failure 500 {object} map[string]interface{} "Ошибка удаления"
// @Router /api/persons/{id} [delete]
func DeletePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("[INFO] DeletePerson: Удаление пользователя с ID %s\n", id)

		if err := db.Delete(&models.Person{}, id).Error; err != nil {
			log.Printf("[ERROR] DeletePerson: Ошибка удаления: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
			return
		}

		log.Printf("[INFO] DeletePerson: Пользователь с ID %s удалён\n", id)
		c.JSON(http.StatusOK, gin.H{"message": "Удалено успешно"})
	}
}
