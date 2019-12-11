package language

type Language struct {
	code     string
	messages map[string]Message
}

func (l *Language) Code() string {
	return l.code
}

func (l *Language) SetCode(in string) string {
	l.code = in
	return l.code
}

func (l *Language) getInt(in uint) (bool, string) {
	for _, v := range l.messages {
		if v.index == in {
			return true, v.Message
		}
	}
	return false, ""
}

func (l *Language) getID(in string) (bool, string) {
	for _, v := range l.messages {
		if v.ID == in {
			return true, v.Message
		}
	}
	return false, ""
}

func (l *Language) getLang(in Message) (bool, string) {
	for _, v := range l.messages {
		if v == in {
			return true, v.Message
		}
	}
	return false, ""
}

func (l *Language) Get(id interface{}) (bool, string) {
	switch v := id.(type) {
	case int:
		return l.getInt(uint(v))
	case uint:
		return l.getInt(v)
	case string:
		return l.getID(v)
	case Message:
		return l.getLang(v)
	default:
		return false, ""
	}
}
