package libs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var BlockedByGoogle = false

func CheckQuery(url string, delay int, c chan *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()

	if !BlockedByGoogle {
		resp := HttpRequest(url)
		if resp.StatusCode == 429 {
			BlockedByGoogle = true
			log.Error().Msg("[!] Google is blocking you for making too many requests. Make sure the requests are rotated properly! Whatever you are doing is probably being fingerprinting...Please try to increase the delay to bypass the bot detection and mitigation system !")
		} else if resp.StatusCode != 200 {
			log.Debug().Msg("[!] HTTP " + fmt.Sprint(resp.StatusCode) + " Error !")
		}
		c <- resp
	}
}

func ProcessQuery(pattern string, rawonly bool, urls []string, delay int, maxworkers int, format string, linksonly bool) []string {
	c := make(chan *http.Response, maxworkers)
	var wg sync.WaitGroup
	var dataResults []string

	for _, url := range urls {
		wg.Add(1)
		log.Debug().Msg("[+] Processing '" + url + "'...")
		go CheckQuery(url, delay, c, &wg)
		time.Sleep(time.Duration(delay) * time.Second)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for resp := range c {
		if rawonly {
			rawData := GetHttpResponse(resp)
			dataResults = ExtractValues(pattern, rawData)
			for i := range dataResults {
				dataResults[i] += "\n"
			}
		} else {
			jsonData, txtData := GenResults(ParseHtml(resp), pattern, format, linksonly)
			if len(jsonData) > 0 {
				jsonResults, err := json.Marshal(jsonData)
				CheckErrors(err)
				dataResults = append(dataResults, string(jsonResults))
			} else {
				dataResults = append(dataResults, strings.Join(txtData, ""))
			}
		}
	}
	return dataResults
}

func GenQueries(dork string, maxpages int, resultsperpage int, proxyurl string) []string {
	var dorksQueries []string

	for count := 0; count <= maxpages; count += resultsperpage {
		dorkQuery := fmt.Sprintf("%s/search?q=%s&start=%d&num=%d", proxyurl, url.QueryEscape(dork), count, resultsperpage)
		dorksQueries = append(dorksQueries, dorkQuery)
	}
	return dorksQueries
}

func GenResults(results []map[string]interface{}, pattern string, format string, linksonly bool) ([]map[string]interface{}, []string) {
	var txtData []string
	var jsonData []map[string]interface{}

	for _, item := range results {
		if SearchPattern(pattern, item["Description"].(string)) || SearchPattern(pattern, item["Title"].(string)) || SearchPattern(pattern, item["Link"].(string)) {
			switch format {
			case "txt":
				if linksonly {
					txtData = append(txtData, string(item["Link"].(string)+"\n"))
				} else {
					txtData = append(txtData, fmt.Sprintf("Title: %s\n Link: %s\n Description: %s", item["Title"].(string), item["Link"].(string), item["Description"].(string)+"\n"))
				}
			case "json":
				if linksonly {
					link := map[string]interface{}{
						"Link": item["Link"].(string),
					}
					jsonData = append(jsonData, link)
				} else {
					jsonData = append(jsonData, item)
				}
			case "csv":
				if linksonly {
					txtData = append(txtData, string(`"`+item["Link"].(string)+`"`+"\n"))
				} else {
					txtData = append(txtData, `"`+item["Title"].(string)+`",`, `"`+item["Link"].(string)+`",`, `"`+item["Description"].(string)+`"`+"\n")
				}
			}
		}
	}
	return jsonData, txtData
}
