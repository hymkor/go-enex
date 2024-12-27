- Non-images attachments were embedded with `<img>` tags into output-html as if they were images. Modified to use `<a>` to link attachments
- The `-h` option shows the version of program, OS and CPU-architecture now

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
