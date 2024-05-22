package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	folders          [][]string
	nFolders         int
	updateInProgress bool
	currentFolder    int
	quit             bool
)

func main() {

	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

func newModel() model {
	return model{
		fName:     0,
		mark:      1,
		time:      2,
		xMark:     "",
		checkMark: "󰗠",
		dash:      "",
		spinner:   spinner.Model{Spinner: randomSpinner()},
	}
}

type model struct {
	fName int
	mark  int
	time  int // indices for current row

	xMark     string
	checkMark string
	dash      string
	spinner   spinner.Model
}

func (m model) Init() tea.Cmd {
	m.getFolders()
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if !quit && currentFolder < nFolders {
		var cmd tea.Cmd
		if updateInProgress == false {
			go m.updateCurrentDir()
		}

		m.setMarker(m.spinnerString())
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}

	return m, tea.Quit

}

func (m model) View() string {

    if nFolders == 0 {
        quit = true
        return "no git folders found" + "\n"
    }

	return m.tableString() + "\n"
}

func (m model) getFolders() {

	cwd, _ := os.Getwd()
	entries, _ := os.ReadDir(cwd)

	for _, e := range entries {
		if e.IsDir() && isGitDir(e.Name()) {
			folders = append(folders, []string{e.Name(), m.dash, m.dash})
			nFolders++
		}
	}

}

func (m model) updateCurrentDir() {

	updateInProgress = true
	cwd, _ := os.Getwd()
	os.Chdir(m.getFolderName())

	currTime := time.Now()
	out, err := exec.Command("git", "pull").Output()
	timeTook := time.Since(currTime)

	os.Chdir(cwd)
	m.setMarker(m.decideMarker(out, err))
	m.setTime(timeTook.String())
	currentFolder, updateInProgress = currentFolder+1, false

}

func (m model) tableString() string {

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Faint(true)).
		StyleFunc(func(row, col int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)
			if col == 1 {
				style = style.Align(lipgloss.Center)
			}
			return style
		}).
		Headers(" Directory ", " Update ", " Time ").
		Rows(folders...)

	return t.Render()
}

func (m model) getFolderName() string {
	return folders[currentFolder][m.fName]
}

func (m model) setMarker(marker string) {
	folders[currentFolder][m.mark] = marker
}

func (m model) setTime(time string) {
	folders[currentFolder][m.time] = time
}

func (m model) spinnerString() string {
	return m.spinner.View()
}

func (m model) decideMarker(out []byte, err error) string {

	if string(out) == "Already up to date.\n" {
		return "up-to-date"
	}

	if err == nil {
		return m.checkMark
	} else {
		return m.xMark
	}

}

func isGitDir(name string) bool {

	entries, _ := os.ReadDir(name)
	for _, e := range entries {
		if e.IsDir() && e.Name() == ".git" {
			return true
		}
	}

	return false

}

func logToFile(line ...any) {

	f, err := os.OpenFile("updateLogFile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(line...)

}

func randomSpinner() spinner.Spinner {
	spinners := [...]spinner.Spinner{
		// spinner.Dot,
		// spinner.Jump,
		// spinner.Line,
		// spinner.Globe,
		spinner.Moon,
		spinner.Meter,
		// spinner.Pulse,
		spinner.Monkey,
		spinner.Points,
		// spinner.Hamburger,
		// spinner.MiniDot,
		// spinner.Ellipsis,
		gamePadSpinner(),
		circleSliceSpinner(),
	}

	return spinners[rand.Intn(len(spinners))]

}

func circleSliceSpinner() spinner.Spinner {
	return spinner.Spinner{
		Frames: []string{"󰪞", "󰪟", "󰪠", "󰪡", "󰪢", "󰪣", "󰪤", "󰪥"},
		FPS:    time.Second / 8,
	}
}

func gamePadSpinner() spinner.Spinner {
	return spinner.Spinner{
		Frames: []string{"󰸴", "󰸵", "󰸸", "󰸷"},
		FPS:    time.Second / 4,
	}
}
