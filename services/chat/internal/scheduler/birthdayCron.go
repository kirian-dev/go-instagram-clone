package scheduler

import (
	"encoding/base64"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/helpers"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func RunBirthdayCron(db *gorm.DB, log *logger.ZapLogger, cfg *config.Config) error {
	c := cron.New()

	_, err := c.AddFunc("0 0 * * *", func() {
		// Check for users with birthdays today
		usersWithBirthdays, err := getUsersWithBirthdaysToday(db)
		if err != nil {
			log.Error("Error retrieving users with birthdays:", err)
			return
		}

		for _, user := range usersWithBirthdays {
			convertUser := helpers.ConvertToResponseUser(user)
			if err := createBirthdayEmail(convertUser, log, cfg); err != nil {
				log.Error("Error creating birthday email:", err)
			}
		}
	})
	if err != nil {
		log.Fatal("Error scheduling birthday check:", err)
	}

	c.Start()

	log.Info("Birthday cron job started successfully")
	select {}
}

func getUsersWithBirthdaysToday(db *gorm.DB) ([]*models.User, error) {
	today := time.Now()
	formattedToday := today.Format("2006-01-02")

	var users []*models.User
	if err := db.Where("birthday::date = ?", formattedToday).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func createBirthdayEmail(user *models.UserResponse, log *logger.ZapLogger, cfg *config.Config) error {
	var firstName = user.FirstName
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[0]
	}
	imagePath := "public/img/birthday.png"
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		log.Error(err)
		return err
	}
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
	imageName := "birthday-image"

	emailData := helpers.EmailData{
		URL:       cfg.ClientOrigin + "/dashboard",
		FirstName: firstName,
		Subject:   "Happy Birthday",
	}

	helpers.SendEmail(user, &emailData, "birthday.html", cfg, log, imageBytes, imageBase64, imageName)
	log.Info("Birthday email sent successfully to user:", user.Email)

	return nil
}
