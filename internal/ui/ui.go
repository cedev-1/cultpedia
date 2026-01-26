package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cultpedia/internal/actions"
	"cultpedia/internal/models"
	"cultpedia/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	keyCtrlC = "ctrl+c"
	keyEnter = "enter"
	keyEsc   = "esc"

	qtypeSingleChoice = "single_choice"
	qtypeTrueFalse    = "true_false"
)

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("5")).
	Bold(true).
	MarginTop(2)

var versionStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("6")).
	Italic(true)

var successStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("2")).
	Bold(true)

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("1")).
	Bold(true)

var infoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("4"))

var boxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1).
	MarginTop(1).
	MarginBottom(1)

type confirmModel struct {
	question     models.Question
	cursor       int
	choices      []string
	version      string
	questionType string
}

type previewModel struct {
	question      models.Question
	version       string
	languageIndex int
	languages     []string
	questionType  string
}

type helpModel struct {
	version      string
	questionType string
}

type versionCheckModel struct {
	version string
}

type warningModel struct {
	message      string
	version      string
	questionType string
}

type versionCheckMsg struct {
	version       string
	remoteVersion string
	err           error
}

type mainMenuModel struct {
	cursor       int
	choices      []string
	message      string
	version      string
	questionType string
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m previewModel) Init() tea.Cmd {
	return nil
}

func (m helpModel) Init() tea.Cmd {
	return nil
}

func (m mainMenuModel) Init() tea.Cmd {
	return nil
}

func (m versionCheckModel) Init() tea.Cmd {
	return func() tea.Msg {
		remoteVersion, err := actions.GetRemoteVersion()
		return versionCheckMsg{version: m.version, remoteVersion: remoteVersion, err: err}
	}
}

func (m warningModel) Init() tea.Cmd {
	return nil
}

func InitialMainModel() versionCheckModel {

	version := "unknown"
	if data, err := os.ReadFile(utils.ManifestFile); err == nil {
		var manifest models.Manifest
		if json.Unmarshal(data, &manifest) == nil {
			version = manifest.Version
		}
	}
	return versionCheckModel{
		version: version,
	}
}

func newMainMenuModel(version string, questionType string, message string) mainMenuModel {
	return mainMenuModel{
		cursor: 0,
		choices: []string{
			"Validate new question",
			"Preview question",
			"Add question to dataset",
			"Switch question type",
		},
		message:      message,
		version:      version,
		questionType: questionType,
	}
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", keyCtrlC:
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case keyEnter:
			if m.cursor == 0 {
				message := actions.AddValidatedQuestion(m.question)
				// Reset template after successful add
				if !strings.Contains(message, "error") {
					if err := actions.ResetTemplate(m.questionType); err == nil {
						message += "\n\n✔ Template file has been reset for your next question."
					}
				}
				return newMainMenuModel(m.version, m.questionType, message), nil
			} else {
				return newMainMenuModel(m.version, m.questionType, ""), nil
			}
		}
	}
	return m, nil
}

func (m previewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left", "h":
			if m.languageIndex > 0 {
				m.languageIndex--
			}
		case "right", "l":
			if m.languageIndex < len(m.languages)-1 {
				m.languageIndex++
			}
		case keyEnter, keyEsc:
			return newMainMenuModel(m.version, m.questionType, ""), nil
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	s := titleStyle.Render("Confirm Addition") + "\n\n"

	s += "Ready to add this question?\n\n"

	s += boxStyle.Render(fmt.Sprintf("Slug: %s\n\nTheme: %s\nDifficulty: %s\nPoints: %.1f\nType: %s\n\nLanguages: ✓ fr ✓ en ✓ es\nAnswers: %d\nSources: %d",
		m.question.Slug,
		m.question.Theme.Slug,
		m.question.Difficulty,
		m.question.Points,
		m.question.Qtype,
		len(m.question.Answers),
		len(m.question.Sources)))

	s += "\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		icon := "✔"
		if choice == "No" {
			icon = "✘"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, icon, choice)
	}

	s += "\n\n" + infoStyle.Render("Commands: [↑↓] Navigate | [Enter] Confirm | [?] Help | [q] Quit")
	return s
}

