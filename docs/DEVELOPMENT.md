# Development Guide

## セットアップ

### 初回セットアップ
```bash
# 開発環境のセットアップ（自動）
./scripts/dev-setup.sh
```

### 手動セットアップ
```bash
# CircleCI CLIのインストール
brew install circleci  # macOS
# または
curl -fLSs https://circle.ci/cli | bash

# pre-commitのインストール
pip install pre-commit
pre-commit install
```

## 開発ワークフロー

### コミット前の検証
```bash
# すべての検証を実行
make validate-all

# 個別の検証
make validate-ci  # CircleCI設定の検証
make lint        # Goコードのリント
make test        # テストの実行
```

### よく使用するコマンド
```bash
# 開発サーバーの起動
make server

# テストの実行
make test

# カバレッジレポートの生成
make test-coverage

# 依存関係の更新
make install-go-modules
```

## 品質保証

### 自動チェック
- **Pre-commit hooks**: コミット時に自動でYAML構文チェック、リント、フォーマット
- **CircleCI**: プッシュ時に自動でテスト、リント、セキュリティスキャン

### 手動チェック
```bash
# CircleCI設定の検証
circleci config validate

# ローカルでのCircleCI実行（オプション）
circleci local execute
```

## トラブルシューティング

### YAML構文エラー
```bash
# 設定ファイルの検証
make validate-ci

# 詳細なエラー情報
circleci config validate --verbose
```

### リントエラー
```bash
# 自動修正可能なエラーの修正
make lint-fix

# 詳細なリント情報
golangci-lint run --verbose
```

### Pre-commitフックの無効化（緊急時）
```bash
# 一時的に無効化
git commit --no-verify -m "your message"

# 完全に無効化
pre-commit uninstall
```

## ベストプラクティス

1. **コミット前**: 必ず `make validate-all` を実行
2. **YAML編集**: VS CodeのYAML拡張機能を使用
3. **設定変更**: CircleCI設定変更後は必ず `make validate-ci` で検証
4. **プッシュ前**: ローカルでテストが通ることを確認

## 推奨ツール

### VS Code拡張機能
- Red Hat YAML Language Support
- Go for Visual Studio Code
- Prettier - Code formatter

### コマンドラインツール
- CircleCI CLI
- golangci-lint
- pre-commit
