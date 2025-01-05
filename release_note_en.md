[Top](./README.md) / English / [Japanese](./release_note_ja.md)

v0.3.6
======
Jan 5, 2024

- To ensure proper separation between adjacent links, the `<a>...</a>` tag for attached files is now wrapped in a `<div class="goenex-attachment-link">...</div>` element.
- The `<img>` tag for attachment images is now wrapped in a `<span class="goenex-attachment-image">...</span>` element.

v0.3.5
======
Jan 4, 2024

- Fixed an issue where ENEX file names ending with `..enex` caused path inconsistencies on Windows due to attempting to create a directory name ending with `.`. The trailing `.` is now removed automatically.

v0.3.4
======
Jan 2, 2024

- Generated data is now stored in a three-level directory structure: (root) → (directory named after the ENEX file) → (directory for attachments of each note).
    - If no ENEX file name is provided and ENEX data is received from standard input, the structure will be limited to two levels: (root) → (directory named after the note).
    - A index.html or README.md file is placed in the root directory, listing links to the index.html or README.md files in each ENEX file's directory.
    - Each index.html or README.md in the ENEX file directories includes a heading with the ENEX file name.
- Added the `-d DIR` option to specify the directory for file output.

v0.3.3
======
Jan 1, 2024

- unenex: Removed `lang="ja"` from the `<html>` tag in index.html.
- exstyle: Display an error if no stylesheets are found in the given HTML.
- For attachments without filenames:
    - Image files are saved with names like `image.png` (the extension is determined by the MIME type).
    - Non-image files are saved with the name `Evernote`.
    - The link text for non-image files is now the actual saved filename instead of `(Untitled Attachment)`.
- Sequential numbering for duplicate filenames now starts from (1).
- Removed code that shortened HTML by eliminating tags like `</div><div>`.
- unenex: Add the `-st` option to specify the stylesheet text inline.  
    e.g., `unenex -st "div{line-height:2.0!important}" *.enex` (CMD.EXE)  
    or    `unenex -st 'div{line-height:2.0!important}' *.enex` (bash)
- unenex -markdown: `%20` is used instead of `+` for SPACE on URL of README.md

[#2]: https://github.com/hymkor/go-enex/issues/2

v0.3.2
======
Dec 31, 2024

## Specifying Style-Sheets

- Added support for loading file data as a stylesheet with the `-sf` option.
- Added a tool called `exstyle` to extract only the stylesheet from HTML files directly exported by the Evernote desktop client (for use with `-sf`).

Since Evernote's stylesheet is copyrighted material owned by Evernote Corporation, it cannot be included with this tool. Therefore, if you wish to achieve a representation closer to the original, you will need to extract the stylesheet themselves using the following steps:

1. Use the Evernote desktop client to export any page in HTML format.
2. Extract only the stylesheet from the exported HTML using:  
    `exstyle some-directly-exported.html > common-enex.css`
3. Specify the extracted stylesheet when using `unenex` from now on:  
    `unenex -sf common-enex.css source.enex`

## Miscellaneous

- Inserted the note title at the beginning using the `<h1>` tag.
- Stopped using the original size of images or the values of `--naturalWidth` and `--naturalHeight` from the `<en-media>` tag for specifying image dimensions.
- Fix: Insert `: ` always after `Create xxxx` on log message

v0.3.1
======
Dec 30, 2024

- Fix: the link text for the attachment was encoded as URL that should not [#2]
- Attachments without file names are now given the link text `(Untitled Attachment)`.
- Resolved an issue where generating attachments without file names would fail by assigning substitute file names like `Untitled`, `Untitled (2)`.
- Changed the method for determining whether an attachment is an image from relying on the file name to using the MIME type.
- File names that differ only in uppercase or lowercase letters are now considered to be duplicates.

v0.3.0
======
Dec 29, 2024

- ([#2]) Do not rename filenames as possible
    - ([#2]-2) Do not append serial-numbers to filenames unless filenames are duplicated
    - ([#2]-3) Do not replace SPACE to UNDERSCORE
- `unenex` can read multiple enex-files and support wildcards now
- ([#2]-5) Fix: images were always expanded to full size

v0.2.1
======
Dec 28, 2024

- [#2] Non-images attachments were embedded with `<img>` tags into output-html as if they were images. Modified to use `<a>` to link attachments
- unenex: `-h` option shows the version of program, OS and CPU-architecture now

v0.2.0
======
May 3,2024

- New executable `unenex`
    - It supports the enex files containing multi-notes
    - HTML and Markdown files are named as (NOTE-TITLE){.html OR .md}
    - Image files are put on the folder named as (NOTE-TITLE).files
    - The characters that can not used on filesystems are replaced to other characters (`<`→ `＜`, `>`→ `＞`, `"`→ `”`, `/`→ `／`, `\`→ `＼`, `|`→ `｜`, `?`→ `？`, `*`→ `＊`, `:`→ `：`, `(`→ `（`, `)`→ `）`, ` `→ `_`)
    - Abandon options `-embed` and `-prefix`. They are kept on `enexToHtml`
    - Executable `enexToHtml` is deprecated, but the source files are kept as an example of the package.
- Fix the problem decode failed when the encoding text has spaces as prefix or postfix

v0.1.1
======
Aug 21, 2023

- Read the hash codes of images from `<recoIndex objID="...">` not only source-url.
- Add `-v` option (verbose)

v0.1.0
======
Mar 23, 2023

- The first release.
