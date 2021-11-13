// 动态JSON 可以转 map[string]interface{}。正如你所知，map 是无序的，不能按照 JSON key 原来的顺序遍历 map
// jsonkeys 包的目的正是 获取动态JSON 中 key 的先后顺序

package jsonkeys

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type KeySlice []string
type JsonKeysMap map[string]KeySlice

type scanner struct {
	step func(*scanner, byte) (int, error)
	// 存放叠加状态，先进后出，"{" "}" 必须成对删除
	states []int
	// 存放正在扫描的 key 的所有字符，扫描结束后合成字符串并清空之
	keyCharacters []byte
	// 存在上一个扫描到的 key，如果它的值是对象"{}"，那么把它插入到 keyPath 中
	// 扫描完它的所有子 key 后将其从 keyPath 中移除
	keyName string
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
	scanStateEndObject     // 对象的结束标志 "}"
	scanStateEndRootObject // 根对象的结束标志 "}"
)

var RootPathName = "root"

var (
	keys JsonKeysMap
	scan = new(scanner)
)

func (jkm JsonKeysMap) Get(keyPath string) (KeySlice, error) {
	if value, ok := jkm[keyPath]; ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("no key (%v)", keyPath)
	}
}

// 扫描数据
func scanData(data []byte) error {
	dataLen := len(data)
	for i := 0; i < dataLen; i++ {
		// fmt.Printf("%v - %v - %v\n", string(data[i]), dataLen, i)

		state, err := scan.step(scan, data[i])
		if err != nil {
			return err
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
		case scanStateBeginValueWithArray:
			continue
		case scanStateBeginValueWithoutQuote:
			i--
			continue
		case scanStateValueCharacter:
			continue
		case scanStateEndValue:
			continue
		case scanStateEndValueWithArray:
			continue
		case scanStateEndValueWithoutQuote:
			i--
			continue
		case scanStateEndObject:
			continue
		case scanStateEndRootObject:
			continue
		}
	}

	return nil
}

// 1. 从JSON数据解析key
func ParseFromData(data []byte) (JsonKeysMap, error) {
	// 初始化
	keys = JsonKeysMap{RootPathName: []string{}}
	scan.reset()

	err := scanData(data)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// 2. 从JSON文件解析key
func ParseFromFile(file string) (JsonKeysMap, error) {
	fileObj, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer fileObj.Close()

	reader := bufio.NewReader(fileObj)
	bufLen := 1024
	buf := make([]byte, bufLen)

	// 初始化
	keys = JsonKeysMap{RootPathName: []string{}}
	scan.reset()

	for {
		n, err := reader.Read(buf)

		if err != nil && err != io.EOF {
			return nil, err
		}

		if n > 0 {
			scanErr := scanData(buf[:n])

			if scanErr != nil {
				return nil, err
			}
		}

		if err == io.EOF {
			break
		}
	}

	return keys, nil
}

func (s *scanner) reset() {
	s.step = stepBeginObject
	s.states = []int{}
	s.keyCharacters = []byte{}
	s.keyName = RootPathName
	s.keyPath = []string{}
}

func (s *scanner) appendKeyPath() {
	s.keyPath = append(s.keyPath, s.keyName)
}

func (s *scanner) deleteLastKeyPath() {
	s.keyPath = s.keyPath[:len(s.keyPath)-1]
}

func (s *scanner) isEndRootPath() bool {
	return len(s.keyPath) == 0
}

func (s *scanner) getFullKeyPath() string {
	fullKeyPath := s.keyPath[0]
	for i := 1; i < len(s.keyPath); i++ {
		fullKeyPath += ("." + s.keyPath[i])
	}

	return fullKeyPath
}

func (s *scanner) setOneKey() {
	key := string(s.keyCharacters)
	// 清空 keyCharacters，准备保存下一个 key
	s.keyCharacters = s.keyCharacters[0:0]
	// 保存 key
	s.keyName = key

	getFullKeyPath := s.getFullKeyPath()
	// 保存一个 key 到 全局变量 keys
	if _, ok := keys[getFullKeyPath]; !ok {
		keys[getFullKeyPath] = KeySlice{key}
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
		scan.appendKeyPath()
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
		if scan.isLastStateBackslash() {
			scan.deleteLastState()
		}
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
// 4. value 以 '{' 开始，以 '}' 结尾, 对象类型
// 判断是否是一个 value 的开始字符
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
	// 4. value 以 '{' 开始，以 '}' 结尾, 对象类型
	case '{':
		scan.step = stepBeginKey
		scan.states = append(scan.states, scanStateBeginObject)
		scan.appendKeyPath()
		return scanStateBeginObject, nil
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
			scan.deleteLastState()
			return scanStateValueCharacter, nil
		} else {
			scan.step = stepEndValue
			scan.states = append(scan.states, scanStateEndValue)
			return scanStateEndValue, nil
		}
	case '\\':
		if scan.isLastStateBackslash() {
			scan.deleteLastState()
			return scanStateValueCharacter, nil
		} else {
			scan.states = append(scan.states, scanStateBackslash)
			return scanStateValueCharacter, nil
		}
	default:
		if scan.isLastStateBackslash() {
			scan.deleteLastState()
		}
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
	// 结束对象 标志
	case '}':
		// scan.step = stepEndObject
		scan.states = append(scan.states, scanStateEndObject)
		scan.deleteLastKeyPath()
		if scan.isEndRootPath() {
			scan.step = stepEndRootObject
			return scanStateEndRootObject, nil
		} else {
			return scanStateEndObject, nil
		}
	default:
		return -1, fmt.Errorf("error json format, See stepEndValue: character(%v)", string(c))
	}
}

// 忽略根对象后面多余的字符
func stepEndRootObject(s *scanner, c byte) (int, error) {
	return 0, nil
}
