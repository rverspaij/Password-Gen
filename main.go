package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
	"gopkg.in/yaml.v3"
)

var database *sql.DB
var length int
var digits bool
var symbols bool
var lower bool
var upper bool

type Config struct {
	Server struct {
		DbName string `yaml:"host"`
	} `yaml:"server"`

	Db struct {
		DbUser string `yaml:"user"`
		DbPass string `yaml:"pass"`
	} `yaml:"db"`
}

func init() {
	flag.IntVar(&length, "l", 0, "Fill in length password!")
	flag.BoolVar(&digits, "digits", true, "Do you want digits?")
	flag.BoolVar(&symbols, "symbols", true, "Do you want digits?")
	flag.BoolVar(&lower, "lower", true, "Do you want digits?")
	flag.BoolVar(&upper, "upper", true, "Do you want digits?")
	flag.Parse()
	if length < 1 {
		log.Fatal("Make sure your passlength is greater than 0!")
	}
}

func main() {
	var config Config

	password, err := genPassword(length, digits, symbols, lower, upper)
	errorHandler(err)

	err = readConfig(&config)
	errorHandler(err)

	err = connectDB(&config)
	errorHandler(err)

	err = createTable()
	errorHandler(err)

	passCheck, err := checkExistense(password)

	if passCheck {
		err := errors.New("password already exists")
		errorHandler(err)
	} else {
		if err := addPass(password); err != nil {
			errorHandler(err)
		}
	}
}

func errorHandler(err error) {
	file, err1 := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err1 != nil {
		log.Fatal(err1)
	}
	defer file.Close()

	logger := log.New(file, "Custom Log", log.LstdFlags)
	if err != nil {
		logger.Fatal(err)
	}
}

func readConfig(cfg *Config) error {
	conf, err := os.ReadFile("conf.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(conf), &cfg)
	if err != nil {
		return err
	}
	if cfg.Db.DbPass == "" || cfg.Db.DbUser == "" || cfg.Server.DbName == "" {
		err := errors.New("empty fields")
		return err
	}
	return err
}

func genPassword(passLength int, digits bool, symbols bool, lower bool, upper bool) (string, error) {
	rand.Seed(time.Now().UnixNano())
	passDigits := rand.Intn(0 + passLength)
	characters := passLength - passDigits
	passSymbols := rand.Intn(0 + characters)
	password, err := password.Generate(passLength, passDigits, passSymbols, lower, upper)
	if err != nil {
		return password, err
	}
	return password, err
}

func connectDB(cfg *Config) error {
	db, err := sql.Open("postgres", "dbname="+cfg.Server.DbName+" user="+cfg.Db.DbUser+" password="+cfg.Db.DbPass+" sslmode=disable")
	if err != nil {
		return err
	}

	database = db
	return err
}

func createTable() error {
	query := `CREATE TABLE IF NOT EXISTS password (
		id		SERIAL	PRIMARY KEY,
		content	TEXT	NOT NULL
	);`
	_, err := database.Exec(query)
	return err
}

func checkExistense(content string) (bool, error) {
	var check bool
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM password WHERE content = '%s')`, content)
	err := database.QueryRow(query).Scan(&check)
	if err != nil {
		return check, err
	}
	fmt.Println(check)
	return check, nil
}

func addPass(content string) error {
	query := `INSERT INTO password (content)
	VALUES($1)
;`
	_, err := database.Exec(query, content)
	fmt.Println("Password is added")
	return err
}
