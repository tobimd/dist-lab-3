package fulcrum

import "fmt"

var (
	id int
	f  = ""
)

func Run(fulcrumId int) {
	id = fulcrumId
	f = fmt.Sprintf("fulcrum_%d", id)

}
