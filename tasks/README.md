# Mockport Task Board

この `tasks/` ディレクトリは、Mockport を Phase 0 から公開 OSS 運用まで TDD ベースで実装するための作業ボードです。

## 方針

- Stripe first で進める。
- Minimal MVP までは Go + Docker + built-in adapter に集中する。
- npm wrapper は後回しにする。
- Rust component は後回しにする。
- 動的 plugin system、複雑な DSL、adapter-specific image は Minimal MVP では作らない。
- Phase 11 以降は、provider 内部再現ではなく、公開 API / 公式 SDK / 主要 workflow / fake state / 主要 error shape の高忠実度互換を目標にする。
- 「Full-compatible」は全 provider 内部や全 undocumented behavior の再現ではなく、測定可能な provider-compatible local API として定義する。
- production code は、対応する failing test を先に作って失敗を確認してから実装する。

## ファイル構成

```txt
tasks/
  README.md
  status.md
  phase0_baseline.md
  phase1_stripe_minimal_mvp.md
  phase2_cli_ux.md
  phase3_ai_safe_mode.md
  phase4_trust_reports.md
  phase5_additional_adapters.md
  phase6_distribution.md
  phase7_public_oss_hardening.md
  phase8_public_env_safety.md
  phase9_public_docs_and_discovery.md
  phase10_public_preview_release.md
  phase11_community_and_maintenance.md
  phase12_compatibility_engine.md
  phase13_fixture_spec_policy.md
  phase14_sdk_contract_harness_foundation.md
  phase15_stateful_provider_workflows.md
  phase16_stripe_provider_compatibility.md
  phase17_openai_provider_compatibility.md
  phase18_github_oauth_provider_compatibility.md
  phase19_slack_provider_compatibility.md
  phase20_provider_compatible_release_track.md
```

## Status の意味

```txt
pending      未着手
in_progress 進行中
blocked      外部要因または設計判断待ち
done         完了、検証済み
```

進捗の正本は `tasks/status.md` です。Phase ファイルは実行計画として使い、作業開始時に `tasks/status.md` の status を `in_progress` に変更し、検証コマンドが通ったら `done` にします。Phase ファイル内の `Status` とチェックボックスは、実行時に可能な範囲で同期します。

## TDD 作成方法

各実装タスクは次の順番で進めます。

1. RED: 期待する振る舞いを示すテストを書く。
2. RED verify: 対象テストだけを実行し、期待どおり失敗することを確認する。
3. GREEN: テストを通す最小実装を書く。
4. GREEN verify: 対象テストを再実行し、成功を確認する。
5. REFACTOR: 重複や命名を整理する。振る舞いは増やさない。
6. FULL verify: `go test ./...` など、その Phase の検証コマンドを通す。
7. Status update: `tasks/status.md` と Phase ファイルのチェックを更新する。
8. Commit checkpoint: Git リポジトリ化後は、Phase または小さな機能単位で commit する。

## Phase の境界

Phase 0 は「Mockport の実装基盤が起動して health check に答える」までです。

Phase 1 は「Stripe-like adapter を Docker で起動し、成功/失敗/認証エラー/rate limit/timeout/webhook/report/AI-safe warning を Minimal MVP の範囲で扱える」までです。

Phase 2 は「空ディレクトリから `mockport init` と `mockport up/run` で 2 分以内に起動できる CLI UX」までです。

Phase 3 は「AI-safe mode が警告、strict fail、redaction、report、docs で明確に差別化される」までです。

Phase 4 は「adapter を増やす前に、対応範囲、未対応 endpoint、scenario coverage、maturity を report で説明できる trust foundation」を作るまでです。

Phase 5 は「Phase 4 の metadata/report contract に沿って、OpenAI-compatible、GitHub OAuth-like、Slack-like の built-in adapters を追加する」までです。

Phase 6 は「GHCR、GitHub release binaries、Homebrew tap 下準備、npm wrapper 下準備、docs site の配布面を検証できる」までです。

Phase 7 は「公開前に必要な LICENSE、SECURITY、CONTRIBUTING、issue/PR templates、public CI gate を揃える」までです。

Phase 8 は「Mockport 用 `.env` を公開してもよい条件を scanner、docs、examples、CI で保証する」までです。

Phase 9 は「公開 docs、support matrix、limitations、examples を整え、初見ユーザーが導入判断できる」までです。

Phase 10 は「初回 public preview release として GitHub Release、checksums、GHCR image を実際に公開し、clean install を検証する」までです。

Phase 11 は「Dependabot、maintainer guide、roadmap、adapter contribution quality bar を整え、継続運用できる」までです。

Phase 12 は「互換性を manifest、score、report として測れる Compatibility Engine を作る」までです。

Phase 13 は「fixture/spec snapshot の安全性、出典、更新ルールを定義する」までです。

Phase 14 は「公式 SDK / 実 client contract を local Mockport に向けて通す harness foundation を作る」までです。

Phase 15 は「主要 workflow を fake state 上で create/retrieve/list/update できる stateful API にする」までです。

Phase 16 は「Stripe adapter を SDK/workflow-compatible local API に引き上げる」までです。

Phase 17 は「OpenAI adapter を SDK/workflow-compatible local API に引き上げる」までです。

Phase 18 は「GitHub OAuth adapter を client/workflow-compatible local API に引き上げる」までです。

Phase 19 は「Slack adapter を client/workflow-compatible local API に引き上げる」までです。

Phase 20 は「provider-compatible release track と compatibility report を継続運用する」までです。

## 参照資料

- `mockport_docs/docs/02_mvp_scope.md`
- `mockport_docs/docs/03_architecture.md`
- `mockport_docs/docs/04_go_engineering_guide.md`
- `mockport_docs/docs/06_config_spec.md`
- `mockport_docs/docs/07_adapter_design.md`
- `mockport_docs/docs/08_cli_spec.md`
- `mockport_docs/docs/09_docker_runtime.md`
- `mockport_docs/docs/10_ai_safe_security.md`
- `mockport_docs/docs/11_testing_strategy.md`
- `mockport_docs/docs/12_roadmap.md`
