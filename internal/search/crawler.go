package search

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type CrawlData struct {
	Url          string
	Success      bool
	ResponseCode int
	CrawlData    ParsedBody
}

type ParsedBody struct {
	CrawlTime       time.Duration
	PageTitle       string
	PageDescription string
	Heading         string
	Links           Links
}

type Links struct {
	Internal []string
	External []string
}

func runCrawl(inputUrl string) CrawlData {
	resp, err := http.Get(inputUrl)
	baseUrl, err := url.Parse(inputUrl)

	if err != nil || resp == nil {
		fmt.Println("somethign went wrong fetching body")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: 0, CrawlData: ParsedBody{}}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("code 200 not found")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}
	}

	// Check HTML
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		// Response is HTML
		data, err := parseBody(resp.Body, baseUrl)
		if err != nil {
			return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}
		}
		return CrawlData{Url: inputUrl, Success: true, ResponseCode: resp.StatusCode, CrawlData: data}
	} else {
		fmt.Println("html response not found")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}
	}
}

func parseBody(body io.Reader, baseUrl *url.URL) (ParsedBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return ParsedBody{}, err
	}
	start := time.Now()

	// Get Links
	links := getLinks(doc, baseUrl)
	// Get Title and Desciption
	title, desc := getPageData(doc)
	// Get H1 Tags
	heading := getPageHeading(doc)

	end := time.Now()
	return ParsedBody{
		CrawlTime:       end.Sub(start),
		PageTitle:       title,
		PageDescription: desc,
		Heading:         heading,
		Links:           links,
	}, nil
}

// DFS (depth first search) Recursive function for scanning HTML tree
func getLinks(n *html.Node, baseUrl *url.URL) Links {
	links := Links{}
	if n == nil {
		return links
	}

	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url, err := url.Parse(attr.Val)
					if err != nil || strings.HasPrefix(url.String(), "#") || strings.HasPrefix(url.String(), "mail") || strings.HasPrefix(url.String(), "tel") || strings.HasPrefix(url.String(), "javascript") || strings.HasPrefix(url.String(), ".pdf") || strings.HasPrefix(url.String(), ".md") {
						continue
					}
					if url.IsAbs() {
						if isSameHost(url.String(), baseUrl.String()) {
							links.Internal = append(links.Internal, url.String())
						} else {
							links.External = append(links.External, url.String())
						}
					} else {
						rel := baseUrl.ResolveReference(url)
						links.Internal = append(links.Internal, rel.String())
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			findLinks(child)
		}
	}
	findLinks(n)
	return links
}

func isSameHost(absoluteURL, baseURL string) bool {
	absURL, err := url.Parse(absoluteURL)
	if err != nil {
		return false
	}
	baseURLPased, err := url.Parse(baseURL)
	if err != nil {
		return false
	}
	return absURL.Host == baseURLPased.Host
}

func getPageData(n *html.Node) (string, string) {
	if n == nil {
		return "", ""
	}

	title, desc := "", ""
	var findMetaAndTitle func(*html.Node)

	findMetaAndTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			// Check if empty
			if n.FirstChild == nil {
				title = " "
			} else {
				title = n.FirstChild.Data
			}
		} else if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string
			for _, attr := range n.Attr {
				if attr.Key == "name" {
					name = attr.Val
				} else if attr.Key == "content" {
					content = attr.Val
				}
			}
			if name == "description" {
				desc = content
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		findMetaAndTitle(child)
	}
	findMetaAndTitle(n)
	return title, desc
}

func getPageHeading(n *html.Node) string {
	if n == nil {
		return ""
	}
	var heading strings.Builder
	var findH1 func(*html.Node)

	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			// Check if node is empty
			if n.FirstChild != nil {
				heading.WriteString(n.FirstChild.Data)
				heading.WriteString(", ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findH1(c)
		}
	}
	return strings.TrimSuffix(heading.String(), ",")
}
