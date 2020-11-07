package common

// Status generalizes operation result.
type Status struct {
	Success   bool
	Operation string
}

// StatusList holds values of statuses for a specific
// key-value mapping.
type StatusList struct {
	Data map[string]*Status
}

// NewStatusList initializes a new list.
func NewStatusList() *StatusList {
	return &StatusList{
		Data: make(map[string]*Status),
	}
}

// Get returns a key from status list.
func (list *StatusList) Get(key string) (*Status, bool) {
	val, ok := list.Data[key]
	return val, ok
}

// HasKey checks if key is inside status list.
func (list *StatusList) HasKey(key string) bool {
	_, ok := list.Data[key]
	return ok
}

// Insert adds a new key, value pair to status list.
func (list *StatusList) Insert(key string, status *Status) {
	list.Data[key] = status
}

// Size returns the length of list.
func (list *StatusList) Size() int {
	return len(list.Data)
}

// Successes counts number of valid statuses inside list.
func (list *StatusList) Successes() int {
	success := 0
	for _, v := range list.Data {
		if v.Success {
			success++
		}
	}
	return success
}

// ToGeneric converts list to map of generic objects.
func (list *StatusList) ToGeneric() map[string]*Generic {
	data := make(map[string]*Generic)
	for key, value := range list.Data {
		data[key] = ObjectToGeneric(value)
	}
	return data
}
