package sensitive_words_filter

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
	var text = "这是一个阳光大男孩"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterManager := NewFilter(tt.args.storeOption, tt.args.filterOption)

			if err := filterManager.LoadDictPath("./dict/test_dict.txt"); err != nil {
				t.Errorf("读入失败:%v", err)
			}

			isSensitive := filterManager.IsSensitive(text)
			t.Logf("IsSensitive= %v", isSensitive)

			matchedAll := filterManager.FindAll(text)
			t.Logf("所有敏感词为：%v", matchedAll)

		})
	}
}
