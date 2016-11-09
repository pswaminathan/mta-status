package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var (
	th       = `<table><tr><th>Line Name</th><th>Status</th><th>Info</th></tr>`
	tf       = `</table>`
	port     int
	certPath string
	keyPath  string
)

func init() {
	flag.IntVar(&port, "port", 443, "Port for server")
	flag.StringVar(&certPath, "cert", "/etc/letsencrypt/live/pswa.me/fullchain.pem", "SSL certificate Location")
	flag.StringVar(&keyPath, "key", "/etc/letsencrypt/keys/0000_key-certbot.pem", "SSL key Location")
}

type Line struct {
	Name   string `xml:"name"`
	Status string `xml:"status"`
	Text   string `xml:"text"`
	Date   string `xml:"Date"`
	Time   string `xml:"Time"`
}

type Result struct {
	XMLName      xml.Name `xml:"service"`
	ResponseCode int      `xml:"responsecode"`
	Timestamp    string   `xml:"timestamp"`
	Subway       []Line   `xml:"subway>line"`
	Bus          []Line   `xml:"bus>line"`
	BT           []Line   `xml:"BT>line"`
	LIRR         []Line   `xml:"LIRR>line"`
	MetroNorth   []Line   `xml:"MetroNorth>line"`
}

func main() {
	flag.Parse()
	u, err := url.Parse("http://web.mta.info")
	if err != nil {
		log.Fatal(err)
	}
	rp := NewReverseProxy(u)
	http.Handle("/status/serviceStatus.txt", rp)
	http.HandleFunc("/service/", getServiceData)
	log.Fatal(http.ListenAndServeTLS(
		":"+strconv.Itoa(port),
		certPath,
		keyPath,
		nil,
	))
}

func getServiceData(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Path[9:]
	body, _ := getMTAData()
	defer body.Close()
	var res Result
	xml.NewDecoder(body).Decode(&res)
	buf := generateTable(service, &res)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	buf.WriteTo(w)
}

func getMTAData() (io.ReadCloser, error) {
	resp, err := http.Get("http://web.mta.info/status/serviceStatus.txt")
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func generateTable(service string, result *Result) *bytes.Buffer {
	buf := new(bytes.Buffer)
	lines := getLines(service, result)
	if lines == nil {
		return buf
	}
	buf.WriteString(th)
	for _, line := range lines {
		buf.WriteString("<tr>")
		buf.WriteString("<td class=\"name\">")
		buf.WriteString(line.Name)
		buf.WriteString("</td>")
		buf.WriteString("<td class=\"status\">")
		buf.WriteString(line.Status)
		buf.WriteString("</td>")
		buf.WriteString("<td class=\"text\">")
		buf.WriteString(line.Text)
		buf.WriteString("</td>")
		buf.WriteString("</tr>")
	}
	buf.WriteString(tf)
	return buf
}

func getLines(service string, result *Result) []Line {
	var lines []Line
	switch service {
	case "subway":
		lines = result.Subway
	case "bus":
		lines = result.Bus
	case "Subway":
		lines = result.Subway
	case "Bus":
		lines = result.Bus
	case "BT":
		lines = result.BT
	case "LIRR":
		lines = result.LIRR
	case "MetroNorth":
		lines = result.MetroNorth
	}
	return lines
}
