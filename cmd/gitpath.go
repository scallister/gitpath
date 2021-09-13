package cmd

/*
Copyright © 2021 Steven Callister

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/yookoala/realpath"
)

func GitPathCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("received %d instead of 1 arguments", len(args))
	}
	fPath := args[0]

	// Setup git library
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to setup repository")
		return err
	}

	// Get Branch Name
	branchName, err := currentBranchName(repo)
	if err != nil {
		return err
	}
	log.Trace().Str("branchName", branchName).Msg("")

	// Determine base url
	baseUrl, err := getBaseUrl(repo)
	if err != nil {
		return err
	}
	log.Trace().Str("baseURL", baseUrl).Msg("")

	// Get relative filepath - to determine what we append to the URL
	repoPath, err := getRepoRoot(repo)
	if err != nil {
		return err
	}
	log.Trace().Str("repoPath", repoPath).Msg("")
	fullPath, err := realpath.Realpath(fPath)
	if err != nil {
		return nil
	}
	relativePath := getRelativePath(fullPath, repoPath)

	// Generate/Print URL
	url := createURL(baseUrl, branchName, relativePath)
	fmt.Printf("%s\n", url)

	return nil
}

// relativePath - Take a full path to a file, and a repository root, and return the relative repository path to the file
func getRelativePath(fullPath, repoPath string) string {
	relativePath := strings.TrimPrefix(fullPath, repoPath)
	relativePath = strings.TrimPrefix(relativePath, "/")
	log.Trace().Str("fullPath", fullPath).
		Str("repoPath", repoPath).
		Str("relativePath", relativePath).
		Msg("")
	return relativePath
}

// createURL - Takes several items and creates a working GitHub / GitLab url
func createURL(baseURL, branchName, relativePath string) string {
	var url string
	if strings.Contains(baseURL, "github") {
		url = fmt.Sprintf("%s/blob/%s/%s", baseURL, branchName, relativePath)
	} else if strings.Contains(baseURL, "gitlab") {
		url = fmt.Sprintf("%s/-/blob/%s/%s", baseURL, branchName, relativePath)
	} else {
		err := fmt.Errorf("unsure if github or gitlab repository, baseURL is %s", baseURL)
		log.Error().Err(err).Msg("Unsure if this is a github or gitlab repository. Please rerun with --verbose and file a GitHub issue.")
		panic(err)
	}
	return url
}

func getRepoRoot(repo *git.Repository) (string, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	return wt.Filesystem.Root(), nil
}

func getBaseUrl(repo *git.Repository) (string, error) {
	remote, err := repo.Remote("origin")
	urls := remote.Config().URLs
	log.Trace().Strs("remoteURLS", urls).Msg("")
	if len(urls) != 1 {
		err := fmt.Errorf("expected an array containig one remote URL, received %s instead", urls)
		return "", err
	}
	baseUrl := trimUrlToBase(urls[0])

	return baseUrl, err
}

func trimUrlToBase(url string) string {
	urlTrimmed := url
	if strings.Contains(urlTrimmed, "@") {
		urlTrimmed = strings.TrimPrefix(url, "git@")
		urlTrimmed = strings.Replace(urlTrimmed, ":", "/", 1)
	}
	urlTrimmed = strings.TrimSuffix(urlTrimmed, ".git")
	if !strings.HasPrefix(urlTrimmed, "https://") {
		urlTrimmed = "https://" + urlTrimmed
	}
	return urlTrimmed
}

func currentBranchName(repo *git.Repository) (string, error) {
	head, err := repo.Head()
	if err != nil {
		err = errors.Wrap(err, "failed to get repository Head")
		return "", err
	}

	// Check if this is a branch
	headName := head.Name()
	if !headName.IsBranch() {
		err = errors.New("GitPath only functions when HEAD is pointed at a branch")
		return "", err
	}

	currentBranchName := strings.TrimPrefix(headName.Short(), "refs/heads/")

	return currentBranchName, nil

}
