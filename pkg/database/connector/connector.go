package connector

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connector() {
	//err := godotenv.Load()
	//if err != nil {
	//	fmt.Println("Error loading .env file")
	//}
	dbUrl := os.Getenv("DATABASE_URL")

	db, sqlErr := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if sqlErr != nil {
		panic(sqlErr)
	}

	DB = db
}
