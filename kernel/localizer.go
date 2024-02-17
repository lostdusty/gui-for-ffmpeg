package kernel

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"path/filepath"
	"sort"
)

type LocalizerContract interface {
	GetLanguages() []Lang
	GetMessage(localizeConfig *i18n.LocalizeConfig) string
	SetCurrentLanguage(lang Lang) error
	SetCurrentLanguageByCode(code string) error
	GetCurrentLanguage() *CurrentLanguage
	AddListener(listener LocalizerListenerContract)
}

type LocalizerListenerContract interface {
	Change(localizerService LocalizerContract)
}

type Lang struct {
	Code  string
	Title string
}

type CurrentLanguage struct {
	Lang             Lang
	localizer        *i18n.Localizer
	localizerDefault *i18n.Localizer
}

type Localizer struct {
	bundle            *i18n.Bundle
	languages         []Lang
	currentLanguage   *CurrentLanguage
	localizerListener map[int]LocalizerListenerContract
}

func NewLocalizer(directory string, languageDefault language.Tag) (*Localizer, error) {
	bundle := i18n.NewBundle(languageDefault)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	languages, err := initLanguages(directory, bundle)
	if err != nil {
		return nil, err
	}

	localizerDefault := i18n.NewLocalizer(bundle, languageDefault.String())

	return &Localizer{
		bundle:    bundle,
		languages: languages,
		currentLanguage: &CurrentLanguage{
			Lang: Lang{
				Code:  languageDefault.String(),
				Title: cases.Title(languageDefault).String(display.Self.Name(languageDefault)),
			},
			localizer:        localizerDefault,
			localizerDefault: localizerDefault,
		},
		localizerListener: map[int]LocalizerListenerContract{},
	}, nil
}

func initLanguages(directory string, bundle *i18n.Bundle) ([]Lang, error) {
	var languages []Lang

	files, err := filepath.Glob(directory + "/active.*.toml")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		lang, err := bundle.LoadMessageFile(file)
		if err != nil {
			return nil, err
		}
		title := cases.Title(lang.Tag).String(display.Self.Name(lang.Tag))
		languages = append(languages, Lang{Code: lang.Tag.String(), Title: title})
	}

	sort.Sort(languagesSort(languages))

	return languages, nil
}

func (l Localizer) GetLanguages() []Lang {
	return l.languages
}

func (l Localizer) GetMessage(localizeConfig *i18n.LocalizeConfig) string {
	message, err := l.GetCurrentLanguage().localizer.Localize(localizeConfig)
	if err != nil {
		message, err = l.GetCurrentLanguage().localizerDefault.Localize(localizeConfig)
		if err != nil {
			return err.Error()
		}
	}
	return message
}

func (l Localizer) SetCurrentLanguage(lang Lang) error {
	l.currentLanguage.Lang = lang
	l.currentLanguage.localizer = i18n.NewLocalizer(l.bundle, lang.Code)
	l.eventSetCurrentLanguage()
	return nil
}

func (l Localizer) SetCurrentLanguageByCode(code string) error {
	lang, err := language.Parse(code)
	if err != nil {
		return err
	}
	title := cases.Title(lang).String(display.Self.Name(lang))
	return l.SetCurrentLanguage(Lang{Code: lang.String(), Title: title})
}

func (l Localizer) GetCurrentLanguage() *CurrentLanguage {
	return l.currentLanguage
}

func (l Localizer) AddListener(listener LocalizerListenerContract) {
	l.localizerListener[len(l.localizerListener)] = listener
}

func (l Localizer) eventSetCurrentLanguage() {
	for _, listener := range l.localizerListener {
		listener.Change(l)
	}
}

type languagesSort []Lang

func (l languagesSort) Len() int      { return len(l) }
func (l languagesSort) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l languagesSort) Less(i, j int) bool {
	return languagePriority(l[i]) < languagePriority(l[j])
}
func languagePriority(l Lang) int {
	priority := 0

	switch l.Code {
	case "ru":
		priority = -3
	case "kk":
		priority = -2
	case "en":
		priority = -1
	}

	return priority
}
