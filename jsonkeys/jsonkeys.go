// 动态JSON 可以转 map[string]interface{}。正如你所知，map 是无序的，不能按照 JSON key 原来的顺序遍历 map
// jsonkeys 包的目的正是 获取动态JSON 中 key 的先后顺序

package jsonkeys

import (
	"bufio"
	"fmt"
	"os"
)

type keySlice []string
type jsonKeysMap map[string]keySlice

type scanner struct {
	step func(*scanner, byte) (int, error)
	// 存放叠加状态，先进后出，"{" "}" 必须成对删除
	states []int
	// 存放正在扫描的 key 的所有字符，扫描结束后合成字符串并清空之
	keyCharacters []byte
	// 存放正在扫描 JSON 的完整路径。
	// 最外层默认是 {RootPathName}, 之后每次扫描到 "{", 都应将其父 key 插入到 keyPath 中
	// 每次 "{" "}" 成对删除，都必须删除 keyPath 最后一个元素
	keyPath []string
}

const (
	scanStateIgnore                 = iota // 没意义的字节，直接忽略
	scanStateBeginObject                   // 对象的开始标志 "{"
	scanStateBeginKey                      // key 的开始标志 '"'
	scanStateKeyCharacter                  // key 的内容字节，支持任何字符，包括与 key 结束标志相同的 '"'
	scanStateBackslash                     // 由于 key 内容支持 '"', 必须判断读取到的 '"' 是否是 key 的结束标志
	scanStateEndKey                        // key 的结束标志 '"'
	scanStateKeyValueSeparator             // key 的结束标志 '"'
	scanStateBeginValue                    // value 的开始标志 '"', 它与 key 的结束标志之间必须有 ":"
	scanStateBeginValueWithArray           // value 以 '[' 开始，以 ']' 结尾, 数组类型
	scanStateBeginValueWithoutQuote        // value 没有以 '"' 为开始标志, 比如 数字，布尔值等
	scanStateValueCharacter                // 由于 value 内容支持 '"', 必须判断读取到的 '"' 是否是 value 的结束标志
	scanStateEndValue                      // value 的结束标志 '"'
	scanStateEndValueWithoutQuote          // value 没有以 '"' 为结束标志, 比如 数字，布尔值等
	scanStateEndValueWithArray             // value 以 '[' 开始，以 ']' 结尾, 数组类型
	// scanStateKeyKeySeparator               // 对象包含下一个 key 的标志 ","
	scanStateEndObject // 对象的结束标志 "}"
)

var RootPathName = "root"

var (
	keys jsonKeysMap
	scan = new(scanner)
)

func GetKeys() jsonKeysMap {
	return keys
}

func ParseFromData(data []byte) (jsonKeysMap, error) {
	keys = jsonKeysMap{RootPathName: []string{}}
	scan.reset()

	dataLen := len(data)
	for i := 0; i < dataLen; i++ {
		// fmt.Printf("%v", string(data[i]))

		state, err := scan.step(scan, data[i])
		if err != nil {
			return nil, err
		}

		switch state {
		case scanStateIgnore:
			continue
		case scanStateBeginObject:
			continue
		case scanStateBeginKey:
			continue
		case scanStateKeyCharacter:
			continue
		case scanStateEndKey:
			continue
		case scanStateKeyValueSeparator:
			continue
		case scanStateBeginValue:
			continue
		case scanStateBeginValueWithoutQuote:
			i--
			continue
		case scanStateValueCharacter:
			continue
		case scanStateEndValue:
			continue
		case scanStateEndValueWithoutQuote:
			i--
			continue
		case scanStateEndObject:
			continue
		}
	}

	return keys, nil
}

func ParseFromFile(file string) (jsonKeysMap, error) {
	fileObj, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer fileObj.Close()

	reader := bufio.NewReader(fileObj)
	buf := make([]byte, 1024)
	_, err = reader.Read(buf)

	if err != nil {
		return nil, err
	}

	return ParseFromData(buf)
}

func (s *scanner) reset() {
	s.step = stepBeginObject
	s.states = []int{}
	s.keyCharacters = []byte{}
	s.keyPath = []string{RootPathName}
}

