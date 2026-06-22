# Limitations 日本語版

[English](limitations.md)

Mockport は selected workflow の local integration test を目的にしており、provider の内部実装や未公開 behavior は再現しません。

## 対象外

- 実 payment processing、fraud、settlement、billing network。
- 実 AI inference、tokenization parity、private scheduling。
- GitHub/Slack/LINE の enterprise policy や full directory state。
- provider sandbox や production validation の完全な代替。

## 未実装の設定ブロック

`mockport.yml` の `scenarios:` ブロックはパースされますが**未実装**です。ランタイムでは無視されます。
このブロックが存在する場合、起動時（`--check` 出力および `/_mockport/report`）に警告が出力されます。

レスポンスの切り替えやエラーケースのシミュレーションには以下を使用してください：

- アダプタの `scenario:` フィールドによる built-in シナリオ
- `X-Mockport-Scenario` リクエストヘッダ（issue #80 参照）

ユーザー定義シナリオの将来の計画については [scenario-policy.md](../scenario-policy.md) を参照してください。

採用前に support matrix、report、adapter examples を確認してください。
