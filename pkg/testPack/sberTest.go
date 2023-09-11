package testPack

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

var infoLog = log.New(os.Stdout, "\033[33mINFO\033[0m\t", log.Ldate|log.Ltime)
var errors = log.New(os.Stderr, "\033[31mERROR\u001B[0m\t", log.Ldate|log.Ltime)

const initialURL = "http://147.78.65.149/start"

func createClient() (client *http.Client, redirectedURL *url.URL) {
	var redirUrl string
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		errors.Panic(err)
	}
	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirUrl = req.URL.String()
			return nil
		},
		Jar: cookieJar,
	}

	u, err := url.Parse(initialURL)
	if err != nil {
		errors.Panic(err)
	}
	doReqGet(client, u)

	redirectedURL, err = url.Parse(redirUrl)
	if err != nil {
		errors.Panic(err)
	}

	return client, redirectedURL
}
func doReqGet(client *http.Client, redirectedURL *url.URL) *http.Response {
	req := &http.Request{
		URL:    redirectedURL,
		Header: make(http.Header),
		Method: "GET",
	}

	resp, err := client.Do(req)
	if err != nil {
		errors.Panic(err)
	}

	infoLog.Printf("Client:\033[32m%s\u001B[0m go to:%s", getSid(client, redirectedURL), redirectedURL)
	return resp
}
func doReqPost(client *http.Client, redirectedURL *url.URL) (*http.Response, *url.URL) {
	req2 := &http.Request{
		URL:    redirectedURL,
		Header: make(http.Header),
		Method: "POST",
	}
	resp, err := client.Do(req2)
	if err != nil {
		errors.Panic(err)
	}
	//defer resp.Body.Close()

	infoLog.Printf("Client:\033[32m%s\033[0m go to:%s", getSid(client, redirectedURL), redirectedURL)
	return resp, resp.Request.URL
}

func RunTest() {
	client, redirectedURL := createClient()
	dictionary := make(map[string]string, 0)

	resp := doReqGet(client, redirectedURL)

	passed := fillMap(resp, dictionary)

	for !passed {
		redirectedURL = createUrl(dictionary, redirectedURL)
		time.Sleep(3 * time.Second)

		resp, redirectedURL = doReqPost(client, redirectedURL)

		passed = fillMap(resp, dictionary)

	}

	infoLog.Printf("Client:\u001B[32m%s\u001B[0m DONE TEST", getSid(client, redirectedURL))

}

func createUrl(dict map[string]string, redirectUrl *url.URL) *url.URL {
	var ansUrl *url.URL
	var err error
	if len(dict) > 0 {
		var answer = "?"

		for key, val := range dict {
			answer += fmt.Sprintf("%s=%s&", key, val)
		}
		answer = strings.TrimSuffix(answer, "&")
		clearMap(dict)
		ansUrl, err = url.Parse(redirectUrl.String() + answer)
		if err != nil {
			errors.Panic(err)
		}
	}

	return ansUrl
}

func getSid(client *http.Client, u *url.URL) string {
	var res string
	cookie := client.Jar.Cookies(u)[0]
	res = cookie.Value
	return res
}

func fillMap(response *http.Response, di map[string]string) bool {
	var passed = false
	doc, err := html.Parse(response.Body)
	if err != nil {
		errors.Panic(err)
	}
	var findElements func(*html.Node)
	findElements = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "select":
				extractSelectInfo(n, di)
			case "input":
				extractInputInfo(n, di)
			case "h1":
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode && n.FirstChild.Data == "Passed" {
					passed = true
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findElements(c)

		}
	}

	findElements(doc)
	return passed
}

func clearMap(src map[string]string) {
	for key := range src {
		delete(src, key)
	}
}

func extractInputInfo(n *html.Node, res map[string]string) {

	if n.Type == html.ElementNode && n.Data == "input" {
		var name string
		var isRadio = false

		for _, attr := range n.Attr {

			if attr.Key == "type" && attr.Val == "radio" {
				isRadio = true
			}

			if attr.Key == "name" {
				name = attr.Val
				if _, ok := res[name]; !ok {
					res[attr.Val] = ""
				}
			}

			if attr.Key == "value" && isRadio {
				if val, _ := res[name]; len(val) < len(attr.Val) {
					res[name] = attr.Val
				}

			}

		}

		if !isRadio {
			if val, _ := res[name]; val == "" {
				res[name] = "test"
			}
		}

	}
}

func extractSelectInfo(n *html.Node, res map[string]string) {
	if n.Type == html.ElementNode && n.Data == "select" {

		var name string
		for _, attr := range n.Attr {
			if attr.Key == "name" {
				name = attr.Val
				if _, ok := res[name]; !ok {
					res[attr.Val] = ""
				}
				break
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "option" {
				for _, attr := range c.Attr {
					if attr.Key == "value" {
						if val, _ := res[name]; len(val) < len(attr.Val) {
							res[name] = attr.Val
						}
						break
					}
				}
			}
		}
	}

}
