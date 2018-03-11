package authentication

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	"dorian/http-gui/src/api.jwt.auth/core/redis"
	"dorian/http-gui/src/api.jwt.auth/services/models"
	"dorian/http-gui/src/api.jwt.auth/settings"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var authBackendInstance *JWTAuthenticationBackend = nil

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

func (backend *JWTAuthenticationBackend) GenerateToken(userUUID string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix(),
		"iat": time.Now().Unix(),
		"sub": userUUID,
	}
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

func (backend *JWTAuthenticationBackend) Authenticate(user *models.User) bool {
	var username string
	var password string
	var enabled bool

	db := settings.GetDB()
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", user.Username)
	err := row.Scan(&username, &password, &enabled)
	if err != nil {
		// Following temporary code is only for testing: whenever the username does not exist, then insert it
		// into db using its password but in hashed form
		// TODO: Must be deleted later
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		_, err := db.Exec(
			"INSERT INTO gral_usr (titulo, password, enabled) VALUES ($1, $2, $3)",
			user.Username,
			string(hashedPassword),
			true,
		)
		if err != nil {

		} // end of temporary code
		return false
	}

	if !enabled {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		return false
	}
	return true
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	redisConn := redis.Connect()
	claims := token.Claims.(jwt.MapClaims)
	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(claims["exp"]))
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisConn := redis.Connect()
	redisToken, _ := redisConn.GetValue(token)

	if redisToken == nil {
		return false
	}

	return true
}

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(settings.Get().PrivateKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(settings.Get().PublicKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}
