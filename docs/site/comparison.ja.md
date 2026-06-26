# Comparison 日本語版

[English](comparison.md)

Mockport は selected SaaS workflow 向けの local / Docker-first emulator です。すべての mock server、provider sandbox、contract test tool の置き換えではありません。開発者が実際に検討する代替手段との比較です。

## まず結論（How to choose）

- **Stripe だけ・状態不要・OpenAPI 由来の網羅性が欲しい** → [stripe-mock](https://github.com/stripe/stripe-mock)（Stripe 公式モック）
- **AWS をローカルで動かしたい** → [LocalStack](https://localstack.cloud/)（Mockport とは守備範囲が重ならない — AWS API 向けで Stripe 等の SaaS は対象外）
- **OpenAPI スペックがあり手早く mock したい** → [Stoplight Prism](https://github.com/stoplightio/prism)
- **JavaScript テスト内で完結させたい（browser / Node.js）** → [MSW](https://github.com/mswjs/msw)
- **決済 + LLM + チャット + OAuth など複数 SaaS を 1 プロセスで、状態・シナリオ・webhook 配送・シークレット安全性込みで** → Mockport
- **HTTP スタブをすべて自前で持つ** → WireMock または手書きダブル（下記）
- **本番前の provider 挙動の正本が欲しい** → provider sandbox（下記）

## 機能比較表

| | Mockport | stripe-mock | Prism | MSW | WireMock |
| --- | --- | --- | --- | --- | --- |
| 対象 | Stripe / OpenAI / Slack / GitHub OAuth / LINE / Zoho OAuth | Stripe のみ | OpenAPI（または Postman）がある任意 API | 任意 HTTP（JavaScript ランタイムのみ） | 任意 HTTP（スタブ手書き） |
| 状態保持（作成→取得の往復） | ✅ | ❌ | ❌ | ハンドラ次第 | スタブ次第 |
| エラーシナリオ切替 | ✅（`scenario:`、`X-Mockport-Scenario`） | ❌ | example / response 選択のみ | 手書き | 手書き |
| webhook / イベント配送 | ✅（Stripe / LINE / Slack ヘルパー） | ❌ | ❌ | — | ✅（設定次第） |
| 言語非依存（HTTP） | ✅ | ✅ | ✅ | ❌ JS のみ | ✅ |
| OpenAPI からの自動生成 | ❌ | ✅（Stripe OpenAPI 由来） | ✅ | — | ❌ |
| シークレット安全ポリシー | ✅（`ai-safe` mode） | — | — | — | — |

`—` は今回の一次情報確認で該当ツールの項目を確定できなかったセルです。

### 根拠（一次情報）

| ツール | 参照 |
| --- | --- |
| stripe-mock | [stripe/stripe-mock README](https://github.com/stripe/stripe-mock/blob/master/README.md) — ステートレス、OpenAPI 由来応答、エラーシナリオ設定不可、webhook 送信なし |
| LocalStack | [LocalStack AWS services docs](https://docs.localstack.cloud/aws/services/) — AWS サービス emulation。Stripe 等 SaaS API は対象外 |
| Prism | [stoplightio/prism README](https://github.com/stoplightio/prism/blob/master/README.md) — OpenAPI/Postman mock server。spec example ベースの動的応答。永続 state なし |
| MSW | [mswjs/msw README](https://github.com/mswjs/msw/blob/main/README.md) — browser / Node.js でのリクエストインターセプト。JavaScript/TypeScript |
| WireMock | [WireMock stubbing docs](https://wiremock.org/docs/stubbing/)、[webhooks and callbacks](https://wiremock.org/docs/webhooks-and-callbacks/) — 汎用スタブ。webhook は設定次第 |

## Mockport を選ばない方がよい場合

- **Stripe の全エンドポイント網羅が必要** — Mockport は workflow 単位のサポート。[support matrix](support-matrix.ja.md) と adapter spec を確認。OpenAPI 網羅は stripe-mock 向き。
- **課金計算・fraud・本物の認可ロジックの再現が必要** — provider 内部は再現しない（[ROADMAP](../../ROADMAP.ja.md) Non-Goals）。
- **OpenAPI から adapter なしで mock を生成したい** — Prism 等を使う。Mockport は汎用 spec コンパイラではない。
- **JavaScript のみで in-process に完結させたい** — サイドカー HTTP より MSW が単純。
- **provider の正本挙動が必要** — 最終確認は実 sandbox。Mockport は local / CI integration 向け。
- **AWS サービスの local emulation が必要** — LocalStack 向け。Mockport は AWS API を emulated しない。

## Mockport vs WireMock

WireMock は汎用 HTTP mock です。スタブ・マッチャー・webhook callback を自前定義します。Mockport は provider-shaped — adapter が selected workflow、fake credential、安全警告、report、support matrix を共有実装します。

## Mockport vs 手書き Test Double

手書きダブルは 1 コードベースでは速いですが drift しやすい。Mockport は adapter 挙動・example・report・compatibility evidence を集中管理し、複数アプリと CI が同じ local provider API を共有できます。

## Mockport vs Provider Sandbox

Provider sandbox は provider 挙動の正本です。Mockport は local / Docker-first / secret-free / deterministic。高速な local と CI integration test に使い、本番前に critical path を実 sandbox で確認してください。

## 位置づけまとめ

Mockport が向く条件:

- 複数 SaaS を 1 プロセスの Docker-first local API で動かしたい
- commit してよい fake env が欲しい
- CI で deterministic な外部サービスシナリオと unsupported behavior の明示が欲しい
- 採用前に adapter coverage と known gap を report で確認したい

Mockport は provider 内部の完全クローンではありません。
