package dataloaders

import (
	"log"
	"strconv"
	"github.com/jinzhu/gorm"
)

const PerPage = 250

// Fetcher is an interface for fetching JSON data
type Fetcher interface {
	BaseUrl() string
	ParseResults(body []byte) (int, error)
	CreateOrSave(db *gorm.DB, index int) error
}

// Fetches all pages using a fetcher
func FetchAll(ch chan<- bool, db *gorm.DB, f Fetcher) {
	log.Println("Started:", f)
	total := 0
	page := 1
	for {
		c, _ := Fetch(db, f, page)
		if c < 1 {
			break
		}
		total += c
		page += 1
	}
	log.Println(f,"Finished:", total, "total")
	ch <- true
}

// Fetches one page with the help of a Fetcher
func Fetch(db *gorm.DB, f Fetcher, page int) (int, error) {
	url := f.BaseUrl()  + "&page=" + strconv.Itoa(page) + "&per_page=" + strconv.Itoa(PerPage)
	body, err := ReadUrl(url)

	if err != nil {
		return 0, err
	} 

	c, err := f.ParseResults(body)

	if err != nil {
		return 0, err
	}

	tx := db.Begin()

	// log.Println("COUNT: ", c)
	for i := 0; i < c; i++ {
		err = f.CreateOrSave(db, i)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return c, nil
}