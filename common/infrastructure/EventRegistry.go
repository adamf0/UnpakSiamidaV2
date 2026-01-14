package infrastructure

import (
	"reflect"
	"sync"
)

var (
	eventRegistry = make(map[string]reflect.Type)
	mu            sync.RWMutex
)

// CanonicalTypeName = PkgPath + "." + Name
func CanonicalTypeName(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath() + "." + t.Name()
}

// RegisterDomainEvent WAJIB pointer
func RegisterDomainEvent(event any) {
	t := reflect.TypeOf(event)
	if t.Kind() != reflect.Pointer {
		panic("domain event must be pointer")
	}

	elem := t.Elem()
	key := elem.PkgPath() + "." + elem.Name()

	mu.Lock()
	defer mu.Unlock()

	eventRegistry[key] = elem
}

// resolveType dipakai OutboxProcessor
func resolveType(typeName string) reflect.Type {
	mu.RLock()
	defer mu.RUnlock()

	t, ok := eventRegistry[typeName]
	if !ok {
		panic("domain event not registered: " + typeName)
	}
	return t
}
