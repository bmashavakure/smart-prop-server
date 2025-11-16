package auth_utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// claims struct for generating jwt tokens
type JWTClaims struct {
	DateTime  string `json:"date_time"`
	UserEmail string `json:"user_email"`

	jwt.RegisteredClaims
}

func HashPassword(pwd string) string {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Println("Error loading .env file")
	//}

	hashingCost := os.Getenv("HashCost")
	intCost, intErr := strconv.Atoi(hashingCost)
	if intErr != nil {
		log.Println("Failed to convert string to int")
	}

	pwdHash, hashErr := bcrypt.GenerateFromPassword([]byte(pwd), intCost)
	if hashErr != nil {
		fmt.Printf("Error trying to hash password: %s\n", hashErr)
	}
	return string(pwdHash)
}

// compare user password and stored hash
func ComparePasswordAndHash(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// hash strings for custom token
func HashString(input string) string {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Println("Error loading .env file")
	//}
	salt := os.Getenv("HASH_SALT")
	byteInput := []byte(input + salt)

	md5hash := md5.Sum(byteInput)

	return hex.EncodeToString(md5hash[:])
}

// generate jwt token func
func GenerateJWTToken(dateTime string, userEmail string) (string, error) {
	//err := godotenv.Load()
	//if err != nil {
	//	fmt.Printf("Error loading .env file")
	//	return "", err
	//}

	jwtKey := []byte(os.Getenv("JWT_KEY"))
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaims{
		DateTime:  dateTime,
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "smart-prop-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
