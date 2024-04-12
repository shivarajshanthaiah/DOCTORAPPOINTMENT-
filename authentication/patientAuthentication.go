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

// Generating jwt token for patient
func GeneratePatientToken(phone string) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.PatientClaims{

		Phone:          phone,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)

}

func AuthenticatePatient(signedStringToken string) (string, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(signedStringToken, &models.PatientClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	//type assert the claims from the token object
	if claims, ok := token.Claims.(*models.PatientClaims); ok && token.Valid {
		return claims.Phone, nil
	}

	return "", errors.New("invalid token")
}

func PatientAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		patientsToken := c.GetHeader("Authorization")

		if patientsToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(patientsToken, "Bearer"))

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
