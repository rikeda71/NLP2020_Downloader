# NLP2020_Downloader

NLP2020 の論文とポスターをダウンロードするためのリポジトリ

## Requirements

- golang
- NLP2020 のオンライン開催版プログラムへのアクセスに必要なユーザ名とパスワード
  - NLP2020 の参加者にのみ通達（公開禁止）
  - `事前参加登録手続き完了のお知らせ`というタイトルのメールに記載

## Usage

1. `cp .env.sample .env`

2. `.env`を編集

   ```text
   USERNAME=ユーザ名
   PASSWORD=パスワード
   ```

3. 以下を実行

   ```shell script
   go build
   ./nlp2020-downloader
   ```
