package main

import (
	"fmt"

	"github.com/gitHusband/goutils/jsonkeys"
)

const file = "./testeasy.json"

func main() {
	testjsonKeys()
}

func testjsonKeys() {
	jsonStr := `{"na\\me":"Tom\\", "age":"18","is\"Cop":false}`

	jsonDataKeysMap, err := jsonkeys.ParseFromData([]byte(jsonStr))
	if err != nil {
		fmt.Printf("\nJSON Data error： %v\n", err)
	}
	fmt.Printf("JSON Data key 顺序： \n%#v\n", jsonDataKeysMap)

	jsonFileKeysMap, err := jsonkeys.ParseFromFile(file)
	if err != nil {
		fmt.Printf("\nJSON File error： %v\n", err)
	}
	fmt.Printf("JSON File key 顺序： \n%#v\n", jsonFileKeysMap)
}
