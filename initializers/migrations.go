package initializers

import "instawham/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
