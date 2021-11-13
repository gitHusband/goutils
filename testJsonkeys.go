package main

import (
	"fmt"

	"github.com/gitHusband/goutils/jsonkeys"
)

func main() {
	testjsonKeys()
}

func testjsonKeys() {
	jsonStr := `
{
	"name": "Tom",
	"age": 25,
	"is\"Cop": false,
	"favoriteFruits": {
		"bannana": "yellow",
		"apple": "red",
		"peach": "pink"
	},
	"familyMembers": [
		"David",
		"Sammy"
	],
	"codeLanguage": {
		"Golange": "21 世纪的\"C\"语言",
		"Javascript": "Web 页面脚本语言",
		"PHP": "世界上最好的语言"
	}
}
`
	// 1. 从JSON数据解析key
	jsonDataKeysMap, err := jsonkeys.ParseFromData([]byte(jsonStr))
	if err != nil {
		fmt.Printf("\nJSON Data error： %v\n", err)
	}
	fmt.Printf("JSON Data key 顺序： \n%#v\n", jsonDataKeysMap)
	for dk, dv := range jsonDataKeysMap {
		fmt.Printf("- data key 顺序：%v: %#v\n", dk, dv)
	}

	file := "./testeasy.json"
	// 2. 从JSON文件解析key
	jsonFileKeysMap, err := jsonkeys.ParseFromFile(file)
	if err != nil {
		fmt.Printf("\nJSON File error： %v\n", err)
	}
	fmt.Printf("JSON File key 顺序： \n%#v\n", jsonFileKeysMap)
	for key, value := range jsonFileKeysMap {
		fmt.Printf("- file key 顺序：%v: %v\n", key, value)
	}
}
