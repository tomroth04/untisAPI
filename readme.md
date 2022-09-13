# UntisApi
## Description
This is a wrapper for the untis api in go. 
Warning this is a work in progress
Currently, it is very limited and supports only authentication using Username/Secret.

The wrapper uses a few external libraries including [Resty](https://github.com/go-resty/resty), [OTP](https://https://github.com/pquerna/otp) and [Gjson](https://github.com/tidwall/gjson)


## Installation
```
go get github.com/tomroth04/untisAPI-go
```

## Getting started
```go
package main
import (
	"fmt"
	"untisAPI"
)

func main() {
	client := untisAPI.NewClient("https://demo.untis.at", "demo", "Awesome", "max", "maxspasswort")

	if err := client.Login(); err != nil {
		panic(err)
	}

	// Get all classes
	classes, err := client.GetClasses(false)
	if err != nil {
		panic(err)
	}
	fmt.Println(classes)

}

```

## Considerations
Due to me not able to access every endpoint of the untis api, I can't guarantee that this wrapper works for every endpoint. If you find any bugs or have any suggestions, feel free to open an issue or a pull request.
Therefore also not all the structs are complete. If you need a struct, feel free to open an issue or a pull request.

## Credits

[Javascript webuntis Api](https://github.com/SchoolUtils/WebUntis) for their documentation of the untis api
