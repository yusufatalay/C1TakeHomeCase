package user

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yusufatalay/C1TakeHomeCase/assessment"
	"github.com/yusufatalay/C1TakeHomeCase/database"
	"github.com/yusufatalay/C1TakeHomeCase/userassessment"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time `json:"-"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
}

func init() {
	database.DBConn.AutoMigrate(&User{})
}

// CreateUser creates the user with given information (FirstName, LastName, email) this informations are gathered from icoming request's body (POST)
func CreateUser(c *fiber.Ctx) error {
	// unmarshall the reques's body into the payload (User object (kind of...))
	var payload User
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	// write user to the database
	database.DBConn.Create(&payload)
	//
	err := c.Status(http.StatusCreated).JSON(payload)

	if err != nil {
		log.Println(err)
	}

	// return a response according to the resul of the process
	response := struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
		Type    bool        `json:"type"`
	}{
		Data:    payload,
		Message: "User has created and inserted to database",
		Type:    true,
	}

	return c.JSON(response)
}

// AfterCreate is a trigger which gets triggered after creation of a User instance,
// This particular trigger will create an assessment and a new entry at UserAssessments table
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	// new assessment has created
	assessmentuuid, err := assessment.CreateAssessment()
	if err != nil {
		return err
	}

	// create new userassessment entry
	err = userassessment.CreateUserAssessment(int(u.ID), assessmentuuid)

	return err

}
