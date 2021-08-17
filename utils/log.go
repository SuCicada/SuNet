package utils

import (
	"fmt"
	"log"
)

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
func LogErr(err error, v ...interface{}) {
	if err != nil {
		if typeof(v[1]) == "func(error)" {
			v[1].(func(error))(err)
		} else {
			log.Fatalf(fmt.Sprintln(v[1]), v[2:]...)
		}
		//switch len(err) {
		//case 1:
		//	log.Fatalln(err[0])
		//default:
		//}
	}
}
