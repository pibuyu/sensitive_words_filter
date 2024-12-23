package sensitive_words_filter

import (
	"embed"
	"fmt"
	"github.com/pibuyu/sensitive_words_filter/filter"
	"github.com/pibuyu/sensitive_words_filter/store"
)

type Manager struct {
	store.Store
	filter.Filter
}

//go:embed dict/default_dict.txt
var defaultDict embed.FS

// NewFilter 原本的方法签名 NewFilter(storeOption StoreOption, filterOption FilterOption)
// 现在仅支持这俩默认选项，那就先不要求用户传递参数过来，后续如果有了更多可选的选项，可以提供一个NewFilterWithOptions方法
func NewFilter() *Manager {
	var filterStore store.Store
	var myFilter filter.Filter

	filterStore = store.NewMemoryModel()

	dfaModel := filter.NewDfaModel()
	go dfaModel.Listen(filterStore.GetAddChan(), filterStore.GetDelChan())
	myFilter = dfaModel

	//switch storeOption.Type {
	//case StoreMemory:
	//	filterStore = store.NewMemoryModel()
	//
	//default:
	//	panic("invalid store type")
	//}
	//
	//switch filterOption.Type {
	//case FilterDfa:
	//	dfaModel := filter.NewDfaModel()
	//
	//	go dfaModel.Listen(filterStore.GetAddChan(), filterStore.GetDelChan())
	//
	//	myFilter = dfaModel

	//default:
	//	panic("invalid filter type")
	//}

	//初始化Filter对象时读入默认dict文件
	//将txt文件静态嵌入到项目中
	dictFile, err := defaultDict.Open("dict/default_dict.txt")
	if err != nil {
		panic(fmt.Sprintf("failed to open embedded dictionary: %v", err))
	}
	defer dictFile.Close()

	// 使用 dictFile (io.Reader) 加载字典
	if err := filterStore.LoadDict(dictFile); err != nil {
		panic(fmt.Sprintf("failed to load dictionary: %v", err))
	}

	return &Manager{
		Store:  filterStore,
		Filter: myFilter,
	}
}
