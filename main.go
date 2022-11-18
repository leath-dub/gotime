package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type body_t struct {
    ViewOptions struct {
        Days []struct {
            Name string
            DayOfWeek int
            IsDefault bool
        }
        Weeks []struct {
            WeekNumber int
            WeekLabel int
            FirstDayInWeek string
        }
        TimePeriods []struct {
            Description string
            StartTime string
            EndTime string
            IsDefault bool
        }
        DatePeriods []struct {
            Description string
            StartDateTime string
            EndDateTime string
            IsDefault bool
            IsThisWeek bool
            IsNextWeek bool
            Type string
        }
        LegendItems []any
        InstitutionConfig struct {}
        DateConfig struct {
          FirstDayInWeek int
          StartDate string
          EndDate string
        }
    }
    CategoryIdentities []string
}

type response_t struct {
    CategoryTypeIdentity string
    CategoryTypeName string
    CategoryEvents []struct {
        EventIdentity string
        HostKey string
        Description string
        EndDateTime string
        EventType string
        IsPublished string
        Location string
        Owner string
        StartDateTime string
        IsDeleted bool
        LastModified string
        ExtraProperties []struct {
            Name string
            DisplayName string
            Value string
            Rank int
        }
    }
}

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

func construct_json_body(filename string, body *body_t, date time.Time, id []string) {
    buf, err := ioutil.ReadFile(filename)
    die(err)

    json.Unmarshal(buf, &body);

    start := start_of_week(date)

    body.ViewOptions.Weeks[0].FirstDayInWeek = start.Format(time.RFC3339)
    body.ViewOptions.Days[0].Name = start.Weekday().String()
    body.ViewOptions.Days[0].DayOfWeek = int(start.Weekday())
    body.CategoryIdentities = id
}

func main() {
    /* -- Make id request -- */
    type identity_t struct {
        TotalPages int
        Results []struct {Identity string `json:"Identity"`} `json:"Results"`
    }

    var req opentimetable_request_t = opentimetable_request_t {
        category: def.categories["Location"],
        prefix: def.prefix,
        request: def.requests["id"],
        headers: def.headers,
        extra: "lg25",
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

    var jsl body_t
    identities := []string {
        msg.Results[0].Identity,
    }
    construct_json_body("body.json", &jsl, time.Now(), identities)

    var other_msg []response_t
    send, err := json.Marshal(&jsl)
    die(err)

    other_req.do(bytes.NewBuffer(send)).Decode(&other_msg)

    fmt.Println(other_msg[0].CategoryEvents[0].ExtraProperties[0].DisplayName)
    for i := 0; i < len(other_msg[0].CategoryEvents); i++ {
        fmt.Println(other_msg[0].CategoryEvents[i].ExtraProperties[0].Value)
    }
    /* -- -- */
}
