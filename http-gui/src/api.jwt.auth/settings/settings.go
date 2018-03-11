package settings

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
)

var environments = map[string]string{
	"production":    "/home/j2eeserver/maxima/config/enviroments/prod.json",
	"preproduction": "/home/dorian/go/src/dorian/resources/profiles/default.json",
	"tests":         "/home/j2eeserver/maxima/config/enviroments/tests.json",
}

type Settings struct {
	PrivateKeyPath     string
	PublicKeyPath      string
	JWTExpirationDelta int
	DbHost             string
	DbPort             string
	DbName             string
	DbUser             string
	DbPasswd           string
}

var settings = Settings{}
var env = "preproduction"
var db *sql.DB

func Init() {
	env = os.Getenv("GO_ENV")
	if env == "" {
		fmt.Println("Warning: Setting preproduction environment due to lack of GO_ENV value")
		env = "preproduction"
	}
	LoadSettingsByEnv(env)
}

func LoadSettingsByEnv(env string) {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)
	}

	jsonErr := json.Unmarshal(content, &settings)
	if jsonErr != nil {
		fmt.Println("Error while parsing config file", jsonErr)
	}
}

func GetEnvironment() string {
	return env
}

func Get() Settings {
	if &settings == nil {
		Init()
	}
	return settings
}

func IsTestEnvironment() bool {
	return env == "tests"
}

func GetDB() *sql.DB {
	if db == nil {
		var err error
		connStr := "host=" + settings.DbHost + " port=" + settings.DbPort + " dbname=" + settings.DbName + " user=" + settings.DbUser + " password=" + settings.DbPasswd + " sslmode=require"
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
	}
	return db
}
