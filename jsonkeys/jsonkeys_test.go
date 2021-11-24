package jsonkeys

import (
	"reflect"
	"testing"
)

func TestFormString(t *testing.T) {
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
	var expected = JsonKeysMap{
		"root":                KeySlice{"name", "age", "is\"Cop", "favoriteFruits", "familyMembers", "codeLanguage"},
		"root.codeLanguage":   KeySlice{"Golange", "Javascript", "PHP"},
		"root.favoriteFruits": KeySlice{"bannana", "apple", "peach"},
	}

	// 1. 从JSON数据解析key
	result, err := ParseFromData([]byte(jsonStr))
	if err != nil {
		t.Errorf("\nJSON Data error： %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("\nResult： %v\nExpected： %v\n", result, expected)
	}
}

func TestFormFile(t *testing.T) {
	file := "./jsonkeys_test.json"

	var expected = JsonKeysMap{
		"root":                KeySlice{"name", "age", "is\"Cop", "favoriteFruits", "familyMembers", "codeLanguage"},
		"root.codeLanguage":   KeySlice{"Golange", "Javascript", "PHP"},
		"root.favoriteFruits": KeySlice{"bannana", "apple", "peach"},
	}

	// 2. 从JSON文件解析key
	result, err := ParseFromFile(file)

	if err != nil {
		t.Errorf("\nJSON Data error： %v\n", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("\nResult： %v\nExpected： %v\n", result, expected)
	}
}
