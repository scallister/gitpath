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

import "testing"

func Test_trimUrlToBase(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "SSH GitLab Url",
			args: args{url: "git@gitlab.com:scallister/gitpath.git"},
			want: "https://gitlab.com/scallister/gitpath",
		},
		{
			name: "HTTPS GitHub Url",
			args: args{url: "https://github.com/scallister/gitpath.git"},
			want: "https://github.com/scallister/gitpath",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimUrlToBase(tt.args.url); got != tt.want {
				t.Errorf("trimUrlToBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRelativePath(t *testing.T) {
	type args struct {
		fullPath string
		repoPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "File in root",
			args: args{
				fullPath: "/mnt/c/Users/steven/IdeaProjects/gitpath/README.md",
				repoPath: "/mnt/c/Users/steven/IdeaProjects/gitpath",
			},
			want: "README.md",
		},
		{
			name: "File in folder",
			args: args{
				fullPath: "/mnt/c/Users/steven/IdeaProjects/gitpath/cmd/gitpath.go",
				repoPath: "/mnt/c/Users/steven/IdeaProjects/gitpath",
			},
			want: "cmd/gitpath.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRelativePath(tt.args.fullPath, tt.args.repoPath); got != tt.want {
				t.Errorf("getRelativePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createURL(t *testing.T) {
	type args struct {
		baseURL      string
		branchName   string
		relativePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "github",
			args: args{
				baseURL:      "https://github.com/scallister/gitpath",
				branchName:   "scallister/initial",
				relativePath: "cmd/gitpath.go",
			},
			want: "https://github.com/scallister/gitpath/blob/scallister/initial/cmd/gitpath.go",
		},
		{
			name: "gitlab",
			args: args{
				baseURL:      "https://gitlab.com/scallister/dotfiles",
				branchName:   "master",
				relativePath: "todo.md",
			},
			want: "https://gitlab.com/scallister/dotfiles/-/blob/master/todo.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createURL(tt.args.baseURL, tt.args.branchName, tt.args.relativePath); got != tt.want {
				t.Errorf("createURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
