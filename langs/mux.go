package language

type LanguageLibrary struct {
	current  string
	language map[string]*Language
}

func (l *LanguageLibrary) Add(lang *Language) (string, error) {
	if l.language[lang.Code()] == nil {
		l.language[lang.Code()] = lang
	}
	return "", nil
}
func (l *LanguageLibrary) Del(lang Language) error {
	if l.language[lang.Code()] != nil {
		delete(l.language, lang.Code())
	}
	return nil
}
func (l *LanguageLibrary) DelID(lang string) error {
	if l.language[lang] != nil {
		delete(l.language, lang)
	}
	return nil
}
func (l *LanguageLibrary) SetLang(lang string) error {
	l.current = lang
	return nil
}
func (l *LanguageLibrary) Get(id string) Language {
	return *l.language[id]
}
