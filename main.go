package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger" // fiber-swagger middleware
	_ "go.mod/docs"
	"go.mod/models"
	"go.mod/storage"
	"gorm.io/gorm"
)

// @title Swagger Example API
// @version 1.0
// @description This is a Pharmacy API.
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

type Pharmacy struct {
	Owner   string `json:"owner"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Repository struct {
	DB *gorm.DB
}

// ShowPharmacies godoc
// @Summary      Create a pharmacy
// @Description  Adds a new pharmacy entity to DB
// @Tags         Pharmacy
// @Accept       json
// @Produce      json
// @Router       /create_pharmacy [post]
// @Success      200  "pharmacy created"
// @Failure      422  "repquest failed"
// @Param		 pharmacy	body		models.Pharmacy	true	"Add Pharmacy"
func (r *Repository) CreatePharmacy(context *fiber.Ctx) error {
	pharmacy := Pharmacy{}

	err := context.BodyParser(&pharmacy)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "repquest failed"})
		return err
	}
	err = r.DB.Create(&pharmacy).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create pharmacy"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "pharmacy created"})

	return nil
}

// ShowPharmacies godoc
// @Summary      Retrieve all pharmacies
// @Description  Retrieves all pharmacies from DB
// @Tags         Pharmacy
// @Accept       json
// @Produce      json
// @Router       /get_pharmacies [get]
// @Success      200  {array}  models.Pharmacy
// @Failure      400  "Could not get pharmacies"
func (r *Repository) GetPharmacies(context *fiber.Ctx) error {
	pharmacyModels := &[]models.Pharmacy{}

	err := r.DB.Find(pharmacyModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get pharmacies"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "pharmacies fetched",
		"data":    pharmacyModels,
	})

	return nil
}

// ShowPharmacy godoc
// @Summary      	Retrieve a pharmacy
// @Description  	Retrieves a single pharmacy by ID
// @Tags         	Pharmacy
// @Accept       	json
// @Produce     	json
// @Param			id	path		int	true	"Pharmacy ID"
// @Success      	200  {array}  models.Pharmacy
// @Failure      	400  "ID is required"
// @Failure      	404  "Could not get pharmacy"
// @Router       	/get_pharmacy/{id} [get]
func (r *Repository) GetPharmacyByID(context *fiber.Ctx) error {
	id := context.Params("id")
	pharmacyModel := &models.Pharmacy{}
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "id is required"})
		return nil
	}

	fmt.Println("the Id is: ", id)

	result := r.DB.Where("id = ?", id).First(pharmacyModel)
	if result.Error != nil || result.RowsAffected == 0 {
		context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "could not get pharmacy"})
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "pharmacy fetched",
		"data":    pharmacyModel,
	})

	return nil
}

// DeletePharmacy godoc
// @Summary      	Delete a pharmacy
// @Description  	Delete a single pharmacy by ID
// @Tags         	Pharmacy
// @Accept       	json
// @Produce     	json
// @Param			id	path		int	true	"Pharmacy ID"
// @Success      	200  "pharmacy deleted"
// @Failure      	400  "ID is required"
// @Failure      	404  "Could not delete pharmacy"
// @Router       	/delete_pharmacy/{id} [delete]
func (r *Repository) DeletePharmacy(context *fiber.Ctx) error {
	pharmacyModel := &[]models.Pharmacy{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "id is required",
			})
		return nil
	}

	result := r.DB.Delete(pharmacyModel, id)

	if result.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete pharmacy"})
		return result.Error
	}

	if result.RowsAffected == 0 {
		// No rows were affected, meaning the record with the given ID was not found
		context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "pharmacy not found for deletion"})
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "pharmacy deleted"})

	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_pharmacy", r.CreatePharmacy)
	api.Delete("delete_pharmacy/:id", r.DeletePharmacy)
	api.Get("/get_pharmacy/:id", r.GetPharmacyByID)
	api.Get("/get_pharmacies", r.GetPharmacies)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not connect to the database")
	}

	err = models.MigratePharmacy(db)
	if err != nil {
		log.Fatal("Could not migrate the database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	r.SetupRoutes(app)
	err = app.Listen(":8080")
	if err != nil {
		log.Fatalf("fiber.Listen failed %s", err)
	}

}
