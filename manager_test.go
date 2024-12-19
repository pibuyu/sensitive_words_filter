package sensitive_words_filter

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestNewFilter(t *testing.T) {
	type args struct {
		storeOption  StoreOption
		filterOption FilterOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "memory+dfa",
			args: args{
				storeOption: StoreOption{
					Type: StoreMemory,
				},
				filterOption: FilterOption{
					Type: FilterDfa,
				},
			},
		},
	}
	var text = "我是一个阳光开朗大男孩，兴趣爱好有打羽毛球，听女声毒物音乐"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var memBefore, memAfter runtime.MemStats
			runtime.ReadMemStats(&memBefore)
			filterManager := NewFilter()
			runtime.ReadMemStats(&memAfter)

			// 计算内存差值
			dfaMemoryUsage := (memAfter.Alloc - memBefore.Alloc) / (1024 * 1024)
			t.Logf("DFA树占用的内存: %d MB", dfaMemoryUsage)

			isSensitive := filterManager.IsSensitive(text)
			t.Logf("IsSensitive= %v", isSensitive)

			matchedAll := filterManager.FindAll(text)
			t.Logf("所有敏感词为：%v", matchedAll)

			replaceResult := filterManager.Replace(text, '*')
			t.Logf("替换结果为:%v", replaceResult)

			//读入新的dict
			if err := filterManager.LoadDictPath("./dict/test_dict.txt"); err != nil {
				t.Errorf("读入失败:%v", err)
			}

			matchedAllAfterLoad := filterManager.FindAll(text)
			t.Logf("所有敏感词为：%v", matchedAllAfterLoad)
		})
	}
}

func TestFilterOutputFile(t *testing.T) {
	type args struct {
		storeOption  StoreOption
		filterOption FilterOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "memory+dfa",
			args: args{
				storeOption: StoreOption{
					Type: StoreMemory,
				},
				filterOption: FilterOption{
					Type: FilterDfa,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 读取 output.txt 文件内容
			file, err := os.Open("./dict/output.txt")
			if err != nil {
				t.Fatalf("打开文件失败: %v", err)
			}
			defer file.Close()

			reader := bufio.NewReader(file)
			text, err := reader.ReadString(0) // 将整个文件内容读取为字符串
			if err != nil && err.Error() != "EOF" {
				t.Fatalf("读取文件失败: %v", err)
			}

			// 初始化敏感词过滤器
			var memBefore, memAfter runtime.MemStats
			runtime.ReadMemStats(&memBefore)
			filterManager := NewFilter()
			runtime.ReadMemStats(&memAfter)

			// 计算敏感词过滤器初始化占用内存
			dfaMemoryUsage := (memAfter.Alloc - memBefore.Alloc) / (1024 * 1024)
			t.Logf("DFA树占用的内存: %d MB", dfaMemoryUsage)

			// 开始检测
			startTime := time.Now()
			matchedAll := filterManager.FindAll(text)
			endTime := time.Now()

			// 计算检测所用时间和字符检测速度
			elapsedTime := endTime.Sub(startTime).Seconds()
			totalChars := len(text)
			charsPerSecond := float64(totalChars) / elapsedTime

			t.Logf("检测到的敏感词为：%v", matchedAll)
			t.Logf("总检测时间为: %.2f 秒", elapsedTime)
			t.Logf("检测字符数: %d", totalChars)
			t.Logf("每秒检测字符数: %.2f", charsPerSecond)
		})
	}
}

// benchmark
func BenchmarkIsSensitivePerformance(b *testing.B) {
	filter := NewFilter()
	text := "我是一个阳光开朗大男孩，兴趣爱好有打羽毛球，听女声毒物音乐"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.IsSensitive(text)
	}
	b.StopTimer()
}

func BenchmarkReplacePerformance(b *testing.B) {
	filter := NewFilter()
	text := "我是一个阳光开朗大男孩，兴趣爱好有打羽毛球，听女声毒物音乐"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Replace(text, '*')
	}
	b.StopTimer()
}

func TestInsert(t *testing.T) {
	testFilePath := "./dict/test.txt"
	dictFilePath := "./dict/default_dict.txt"
	outputFilePath := "./dict/output.txt"

	// 读取 test.txt 文件内容
	originalText, err := readFile(testFilePath)
	if err != nil {
		fmt.Printf("读取文件 %s 失败: %v\n", testFilePath, err)
		return
	}

	// 读取 default_dict.txt 文件中的敏感词
	sensitiveWords, err := readSensitiveWords(dictFilePath)
	if err != nil {
		fmt.Printf("读取文件 %s 失败: %v\n", dictFilePath, err)
		return
	}

	// 将敏感词随机插入到文本中
	modifiedText := insertSensitiveWords(originalText, sensitiveWords)

	// 将结果写入到新的文件中
	err = writeFile(outputFilePath, modifiedText)
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}

	fmt.Printf("成功将敏感词插入到文本中，结果已保存到 %s\n", outputFilePath)

}

// 读取文本文件内容
func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content.String(), nil
}

// 读取敏感词文件
func readSensitiveWords(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// 随机插入敏感词
func insertSensitiveWords(text string, sensitiveWords []string) string {
	rand.Seed(time.Now().UnixNano())
	words := strings.Fields(text) // 将原始文本分割成单词
	totalWords := len(words)

	// 决定插入敏感词的位置
	insertPositions := rand.Perm(totalWords)
	if len(sensitiveWords) > totalWords {
		insertPositions = insertPositions[:totalWords]
	} else {
		insertPositions = insertPositions[:len(sensitiveWords)]
	}

	for i, pos := range insertPositions {
		if pos >= totalWords {
			pos = totalWords - 1
		}
		words = append(words[:pos+1], words[pos:]...) // 防止越界插入
		words[pos+1] = sensitiveWords[i%len(sensitiveWords)]
	}

	return strings.Join(words, " ")
}

// 写入文件
func writeFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}

	writer.Flush()
	return nil
}
