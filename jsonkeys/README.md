# jsonkeys - 获取 json key 的先后顺序


当使用GO 标准库 `encoding/json` 解析动态JSON 的时候，我们将结果解析为 `map[string]interface{}`。

而 GO `map` 类型的key 是无序的，也就是说你不能确定JSON key 的先后顺序。

如果你需要确定 JSON key 的顺序，可以使用 `jsonkeys` 包。

> 目录

[TOC]

## 1. 使用示例

```
testJsonkeys.go 文件
```

```
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
    
	for dk, dv := range jsonDataKeysMap {
		fmt.Printf("- data key 顺序：%v: %#v\n", dk, dv)
	}
}
```
**执行结果**：
```
% go run testJsonkeys.go
- data key 顺序：root.favoriteFruits: jsonkeys.keySlice{"bannana", "apple", "peach"}
- data key 顺序：root.codeLanguage: jsonkeys.keySlice{"Golange", "Javascript", "PHP"}
- data key 顺序：root: jsonkeys.keySlice{"name", "age", "is\"Cop", "favoriteFruits", "familyMembers", "codeLanguage"}
```
## 2. API
<table style="width:100%">
<thead>
	<th>函数</th>
	<th>释义</th>
</thead>
<tbody>
	<tr align="left">
		<td>ParseFromData</td>
		<td><strong style="font-size: 15px">从JSON数据解析key</strong><br/> <em style="color: #888888">func ParseFromData(data []byte) (JsonKeysMap, error)</em></td>
	</tr>
	<tr align="left">
		<td>ParseFromFile</td>
		<td><strong style="font-size: 15px">从JSON文件解析key</strong><br/> <em style="color: #888888">func ParseFromFile(file string) (JsonKeysMap, error)</em></td>
	</tr>
	<tr align="left">
		<td>Get</td>
		<td><strong style="font-size: 15px">获取 key 排序切片</strong><br/> <em style="color: #888888">func (jkm JsonKeysMap) Get(keyPath string) (StringSlice, error)</em></td>
	</tr>
</tbody>
</table>