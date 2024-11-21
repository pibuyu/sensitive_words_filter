<h1>基于DFA+本地内存的敏感词过滤器实现。</h1>

<h2>How to run?</h2>

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
