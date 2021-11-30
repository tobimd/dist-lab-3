package informant

import (
	"fmt"
)

var (
	id int
	f  = ""
)

func Run(informantId int) {
	id = informantId
	f = fmt.Sprintf("informant_%d", id)
}
