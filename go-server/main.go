package main

import (
	"encoding/xml"
	"flag"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	port     int
	certPath string
	keyPath  string
	insecure bool
)

func init() {
	flag.IntVar(&port, "port", 443, "Port for server")
	flag.StringVar(&certPath, "cert", "/etc/letsencrypt/live/pswa.me/fullchain.pem", "SSL certificate Location")
	flag.StringVar(&keyPath, "key", "/etc/letsencrypt/keys/0000_key-certbot.pem", "SSL key Location")
	flag.BoolVar(&insecure, "insecure", false, "Run over http (useful for local testing)")
}

type Line struct {
	Name   string        `xml:"name"`
	Status string        `xml:"status"`
	Text   template.HTML `xml:"text"`
	Date   string        `xml:"Date"`
	Time   string        `xml:"Time"`
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

	addr := ":" + strconv.Itoa(port)
	if insecure {
		log.Fatal(http.ListenAndServe(addr, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(addr, certPath, keyPath, nil))
	}
}

func getServiceData(rw http.ResponseWriter, r *http.Request) {
	var res Result
	// TODO do some real error handling here
	_ = getMTAData(&res)

	service := r.URL.Path[9:]
	lines := getLines(service, &res)
	// TODO nil check here
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	serializeToHTML(rw, lines)
}

func getMTAData(data interface{}) error {
	resp, err := http.Get("http://web.mta.info/status/serviceStatus.txt")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return xml.NewDecoder(resp.Body).Decode(data)
}

func serializeToHTML(rw http.ResponseWriter, lines []Line) {
	t := template.Must(template.New("status_table").Parse(tableTemplate))
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(rw, lines)
}

func getLines(service string, result *Result) []Line {
	var lines []Line
	switch strings.ToLower(service) {
	case "subway":
		lines = result.Subway
	case "bus":
		lines = result.Bus
	case "bt":
		lines = result.BT
	case "lirr":
		lines = result.LIRR
	case "metronorth":
		lines = result.MetroNorth
	}
	return lines
}

var (
	tableTemplate = `<table class="status-table">
	<tr class="status-headers"><th>Line Name</th><th>Status</th><th>More Information</th></tr>
	{{range $i, $e := .}}<tr class="status-row" id="status-row-{{$i}}"><td class="line-name">{{$e.Name}}</td> <td class="line-status">{{$e.Status}}</td><td class="line-text">{{$e.Text}}</td></tr>{{end}}
</table>
`
)
