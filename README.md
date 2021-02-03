# UrlSelector

## URLの設定方法

### JSONファイルの準備

以下のように、`name` と `url` を書いたJSONファイルを作成する。

```UrlSelector.json
[
  {
    "name": "Google",
    "url": "https://google.com"
  },
  {
    "name": "Yahoo!",
    "url": "https://yahoo.co.jp"
  }
]
```

### JSONファイルの指定方法

以下の優先順位で設定用JSONファイルが読み込まれます。読み込まれるJSONファイルは一つのみです。

1. 引数指定
1. 実行コマンドの実態と同じディレクトリの `UrlSelector.json`
1. 環境変数 `URLSELECTOR_CONFIG` で指定されたファイルパス
1. ホームディレクトリの `UrlSelector.json`