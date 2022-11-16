package main

import (
    "fmt"
    //"io/ioutil"
    "net/http"
    "os"
    "time"
    "encoding/json"
)

/* some literals/constants */
const id_request string = "https://opentimetable.dcu.ie/broker/api/CategoryTypes/241e4d36-60e0-49f8-b27e-99416745d98d/Categories/Filter?pageNumber=1&query="
var headers http.Header = http.Header {
        "Authorization": {"basic T64Mdy7m["},
        "Content-Type": {"application/json; charset=utf-8"},
        "Accept": {"application/json; charset=utf-8"},
        "credentials": {"include"},
        "Origin": {"https://opentimetable.dcu.ie/"},
        "Referer": {"https://opentimetable.dcu.ie/"},
};

/* request id from module code */
func module_code_to_id(module string) (string) {
    var url string = id_request + module;

    client := http.Client {Timeout: time.Duration(3) * time.Second};

    req, err := http.NewRequest("POST", url, nil);
    if err != nil {
        fmt.Fprintf(os.Stderr, "error in constructing http request: %s\n", err);
    }

    req.Header = headers;
    res, err := client.Do(req);

    var decoder *json.Decoder;
    decoder = json.NewDecoder(res.Body);

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

    var msg message_t;
    err = decoder.Decode(&msg);
    return msg.Results[0].Identity;
}

func main() {
    fmt.Printf("%s\n", module_code_to_id("comsci2"));
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
