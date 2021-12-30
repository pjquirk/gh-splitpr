package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh"
)

type SplitSettings struct {
	CloneDir       string
	NewBranchName  string
	BaseBranchName string
	NewPrTitle     string
	BodyFile       string
}

type prSettingsRaw struct {
	BaseRefName string `json:"baseRefName"`
	HeadRefName string `json:"headRefName"`
	Title       string `json:"title"`
	Body        string `json:"body"`
}

func GetSplitSettings(repository gh.Repository, pullRequestId int, commits []Commit) tea.Msg {
	nwo := ToNwo(repository)
	stdOut, _, err := gh.Exec("pr", "view", strconv.Itoa(pullRequestId), "--repo", nwo, "--json", "baseRefName,headRefName,title,body")
	if err != nil {
		return ErrMsg{err}
	}

	var prSettings prSettingsRaw
	err = json.Unmarshal(stdOut.Bytes(), &prSettings)
	if err != nil {
		return ErrMsg{err}
	}

	if len(commits) < 1 {
		return ErrMsg{errors.New("Cannot split a PR without any commits selected")}
	}
	newBranchName := fmt.Sprintf("%s-split-%.8s", prSettings.HeadRefName, commits[0].Sha)

	cloneDir, err := os.MkdirTemp("", "gh-splitpr")
	if err != nil {
		return ErrMsg{err}
	}

	bodyFile, err := os.CreateTemp(cloneDir, ".gh-splitpr-body")
	if err != nil {
		return ErrMsg{err}
	}
	defer closeFile(bodyFile)
	bodyFile.WriteString(fmt.Sprintf("Split from #%d\n\n", pullRequestId) + prSettings.Body)

	return SplitSettings{
		CloneDir:       cloneDir,
		NewBranchName:  newBranchName,
		BaseBranchName: prSettings.BaseRefName,
		NewPrTitle:     fmt.Sprintf("Split from #%d - %s", pullRequestId, prSettings.Title),
		BodyFile:       bodyFile.Name(),
	}
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		panic(err)
	}
}
