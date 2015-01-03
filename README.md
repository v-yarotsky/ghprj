gh-prj
=======================

## Disclaimer

The thing is currently WIP,  I'm learning myself some golang, so it's going to be ugly.

## Summary

There's a command-line app `gh` and an [Alfred Extension](https://github.com/v-yarotsky/gh-prj/blob/master/assets/Github%20Prj.alfredworkflow?raw=true).
This alfred extension helps you open github repos in your web browser with a few keystrokes.

[Download latest](https://github.com/v-yarotsky/gh-prj/blob/master/assets/Github%20Prj.alfredworkflow?raw=true)

## Installation

Dependencies: go

```
brew install go
```

```
git clone https://github.com/v-yarotsky/gh-prj.git
cd gh-prj
make install
```

## Initial setup

Log in into github account using the `gh` commaind in terminal.

```
$ gh
# Username: <your username>
# Password: <your password>
```

## Usage

Just type `gh <repo name substring>` in alfred and hit Enter.
To reload stale repository list, type `ghreload` and hit Enter.

## TODO

- Better fuzzy matching
- More functional command-line-only tool

## License

[WTFPL](https://github.com/v-yarotsky/gh-prj/blob/master/LICENSE.txt?raw=true)