func (m previewModel) View() string {
	s := titleStyle.Render("Question Preview") + "\n\n"
	currentLang := m.languages[m.languageIndex]

	langBar := "Languages: "
	for i, lang := range m.languages {
		if i == m.languageIndex {
			langBar += successStyle.Render("[" + strings.ToUpper(lang) + "]")
		} else {
			langBar += "[" + lang + "]"
		}
		if i < len(m.languages)-1 {
			langBar += " "
		}
	}
	s += langBar + "\n\n"

	s += fmt.Sprintf("  Slug: %s\n", infoStyle.Render(m.question.Slug))
	s += fmt.Sprintf("  Theme: %s\n", infoStyle.Render(m.question.Theme.Slug))

	if len(m.question.Subthemes) > 0 {
		subthemes := make([]string, len(m.question.Subthemes))
		for i, sub := range m.question.Subthemes {
			subthemes[i] = sub.Slug
		}
		s += fmt.Sprintf("  Subthemes: %s\n", infoStyle.Render(strings.Join(subthemes, ", ")))
	}

	if len(m.question.Tags) > 0 {
		tags := make([]string, len(m.question.Tags))
		for i, tag := range m.question.Tags {
			tags[i] = tag.Slug
		}
		s += fmt.Sprintf("  Tags: %s\n", infoStyle.Render(strings.Join(tags, ", ")))
	}

	s += fmt.Sprintf("  Difficulty: %s | Points: %.1f | Type: %s\n", m.question.Difficulty, m.question.Points, m.question.Qtype)
	s += fmt.Sprintf("  Languages: ✓ fr ✓ en ✓ es | Answers: %d | Sources: %d\n", len(m.question.Answers), len(m.question.Sources))

	content := m.question.I18n[currentLang]
	s += "\n" + boxStyle.Render(fmt.Sprintf("Title (%s): %s\n\nQuestion: %s\n\nExplanation: %s", strings.ToUpper(currentLang), content.Title, content.Stem, content.Explanation))

	s += "\n\nAnswers:\n"
	for _, choice := range m.question.Answers {
		correctMark := " "
		if choice.IsCorrect {
			correctMark = "✓"
		}
		s += fmt.Sprintf("  [%s] %s\n", correctMark, choice.I18n[currentLang].Label)
	}

	s += "\n\n" + infoStyle.Render("Commands: [←→] Change Language | [Enter/Esc] Back | [?] Help | [q] Quit")
	return s
}

func (m helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", keyCtrlC:
			return m, tea.Quit
		case keyEnter, keyEsc:
			return newMainMenuModel(m.version, m.questionType, ""), nil
		}
	}
	return m, nil
}

func (m helpModel) View() string {
	s := titleStyle.Render("Help") + "\n\n"

	themesStr := "science, history, geography, sports, gaming"
	if themes, err := actions.GetAvailableThemes(); err == nil && len(themes) > 0 {
		themesStr = strings.Join(themes, ", ")
	}

	helpContent := fmt.Sprintf(`Navigation:
  ↑/↓ or k/j              Move up/down in menus
  ←/→ or h/l              Switch languages (in preview)
  Enter                   Confirm selection
  Esc                     Go back to menu

Actions:
  ?                       Show this help
  q / Ctrl+C              Exit (from main menu)

Available Themes:
  %s

Tips:
  • Edit the template file matching your question type
  • All 3 languages required: fr, en, es
  • Each language needs: title, stem, explanation
  • Minimum text lengths: stem 10 chars, explanation 20 chars
  • Always provide at least 1 source URL
  • Use slug format: {theme}-{subtheme}-{key}-{detail}

For more info, see CONTRIBUTING.md in the docs/ folder.
Or visit:
  https://docs.culturae.me/cultpedia/

Thank you for contributing to Cultpedia!`, themesStr)

	s += boxStyle.Render(helpContent)

	s += "\n\n" + infoStyle.Render("Press any key to close")
	return s
}

