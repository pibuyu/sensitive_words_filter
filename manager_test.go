package sensitive_words_filter

import (
	"fmt"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterManager := NewFilter(tt.args.storeOption, tt.args.filterOption)

			err := filterManager.LoadDictPath("./dict/default_dict.txt")
			if err != nil {
				fmt.Println(err)
			}

			//err = filterManager.AddWord("小行家", "敏感词2", "敏感词3")
			//if err != nil {
			//	t.Errorf("add sensitive word failed, err: %v", err)
			//}

			isSensitive := filterManager.IsSensitive("啦啦啦，我是快乐的小行家")
			t.Logf("IsSensitive= %v", isSensitive)

			//filtered := filterManager.Remove("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词")
			//if !reflect.DeepEqual(filtered, "这是,这是,这是,这是,这里没有敏感词") {
			//	t.Errorf("Remove() = %v, want %v", filtered, "这是,这是,这是,这是,这里没有敏感词")
			//}
			//
			//replaced := filterManager.Replace("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词", '*')
			//if !reflect.DeepEqual(replaced, "这是****,这是****,这是****,这是****,这里没有敏感词") {
			//	t.Errorf("Replace() = %v, want %v", replaced, "这是****,这是****,这是****,这是****,这里没有敏感词")
			//}
			//
			//matchedOne := filterManager.FindOne("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词")
			//if !reflect.DeepEqual(matchedOne, "敏感词1") {
			//	t.Errorf("FindOne() = %v, want %v", matchedOne, "敏感词1")
			//}
			//
			matchedAll := filterManager.FindAll("啦啦啦")
			t.Logf("所有敏感词为：%v", matchedAll)
			//if !reflect.DeepEqual(matchedAll, []string{"敏感词1", "敏感词2", "敏感词3"}) {
			//	t.Errorf("FindAll() = %v, want %v", matchedAll, []string{"敏感词1", "敏感词2", "敏感词3"})
			//}
			//
			//matchedMap := filterManager.FindAllCount("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词")
			//t.Logf("查找敏感词数量的map为：%v", matchedMap)
			//if !reflect.DeepEqual(matchedMap, map[string]int{"敏感词1": 2, "敏感词2": 1, "敏感词3": 1}) {
			//	t.Errorf("FindAllCount() = %v, want %v", matchedMap, map[string]int{"敏感词1": 2, "敏感词2": 1, "敏感词3": 1})
			//}
		})
	}
}
