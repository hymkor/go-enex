[#2]: https://github.com/hymkor/go-enex/issues/2

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
Dec 20, 2024

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
