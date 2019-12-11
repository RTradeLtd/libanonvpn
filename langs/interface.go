package language

type Message struct {
	ID      string
	Message string
	index   uint
}

type Lang interface {
	Get(id interface{}) (bool, string)
	Code() string
}

type LanguageMux interface {
	Add(lang Language) (string, error)
	Del(lang Language) error
	DelID(lang string) error
	SetLang(lang string) error
	Get(id string) Language
}
