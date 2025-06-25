# Nextver CLI

Evaluate next version number with SemVer format based on commit messages.

This CLI requires the commit messages to be written in [conventional commit](https://www.conventionalcommits.org/) style.

## Installation

### Curl

```
curl -sL https://raw.githubusercontent.com/chonla/nextver/refs/heads/main/install.sh | sh
```

### Homebrew/Linuxbrew

```
brew tap chonla/universe
brew install nextver
```

### From source

```
go get github.com/chonla/nextver
```

### Upgrade

```
brew upgrade
```

## Usage

```
Usage of nextver:

  nextver [options...] [dir]

Options:
  -d  Debug mode, print considering steps.
  -e  Suppress trailing new line. Print only version out.
  -n  Version is not prefixed by v, for example, 1.0.0.
  -t  Show detected latest version.
  -v  Show version of nextver.
```

## Example

```
$ nextver
v1.1.0

$ nextver -e
v1.1.0

$ nextver -d
Detected current version: v1.0.0
Version is prefixed by v.
HEAD commit ID=e6f55889e779abc2b4bb616ffd0cf8b3aca6c05a
Latest tag commit ID=baf76a97095b0eeae8176bb12c4b89eef0e50380
1 commit(s) since latest tag
============
Commit stats
------------
Major change(s) = 0
Minor change(s) = 1
Revision change(s) = 0
============
Estimated next version: v1.1.0

$ nextver -t
v1.0.0
```

## References

* [Conventional Commit](https://www.conventionalcommits.org/)

## License

[MIT](LICENSE.txt)
