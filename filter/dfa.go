package filter

import "sync"

type dfaNode struct {
	children map[rune]*dfaNode
	isLeaf   bool
}

func newDfaNode() *dfaNode {
	return &dfaNode{
		children: make(map[rune]*dfaNode),
		isLeaf:   false,
	}
}

type DfaModel struct {
	root *dfaNode
	mu   sync.RWMutex //加一个读写互斥锁
}

func NewDfaModel() *DfaModel {
	return &DfaModel{
		root: newDfaNode(),
	}
}

func (m *DfaModel) AddWords(words ...string) {
	//m.mu.Lock()
	//defer m.mu.Unlock() //批量写入之前先获取锁
	for _, word := range words {
		m.AddWord(word)
	}
}

func (m *DfaModel) AddWord(word string) {
	now := m.root
	runes := []rune(word)

	//遍历每个字符，如果字符存在于children中，移动到对应的子节点；否则创建新的子节点并添加到children中
	for _, r := range runes {
		if next, ok := now.children[r]; ok {
			now = next
		} else {
			next = newDfaNode()
			now.children[r] = next
			now = next
		}
	}

	now.isLeaf = true
}

func (m *DfaModel) DelWords(words ...string) {
	//批量删除之前同样获取锁
	//m.mu.Lock()
	//defer m.mu.Unlock()
	for _, word := range words {
		m.DelWord(word)
	}
}

func (m *DfaModel) DelWord(word string) {
	var lastLeaf *dfaNode
	var lastLeafNextRune rune
	now := m.root
	runes := []rune(word)

	//如果字符不存在于children中，说明这个单词不在DFA树中，直接返回；如果存在，移动到对应的子节点
	for _, r := range runes {
		if next, ok := now.children[r]; !ok {
			return
		} else {
			//记录遍历到的最后一个叶子节点，及其对应的下一个字符。
			//遍历结束之后，从最后一个叶子节点的children中删除lastLeafNextRune，移除该单词
			if now.isLeaf {
				lastLeaf = now
				lastLeafNextRune = r
			}
			now = next
		}
	}

	delete(lastLeaf.children, lastLeafNextRune)
}

func (m *DfaModel) Listen(addChan, delChan <-chan string) {
	go func() {
		for word := range addChan {
			m.AddWord(word)
		}
	}()

	go func() {
		for word := range delChan {
			m.DelWord(word)
		}
	}()
}

func (m *DfaModel) FindAll(text string) []string {
	var matches []string // stores words that match in dict
	var found bool       // if current rune in node's map
	var now *dfaNode     // current node

	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]

		if !found {
			parent = m.root
			pos = start
			start++
			continue
		}

		if now.isLeaf && start <= pos {
			matches = append(matches, string(runes[start:pos+1]))
		}

		if pos == length-1 {
			parent = m.root
			pos = start
			start++
			continue
		}

		parent = now
	}

	var res []string
	set := make(map[string]struct{})

	for _, word := range matches {
		if _, ok := set[word]; !ok {
			set[word] = struct{}{}
			res = append(res, word)
		}
	}

	return res
}

func (m *DfaModel) FindAllCount(text string) map[string]int {
	res := make(map[string]int)
	var found bool
	var now *dfaNode

	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]

		//如果当前字符无法匹配，重置匹配状态，从起始位置的下一个字符重新开始
		if !found {
			parent = m.root
			pos = start
			start++
			continue
		}

		//如果当前字符是叶子节点，并且起始位置小于等于当前位置，则将当前单词加入结果
		if now.isLeaf && start <= pos {
			res[string(runes[start:pos+1])]++
		}

		if pos == length-1 {
			parent = m.root
			pos = start
			start++
			continue
		}

		parent = now
	}

	return res
}

func (m *DfaModel) FindOne(text string) string {
	var found bool
	var now *dfaNode

	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]

		if !found || (!now.isLeaf && pos == length-1) {
			parent = m.root
			pos = start
			start++
			continue
		}
		//和findAll逻辑基本一致，不过这里不记录具体单词，旨在isLeaf==true时让count++
		if now.isLeaf && start <= pos {
			return string(runes[start : pos+1])
		}

		parent = now
	}

	return ""
}

func (m *DfaModel) IsSensitive(text string) bool {
	return m.FindOne(text) != ""
}

func (m *DfaModel) Replace(text string, repl rune) string {
	var found bool
	var now *dfaNode

	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]

		if !found || (!now.isLeaf && pos == length-1) {
			parent = m.root
			pos = start
			start++
			continue
		}

		//如果匹配上了，从start到pos中间的字符都得换成repl
		if now.isLeaf && start <= pos {
			for i := start; i <= pos; i++ {
				runes[i] = repl
			}
		}

		parent = now
	}

	return string(runes)
}

func (m *DfaModel) Remove(text string) string {
	var found bool
	var now *dfaNode

	start := 0 // 从文本的第几个文字开始匹配
	parent := m.root
	runes := []rune(text)
	length := len(runes)
	filtered := make([]rune, 0, length)

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]

		if !found || (!now.isLeaf && pos == length-1) {
			filtered = append(filtered, runes[start])
			parent = m.root
			pos = start
			start++
			continue
		}

		if now.isLeaf {
			start = pos + 1
			parent = m.root
		} else {
			parent = now
		}
	}
	//	如果匹配上，则跳过当前单词，不将这个单词加入到filtered结果中
	filtered = append(filtered, runes[start:]...)

	return string(filtered)
}
