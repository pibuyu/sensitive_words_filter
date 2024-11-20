package sensitive_words_filter

import (
	"sensitive_words_filter/filter"
	"sensitive_words_filter/store"
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

	return &Manager{
		Store:  filterStore,
		Filter: myFilter,
	}
}
