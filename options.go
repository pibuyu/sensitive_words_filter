package sensitive_words_filter

const (
	StoreMemory = iota
)

const (
	FilterDfa = iota
)

type StoreOption struct {
	Type uint32
}

type FilterOption struct {
	Type uint32
}