func (s *scanner) getFullKeyPath() string {
	fullKeyPath := s.keyPath[0]
	for i := 1; i < len(s.keyPath); i++ {
		fullKeyPath += ("." + fullKeyPath)
	}

	return fullKeyPath
}

func (s *scanner) setOneKey() {
	key := string(s.keyCharacters)
	// 清空 keyCharacters，准备保存下一个 key
	s.keyCharacters = s.keyCharacters[0:0]
	getFullKeyPath := s.getFullKeyPath()

	// 保存一个 key 到 全局变量 keys
	if _, ok := keys[getFullKeyPath]; !ok {
		keys[getFullKeyPath] = keySlice{key}
	} else {
		keys[getFullKeyPath] = append(keys[getFullKeyPath], key)
	}
}

// 获取最后一个 state
func (s *scanner) getLastState() int {
	return s.states[len(s.states)-1]
}

// 删除最后一个 state
func (s *scanner) deleteLastState() {
	s.states = s.states[:len(s.states)-1]
}

// 判断上一个字符是否是 "\"
func (s *scanner) isLastStateBackslash() bool {
	return s.getLastState() == scanStateBackslash
}

func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\t' || c == '\r' || c == '\n')
}

// 结构体 以 “{” 开始，以 "}" 结尾
// 判断是否是一个结构体的开始字符 “{”
// 标志着 下一个内容是 字段名（key）
func stepBeginObject(s *scanner, c byte) (int, error) {
	if isSpace(c) {
		return scanStateIgnore, nil
	}
	switch c {
	case '{':
		scan.step = stepBeginKey
		// 保存 scanStateBeginObject
		scan.states = append(scan.states, scanStateBeginObject)
		return scanStateBeginObject, nil
	default:
		return -1, fmt.Errorf("error json format, See stepBeginObject: character(%v)", string(c))
	}
}

// key 以 '"' 开始，以 '"' 结尾
// 判断是否是一个 key 的开始字符 '"'
// 在双引号之间的所有字符将合成 key 字符串保存起来。
func stepBeginKey(s *scanner, c byte) (int, error) {
	if isSpace(c) {
		return scanStateIgnore, nil
	}
	switch c {
	// key 开始标志
	case '"':
		scan.step = stepKeyCharacter
		scan.states = append(scan.states, scanStateBeginKey)
		return scanStateBeginKey, nil
	default:
		return -1, fmt.Errorf("error json format, See stepBeginKey: character(%v)", string(c))
	}
}

// 任何字符都能作为 key 值，包括 '"'
// 那么必须区分 '"' 是否是 key 的结尾
func stepKeyCharacter(s *scanner, c byte) (int, error) {
	switch c {
	// key 结束标志
	case '"':
		if scan.isLastStateBackslash() {
			scan.keyCharacters = append(scan.keyCharacters, c)
			scan.deleteLastState()
			return scanStateKeyCharacter, nil
		} else {
			scan.step = stepEndKey
			scan.states = append(scan.states, scanStateEndKey)
			scan.setOneKey()
			return scanStateEndKey, nil
		}
	case '\\':
		if scan.isLastStateBackslash() {
			scan.keyCharacters = append(scan.keyCharacters, c)
			scan.deleteLastState()
			return scanStateKeyCharacter, nil
		} else {
			scan.states = append(scan.states, scanStateBackslash)
			return scanStateKeyCharacter, nil
		}
	default:
		scan.keyCharacters = append(scan.keyCharacters, c)
		return scanStateKeyCharacter, nil
	}
}

// key 以 '"' 开始，以 '"' 结尾
// 判断是否是一个 key 的结束字符 '"'
// 标志着 下一个内容是 字段值（value）
func stepEndKey(s *scanner, c byte) (int, error) {
	if isSpace(c) {
		return scanStateIgnore, nil
	}
	switch c {
	// key 与 value 的分隔符
	case ':':
		scan.step = stepBeginValue
		scan.states = append(scan.states, scanStateKeyValueSeparator)
		return scanStateKeyValueSeparator, nil
	default:
		return -1, fmt.Errorf("error json format, See stepEndKey: character(%v)", string(c))
	}
}

