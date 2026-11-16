package infrastructure

import (
	"reflect"
	"strings"
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
	simpleKey := elem.Name()

	mu.Lock()
	defer mu.Unlock()

	eventRegistry[key] = elem
	eventRegistry[simpleKey] = elem
}

// resolveType dipakai OutboxProcessor
func resolveType(typeName string) (reflect.Type, bool) {
	mu.RLock()
	defer mu.RUnlock()

	t, ok := eventRegistry[typeName]
	if ok {
		return t, true
	}

	// Try short/simple name
	if idx := strings.LastIndex(typeName, "."); idx >= 0 && idx < len(typeName)-1 {
		short := typeName[idx+1:]
		if tShort, okShort := eventRegistry[short]; okShort {
			return tShort, true
		}
	}

	return nil, false
}
