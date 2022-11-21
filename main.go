package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotime/src"
	"html/template"
	"io"

	"net/http"
	"os"
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

type opentimetable_api_t struct {// {{{ OLD
    categories map[string]string
    requests map[string]string
    prefix string
    headers http.Header
}

type opentimetable_request_t struct {
    category string
    request string
    prefix string
    headers http.Header
    extra string
}

var def opentimetable_api_t = opentimetable_api_t {
    categories: map[string]string {
        "Programmes of Study": "241e4d36-60e0-49f8-b27e-99416745d98d",
        "Module": "525fe79b-73c3-4b5c-8186-83c652b3adcc",
        "Location": "1e042cb1-547d-41d4-ae93-a1f2c3d34538",
    },
    requests: map[string]string {
        "id": "/categories/filter?pagenumber=1&query=",
        "timetable": "/categories/events/filter",
    },
    prefix: "https://opentimetable.dcu.ie/broker/api/categorytypes/",
    headers: http.Header {
        "Authorization": {"basic T64Mdy7m["},
        "Content-Type": {"application/json; charset=utf-8"},
        "Accept": {"application/json; charset=utf-8"},
        "credentials": {"include"},
        "Origin": {"https://opentimetable.dcu.ie/"},
        "Referer": {"https://opentimetable.dcu.ie/"},
    },
}// }}}

func die(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %s", err)
        os.Exit(1)
    }
}

func (self opentimetable_request_t) do(body io.Reader) *json.Decoder {
    url := self.prefix + self.category + self.request + self.extra

    req, err := http.NewRequest("POST", url, body)
    die(err)

    req.Header = self.headers
    res, err := http.DefaultClient.Do(req)
    if res.StatusCode != 200 {
        fmt.Fprintf(os.Stderr, "http status code was not 200: status %d\n", res.StatusCode)
        os.Exit(1)
    }
    die(err)

    return json.NewDecoder(res.Body)
}

func start_of_week(date time.Time) time.Time {
    offset := int(date.Weekday()) - 1
    return date.AddDate(0, 0, -offset)
}

func construct_json_body(filename string, date time.Time, id []string, buf *bytes.Buffer) {
    start := start_of_week(date)
    body := lib.BodyTemplate {
        FirstDayInWeek: start.Format(time.RFC3339),
        Name: date.Weekday().String(),
        DayOfWeek: int(date.Weekday()),
        CategoryIdentities: id,
    }

    tmpl, err := template.ParseFiles(filename)
    die(err)

    err = tmpl.Execute(buf, body)
    die(err)
}


func main() {
    /*
    type identity_t struct {
        TotalPages int
        Results []struct {Identity string `json:"Identity"`} `json:"Results"`
    }

    var req opentimetable_request_t = opentimetable_request_t {
        category: def.categories["Programmes of Study"],
        prefix: def.prefix,
        request: def.requests["id"],
        headers: def.headers,
        extra: "comsci2",
    }

    var msg identity_t
    req.do(nil).Decode(&msg)

    var other_req opentimetable_request_t = opentimetable_request_t {
        category: def.categories["Programmes of Study"],
        prefix: def.prefix,
        request: def.requests["timetable"],
        headers: def.headers,
        extra: "",
    }

    identities := []string {
        msg.Results[0].Identity,
    }

    var buf bytes.Buffer
    construct_json_body("body.json", time.Now().AddDate(0, 0, 2), identities, &buf)

    var other_msg []lib.ResponseTemplate

    other_req.do(&buf).Decode(&other_msg)

    fmt.Println(other_msg)
    */

    var newRequest = Request {
        Category: Module,
        Type: Timetable,
        Query: "ca116",
    }
    newRequest.Do()
}
