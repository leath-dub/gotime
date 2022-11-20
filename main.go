package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotime/src"
	"html/template"
	"io"

	// "io/ioutil"
	"net/http"
	"os"
	"time"

)

/* some literals/constants */

type opentimetable_api_t struct {
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
}

func die(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %s", err)
        os.Exit(1)
    }
}

func (self opentimetable_request_t) do(body io.Reader) *json.Decoder {
    client := http.Client {Timeout: time.Duration(10) * time.Second}
    url := self.prefix + self.category + self.request + self.extra

    req, err := http.NewRequest("POST", url, body)
    die(err)

    req.Header = self.headers
    res, err := client.Do(req)
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
    /* -- Make id request -- */
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
    /* -- -- */

    /* -- Make timetable request -- */
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
    fmt.Println(buf.String())

    var other_msg []lib.ResponseTemplate

    other_req.do(&buf).Decode(&other_msg)

    fmt.Println(other_msg)
    /*
    fmt.Println(other_msg[0].CategoryEvents[0].ExtraProperties[0].DisplayName)
    for i := 0; i < len(other_msg[0].CategoryEvents); i++ {
        fmt.Println(other_msg[0].CategoryEvents[i].ExtraProperties[0].Value)
    }
    */
    /* -- -- */
}
