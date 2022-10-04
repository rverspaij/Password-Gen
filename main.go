package main

import (
	"database/sql"
	"flag"
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
	readConfig(&config)
	connectDB(&config)
	createTable()
	checkExistense(password)
	addPass(password)
}

func errorHandler(error) {

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
		log.Fatal("empty fields")
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
	log.Printf(password)
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

func checkExistense(content string) bool {
	query := `SELECT * FROM password WHERE EXISTS (content);`
	_, err := database.Exec(query)
	if err != nil {
		log.Fatal("no")
		return false
	} else {
		return true
	}
}

func addPass(content string) error {
	query := `INSERT INTO password (content)
	VALUES($1)
;`
	_, err := database.Exec(query, content)
	return err
}
