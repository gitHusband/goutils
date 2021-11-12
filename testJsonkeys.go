package main

import (
	"fmt"

	"github.com/gitHusband/goutils/jsonkeys"
)

func main() {
	testjsonKeys()
}

func testjsonKeys() {
	jsonStr := `{"na\\me":"Tom\\", "age":"18","is\"Cop":false}`

	jsonKeysMap := jsonkeys.ParseFromData([]byte(jsonStr))
	fmt.Printf("JSON key 顺序： \n%#v\n", jsonKeysMap)
}
