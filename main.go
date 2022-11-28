package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gotime/src"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

/* default body file name */
const TemplateFile = "body.json"
var offset = 0

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
    err = json.NewDecoder(response.Body).Decode(&responseStruct)
    dieError("error in decoding response: ", err)

    if len(responseStruct.Results) > 0 {
        return responseStruct.Results[0].Identity
    }
    log.Fatalln("empty response")
    return ""
}

func max(x, y int) int {
    if x < y {
        return y
    }
    return x
}

func maxCategoryEvent(events *lib.CategoryEvents) int {
    var localMax int
    globalMax := 0
    for _, event := range *events {
        localMax = max(len(event.ExtraProperties[0].Value), len(event.Location) + 10)
        if localMax > globalMax {
            globalMax = localMax
        }
    }
    return globalMax
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

    /* sort CategoryEvents based on time -- see lib.go for more details */
    if len(decode) == 0 {
        return
    }
    events := decode[0].CategoryEvents

    lineWidth := maxCategoryEvent(&events)

    sort.Stable(lib.CategoryEvents(events))

    var responseTime time.Time
    for _, v := range events {
        if len(v.ExtraProperties) < 2 {
            continue
        }
        fmt.Println(strings.Repeat("─", lineWidth))
        fmt.Print("\033[3m", v.ExtraProperties[0].Value, "\033[m\n")
        fmt.Println("\033[35mLocation:\033[m", v.Location)
        fmt.Println("\033[36mLecturer:\033[m", v.ExtraProperties[1].Value)
        responseTime, err = time.Parse(time.RFC3339, v.StartDateTime)
        fmt.Printf("\033[34mTime:\033[m %d:%02d-", responseTime.Hour(), responseTime.Minute())
        responseTime, err = time.Parse(time.RFC3339, v.EndDateTime)
        fmt.Printf("%d:%02d\n", responseTime.Hour(), responseTime.Minute())
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
            getTimetable(&self, time.Now().AddDate(0, 0, offset))
    }
}

func getOpt(n int) (int, error) {
    var char []byte = make([]byte, 1)
    var pos int
    for {
        os.Stdin.Read(char)
        switch char[0] {
            case 66: fallthrough /* down arrow */
            case 'j': if pos < n - 1 {
                pos++
                fmt.Print("\033[1B")
            }
            case 65: fallthrough /* up arrow */
            case 'k': if pos > 0 {
                pos--
                fmt.Print("\033[1A")
            }
            case 'q': {
                return 0, errors.New("error: user exited")
            }
            case '\n': {
                fmt.Print("*\033[2D")
                fmt.Printf("\033[%dB", n - pos)
                return pos, nil
            }
        }
    }
}

func getInput() Request {
    var result Request

    /* disable input buffering */
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    /* disable echo of input */
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

    options := []string { ProgrammesOfStudy, Location, Module }

    /* clear screen and move cursor to HOME */
    fmt.Print("\033[2J\033[H")
    fmt.Println("Choose a category:")
    fmt.Println("[ ] Programmes of Study")
    fmt.Println("[ ] Location")
    fmt.Println("[ ] Module")
    fmt.Print("\033[H\033[1B\033[1C")

    pos, err := getOpt(3)
    dieError("", err)

    result.Category = options[pos]

    opts := []bool { Timetable, Id }
    /* clear screen and move cursor to HOME */
    fmt.Print("\033[2J\033[H")
    fmt.Println("Choose Request Type:")
    fmt.Println("[ ] Timetable")
    fmt.Println("[ ] Id")
    fmt.Print("\033[H\033[1B\033[1C")

    pos, err = getOpt(2)
    dieError("", err)

    result.Type = opts[pos]

    exec.Command("stty", "-F", "/dev/tty", "-cbreak").Run()
    exec.Command("stty", "-F", "/dev/tty", "echo").Run()

    fmt.Print("Query: ")
    var str string
    fmt.Scanf("%s\n", &str)
    fmt.Print("\033[2J\033[H")

    result.Query = str
    return result
}

func main() {
    /* temporary: specify argument -- how many days from now */
    if len(os.Args) > 1 {
        offset, _ = strconv.Atoi(os.Args[1])
    }
    newRequest := getInput()
    newRequest.Do()
}
