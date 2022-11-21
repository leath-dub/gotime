package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotime/src"
	"html/template"
	"net/http"
	"time"
    "log"
)

/* default body file name */
const TemplateFile = "body.json"

/* Go template for url */
const UrlTemplate =
"https://opentimetable.dcu.ie/broker/api/categorytypes/{{ .Category }}/categories{{ if .Type }}/events/filter{{ else }}/filter?pagenumber=1&query={{ end }}{{ if not .Type }}{{ .Query }}{{ end }}"

type IdResponse struct {
    TotalPages int
    Results []struct {Identity string `json:"Identity"`} `json:"Results"`
}

type Request struct {
    Category string
    Type bool
    Query string
    Id []string
}

/* Type field will be one of these */
const (
    Id = false
    Timetable = true
)

/* Category field will be one of these */
const (
    ProgrammesOfStudy = "241e4d36-60e0-49f8-b27e-99416745d98d"
    Module = "525fe79b-73c3-4b5c-8186-83c652b3adcc"
    Location = "1e042cb1-547d-41d4-ae93-a1f2c3d34538"
)

/* Headers for any request */
var DefaultHeaders = http.Header {
    "Authorization": {"basic T64Mdy7m["},
    "Content-Type": {"application/json; charset=utf-8"},
    "Accept": {"application/json; charset=utf-8"},
    "credentials": {"include"},
    "Origin": {"https://opentimetable.dcu.ie/"},
    "Referer": {"https://opentimetable.dcu.ie/"},
}

func dieError(message string, err error) {
    if err != nil {
        log.Fatalln(message, err)
    }
}

func handleResponse(response *http.Response) {
    if response.StatusCode != 200 {
        log.Fatalln("unexpected http response: ", response.Status)
    }
}

/* get url based on Request struct */
func getUrl(self *Request, framework string) string {
    parser, err := template.New("url").Parse(framework)
    dieError("error parsing template: ", err)

    var buf bytes.Buffer
    err = parser.Execute(&buf, self)

    return string(buf.Bytes())
}

/* executer Id request */
func getId(category string, query string) string {
    var requestObject = Request {
        Category: category,
        Type: Id,
        Query: query,
    }

    url := getUrl(&requestObject, UrlTemplate)

    /* fetch id based on query */
    httpRequest, err := http.NewRequest("POST", url, nil)
    dieError("error in http request: ", err)

    /* set headers */
    httpRequest.Header = DefaultHeaders
    response, err := http.DefaultClient.Do(httpRequest)
    dieError("error in http request: ", err)
    handleResponse(response)

    /* Write response into data struct */
    var responseStruct IdResponse
    json.NewDecoder(response.Body).Decode(&responseStruct)

    return responseStruct.Results[0].Identity
}

/* execute timetable request */
func getTimetable(request *Request, date time.Time) {
    body := lib.NewBody(date, request.Id)

    templateParser, err := template.ParseFiles(TemplateFile)
    dieError("error parsing template file: ", err)

    /* parse template into jsonBody buffer */
    jsonBody := new(bytes.Buffer)
    err = templateParser.Execute(jsonBody, body)
    dieError("error in template execution: ", err)

    url := getUrl(request, UrlTemplate)

    /* construct http request */
    httpRequest, err := http.NewRequest("POST", url, jsonBody)

    /* set headers */
    httpRequest.Header = DefaultHeaders

    /* execute request */
    response, err := http.DefaultClient.Do(httpRequest)
    dieError("error in http request: ", err)
    handleResponse(response)


    var decode []lib.ResponseTemplate

    /* write response to stdout */
    json.NewDecoder(response.Body).Decode(&decode)

    var responseTime time.Time
    for i, v := range decode[0].CategoryEvents {
        fmt.Println("Name:", v.ExtraProperties[0].Value)
        fmt.Println("Location:", v.Location)
        fmt.Println("Lecturer:", v.ExtraProperties[1].Value)
        responseTime, err = time.Parse(time.RFC3339, v.StartDateTime)
        fmt.Printf("Time: %d:%02d-", responseTime.Hour(), responseTime.Minute())
        responseTime, err = time.Parse(time.RFC3339, v.EndDateTime)
        fmt.Printf("%d:%02d\n", responseTime.Hour(), responseTime.Minute())
        if i != len(decode[0].CategoryEvents) - 1 {
            fmt.Println("------")
        }
    }
}

/* helper function to execute requests based on Request */
func (self Request) Do() {
    switch self.Type {
        case Id:
            fmt.Println(getId(self.Category, self.Query))
        case Timetable:
            /* ensure id */
            if !(len(self.Id) > 0) {
                self.Id = append(self.Id, getId(self.Category, self.Query))
            }
            getTimetable(&self, time.Now())
    }
}

func main() {
    var newRequest = Request {
        Category: Module,
        Type: Timetable,
        Query: "ca116",
    }
    newRequest.Do()
}
