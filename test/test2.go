package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type event struct {
	signal int16
	value  []byte
}

func main() {

	data := make([]byte, 2) // unsafe.Sizeof(uint16(1))
	binary.BigEndian.PutUint16(data, uint16(123))
	a := &event{
		signal: 771,
		value:  data,
	}
	fmt.Println(a)
	res, _ := json.Marshal(a)
	fmt.Println(res)

	e := &event{}
	json.Unmarshal(res, e)

	fmt.Println(*e)

}
