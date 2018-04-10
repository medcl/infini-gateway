# timeago

[![Build Status](https://travis-ci.org/ararog/timeago.svg?branch=master)](https://travis-ci.org/ararog/timeago)
[![Coverage Status](https://coveralls.io/repos/github/ararog/timeago/badge.svg?branch=master)](https://coveralls.io/github/ararog/timeago?branch=master)

TimeAgo is a library used to calculate how much time has been passed between
two dates, this library is mainly based on time type of go.


## Example

```golang
import (
  "fmt"
  "time"
  timeago "github.com/ararog/timeago"
)

d, _ := time.ParseDuration("-3h")
start := time.Now()
end := time.Now().Add(d)
got, _ := timeago.TimeAgoWithTime(start, end)
fmt.Printf("Output: %s\n", got)
```
