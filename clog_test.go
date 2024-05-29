package clog_test

type Entry map[string]any

func (e Entry) GetAny(key string) (any, bool) {
	v, ok := e[key]
	if !ok {
		return "", false
	}

	return v, true
}

func (e Entry) GetString(key string) (string, bool) {
	v, ok := e.GetAny(key)
	if !ok {
		return "", false
	}

	s, ok := v.(string)
	if !ok {
		return "", false
	}

	return s, true
}

func (e Entry) GetMap(key string) (map[string]any, bool) {
	v, ok := e.GetAny(key)
	if !ok {
		return nil, false
	}

	m, ok := v.(map[string]any)
	if !ok {
		return nil, false
	}
	return m, true
}
