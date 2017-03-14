package googlefinance

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// Store fetches currency exchange data from google.com/finance/converter
type Store struct {
}

// New returns a new instance of Store
func New() *Store {
	return &Store{}
}

// GetExchangeRate returns exchange rate from google.com/finance/converter
// Use from and to for specifying the currencies to convert between
// and date to specify the date of conversion
func (s *Store) GetExchangeRate(from, to string, date string) (float64, error) {
	resp, err := http.Get("https://www.google.com/finance/converter?a=1&from=" + from + "&to=" + to)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return 0, err
	}

	node, found := getNodeByAttr(root, "id", "currency_converter_result")
	if !found {
		return 0, errors.New("Data fetching failed")
	}
	node, found = getNodeByAttr(node, "class", "bld")
	if !found {
		return 0, errors.New("currency " + from + " or " + to + " not found")
	}

	rate, err := strconv.ParseFloat(strings.Split(node.FirstChild.Data, " ")[0], 5)
	if err != nil {
		return 0, errors.New("Data fetching failed")
	}

	return rate, nil
}

func getNodeByAttr(n *html.Node, attrName, val string) (*html.Node, bool) {
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			if attr.Val == val {
				return n, true
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, found := getNodeByAttr(c, attrName, val)
		if found {
			return result, found
		}
	}
	return nil, false
}
