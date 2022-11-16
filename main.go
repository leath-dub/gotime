package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	// "encoding/json"
)

const id_request string = "https://opentimetable.dcu.ie/broker/api/CategoryTypes/241e4d36-60e0-49f8-b27e-99416745d98d/Categories/Filter?pageNumber=1&query="
func module_code_to_id(module string) {
    var url string = id_request + module;

    client := http.Client {Timeout: time.Duration(3) * time.Second};

    req, err := http.NewRequest("POST", url, nil);
    if err != nil {
        fmt.Fprintf(os.Stderr, "error in constructing http request: %s\n", err);
    }

    req.Header.Add("Authorization", "basic T64Mdy7m[");
    req.Header.Add("Content-Type", "application/json; charset=utf-8");
    req.Header.Add("Accept", "application/json; charset=utf-8");
    req.Header.Add("credentials", "include");
    req.Header.Add("Origin", "https://opentimetable.dcu.ie/");
    req.Header.Add("Referer", "https://opentimetable.dcu.ie/");

    res, err := client.Do(req);
    body, err := ioutil.ReadAll(res.Body);
    fmt.Printf(string(body));
}

func main() {
    module_code_to_id("comsci2");
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

