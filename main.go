package main

import (
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"io/ioutil"
	"net/http"
	"os"
)

type Configuration struct {
	DbAddress  string `json:"DB_ADDRESS"`
	DbUser     string `json:"DB_USER"`
	DbPassword string `json:"DB_PASSWORD"`
	DbName     string `json:"DB_NAME"`
}
type BodyInput struct {
	FirstLetter      string `json:"firstLetter"`
	SecondLetter     string `json:"secondLetter"`
	ThirdLetter      string `json:"thirdLetter"`
	FourthLetter     string `json:"fourthLetter"`
	FifthLetter      string `json:"fifthLetter"`
	UsedCharacters   string `json:"usedCharacters"`
	UnusedCharacters string `json:"unusedCharacters"`
}

type Words struct {
	Id   int
	Word string
}

func main() {

	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	db := pg.Connect(&pg.Options{
		Addr:     configuration.DbAddress,
		User:     configuration.DbUser,
		Password: configuration.DbPassword,
		Database: configuration.DbName,
	})
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.SetTrustedProxies([]string{"*"})
	r.Use(cors.Default())

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST"},
	}))
	r.POST("/postPossibleWords", func(c *gin.Context) {
		body := c.Request.Body
		value, err := ioutil.ReadAll(body)

		if err != nil {
			panic(err)
		}
		var bodyConverted BodyInput
		err = json.Unmarshal(value, &bodyConverted)
		if err != nil {
			panic(err)
		}
		query := ""
		queryLength := len(query)
		if len(bodyConverted.FirstLetter) > 0 && len(bodyConverted.FirstLetter) < 2 {
			query += " SUBSTRING(word, 1, 1) = '" + bodyConverted.FirstLetter + "'"
		}
		if len(bodyConverted.SecondLetter) > 0 && len(bodyConverted.SecondLetter) < 2 {
			if len(query) > queryLength {
				query += " AND "
			}
			query += " SUBSTRING(word, 2, 1) = '" + bodyConverted.SecondLetter + "'"
		}
		if len(bodyConverted.ThirdLetter) > 0 && len(bodyConverted.ThirdLetter) < 2 {
			if len(query) > queryLength {
				query += " AND "
			}
			query += " SUBSTRING(word, 3, 1) = '" + bodyConverted.ThirdLetter + "'"
		}
		if len(bodyConverted.FourthLetter) > 0 && len(bodyConverted.FourthLetter) < 2 {
			if len(query) > queryLength {
				query += " AND "
			}
			query += " SUBSTRING(word, 4, 1) = '" + bodyConverted.FourthLetter + "'"
		}
		if len(bodyConverted.FifthLetter) > 0 && len(bodyConverted.FifthLetter) < 2 {
			if len(query) > queryLength {
				query += " AND "
			}
			query += " SUBSTRING(word, 5, 1) = '" + bodyConverted.FifthLetter + "'"
		}
		if len(bodyConverted.UnusedCharacters) > 0 {
			for _, char := range bodyConverted.UnusedCharacters {
				if len(query) > queryLength {
					query += " AND "
				}
				query += " word NOT LIKE '%" + string(char) + "%'"
			}
		}
		if len(bodyConverted.UsedCharacters) > 0 {
			for _, char := range bodyConverted.UsedCharacters {
				if len(query) > queryLength {
					query += " AND "
				}
				query += " word LIKE '%" + string(char) + "%'"
			}
		}

		if len(query) > queryLength {
			query = " AND " + query
		}
		query = "SELECT word FROM public.words WHERE LENGTH(word) = 5 " + query + ";"

		// Select all words.
		var _words []string
		_, err = db.Query(&_words, query)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"return": _words})
	})

	r.Run(":3526")

}
