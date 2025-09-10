# SP API アーキテクチャモダナイゼーション提案

## 現状分析

### 現在のアーキテクチャの問題点

1. **レイヤー中心の構造**: 現在は技術的な層（domain, usecase, infrastructure）で分離されているが、ビジネス機能ごとの分離が不十分
2. **貧弱なドメインモデル**: エンティティが単純なデータ構造体で、ビジネスロジックが不足
3. **古いGoバージョン**: Go 1.16は2021年のバージョンで、Genericsなどの重要な機能が利用できない
4. **基本的なClean Architecture**: DDDの概念（Aggregate, Value Object, Domain Event）が実装されていない
5. **GraphQLスキーマとドメインモデルの密結合**: 自動生成されたモデルがドメインロジックを制約している

## 提案する現代的なアーキテクチャ

### 1. 技術スタックの更新

#### Go バージョンとツール
- **Go 1.22+**: Generics、パフォーマンス改善、セキュリティ更新
- **golangci-lint**: コード品質とベストプラクティスの強制
- **Wire**: 依存性注入（現状維持、優秀な選択）
- **Atlas**: データベースマイグレーション管理

#### 依存性注入の改善
```go
// 現在: 手動での依存性注入
func InitializeEventUseCase(cfg Config) usecase.EventUseCase {
    wire.Build(
        persistence.NewEventPersistence,
        usecase.NewEventUseCase,
    )
    return nil
}

// 提案: より構造化されたDI
var EventModule = fx.Module("event",
    fx.Provide(
        NewEventRepository,
        NewEventService,
        NewEventController,
    ),
    fx.Invoke(RegisterEventRoutes),
)
```

### 2. DDDベースのディレクトリ構造

#### 現在の構造（レイヤー中心）
```
sp_api/
├── domain/
│   ├── model/
│   └── repository/
├── usecase/
├── infrastructure/
└── graph/
```

#### 提案する構造（機能中心）
```
sp_api/
├── cmd/
│   └── api/
│       ├── main.go
│       └── wire.go
├── internal/
│   ├── event/
│   │   ├── aggregate.go          # Event Aggregate Root
│   │   ├── value_objects.go      # Score, MatchTime, etc.
│   │   ├── events.go             # Domain Events
│   │   ├── repository.go         # Repository Interface
│   │   ├── service.go            # Application Service
│   │   └── controller.go         # HTTP/GraphQL Controller
│   ├── team/
│   │   ├── aggregate.go          # Team Aggregate Root
│   │   ├── value_objects.go     # TeamName, FIFA Code, etc.
│   │   ├── events.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── controller.go
│   ├── player/
│   │   ├── aggregate.go
│   │   ├── value_objects.go     # PlayerName, Position, etc.
│   │   ├── events.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── controller.go
│   ├── shared/
│   │   ├── domain/
│   │   │   ├── events.go        # Domain Event Infrastructure
│   │   │   ├── value_object.go  # Base Value Object
│   │   │   └── aggregate.go      # Base Aggregate
│   │   ├── infrastructure/
│   │   │   ├── database/
│   │   │   ├── eventbus/
│   │   │   └── config/
│   │   └── application/
│   │       ├── errors.go
│   │       └── responses.go
│   └── graphql/
│       ├── schema.graphqls
│       ├── resolvers/
│       └── generated/
├── pkg/                          # 外部向けライブラリ
│   ├── logger/
│   ├── validator/
│   └── utils/
├── migrations/
├── docker/
└── docs/
```

### 3. DDDパターンの実装

#### Aggregate Root の例
```go
// internal/event/aggregate.go
type Event struct {
    id          EventID
    name        EventName
    nameJP      EventNameJP
    participants []TeamID
    matches      []Match
    status      EventStatus
    createdAt   time.Time
    updatedAt   time.Time
    
    // Domain Events
    events []DomainEvent
}

// ビジネスロジックを含むメソッド
func (e *Event) AddMatch(match Match) error {
    if e.status != EventStatusActive {
        return ErrEventNotActive
    }
    
    if len(e.matches) >= MaxMatchesPerEvent {
        return ErrTooManyMatches
    }
    
    e.matches = append(e.matches, match)
    e.addEvent(NewMatchAddedEvent(e.id, match.ID()))
    
    return nil
}

func (e *Event) StartEvent() error {
    if e.status != EventStatusScheduled {
        return ErrEventNotScheduled
    }
    
    e.status = EventStatusActive
    e.addEvent(NewEventStartedEvent(e.id))
    
    return nil
}
```

#### Value Objects の例
```go
// internal/event/value_objects.go
type EventName struct {
    value string
}

func NewEventName(name string) (EventName, error) {
    if len(name) < 3 || len(name) > 100 {
        return EventName{}, ErrInvalidEventName
    }
    
    return EventName{value: name}, nil
}

func (n EventName) String() string {
    return n.value
}

type Score struct {
    home int
    away int
}

func NewScore(home, away int) (Score, error) {
    if home < 0 || away < 0 {
        return Score{}, ErrInvalidScore
    }
    
    return Score{home: home, away: away}, nil
}

func (s Score) Home() int { return s.home }
func (s Score) Away() int { return s.away }
func (s Score) IsDraw() bool { return s.home == s.away }
func (s Score) Winner() *TeamID {
    if s.home > s.away {
        return &s.home
    } else if s.away > s.home {
        return &s.away
    }
    return nil
}
```

