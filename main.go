package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	// "encoding/json"
)

func main() {
    var res *http.Response;
    var err error;

    /* make http get request */
    res, err = http.Get("https://opentimetable.dcu.ie");

    /* handle error if any */
    if err != nil {
        fmt.Fprintf(os.Stderr, "error http request: %s\n", err);
        return;
    }

    body, err := ioutil.ReadAll(res.Body);

    fmt.Println(string(body));
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

