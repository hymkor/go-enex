[Top](./README.md) / [English](./release_note_en.md) / Japanese

v0.4.0
------
Jan 7, 2025

## unenex 修正

- `<en-todo>` タグを BALLOT BOX へ置換するようにした (U+2610 or U+2611)
- HTML の先頭に `<!DOCTYPE html>` を挿入するようにした。

## go-enex パッケージ修正

### 以下の変数・関数を廃止した

- `(*Resource) DataBeforeDecoded`
- `Parse`
- `ParseVerbose`
- `(*Export) ExHeader`

### 以下の変数・関数を非公開とした

- `ToSafe` (`defaultSanitizer`)
- `SerialNo` (`_SerialNo`)
- `ShrinkMarkdown` (`shrinkMarkdown`)
- `NewImgSrc` (`newAttachments`)

### 以下の関数・型の名前を変更した

- 関数: `ParseMulti` → 新`Parse`
- 型: `ImgSrc` → `Attachments`
- 型: `Export` → `Note`

### 以下の関数を `unenex` から `go-enex` へ移動

- `ToHtmls`
- `FilesToHtmls`
- `ToMarkdowns`
- `FilesToMarkdowns`

### その他

- `unenex`側のHTML化関数、markdown化関数を go-enex 側へ移動
- `(*Resource) Data` は戻り値でエラーも返すようにした
- `(*Export) ToHtml` の第一引数を interface からコールバック関数へ変更
- `Attachment.baseName` は BaseName として公開フィールドとした
- メソッド: `(*Export) HtmlAndDir` を `(*Note) Extract` へ改名し、引数の`Option`型で細かい指定を行えるようにした。

v0.3.6
======
Jan 5, 2025

- 添付ファイルへリンクが独立した行になるように、タグ: `<a>..</a>` は `<div class="goenex-attachment-link">..</div>` で囲むようにした。
- 添付画像のタグ: `<img>..</img>` は `<span class="goenex-attachment-image">..</span>` で囲むようにした。

Thanks to [@Juelicher-Trainee]

v0.3.5
======
Jan 4, 2025

- enexファイル名の末尾が `..enex` の時、`.` で終わるディレクトリ名を作成しようとするため、Windows ではパス名に不整合が発生する問題を修正。末尾の `.` は削除するようにした。

Thanks to [@Juelicher-Trainee]

v0.3.4
======
Jan 2, 2025

- 生成データは (ルート) → (enexファイル名ディレクトリ) → (ノート名の添付ファイル置き場用ディレクトリ) と3階層のディレクトリに格納するようにした。
    - enexファイル名が与えられず、標準入力から enexファイルデータを受けとった時は、(ルート) → (ノート名ディレクトリ) だけの2階層とする
    - ルートには各enexファイル名ディレクトリの index.html or README.md へのリンクをリストした index.html or README.md を置くようにした。
    - 各enexファイル名ディレクトリの index.html or README.md には enex ファイル名の大見出しを入れるようにした。
- 出力先ディレクトリを指定する `-d DIR` オプションを追加

Thanks to [@Juelicher-Trainee]

v0.3.3
======
Jan 1, 2025

- unenex: `index.html` の `<html>` タグより `lang="ja"` を削除
- exstyle: 与えられたHTML中にスタイルシートが見付からなかった時、エラーを表示
- ファイル名のない添付ファイルについて
    - 画像ファイルは `image.png` (拡張子はMIME型より選定) という名前で保存
    - 非画像ファイルは `Evernote` という名前で保存
    - 非画像ファイルのリンクテキストは `(Untitled Attachment)` ではなく、
      実際の保存ファイル名とした
- 同じファイル名につける通し番号は (1) から始めるようにした
- `</div><div>` を削除するなどのHTMLを短くするコードを削除
- unenex: スタイルシートをインラインで指定するオプション `-st` を追加  
    例: `unenex -st "div{line-height:2.0!important}" *.enex` (CMD.EXE)  
    or  `unenex -st 'div{line-height:2.0!important}' *.enex` (bash)
- unenex -markdown: README.md に記す URL で空白を表現する `+` のかわりに `%20` を使うようにした

Thanks to [@Juelicher-Trainee]

