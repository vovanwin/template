package config

// Store возвращает бэкенд для прямого доступа (используется в UI и тестах)
func (f *Flags) Store() FlagStore {
	return f.store
}
