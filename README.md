# pub

[![](https://pkg.go.dev/badge/github.com/JessebotX/pub)](https://pkg.go.dev/github.com/JessebotX/pub)

- Easy to use **static site generator** that also supports creating RSS feeds, EPUBs, etc. from your Markdown files (WIP)
- Highly **customizable and extensible**: modify the underlying layout to your heart's content using `go`'s advanced text templating engine
- **Zero/minimal configuration necessary**: `pub` comes with great defaults that emphasize user accessibility and system compatibility
- **Local and private**: use the program and browse the documentation (WIP) locally without an internet connection
- **No lock-in** to any service or ecosystem
- **Backward and forward-compatible** (for version >= 1.0.0)

## Description

`pub` is an easy-to-use _static site generator_ specifically **designed for independently publishing your Markdown-based written works on the web**. `pub` also supports _generating other popular open formats such as PDFs and EPUBs (WIP)_, which are supported across almost all machines and e-readers.

`pub` requires _minimal configuration_, allowing you to put all your focus on the act of writing itself. These great defaults emphasize user accessibility and widespread compatibility with many systems, letting you produce works that are comfortable to read in with little effort. Nevertheless, `pub` is also _extremely configurable and extensible_, with an embedded text templating engine (`go`'s `text/template` and `html/template`) and support for custom configuration fields that allow you to develop and use custom themes to modify the overall presentation of the output.

Additionally, `pub` is permissively licensed. You can use the program and browse the documentation locally without internet, and with all your writing being done in open plain-text formats, there is little potential lock-in to any single service/ecosystem.

The upcoming `pub` version 1.0.0 is also committed to not only be _backward-compatible_ but also _forward-compatible_. Old book projects should be able to be successfully built in the future.

## Install

### Go

Install the latest version with `go` itself.

```shell
go install github.com/JessebotX/pub/cmd/pub@latest
```

## Usage

See [docs/getting-started.md](docs/getting-started.md) for more information.

For help using the CLI utility, run program with `-h` or `--help` flag like so:

```shell
pub --help
```

OR access the CLI man page on a unix-like system (WIP):

```shell
man pub
```

