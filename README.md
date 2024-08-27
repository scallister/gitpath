# gitpath

Tool that returns GitHub and GitLab urls for local filepaths.

## Usage
```bash
$ gitpath README.md
https://github.com/scallister/gitpath/blob/scallister/initial/README.md

$ gitpath cmd/root.go
https://github.com/scallister/gitpath/blob/scallister/initial/cmd/root.go
```
## Brew Install
```bash
brew tap scallister/scallister
brew install gitpath
```

## Go Install
```bash
go install github.com/scallister/gitpath@latest
```

## Features
- Generates a url to the branch that is currently being used
- Branch can be overridden with `--main`, `--master`, and `--branch <name>`