#### Domain Events の例
```go
// internal/event/events.go
type EventStarted struct {
    EventID EventID
    StartedAt time.Time
}

func NewEventStartedEvent(eventID EventID) EventStarted {
    return EventStarted{
        EventID: eventID,
        StartedAt: time.Now(),
    }
}

type MatchAdded struct {
    EventID EventID
    MatchID MatchID
    AddedAt time.Time
}

func NewMatchAddedEvent(eventID EventID, matchID MatchID) MatchAdded {
    return MatchAdded{
        EventID: eventID,
        MatchID: matchID,
        AddedAt: time.Now(),
    }
}
```

### 4. Application Service の改善

```go
// internal/event/service.go
type EventService struct {
    repo        EventRepository
    eventBus    DomainEventBus
    teamRepo    TeamRepository
}

func NewEventService(repo EventRepository, eventBus DomainEventBus, teamRepo TeamRepository) *EventService {
    return &EventService{
        repo:     repo,
        eventBus: eventBus,
        teamRepo: teamRepo,
    }
}

func (s *EventService) StartEvent(ctx context.Context, eventID EventID) error {
    // 1. Aggregateを取得
    event, err := s.repo.FindByID(ctx, eventID)
    if err != nil {
        return err
    }
    
    // 2. ビジネスロジックを実行
    if err := event.StartEvent(); err != nil {
        return err
    }
    
    // 3. 変更を保存
    if err := s.repo.Save(ctx, event); err != nil {
        return err
    }
    
    // 4. Domain Eventsを発行
    s.eventBus.Publish(ctx, event.GetEvents()...)
    
    return nil
}
```

### 5. Repository パターンの改善

```go
// internal/event/repository.go
type EventRepository interface {
    FindByID(ctx context.Context, id EventID) (*Event, error)
    FindByName(ctx context.Context, name EventName) (*Event, error)
    FindAll(ctx context.Context) ([]*Event, error)
    Save(ctx context.Context, event *Event) error
    Delete(ctx context.Context, id EventID) error
}

// internal/shared/infrastructure/persistence/event_repository.go
type eventRepository struct {
    db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
    return &eventRepository{db: db}
}

func (r *eventRepository) FindByID(ctx context.Context, id EventID) (*Event, error) {
    // 複雑なクエリでAggregateを再構築
    query := `
        SELECT e.id, e.name, e.name_jp, e.status, e.created_at, e.updated_at,
               m.id, m.stadium_id, m.kickoff, m.category,
               p.event_id, p.team_id
        FROM events e
        LEFT JOIN matches m ON e.id = m.event_id
        LEFT JOIN participants p ON e.id = p.event_id
        WHERE e.id = $1
    `
    
    // クエリ実行とAggregate再構築
    // ...
}
```

### 6. データストアの評価

#### PostgreSQL の継続推奨理由
1. **ACID特性**: トランザクションの整合性が重要
2. **複雑なクエリ**: サッカーデータは関係性が複雑
3. **JSONサポート**: 柔軟なスキーマ変更に対応
4. **パフォーマンス**: 適切なインデックス設計で高速
5. **成熟度**: 安定性とコミュニティサポート

#### 追加検討事項
- **Redis**: キャッシュ層として活用
- **Elasticsearch**: 検索機能の強化
- **Event Sourcing**: 将来の拡張性を考慮

### 7. GraphQL スキーマの改善

```graphql
# 現在のスキーマ（自動生成ベース）
type Event {
    id: ID!
    name: String!
    name_jp: String!
    participants: [Team!]!
    matches: [Match!]!
}

# 提案するスキーマ（ドメイン駆動）
type Event {
    id: ID!
    name: EventName!
    nameJP: EventNameJP!
    status: EventStatus!
    participants: [Team!]!
    matches: [Match!]!
    createdAt: DateTime!
    updatedAt: DateTime!
}

enum EventStatus {
    SCHEDULED
    ACTIVE
    COMPLETED
    CANCELLED
}

input CreateEventInput {
    name: String!
    nameJP: String!
    participantTeamIds: [ID!]!
}

type EventMutation {
    startEvent(eventId: ID!): Event!
    addMatch(eventId: ID!, matchInput: CreateMatchInput!): Event!
    completeEvent(eventId: ID!): Event!
}
```

## 実装計画

### Phase 1: 基盤整備
1. Go 1.22+ へのアップグレード
2. golangci-lint の導入
3. 新しいディレクトリ構造の作成
4. 基本的なDDDパターンの実装

### Phase 2: ドメイン実装
1. Event Aggregate の実装
2. Team Aggregate の実装
3. Player Aggregate の実装
4. Value Objects の実装

### Phase 3: インフラ整備
1. Repository 実装の改善
2. Domain Events の実装
3. Event Bus の実装
4. テストの追加

### Phase 4: API改善
1. GraphQL スキーマの更新
2. Resolver の改善
3. エラーハンドリングの統一
4. ドキュメントの整備

## 期待される効果

1. **保守性の向上**: 機能ごとの分離により、変更の影響範囲が明確
2. **テスタビリティ**: ドメインロジックの分離により、単体テストが容易
3. **拡張性**: 新しい機能の追加が容易
4. **型安全性**: Value Objects により、コンパイル時のエラー検出
5. **ビジネスロジックの明確化**: ドメインエキスパートとの協業が容易

## リスクと対策

1. **学習コスト**: DDDパターンの習得が必要
   - 対策: 段階的な導入とチーム教育
2. **複雑性の増加**: 初期実装が複雑になる可能性
   - 対策: シンプルな実装から開始し、段階的に複雑化
3. **パフォーマンス**: オブジェクトの生成コスト
   - 対策: プロファイリングと最適化

