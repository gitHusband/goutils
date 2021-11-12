package main

import (
	"fmt"

	"github.com/gitHusband/goutils/jsonkeys"
)

func main() {
	testjsonKeys()
}

func testjsonKeys() {
	jsonStr := `{"name":"Tom", "age":"18","isCop":"false"}`

	jsonKeysMap := jsonkeys.ParseKeys([]byte(jsonStr))
	fmt.Printf("JSON key 顺序： \n%#v\n", jsonKeysMap)
}
