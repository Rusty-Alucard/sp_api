# SP API アーキテクチャモダナイゼーション議論記録

## 議論の背景
- 現在のプロジェクトは基本的なClean Architectureを採用
- Go 1.16という古いバージョンを使用
- レイヤー中心の構造でビジネス機能の分離が不十分
- DDDの概念（Aggregate, Value Object, Domain Event）が未実装

## 調査したリソース

### Context7で調査したライブラリ
1. **Go Clean Architecture** (/wesionaryteam/go_clean_architecture)
   - Gin Web Framework + Uber-go/fx + GORM
   - Trust Score: 8.7, Code Snippets: 31
   - 機能中心のモジュール構造を採用

2. **Rich Domain** (/4lessandrodev/rich-domain)
   - DDDパターンの実装例
   - Trust Score: 9.6, Code Snippets: 85
   - Aggregate, Value Object, Domain Eventの実装パターン

### Geminiとの議論
- Go 1.22+へのアップグレードの重要性
- Genericsの活用による型安全性の向上
- golangci-lintの導入によるコード品質向上
- 機能中心のディレクトリ構造への変更
- PostgreSQLの継続使用の妥当性

## 主要な決定事項

### 1. 技術スタックの更新
- **Go 1.22+**: Generics、パフォーマンス改善
- **golangci-lint**: コード品質管理
- **Wire**: 依存性注入（現状維持）
- **Atlas**: データベースマイグレーション

### 2. アーキテクチャパターン
- **DDD (Domain-Driven Design)**: 厳密な実装
- **Aggregate Pattern**: ビジネスロジックの集約
- **Value Objects**: 型安全性とビジネスルールの表現
- **Domain Events**: 副作用の分離

### 3. ディレクトリ構造
- **現在**: レイヤー中心（domain/, usecase/, infrastructure/）
- **提案**: 機能中心（event/, team/, player/）
- **理由**: 高凝集、低結合、スケーラビリティ

### 4. データストア
- **PostgreSQL継続**: ACID特性、複雑なクエリ、JSONサポート
- **Redis追加**: キャッシュ層
- **将来検討**: Elasticsearch（検索）、Event Sourcing

## 実装計画

### Phase 1: 基盤整備
1. Go 1.22+ アップグレード
2. golangci-lint 導入
3. 新ディレクトリ構造作成
4. 基本DDDパターン実装

### Phase 2: ドメイン実装
1. Event Aggregate 実装
2. Team Aggregate 実装
3. Player Aggregate 実装
4. Value Objects 実装

### Phase 3: インフラ整備
1. Repository 実装改善
2. Domain Events 実装
3. Event Bus 実装
4. テスト追加

### Phase 4: API改善
1. GraphQL スキーマ更新
2. Resolver 改善
3. エラーハンドリング統一
4. ドキュメント整備

## 期待される効果
1. **保守性向上**: 機能ごとの分離
2. **テスタビリティ**: ドメインロジック分離
3. **拡張性**: 新機能追加の容易さ
4. **型安全性**: Value Objects活用
5. **ビジネスロジック明確化**: ドメインエキスパート協業

## リスクと対策
1. **学習コスト**: 段階的導入とチーム教育
2. **複雑性増加**: シンプルな実装から開始
3. **パフォーマンス**: プロファイリングと最適化

## 次のステップ
1. 提案書のレビューと承認
2. Phase 1の実装開始
3. チームメンバーへのDDD教育
4. 段階的な移行計画の策定

