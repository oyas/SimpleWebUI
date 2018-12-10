# Simple Web Interface

簡単なWebインターフェースを作成します。
インターフェースはマークダウンで記述できます。
アプリ側の通信は unix socket を使用します。


## 使用方法

1. `mark.md`を書き換えてWebの画面を作ります。
2. `go build` を実行し、`SimpleWebUI` を作成します。
3. `./SimpleWebUI` を実行して、サーバーを起動します。
4. `localhost:8080` へアクセスしてWebの画面を表示します。
5. アプリ側との通信は、`io.socket`を通して行います。


## mark.md の書き方

Markdownのパーサーには、[blackfriday](https://github.com/russross/blackfriday) を使用しています。
ただし、フォームを作成できるようにするために、Markdownパーサーに通す前に、独自に定義した`#`から始まる関数をパースします。

例えば、以下のように書きます。


	text input example

	#Input(name) label_text

	#submit submit_button_text


1行目は通常のテキスト、3行目はテキストボックス、5行目は submitボタン にパースされます。
その結果、以下のような出力が得られます。

```
<p>text input example</p>

<p><label> label_text <input type="text" class="form-control" name="name"></label></p>

<p><button type="submit" class="btn btn-default">submit_button_text</button></p>
```

`h1`などに変換する通常のmarkdownとの区別は、`#`の後にスペースがあるかないかです。
スペースがない場合のみ上記のような変換が行われます。

#### SharpFunc の定義

`#`から始まる関数の定義は、`parser/functions/` にあります。
ここに関数定義を書くことで、新しい関数を好きに追加できます。

例えば、inputタグを出力する`#input`関数は以下のように定義されています。


    func (sf SharpFunc) Input(args []string, body string) string {
        name := ""
        if len(args) > 0 {
            name = args[0]
        }
        return  "<label>" + body + " " +
                "<input type=\"text\" class=\"form-control\" name=\"" + name + "\">" +
                "</label>"
    }



#### Vue.js を使う

Vue.js を使って動的に操作できる部分を作ることができます。

	Vue text: {{ message }}

この`{{ }}`の部分は、アプリ側からjsonデータを送ることで値を変更できます。


## unix socket の通信の例

`server.go` を起動した時にそのディレクトリに`io.socket`が作成されます。
アプリ側とのデータの送受信はこの `io.socket` を使って行います。

`nc` コマンドを使って通信を試せます。

	$ nc -U io.socket

Vue用のデータを送るときは、json形式で改行を含まず1行で入力してください。

	{"message": "Hello, World!"}
