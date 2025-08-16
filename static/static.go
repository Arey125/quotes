package static

import (
	"embed"
	"fmt"
	"time"
)

//go:embed *
var StaticFiles embed.FS

var Timestamp string = fmt.Sprint(time.Now().Unix())

func GetPath(path string) string {
	return "/static/" + Timestamp + path
}
