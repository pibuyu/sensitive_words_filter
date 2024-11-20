package sensitive_words_filter

import (
	"fmt"
	"github.com/pibuyu/sensitive_words_filter/filter"
	"github.com/pibuyu/sensitive_words_filter/store"
)

type Manager struct {
	store.Store
	filter.Filter
}

func NewFilter(storeOption StoreOption, filterOption FilterOption) *Manager {
	var filterStore store.Store
	var myFilter filter.Filter

	switch storeOption.Type {
	case StoreMemory:
		filterStore = store.NewMemoryModel()

	default:
		panic("invalid store type")
	}

	switch filterOption.Type {
	case FilterDfa:
		dfaModel := filter.NewDfaModel()

		go dfaModel.Listen(filterStore.GetAddChan(), filterStore.GetDelChan())

		myFilter = dfaModel

	default:
		panic("invalid filter type")
	}

	//构造manager时读入默认的dict，后续仍然可以读入别的dict
	if err := filterStore.LoadDictPath("./dict/default_dict.txt"); err != nil {
		panic(fmt.Sprintf("failed to load default dictionary file: %v", err))
	}
	return &Manager{
		Store:  filterStore,
		Filter: myFilter,
	}
}
