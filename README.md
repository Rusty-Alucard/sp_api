# SP API - Soccer/Football Data Management System

## 概要

SP APIは、サッカー/フットボールの大会、チーム、選手、試合などのデータを管理するためのGraphQL APIサーバーです。Go言語で開発されており、Clean Architectureの設計原則に従って構築されています。

## アーキテクチャ

### 技術スタック
- **言語**: Go 1.16
- **GraphQL**: gqlgen (99designs/gqlgen)
- **データベース**: PostgreSQL
- **依存性注入**: Google Wire
- **設定管理**: Viper
- **コンテナ**: Docker & Docker Compose
- **その他**: Redis, HashiCorp Vault

### アーキテクチャパターン
Clean Architectureを採用し、以下の層に分離されています：

```
├── domain/          # ドメイン層
│   ├── model/       # エンティティ・値オブジェクト
│   └── repository/  # リポジトリインターフェース
├── usecase/         # ユースケース層
├── infrastructure/  # インフラストラクチャ層
│   ├── persistence/ # データ永続化
│   └── provider/    # 外部サービス連携
└── graph/          # GraphQL層
    ├── generated/   # 自動生成コード
    └── resolver/    # リゾルバー実装
```

## 主要機能

### GraphQLスキーマ
以下の主要なエンティティを管理：

- **Event**: 大会（FIFAワールドカップなど）
- **Team**: チーム（国代表など）
- **Player**: 選手
- **Coach**: コーチ
- **Match**: 試合
- **Stadium**: スタジアム
- **Timeline**: 試合のタイムライン（ゴール、カードなど）

### 主要なクエリ
```graphql
type Query {
    events: [Event!]!           # 全大会取得
    event(id: ID!): Event!      # 特定大会取得
    teams: [Team!]!             # 全チーム取得
    team(fifa: String!): Team!  # FIFAコードでチーム取得
}
```

## データベース設計

### 主要テーブル
- `events`: 大会情報
- `teams`: チーム情報（FIFAコード、連盟など）
- `players`: 選手情報
- `coaches`: コーチ情報
- `matches`: 試合情報
- `stadiums`: スタジアム情報
- `participants`: 大会参加チーム
- `registered_players`: 登録選手
- `match_progresses`: 試合進行状況

### インデックス
- チームのFIFAコード
- 連盟別
- 大会別
- 試合別

## 開発環境

### 前提条件
- Go 1.16以上
- Docker & Docker Compose
- PostgreSQL 13.2
- Redis
- HashiCorp Vault

### セットアップ
```bash
# 依存関係のインストール
make install-all

# GraphQLコード生成
make gqlgen

# サーバー起動
make server
```

### Docker環境
```bash
# インフラストラクチャ起動
docker-compose up -d

# サービス一覧
- PostgreSQL: localhost:25432
- Redis: localhost:26379
- Vault: localhost:28200
```

## 設定

### 環境設定
- `config/environments/local.yml`: ローカル環境設定
- `config/environments/test.yml`: テスト環境設定

### データベース接続設定
```yaml
database:
    driver: postgres
    user: sp
    password: sp
    host: localhost
    port: 25432
    db: sp
```

## プロジェクト構造

```
sp_api/
├── config/                 # 設定管理
├── db/                     # データベース関連
│   └── docker/init/        # 初期化スクリプト
├── domain/                 # ドメイン層
│   ├── model/              # エンティティ
│   └── repository/         # リポジトリインターフェース
├── graph/                  # GraphQL層
│   ├── generated/          # 自動生成コード
│   └── resolver/           # リゾルバー
├── infrastructure/         # インフラストラクチャ層
│   ├── persistence/        # データ永続化
│   └── provider/           # 外部サービス
├── usecase/               # ユースケース層
├── vault/                 # Vault設定
├── docker-compose.yml     # Docker設定
├── gqlgen.yml            # GraphQL設定
├── wire.go               # 依存性注入設定
└── server.go             # エントリーポイント
```

## 開発フロー

### GraphQLスキーマ変更時
1. `graph/schema.graphqls`を編集
2. `make gqlgen`でコード生成
3. リゾルバーを実装

### データベース変更時
1. `db/docker/init/1_ddl.sql`を編集
2. Dockerコンテナを再起動

### 依存性注入
Google Wireを使用して依存性を管理：
- `wire.go`: 依存性定義
- `wire_gen.go`: 自動生成された注入コード

## セキュリティ

- HashiCorp Vaultを使用した機密情報管理
- データベース接続プール設定
- 環境変数による設定の外部化

## 今後の拡張予定

- 認証・認可機能
- リアルタイム更新（WebSocket）
- キャッシュ機能（Redis）
- ログ管理
- メトリクス収集
- API バージョニング

## ライセンス

このプロジェクトのライセンス情報は、プロジェクトのルートディレクトリにあるLICENSEファイルを参照してください。