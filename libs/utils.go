package libs

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
const RequestTimeout = 30

func CheckErrors(err error) {
	if err != nil {
		log.Error().Err(err).Msg("[!] An error occured !")
	}
}

func OutFile(results string, filepath string) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	CheckErrors(err)
	if err == nil {
		_, err = f.Write([]byte(results))
		CheckErrors(err)
	}
	defer f.Close()
}

func OutResults(dataresults []string, format string, linksonly bool, rawonly bool) string {
	var results string

	if rawonly {
		format = "txt"
	}

	switch format {
	case "txt":
		results = strings.Join(dataresults, "")
	case "json":
		results = strings.Join(dataresults, "")
	case "csv":
		if linksonly {
			results = "\"Link\"" + "\n"
		} else {
			results = "\"Title\",\"Link\",\"Description\"" + "\n"
		}
		results = results + strings.Join(dataresults, "")
	}

	return results
}

func SearchPattern(regex string, content string) bool {
	pattern := regexp.MustCompile(regex)
	match := pattern.MatchString(content)
	return match
}

func ExtractValues(regex string, content string) []string {
	var uniqMatches []string
	pattern := regexp.MustCompile(regex)
	matches := pattern.FindAllString(content, -1)
	for _, match := range matches {
		isDuplicate := false
		for _, uniqueMatch := range uniqMatches {
			if match == uniqueMatch {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			uniqMatches = append(uniqMatches, match)
		}
	}
	return uniqMatches
}

func HttpRequest(url string) *http.Response {
	cli := &http.Client{
		Timeout: RequestTimeout * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	CheckErrors(err)

	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := cli.Do(req)
	CheckErrors(err)
	return resp
}

func ParseHtml(resp *http.Response) []map[string]interface{} {
	titleLinkDescList := make([]map[string]interface{}, 0)
	body, err := html.Parse(resp.Body)
	CheckErrors(err)
	var links []string
	var titles []string
	var descs []string
	var node func(*html.Node)
	node = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h3":
				titles = append(titles, GetTextNode(n))
				fallthrough
			case "a":
				for _, a := range n.Attr {
					if a.Key == "href" && (strings.HasPrefix(a.Val, "http://") || strings.HasPrefix(a.Val, "https://")) && !strings.Contains(a.Val, ".google.") {
						links = append(links, a.Val)
					}
				}
			}

			if n.Data == "div" {
				for _, d := range n.Attr {
					if d.Key == "style" && strings.HasPrefix(d.Val, "-webkit-line-clamp:2") {
						descs = append(descs, GetTextNode(n))
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			node(c)
		}
	}
	node(body)

	nblinks := len(links)
	nbtitles := len(titles)
	nbdescs := len(descs)

	minLength := nbtitles
	if nblinks < minLength {
		minLength = nblinks
	}
	if nbdescs < minLength {
		minLength = nbdescs
	}

	if nblinks == 0 || nbtitles == 0 || nbdescs == 0 {
		log.Error().Err(err).Msg("[!] HTML parsing error: title or link or description list empty !")
	} else {
		for idx := 0; idx < minLength; idx++ {
			ntitleLinkDesc := map[string]interface{}{
				"Title":       titles[idx],
				"Link":        links[idx],
				"Description": descs[idx],
			}
			titleLinkDescList = append(titleLinkDescList, ntitleLinkDesc)
		}
	}
	return titleLinkDescList
}

func GetTextNode(n *html.Node) string {
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
		}
	}
	return text
}

func GetHttpResponse(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	CheckErrors(err)
	return string(body)
}