func (m versionCheckModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case versionCheckMsg:
		if msg.err != nil || msg.remoteVersion == "" || msg.remoteVersion == msg.version {
			return newMainMenuModel(msg.version, qtypeSingleChoice, ""), nil
		} else {
			message := fmt.Sprintf("Your local dataset version (%s) is outdated.\nThe latest version on GitHub is %s.\n\nPlease run 'git pull' to update your local repository before proceeding.\n\nPress Enter to continue anyway (not recommended).", msg.version, msg.remoteVersion)
			return warningModel{
				message:      message,
				version:      msg.version,
				questionType: qtypeSingleChoice,
			}, nil
		}
	}
	return m, nil
}

func (m versionCheckModel) View() string {
	return titleStyle.Render("Checking Dataset Version") + "\n\nChecking for updates...\n\n" + infoStyle.Render("Please wait...")
}

func (m warningModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", keyEsc, keyCtrlC:
			return m, tea.Quit
		case keyEnter:
			return newMainMenuModel(m.version, m.questionType, ""), nil
		}
	}
	return m, nil
}

func (m warningModel) View() string {
	s := titleStyle.Render("Version Mismatch") + "\n\n"
	s += boxStyle.Render(m.message)
	s += "\n\n" + infoStyle.Render("Commands: [Enter] Continue Anyway | [Esc/q] Quit")
	return s
}

func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", keyCtrlC:
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case keyEnter:
			switch m.cursor {
			case 0:
				_, err := actions.ValidateNewQuestionWithType(m.questionType)
				if err != nil {
					m.message = err.Error()
				} else {
					m.message = successStyle.Render("✔ New question is valid!")
				}
			case 1:
				question, err := actions.ValidateNewQuestionWithType(m.questionType)
				if err != nil {
					m.message = err.Error()
				} else {
					return previewModel{
						question:      question,
						languageIndex: 1,
						languages:     []string{"fr", "en", "es"},
						version:       m.version,
						questionType:  m.questionType,
					}, nil
				}
			case 2:
				question, err := actions.ValidateNewQuestionWithType(m.questionType)
				if err != nil {
					m.message = err.Error()
				} else {
					return confirmModel{
						question:     question,
						cursor:       0,
						choices:      []string{"Yes", "No"},
						version:      m.version,
						questionType: m.questionType,
					}, nil
				}
			case 3:
				if m.questionType == qtypeSingleChoice {
					m.questionType = qtypeTrueFalse
					m.message = "Switched to True/False mode\nEdit: datasets/new-question-true-false.json"
				} else {
					m.questionType = qtypeSingleChoice
					m.message = "Switched to Single Choice mode\nEdit: datasets/new-question.json"
				}
			}
		case "?":
			return helpModel{version: m.version, questionType: m.questionType}, nil
		}
	}
	return m, nil
}

func (m mainMenuModel) View() string {
	title := titleStyle.Render("Welcome to Cultpedia")
	s := title + "\n"

	versionStr := versionStyle.Render(fmt.Sprintf("Database version: %s", m.version))
	s += versionStr + "\n"

	modeDisplay := qtypeSingleChoice
	templateFile := "datasets/new-question.json"
	if m.questionType == qtypeTrueFalse {
		modeDisplay = qtypeTrueFalse
		templateFile = "datasets/new-question-true-false.json"
	}
	s += infoStyle.Render(fmt.Sprintf("Mode: %s | Template: %s", modeDisplay, templateFile)) + "\n\n"

	s += "What would you like to do?\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	if m.message != "" {
		s += "\n"
		isError := strings.Contains(m.message, "✗") || strings.Contains(m.message, "error") || strings.Contains(m.message, "missing") || strings.Contains(m.message, "duplicate") || strings.Contains(m.message, "Invalid") || strings.Contains(m.message, "too short") || strings.Contains(m.message, "between")

		if isError {
			s += errorStyle.Render("✗ Validation Failed") + "\n"
			s += boxStyle.Render(m.message)
		} else if strings.Contains(m.message, "valid") {
			s += successStyle.Render("✔ Success") + "\n"
			s += boxStyle.Render(m.message)
		} else {
			s += boxStyle.Render(m.message)
		}
	}

	s += "\n\n" + infoStyle.Render("Commands: [↑↓] Navigate | [Enter] Select | [?] Help | [q] Quit")
	return s
}