[#2]: https://github.com/hymkor/go-enex/issues/2

v0.3.2
======
Dec 31, 2024

## スタイルシート指定

- `-sf` でファイルのデータをスタイルシートとして読み込めるようにした
- Evernote デスクトップクライアントで直接エクスポートされたHTMLファイルからスタイルシートだけを抽出するツール `exstyle` を用意した(`-sf`用)

Evernote のスタイルシートは、Evernote社の著作物であるため、本ツールに添付することができない。そのため、もし本物により近い表現にしたい場合、ユーザ各位にてスタイルシートを次の手順で取り出していただく必要がある。

1. Evernote のデスクトップクライアントで、適当なページを HTML 形式でエクスポートする
2. エクスポートされた HTML からスタイルシートのみを抽出する  
    `exstyle some-directly-exported.html > common-enex.css`
3. 以後、unenex を使う時、そのスタイルシートを使うよう指定するようにする  
    `unenex -sf common-enex.css source.enex`

## その他の修正

- ノート名を `<h1>` タグで冒頭に挿入するようにした
- 画像のサイズ指定で、画像の原寸や `<en-media>` タグの`--naturalWidth`, `--naturalHeight` の値を使用するのをやめた
- ログメッセージで `Create xxxx` の後には必ず `: ` を入れるようにした

Thanks to [@Juelicher-Trainee]

v0.3.1
======
Dec 30, 2024

- 添付ファイルのリンクテキストが URL エンコードされたものになっていた不具合を修正 [#2]
- ファイル名がない添付ファイルは、リンクテキストを `(Untitled Attachment)` とするようにした
- ファイル名がない添付ファイルの生成に失敗する問題があったため、`Untitled`, `Untitled (2)` と代替ファイル名を与えるようにした
- 添付ファイルが画像かどうかの判断はファイル名ではなく、MIME の型を使うようにした
- 英大文字・英小文字が違うだけのファイル名は同一ファイル名が重複していると判断するようにした

Thanks to [@Juelicher-Trainee]

v0.3.0
======
Dec 29, 2024

- ([#2]) ファイル名を必要以上にリネームしないようにした。
    - ([#2]-2) 名前が重複しない限り、添付ファイルのファイル名の末尾に通し番号をつけない
    - ([#2]-3) ファイル名の空白をアンダスコアへ置換しない
- `unenex`: 複数の enex ファイル、および、ワイルドカード指定に対応
- ([#2]-5) 画像が常に最大サイズで展開されていた問題を修正

Thanks to [@Juelicher-Trainee]

v0.2.1
======
(2024.12.28)

- ([#2]) 非画像ファイルが `<img>` タグで出力HTMLに埋め込まれていたのを、`<a>` でリンクを張るよう修正した。
- unenex: `-h` オプションで、プログラムバージョン、OS、CPUアーキテクチャも表示するようにした

Thanks to [@Juelicher-Trainee]

v0.2.0
======
(2024.05.03)

- 新しい実行ファイル `unenex` を用意
    - 複数ノートを保持した enex ファイルをサポート
    ーHTMLやmarkdownファイルは (ノート名){.html OR .md} という名前になります
    - イメージファイルは (ノート名).files というフォルダーに配置されます
    - ファイル名に使えない文字は他の文字に置き換えられます(`<`→ `＜`, `>`→ `＞`, `"`→ `”`, `/`→ `／`, `\`→ `＼`, `|`→ `｜`, `?`→ `？`, `*`→ `＊`, `:`→ `：`, `(`→ `（`, `)`→ `）`, ` `→ `_`)
    - `embed` や `-prefix` は廃止となります。それらは `enexToHtml` で維持されます
    - 実行ファイル `enexToHtml` は非推奨となりますが、ソースファイルはパッケージの使用例として残します
- 添付ファイルのテキストエンコード表現の前後に空白がある時、デコードに失敗する問題を修正

v0.1.1
======
(2023.08.21)

- イメージのハッシュコードを、ソースURLだけでなく、`<recoIndex objID="...">` からも読み取るようにした
- `-v` オプションを追加

Thanks to [@Laetgark]

v0.1.0
======
(2023.03.23)

- 初回リリース

[@Juelicher-Trainee]: https://github.com/Juelicher-Trainee
[@Laetgark]: https://github.com/Laetgark
