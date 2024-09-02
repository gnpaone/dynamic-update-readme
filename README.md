# Dynamic Update Readme

[![Releases](https://img.shields.io/github/v/release/gnpaone/dynamic-update-readme?style=flat-square)](https://github.com/gnpaone/dynamic-update-readme/releases)
[![Build](https://img.shields.io/github/actions/workflow/status/gnpaone/dynamic-update-readme/build.yml?style=flat-square&logo=github)](https://github.com/gnpaone/dynamic-update-readme/actions/workflows/build.yml)
[![LICENSE](https://img.shields.io/github/license/gnpaone/dynamic-update-readme?color=green&style=flat-square)](https://github.com/gnpaone/dynamic-update-readme/blob/main/LICENSE)
[![Issues](https://img.shields.io/github/issues/gnpaone/dynamic-update-readme?color=orange&style=flat-square)](https://github.com/gnpaone/dynamic-update-readme/issues)
[![Go](https://img.shields.io/github/go-mod/go-version/gnpaone/dynamic-update-readme?color=maroon&style=flat-square)](https://github.com/gnpaone/dynamic-update-readme/blob/main/go.mod)
[![Godoc](https://pkg.go.dev/badge/github.com/gnpaone/dynamic-update-readme.svg?utm_source=godoc)](https://godoc.org/github.com/gnpaone/dynamic-update-readme)
[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-v1.0.2-undefined.svg?logo=github&logoColor=white&style=flat-square)](https://github.com/marketplace/actions/dynamic-update-readme)

#### As a power user, automate the updating of markdown text with content obtained from other actions in the workflow, as well as a Go module that can update markdown text

---

## Introduction

A Go module that updates markdown text and is integrated as a GitHub Action.

## Parameters

| Parameter        | Workflow default    | Description                                | Required    | Possible values    |
|:----------------:|:-------------------:|:-------------------------------------------|:-----------:|:-------------------|
| readme_path      | ./README.md         | Path of the markdown file                  | No          | *File path relative to the root directory of the GitHub repo*    |
| marker_text      | *null*              | Marker text to replace in markdown file    | Yes         | *Example markers to be added in the markdown:*<br>`<!-- `*marker_text*`_START --><!-- `*marker_text*`_END -->`    |
| markdown_text    | *null*              | Markdown text that needs to be updated     | Yes         | *Any markdown compatible text*    |
| table            | false               | Markdown text is a table                   | No          | true, false    |
| table_options    | *null*              | Alignment for the table                    | No          | align-*Alignment*, col-*Column*-align-*Alignment*, col-*Column*-w-*Min-width*    |

<sup>_Alignment_: `left, right, center`<br>_Column_: Column number of the table (starts with `0`)<br>_Min-width_: Minimum width of a column (in other words, minimum number of hyphens (`-`) in the delimiter row)</sup>

### Notes:
* Make sure to change the following in your GitHub repo settings: `Actions` > `General` > `Workflow permissions` > Choose `Read and write permissions` > Check `Allow GitHub Actions to create and approve pull requests` > `Save`.
* Other parameters `commit_user`, `commit_email`, `commit_message` and `confirm_and_push` related to action workflow are optional.
* If `confirm_and_push` is "false" committer detalis can be accessed via `outputs.commit_user`, `outputs.commit_email` & `outputs.commit_message` for further usage in the workflow.
* The syntax mostly revolves around [GitHub Flavored Markdown parser](https://github.github.com/gfm/).
* The `table` parameter is optional as `markdown_text` supports table markdown syntax by default but this parameter can be used as a simple alternative or any special use case. If `table` is "true":
  - It uses [github.com/willabides/mdtable](https://godoc.org/github.com/willabides/mdtable) module 
  - `markdown_text` contents should follow these conditions:
    + Table row contents are seperated with ";" delimiter. First element will make up the table header.
    + For each element of table rows, table column contents are seperated with "," delimiter.
  - `table_options` can be used only if `table` is "true".
    + Description of options/values:
      1. `align-`: This option aligns the whole table according to the alignment set.
      2. `col-`: This option sets properties of any particular column based on the column number.
    + Order of preference of the alignment options: `col-` > `align-`
* The markdown file must contain start marker `<!-- EXAMPLE_MARKER_START -->` and end marker `<!-- EXAMPLE_MARKER_END -->` where "EXAMPLE_MARKER" is the input of `marker_text` parameter. Note that the `_START` and `_END` part is important.

## Usage
Check the example workflow [here](https://github.com/gnpaone/dynamic-update-readme/blob/main/examples/update.yml).

## Contributing
See [CONTRIBUTING.md](https://github.com/gnpaone/dynamic-update-readme/blob/main/CONTRIBUTING.md) for more details.

## Acknowledgements

This project uses few of the open-source Go modules, without which it would not have been possible. I extend my gratitude to the developers and maintainers of these modules for their valuable contributions to the open-source community:

* [github.com/willabides/mdtable](https://godoc.org/github.com/willabides/mdtable) - Generates markdown tables from string slices.
* [github.com/mattn/go-runewidth](https://godoc.org/github.com/mattn/go-runewidth) - Provides functions to get fixed width of the character or string.
* [github.com/rivo/uniseg](https://godoc.org/github.com/rivo/uniseg) - Implements Unicode Text Segmentation, Word Wrapping, and String Width Calculation.

## License
This project is licensed under GPL-3.0.  
Copyright © 2024, Naveen Prashanth. All Rights Reserved.
