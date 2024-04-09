# UntisAPI

## Description

This is an unofficial wrapper for the untis api in go.
Currently, it is very limited and supports only authentication using Username/Secret.

The library is not perfect, am not a professional go developer, so feel free to contribute.

## Installation

```
go get github.com/tomroth04/untisAPI
```

## Getting started

```golang
import (
    "fmt"
    "log"
    "time"
    "github.com/tomroth04/untisAPI"
)

username := "" // fill in your username
password := "" // fill in your password
school := "" // fill in your school
client := untisAPI.NewClient("antiope.webuntis.com", school, "ANDROID", username, password)

if err := client.Login(); err != nil {
    log.Fatalln(err)
}

endDate := time.Now().Add(time.Hour * 24 * 30)

// check timetable for the next 30 days
lessons, err := client.GetOwnTimetableForRange(time.Now(), endDate, false)
if err != nil {
    log.Fatalln(err)
}

for _, lesson := range lessons {
    fmt.Println(lesson)
}
```
