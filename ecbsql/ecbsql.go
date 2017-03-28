package ecbsql

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/farhan-shahid/exchangerates"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const schema = `
CREATE TABLE ExchangeRate (
    date DATE,
    fromCurr CHAR(3),
    toCurr CHAR(3),
	rate DECIMAL(10,5)
);`
const db = "ExchangeDB"
const dbuser = "root"
const passEnv = "MYSQLPASS"

// Store fetches and stores historical currency exchange data from ecb.europa.eu
type Store struct {
	db *sqlx.DB
}

// New returns a new instance of Store
func New() *Store {
	s := &Store{}
	s.fetchData()
	return s
}

// GetExchangeRate returns exchange rate from the ecb dataset.
// Use from and to for specifying the currencies to convert between
// and date to specify the date of conversion
func (s *Store) GetExchangeRate(from, to string, date string) (float64, error) {
	fromVal, err := s.lookup(from, date)
	if err != nil {
		return 0, err
	}

	toVal, err := s.lookup(to, date)
	if err != nil {
		return 0, err
	}

	rate := calcRate(fromVal, toVal)
	return rate, nil
}

// GetMonthExchangeRates returns a list of exchange rate values for the month specified
func (s *Store) GetMonthExchangeRates(from, to string, year, month int) ([]exchangerates.DateRate, error) {
	rates := make([]exchangerates.DateRate, 0, 31)
	for i := 1; i <= 31; i++ {
		date := strconv.Itoa(year) + "-" + fmt.Sprintf("%02d", month) + "-" + fmt.Sprintf("%02d", i)

		fromVal, err := s.lookup(from, date)
		if err != nil {
			continue
		}
		toVal, err := s.lookup(to, date)
		if err != nil {
			continue
		}

		rate := calcRate(fromVal, toVal)
		t, _ := time.Parse("2006-01-02", date)
		rates = append(rates, exchangerates.DateRate{Rate: rate, Date: t})
	}
	if len(rates) == 0 {
		return nil, errors.New("No data exists")
	}
	return rates, nil
}

func (s *Store) fetchData() (err error) {
	connStr := fmt.Sprintf("%s:%s@/%s", dbuser, os.Getenv("MYSQLPASS"), db)
	s.db, err = sqlx.Connect("mysql", connStr)
	if err != nil {
		return err
	}

	var exists int
	err = s.db.Get(&exists, `SELECT count(*) FROM information_schema.TABLES WHERE (TABLE_SCHEMA = ?) AND (TABLE_NAME = ExchangeRate)`, db)
	if err != nil {
		return err
	}

	if exists != 0 {
		return nil //table already exists so do not fetch data
	}
	_, err = s.db.Exec(schema)
	if err != nil {
		return err
	}

	resp, err := http.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type curr struct {
		Currency string  `xml:"currency,attr"`
		Rate     float64 `xml:"rate,attr"`
	}
	type dateList struct {
		Date string `xml:"time,attr"`
		Curr []curr `xml:"Cube"`
	}
	type data struct {
		Rates []dateList `xml:"Cube>Cube"`
	}
	d := &data{}

	err = xml.Unmarshal(body, d)
	if err != nil {
		return err
	}

	tx := s.db.MustBegin()
	for _, i := range d.Rates {
		for _, j := range i.Curr {

			tx.Exec(`INSERT INTO ExchangeRate (fromCurr, toCurr, date, rate) VALUES (?, ?, ?, ?)`, "EUR", j.Currency, i.Date, j.Rate)
		}
		tx.Exec(`INSERT INTO ExchangeRate (fromCurr, toCurr, date, rate) VALUES (?, ?, ?, ?)`, "EUR", "EUR", i.Date, 1)
	}
	err = tx.Commit()
	return err
}

func (s *Store) lookup(curr string, date string) (float64, error) {
	var value float64
	err := s.db.Get(&value, `SELECT rate FROM ExchangeRate WHERE date=? AND fromCurr="EUR" AND toCurr=?`, date, curr)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func calcRate(fromVal, toVal float64) (rate float64) {
	rate = toVal / fromVal
	rateStr := fmt.Sprintf("%.5f", rate) //reduce precision
	rate, _ = strconv.ParseFloat(rateStr, 5)
	return
}
