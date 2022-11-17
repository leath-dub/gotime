package main

import (
    "fmt"
    "net/http"
    "os"
    "time"
    "encoding/json"
)

/* some literals/constants */

type opentimetable_api_t struct {
    categories map[string]string
    requests map[string]string
    prefix string
    headers http.Header
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

func construct_api_request(
    prefix string, cat string, req string, extra string,
) (string) {
    return prefix + def.categories[cat] + def.requests[req] + extra
}

func api_request(url string) (string) {
    client := http.Client {Timeout: time.Duration(3) * time.Second}

    req, err := http.NewRequest("POST", url, nil)
    die(err)

    req.Header = def.headers
    res, err := client.Do(req)
    if res.StatusCode != 200 {
        fmt.Fprintf(os.Stderr, "http status code was not 200: status %d\n", res.StatusCode)
        os.Exit(1)
    }
    die(err)

    type identity_t struct {
        Results []struct {Identity string `json:"Identity"`} `json:"Results"`
    }

    var msg identity_t
    err = json.NewDecoder(res.Body).Decode(&msg)
    die(err)

    die(err)

    if len(msg.Results) == 0 {
        fmt.Fprintf(
            os.Stderr,
            "Results array is length 0, possibly invalid request, ",
        )
        os.Exit(1)
    }

    return msg.Results[0].Identity
}

func main() {
    fmt.Printf(api_request(
        construct_api_request(
            def.prefix, "Programmes of Study", "id", "comsci2",
        ),
    ))
}

/*
    curl -XPOST \
        -H "Authorization: basic T64Mdy7m[" \
        -H "Content-Type: application/json; charset=utf-8" \
        -H "Accept: application/json; charset=utf-8" \
        -H "credentials: include" \
        -H "Origin: https://opentimetable.dcu.ie/" \
        -H "Referer: https://opentimetable.dcu.ie/" \
        "https://opentimetable.dcu.ie/broker/api/CategoryTypes/241e4d36-60e0-49f8-b27e-99416745d98d/Categories/Filter?pageNumber=1&query=${module}" |
        jq '.Results[0].Identity' |
*/

/* response
{"TotalPages":1,"CurrentPage":1,"Results":[{"ParentCategoryIdentities":["af9505b6-8af2-eae8-a6e3-8
12154330274","7f505ad1-83a2-9bfd-5416-54c52e9de16d","f8c44a18-b544-04e9-134e-db5e84826dbf"],"Categ
oryTypeIdentity":"241e4d36-60e0-49f8-b27e-99416745d98d","CategoryTypeName":null,"CategoryEvents":n
ull,"Name":"COMSCI2","Identity":"3195ffd3-b64c-9a1b-d344-7fc17c57f03d"}],"Count":1}
*/
/* GO equivalent
type result_t struct {
    ParentCategoryIdentities []string
    CategoryTypeIdentity string
    CategoryTypeName string
    CategoryEvents string
    Name string
    Identity string
}

type message_t struct {
    TotalPages int
    CurrentPage int
    Results []result_t
    Count int
}
*/
// "https://opentimetable.dcu.ie/broker/api/categoryTypes/241e4d36-60e0-49f8-b27e-99416745d98d/categories/events/filter")

