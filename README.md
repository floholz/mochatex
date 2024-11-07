# MochaTeX
[![GitHub tag (with filter)](https://img.shields.io/github/v/release/floholz/mochatex?label=latest)](https://github.com/floholz/mochatex/releases/latest)
[![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/floholz/mochatex/go.yml)](https://github.com/floholz/mochatex/actions/workflows/go.yml)
[![GitHub License](https://img.shields.io/github/license/floholz/mochatex)](./LICENSE)
[![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/floholz/mochatex?logo=go&labelColor=gray&label=%20)](https://go.dev/dl/)

> This project is based on [LaTTe](https://github.com/raphaelreyna/latte) by [raphaelreyna](https://github.com/raphaelreyna). 

Generate PDFs using LaTeX templates and JSON. This app provides a [CLI](#toc-cli) as well as a [GUI](#toc-gui).

## Table of Contents
* [About](#toc-about)
* [Install](#toc-install)
* [Usage](#toc-usage)
	* [CLI](#toc-cli)
	* [GUI](#toc-gui)
* [Roadmap](#toc-roadmap)
* [Contributing](#toc-contributing)
* [License](#toc-license)

<a name="toc-about"></a>
## About
MochaTeX helps you generate professional looking PDFs by using your TEX files as templates and filling in the details with JSON.
Under the hood, it uses [pdfTeX / pdfLaTeX](https://tug.org/applications/pdftex) to create a PDF from a TEX file that has been filled in using [Go's templating package](https://golang.org/pkg/text/template/).

<a name="toc-install"></a>
## Install

```shell
go install github.com/floholz/mochatex@latest
```

<a name="toc-usage"></a>
## Usage

You can run MochaTeX as a [cli tool](#toc-cli), or with its [GUI](#toc-gui).

<a name="toc-cli"></a>
### CLI


```shell
mochatex
```

<a name="toc-gui"></a>
### GUI

> coming soon


<a name="toc-roadmap"></a>
### Roadmap
- CLI tool
- GUI

---

<a name="toc-contributing"></a>
### 🤝 Contributing

Contributions, issues and feature requests are welcome!

Feel free to check [issues page](https://github.com/floholz/mochatex/issues).


<a name="toc-license"></a>
### 📝 License

By [floholz](https://github.com/floholz), 2024.

This project is [MIT](./LICENSE) licensed.

---

### Show your support

Give a ⭐ if this project helped you!
