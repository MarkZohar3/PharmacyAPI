package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mod/models"
	"go.mod/storage"
	"gorm.io/gorm"
)

type Pharmacy struct {
	Owner   string `json:"owner"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Repository struct {
	DB *gorm.DB
}

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

func (r *Repository) GetPharmacyByID(context *fiber.Ctx) error {
	id := context.Params("id")
	pharmacyModel := &models.Pharmacy{}
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "id is required"})
		return nil
	}

	fmt.Println("the Id is: ", id)

	err := r.DB.Where("id = ?", id).First(pharmacyModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get pharmacy"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "pharmacy fetched",
		"data":    pharmacyModel,
	})

	return nil
}

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

	err := r.DB.Delete(pharmacyModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete pharmacy"})
		return err.Error
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
	r.SetupRoutes(app)
	app.Listen(":8080")

}