// 1. value 以 '"' 开始，以 '"' 结尾
// 2. value 以 '[' 开始，以 ']' 结尾, 数组类型
// 3. value 不以 '"' 开始，不以 '"' 结尾, 比如 数字，布尔值等
// 判断是否是一个 value 的开始字符 '"'
// 在双引号之间的所有字符将合成 value 字符串保存起来。
func stepBeginValue(s *scanner, c byte) (int, error) {
	if isSpace(c) {
		return scanStateIgnore, nil
	}
	switch c {
	// 1. value 开始标志 '"'
	case '"':
		scan.step = stepValueCharacter
		scan.states = append(scan.states, scanStateBeginValue)
		return scanStateBeginValue, nil
	// 2. value 以 '[' 开始，以 ']' 结尾, 数组类型
	case '[':
		scan.step = stepValueCharacterWithArray
		scan.states = append(scan.states, scanStateBeginValueWithArray)
		return scanStateBeginValueWithArray, nil
	// 3. value 没有开始结束标志
	default:
		scan.step = stepValueCharacterWithoutQuote
		scan.states = append(scan.states, scanStateBeginValueWithoutQuote)
		return scanStateBeginValueWithoutQuote, nil
	}
}

// 1. value 以 '"' 开始，以 '"' 结尾
// 任何字符都能作为 value 值，包括 '"'
// 那么必须区分 '"' 是否是 value 的结尾
// 我们不要 value, 所以就不保存它了
func stepValueCharacter(s *scanner, c byte) (int, error) {
	switch c {
	case '"':
		if scan.isLastStateBackslash() {
			scan.keyCharacters = append(scan.keyCharacters, c)
			scan.deleteLastState()
			return scanStateValueCharacter, nil
		} else {
			scan.step = stepEndValue
			scan.states = append(scan.states, scanStateEndValue)
			return scanStateEndValue, nil
		}
	case '\\':
		if scan.isLastStateBackslash() {
			scan.keyCharacters = append(scan.keyCharacters, c)
			scan.deleteLastState()
			return scanStateValueCharacter, nil
		} else {
			scan.states = append(scan.states, scanStateBackslash)
			return scanStateValueCharacter, nil
		}

	default:
		return scanStateValueCharacter, nil
	}
}

// 2. value 以 '[' 开始，以 ']' 结尾, 数组类型
// 我们不要 value, 所以就不保存它了
func stepValueCharacterWithArray(s *scanner, c byte) (int, error) {
	switch c {
	case ']':
		scan.step = stepEndValue
		scan.states = append(scan.states, scanStateEndValueWithArray)
		return scanStateEndValueWithArray, nil
	default:
		return scanStateValueCharacter, nil
	}
}

// 3. value 不以 '"' 开始，不以 '"' 结尾, 比如 数字，布尔值等
// 我们不要 value, 所以就不保存它了
func stepValueCharacterWithoutQuote(s *scanner, c byte) (int, error) {
	switch c {
	case ',', '}':
		scan.step = stepEndValue
		scan.states = append(scan.states, scanStateEndValueWithoutQuote)
		return scanStateEndValueWithoutQuote, nil
	default:
		return scanStateValueCharacter, nil
	}
}

// value 以 '"' 开始，以 '"' 结尾
// 判断是否是一个 value 的结束字符 '"'
// 标志着 下一个内容是 新的 key(",") 或者 结束("}")
func stepEndValue(s *scanner, c byte) (int, error) {
	if isSpace(c) {
		return scanStateIgnore, nil
	}
	switch c {
	// key 与 key 之间的分割符
	case ',':
		scan.step = stepBeginKey
		scan.states = append(scan.states, scanStateBeginKey)
		return scanStateBeginKey, nil
	// 目前不考虑多级嵌套的JSON
	case '}':
		scan.step = stepEndObject
		scan.states = append(scan.states, scanStateEndObject)
		return scanStateEndObject, nil
	default:
		return -1, fmt.Errorf("error json format, See stepEndValue: character(%v)", string(c))
	}
}

func stepEndObject(s *scanner, c byte) (int, error) {
	return 0, nil
}
