package localizer

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"path/filepath"
	"sort"
)

type ServiceContract interface {
	GetLanguages() []Lang
	GetMessage(localizeConfig *i18n.LocalizeConfig) string
	SetCurrentLanguage(lang Lang) error
	GetCurrentLanguage() *CurrentLanguage
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

type Service struct {
	bundle          *i18n.Bundle
	languages       []Lang
	currentLanguage *CurrentLanguage
}

func NewService(directory string, languageDefault language.Tag) (*Service, error) {
	bundle := i18n.NewBundle(languageDefault)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	languages, err := initLanguages(directory, bundle)
	if err != nil {
		return nil, err
	}

	localizerDefault := i18n.NewLocalizer(bundle, languageDefault.String())

	return &Service{
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
	}, nil
}

func initLanguages(directory string, bundle *i18n.Bundle) ([]Lang, error) {
	var languages []Lang

	files, err := filepath.Glob(directory + "/active.*.toml")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		language, err := bundle.LoadMessageFile(file)
		if err != nil {
			return nil, err
		}
		title := cases.Title(language.Tag).String(display.Self.Name(language.Tag))
		languages = append(languages, Lang{Code: language.Tag.String(), Title: title})
	}

	sort.Sort(languagesSort(languages))

	return languages, nil
}

func (s Service) GetLanguages() []Lang {
	return s.languages
}

func (s Service) GetMessage(localizeConfig *i18n.LocalizeConfig) string {
	message, err := s.GetCurrentLanguage().localizer.Localize(localizeConfig)
	if err != nil {
		message, err = s.GetCurrentLanguage().localizerDefault.Localize(localizeConfig)
		if err != nil {
			return err.Error()
		}
	}
	return message
}

func (s Service) SetCurrentLanguage(lang Lang) error {
	s.currentLanguage.Lang = lang
	s.currentLanguage.localizer = i18n.NewLocalizer(s.bundle, lang.Code)
	return nil
}

func (s Service) GetCurrentLanguage() *CurrentLanguage {
	return s.currentLanguage
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
