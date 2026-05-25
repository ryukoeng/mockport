# Mockport Task Board

この `tasks/` ディレクトリは、Mockport を Phase 0 から Phase 6 まで TDD ベースで実装するための作業ボードです。

## 方針

- Stripe first で進める。
- Minimal MVP までは Go + Docker + built-in adapter に集中する。
- npm wrapper は後回しにする。
- Rust component は後回しにする。
- 動的 plugin system、複雑な DSL、adapter-specific image は Minimal MVP では作らない。
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
  phase4_additional_adapters.md
  phase5_compatibility_reports.md
  phase6_distribution.md
```

## Status の意味

```txt
pending      未着手
in_progress 進行中
blocked      外部要因または設計判断待ち
done         完了、検証済み
```

各タスクは `tasks/status.md` の表と、Phase ファイル内のチェックボックスの両方で管理します。作業開始時に status を `in_progress` に変更し、検証コマンドが通ったら `done` にします。

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

Phase 4 は「OpenAI-compatible、GitHub OAuth-like、Slack-like の built-in adapters を追加し、各 adapter の代表シナリオが動く」までです。

Phase 5 は「report が対応範囲、未対応 endpoint、request replay、adapter maturity を説明できる」までです。

Phase 6 は「GHCR、GitHub release binaries、Homebrew tap 下準備、npm wrapper 下準備、docs site の配布面を検証できる」までです。

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
