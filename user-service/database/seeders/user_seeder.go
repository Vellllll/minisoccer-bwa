package seeders

import (
	"user-service/constants"
	"user-service/domain/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("arvel123"), bcrypt.DefaultCost)

	user := models.User{
		UUID:        uuid.New(),
		Name:        "Arvel",
		Username:    "arvel",
		Password:    string(password),
		PhoneNumber: "08273648237",
		Email:       "arvel@gmail.com",
		RoleID:      constants.Admin,
	}

	err := db.FirstOrCreate(&user, models.User{Username: user.Username}).Error
	if err != nil {
		logrus.Errorf("failed to seed user %v", err)
		panic(err)
	}
	logrus.Printf("seeding user %s", user.Username)
}
