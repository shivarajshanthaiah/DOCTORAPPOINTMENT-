package authentication

import (
	"doctorAppointment/models"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("secretKey")

func GeneratePatienttoken(phone string) (string, error) {
	claims := jwt.MapClaims{
		"phone": phone,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func AuthenticatePatient(signedStringToken string) (string, error) {
	// Parse the token
	var patientClaims models.PatientClaims
	token, err := jwt.ParseWithClaims(signedStringToken, &patientClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return "", err
	}

	// Validate the token
	if !token.Valid {
		return "", errors.New("token is not valid")
	}

	//type assert the claims from the token object
	claims, ok := token.Claims.(*models.PatientClaims)

	if !ok {
		err = errors.New("could't parse claims")
		return "", err
	}
	phone := claims.Phone

	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("token expired")
		return "", err
	}

	return phone, nil
}

func PatientAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		patientsToken := c.GetHeader("Authorization")
		if patientsToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.Replace(patientsToken, "Bearer ", "", 1)

		phone, err := AuthenticatePatient(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// If authentication is successful, set the phone number in the request context
		// so that handlers can access it if needed.
		c.Set("phone", phone)

		c.Next()
	}
}
