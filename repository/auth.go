package repository

import (
	"time"

	"github.com/GilangAndhika/bukuin_be/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.Users) error {
	// mengenkripsi password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// jika role tidak diisi, atur nilainya ke 2
	if user.IdRole == 0 {
		user.IdRole = 2
	}

	// Simpan user ke database
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserByUsername(db *gorm.DB, username string) (*models.Users, error) {
	var user models.Users
	// Cari user berdasarkan username
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByID(db *gorm.DB, IdUser uint) (*models.Users, error) {
	var user models.Users
	// Cari user berdasarkan ID
	result := db.First(&user, IdUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func CreateToken(user *models.Users) (string, error) {
	claims := &models.JWTClaims{
		IdUser: user.IdUser,
		IdRole: user.IdRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(), // token kadaluarsa dalam 1 jam
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
