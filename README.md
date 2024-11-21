<h1>基于DFA+本地内存的敏感词过滤器实现。</h1>

<h2>Example</h2>

```go
import (
	"testing"
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
			filterManager := NewFilter()

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
```
<h2>How to use?</h2>
<code>go get -u "github.com/pibuyu/sensitive_words_filter"</code>
<h3>API</h3>
<h4>LoadDictPath(path ...string)error   	读入自定义敏感词词典</h4>
<h4>IsSensitive(text string)bool        	检测给定文本中是否包含敏感词</h4>
<h4>FindAll(text string)string 			找出给定文本中的所有敏感词</h4>
<h4>Replace(text string,replace rune)string 	将给定文本中的敏感词替换为指定字符</h4>
<h4>Remove(text string)string 			删除给定文本中的所有敏感词</h4>
