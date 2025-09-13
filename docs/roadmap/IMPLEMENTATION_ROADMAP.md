# SP API 実装ロードマップ

## 概要

SP APIプロジェクトの要件定義とアーキテクチャギャップ分析に基づき、段階的な実装計画を策定する。

## 1. 実装戦略

### 1.1 アプローチ
- **段階的開発**: 機能単位での段階的実装
- **MVP重視**: 最小限の機能で価値を提供
- **継続的改善**: フィードバックに基づく改善
- **技術的負債の解消**: 既存の問題の段階的解決

### 1.2 開発原則
- **ドメイン駆動設計**: ビジネスロジックの明確化
- **テスト駆動開発**: 品質の確保
- **APIファースト**: フロントエンドとの分離開発
- **セキュリティファースト**: セキュリティの早期実装

## 2. Phase 1: 基盤機能（2ヶ月）

### 2.1 目標
- ユーザー管理機能の実装
- 予想機能の実装
- 基本的な認証システムの実装
- GraphQL APIの拡張

### 2.2 実装内容

#### 2.2.1 ユーザー管理機能（3週間）

**Week 1: ドメインモデル設計**
```go
// internal/user/aggregate.go
type User struct {
    id        UserID
    email     Email
    password  PasswordHash
    profile   UserProfile
    role      UserRole
    status    UserStatus
    createdAt time.Time
    updatedAt time.Time
    
    events []DomainEvent
}

// internal/user/value_objects.go
type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    if !isValidEmail(email) {
        return Email{}, ErrInvalidEmail
    }
    return Email{value: email}, nil
}

type PasswordHash struct {
    value string
}

func NewPasswordHash(password string) (PasswordHash, error) {
    if len(password) < 8 {
        return PasswordHash{}, ErrPasswordTooShort
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return PasswordHash{}, err
    }
    return PasswordHash{value: string(hash)}, nil
}
```

**Week 2: リポジトリ実装**
```go
// internal/user/repository.go
type UserRepository interface {
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    Save(ctx context.Context, user *User) error
    Delete(ctx context.Context, id UserID) error
}

// infrastructure/persistence/user_repository.go
type userRepository struct {
    db *sql.DB
}

func (r *userRepository) FindByID(ctx context.Context, id UserID) (*User, error) {
    query := `
        SELECT id, email, password_hash, profile, role, status, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
    // 実装...
}
```

**Week 3: サービス実装**
```go
// internal/user/service.go
type UserService struct {
    repo UserRepository
    eventBus DomainEventBus
}

func (s *UserService) RegisterUser(ctx context.Context, input RegisterUserInput) (*User, error) {
    // 1. 入力値検証
    email, err := NewEmail(input.Email)
    if err != nil {
        return nil, err
    }
    
    password, err := NewPasswordHash(input.Password)
    if err != nil {
        return nil, err
    }
    
    // 2. 重複チェック
    existingUser, err := s.repo.FindByEmail(ctx, email)
    if err == nil && existingUser != nil {
        return nil, ErrUserAlreadyExists
    }
    
    // 3. ユーザー作成
    user := &User{
        id:        NewUserID(),
        email:     email,
        password:  password,
        profile:   input.Profile,
        role:      UserRoleUser,
        status:    UserStatusActive,
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }
    
    // 4. 保存
    if err := s.repo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. ドメインイベント発行
    user.addEvent(NewUserRegisteredEvent(user.id, user.email))
    s.eventBus.Publish(ctx, user.GetEvents()...)
    
    return user, nil
}
```

#### 2.2.2 認証システム（2週間）

**Week 1: JWT認証実装**
```go
// internal/auth/service.go
type AuthService struct {
    userRepo UserRepository
    jwtSecret string
    tokenExpiry time.Duration
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthToken, error) {
    // 1. ユーザー検索
    email, err := NewEmail(input.Email)
    if err != nil {
        return nil, err
    }
    
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }
    
    // 2. パスワード検証
    if !user.password.Verify(input.Password) {
        return nil, ErrInvalidCredentials
    }
    
    // 3. JWTトークン生成
    token, err := s.generateJWT(user)
    if err != nil {
        return nil, err
    }
    
    return &AuthToken{
        Token:     token,
        ExpiresAt: time.Now().Add(s.tokenExpiry),
        UserID:    user.id,
    }, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*User, error) {
    // JWTトークン検証
    claims, err := s.parseJWT(token)
    if err != nil {
        return nil, ErrInvalidToken
    }
    
    userID, err := NewUserID(claims.UserID)
    if err != nil {
        return nil, ErrInvalidToken
    }
    
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, ErrUserNotFound
    }
    
    return user, nil
}
```

**Week 2: 認証ミドルウェア実装**
```go
// internal/auth/middleware.go
func AuthMiddleware(authService AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Bearer トークンの抽出
        if strings.HasPrefix(token, "Bearer ") {
            token = token[7:]
        }
        
        user, err := authService.ValidateToken(c.Request.Context(), token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user", user)
        c.Next()
    }
}
```

#### 2.2.3 予想機能（3週間）

**Week 1: 予想ドメインモデル**
```go
// internal/prediction/aggregate.go
type Prediction struct {
    id        PredictionID
    userID    UserID
    matchID   MatchID
    type      PredictionType
    value     PredictionValue
    status    PredictionStatus
    points    int
    createdAt time.Time
    updatedAt time.Time
    
    events []DomainEvent
}

type PredictionType string

const (
    PredictionTypeMatchResult PredictionType = "MATCH_RESULT"
    PredictionTypeScore       PredictionType = "SCORE"
    PredictionTypeWinner      PredictionType = "WINNER"
)

type PredictionValue struct {
    matchResult *MatchResult
    score       *Score
    winner      *TeamID
}

type MatchResult string

const (
    MatchResultHome    MatchResult = "HOME"
    MatchResultAway    MatchResult = "AWAY"
    MatchResultDraw    MatchResult = "DRAW"
)
```

**Week 2: 予想サービス実装**
```go
// internal/prediction/service.go
type PredictionService struct {
    predRepo PredictionRepository
    matchRepo MatchRepository
    userRepo UserRepository
    eventBus DomainEventBus
}

func (s *PredictionService) CreatePrediction(ctx context.Context, input CreatePredictionInput) (*Prediction, error) {
    // 1. ユーザー検証
    user, err := s.userRepo.FindByID(ctx, input.UserID)
    if err != nil {
        return nil, err
    }
    
    // 2. 試合検証
    match, err := s.matchRepo.FindByID(ctx, input.MatchID)
    if err != nil {
        return nil, err
    }
    
    // 3. 予想制限チェック
    if err := s.validatePredictionConstraints(ctx, user.id, match); err != nil {
        return nil, err
    }
    
    // 4. 予想作成
    prediction := &Prediction{
        id:        NewPredictionID(),
        userID:    user.id,
        matchID:   match.id,
        type:      input.Type,
        value:     input.Value,
        status:    PredictionStatusPending,
        points:    0,
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }
    
    // 5. 保存
    if err := s.predRepo.Save(ctx, prediction); err != nil {
        return nil, err
    }
    
    // 6. ドメインイベント発行
    prediction.addEvent(NewPredictionCreatedEvent(prediction.id, prediction.userID, prediction.matchID))
    s.eventBus.Publish(ctx, prediction.GetEvents()...)
    
    return prediction, nil
}

func (s *PredictionService) validatePredictionConstraints(ctx context.Context, userID UserID, match *Match) error {
    // 試合開始前チェック
    if time.Now().After(match.kickoff) {
        return ErrMatchAlreadyStarted
    }
    
    // 既存予想チェック
    existing, err := s.predRepo.FindByUserAndMatch(ctx, userID, match.id)
    if err == nil && existing != nil {
        return ErrPredictionAlreadyExists
    }
    
    return nil
}
```

**Week 3: GraphQL API拡張**
```graphql
# graph/schema.graphqls
type User {
    id: ID!
    email: String!
    profile: UserProfile!
    role: UserRole!
    status: UserStatus!
    createdAt: DateTime!
    updatedAt: DateTime!
}

type UserProfile {
    name: String!
    nickname: String
    birthday: DateTime
    avatar: String
}

type Prediction {
    id: ID!
    user: User!
    match: Match!
    type: PredictionType!
    value: PredictionValue!
    status: PredictionStatus!
    points: Int!
    createdAt: DateTime!
    updatedAt: DateTime!
}

type PredictionValue {
    matchResult: MatchResult
    score: Score
    winner: Team
}

enum PredictionType {
    MATCH_RESULT
    SCORE
    WINNER
}

enum PredictionStatus {
    PENDING
    CORRECT
    INCORRECT
    PARTIAL
}

enum MatchResult {
    HOME
    AWAY
    DRAW
}

type Mutation {
    registerUser(input: RegisterUserInput!): User!
    loginUser(input: LoginInput!): AuthPayload!
    createPrediction(input: CreatePredictionInput!): Prediction!
    updatePrediction(id: ID!, input: UpdatePredictionInput!): Prediction!
    deletePrediction(id: ID!): Boolean!
}

type Query {
    me: User
    predictions(userId: ID, matchId: ID): [Prediction!]!
    prediction(id: ID!): Prediction
}
```

### 2.3 テスト実装

#### 2.3.1 単体テスト
```go
// internal/user/service_test.go
func TestUserService_RegisterUser(t *testing.T) {
    tests := []struct {
        name    string
        input   RegisterUserInput
        wantErr bool
        err     error
    }{
        {
            name: "正常なユーザー登録",
            input: RegisterUserInput{
                Email:    "test@example.com",
                Password: "password123",
                Profile: UserProfile{
                    Name: "Test User",
                },
            },
            wantErr: false,
        },
        {
            name: "無効なメールアドレス",
            input: RegisterUserInput{
                Email:    "invalid-email",
                Password: "password123",
                Profile: UserProfile{
                    Name: "Test User",
                },
            },
            wantErr: true,
            err:     ErrInvalidEmail,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実装
        })
    }
}
```

#### 2.3.2 統合テスト
```go
// integration/user_test.go
func TestUserIntegration(t *testing.T) {
    // テストデータベースのセットアップ
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // サービスの初期化
    userRepo := persistence.NewUserRepository(db)
    userService := user.NewUserService(userRepo, mockEventBus)
    
    // テスト実行
    user, err := userService.RegisterUser(context.Background(), RegisterUserInput{
        Email:    "test@example.com",
        Password: "password123",
        Profile: UserProfile{
            Name: "Test User",
        },
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email.String())
}
```

### 2.4 成果物
- ユーザー管理機能の完全実装
- 認証システムの完全実装
- 予想機能の完全実装
- GraphQL APIの拡張
- 単体テスト・統合テストの実装

## 3. Phase 2: 拡張機能（2ヶ月）

### 3.1 目標
- ポイント・ランキング機能の実装
- 検索機能の実装
- パフォーマンス最適化
- セキュリティ強化

### 3.2 実装内容

#### 3.2.1 ポイント・ランキング機能（3週間）

**Week 1: ポイント計算システム**
```go
// internal/point/service.go
type PointService struct {
    predRepo PredictionRepository
    matchRepo MatchRepository
    userRepo UserRepository
    pointRepo PointRepository
}

func (s *PointService) CalculatePoints(ctx context.Context, matchID MatchID) error {
    // 1. 試合結果取得
    match, err := s.matchRepo.FindByID(ctx, matchID)
    if err != nil {
        return err
    }
    
    // 2. 予想一覧取得
    predictions, err := s.predRepo.FindByMatch(ctx, matchID)
    if err != nil {
        return err
    }
    
    // 3. 各予想のポイント計算
    for _, pred := range predictions {
        points := s.calculatePredictionPoints(pred, match)
        pred.points = points
        pred.status = s.determinePredictionStatus(pred, match)
        
        if err := s.predRepo.Save(ctx, pred); err != nil {
            return err
        }
        
        // 4. ユーザーポイント更新
        if err := s.updateUserPoints(ctx, pred.userID, points); err != nil {
            return err
        }
    }
    
    return nil
}

func (s *PointService) calculatePredictionPoints(pred *Prediction, match *Match) int {
    switch pred.type {
    case PredictionTypeMatchResult:
        return s.calculateMatchResultPoints(pred, match)
    case PredictionTypeScore:
        return s.calculateScorePoints(pred, match)
    case PredictionTypeWinner:
        return s.calculateWinnerPoints(pred, match)
    default:
        return 0
    }
}
```

**Week 2: ランキングシステム**
```go
// internal/ranking/service.go
type RankingService struct {
    userRepo UserRepository
    pointRepo PointRepository
    rankingRepo RankingRepository
}

func (s *RankingService) GetLeaderboard(ctx context.Context, eventID EventID, period LeaderboardPeriod) ([]LeaderboardEntry, error) {
    // 1. ユーザーポイント取得
    userPoints, err := s.pointRepo.GetUserPointsByEvent(ctx, eventID, period)
    if err != nil {
        return nil, err
    }
    
    // 2. ランキング計算
    entries := make([]LeaderboardEntry, 0, len(userPoints))
    for i, up := range userPoints {
        user, err := s.userRepo.FindByID(ctx, up.UserID)
        if err != nil {
            continue
        }
        
        entries = append(entries, LeaderboardEntry{
            Rank:     i + 1,
            User:     user,
            Points:   up.Points,
            Accuracy: up.Accuracy,
        })
    }
    
    return entries, nil
}
```

**Week 3: 統計機能**
```go
// internal/statistics/service.go
type StatisticsService struct {
    predRepo PredictionRepository
    userRepo UserRepository
    matchRepo MatchRepository
}

func (s *StatisticsService) GetUserStats(ctx context.Context, userID UserID) (*UserStats, error) {
    // 1. 予想統計取得
    predictions, err := s.predRepo.FindByUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 2. 統計計算
    stats := &UserStats{
        TotalPredictions: len(predictions),
        CorrectPredictions: 0,
        IncorrectPredictions: 0,
        PartialPredictions: 0,
        TotalPoints: 0,
    }
    
    for _, pred := range predictions {
        stats.TotalPoints += pred.points
        switch pred.status {
        case PredictionStatusCorrect:
            stats.CorrectPredictions++
        case PredictionStatusIncorrect:
            stats.IncorrectPredictions++
        case PredictionStatusPartial:
            stats.PartialPredictions++
        }
    }
    
    // 3. 精度計算
    if stats.TotalPredictions > 0 {
        stats.Accuracy = float64(stats.CorrectPredictions) / float64(stats.TotalPredictions)
    }
    
    return stats, nil
}
```

#### 3.2.2 検索機能（3週間）

**Week 1: Elasticsearch統合**
```go
// infrastructure/search/elasticsearch.go
type ElasticsearchClient struct {
    client *elasticsearch.Client
    index  string
}

func NewElasticsearchClient(config ElasticsearchConfig) (*ElasticsearchClient, error) {
    client, err := elasticsearch.NewClient(elasticsearch.Config{
        Addresses: config.Addresses,
        Username:  config.Username,
        Password:  config.Password,
    })
    if err != nil {
        return nil, err
    }
    
    return &ElasticsearchClient{
        client: client,
        index:  config.Index,
    }, nil
}

func (c *ElasticsearchClient) IndexPlayer(ctx context.Context, player *Player) error {
    doc := map[string]interface{}{
        "id":       player.id.String(),
        "name":     player.name.String(),
        "team":     player.team.name.String(),
        "position": player.position.String(),
        "number":   player.number,
    }
    
    _, err := c.client.Index(
        c.index,
        strings.NewReader(toJSON(doc)),
        c.client.Index.WithDocumentID(player.id.String()),
        c.client.Index.WithContext(ctx),
    )
    
    return err
}
```

**Week 2: 検索サービス実装**
```go
// internal/search/service.go
type SearchService struct {
    esClient *ElasticsearchClient
    playerRepo PlayerRepository
    teamRepo TeamRepository
}

func (s *SearchService) SearchPlayers(ctx context.Context, query string, filters SearchFilters) ([]Player, error) {
    // 1. Elasticsearchクエリ構築
    esQuery := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []interface{}{
                    map[string]interface{}{
                        "multi_match": map[string]interface{}{
                            "query":  query,
                            "fields": []string{"name^2", "team", "position"},
                        },
                    },
                },
            },
        },
    }
    
    // 2. フィルター適用
    if filters.Team != "" {
        esQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = []interface{}{
            map[string]interface{}{
                "term": map[string]interface{}{
                    "team": filters.Team,
                },
            },
        }
    }
    
    // 3. 検索実行
    response, err := s.esClient.client.Search(
        s.esClient.client.Search.WithIndex(s.esClient.index),
        s.esClient.client.Search.WithBody(strings.NewReader(toJSON(esQuery))),
        s.esClient.client.Search.WithContext(ctx),
    )
    if err != nil {
        return nil, err
    }
    
    // 4. 結果処理
    var results []Player
    // レスポンス解析とPlayerオブジェクトの構築
    // ...
    
    return results, nil
}
```

**Week 3: 検索API実装**
```graphql
# graph/schema.graphqls
type SearchResult {
    players: [Player!]!
    teams: [Team!]!
    total: Int!
}

input SearchFilters {
    team: String
    position: Position
    event: ID
}

type Query {
    search(query: String!, filters: SearchFilters): SearchResult!
}
```

#### 3.2.3 パフォーマンス最適化（2週間）

**Week 1: Redis キャッシュ実装**
```go
// infrastructure/cache/redis.go
type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewRedisCache(config RedisConfig) (*RedisCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     config.Addr,
        Password: config.Password,
        DB:       config.DB,
    })
    
    return &RedisCache{
        client: client,
        ttl:    config.TTL,
    }, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
    return c.client.Get(ctx, key).Bytes()
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte) error {
    return c.client.Set(ctx, key, value, c.ttl).Err()
}

// internal/user/service.go
func (s *UserService) FindByID(ctx context.Context, id UserID) (*User, error) {
    // 1. キャッシュから取得を試行
    cacheKey := fmt.Sprintf("user:%s", id.String())
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        var user User
        if err := json.Unmarshal(cached, &user); err == nil {
            return &user, nil
        }
    }
    
    // 2. データベースから取得
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. キャッシュに保存
    if data, err := json.Marshal(user); err == nil {
        s.cache.Set(ctx, cacheKey, data)
    }
    
    return user, nil
}
```

**Week 2: データベース最適化**
```sql
-- インデックス最適化
CREATE INDEX CONCURRENTLY idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY idx_predictions_user_match ON predictions(user_id, match_id);
CREATE INDEX CONCURRENTLY idx_predictions_match ON predictions(match_id);
CREATE INDEX CONCURRENTLY idx_matches_event ON matches(event_id);
CREATE INDEX CONCURRENTLY idx_matches_kickoff ON matches(kickoff);

-- パーティショニング（大量データの場合）
CREATE TABLE predictions_2024 PARTITION OF predictions
FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
```

### 3.3 成果物
- ポイント・ランキング機能の完全実装
- 検索機能の完全実装
- パフォーマンス最適化の実装
- キャッシュシステムの実装

## 4. Phase 3: 高度機能（2ヶ月）

### 4.1 目標
- 通知機能の実装
- マイクロサービス化
- 監視・運用機能の実装
- セキュリティ強化

### 4.2 実装内容

#### 4.2.1 通知機能（3週間）

**Week 1: WebSocket実装**
```go
// internal/notification/hub.go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    userClients map[UserID][]*Client
}

type Client struct {
    hub    *Hub
    conn   *websocket.Conn
    send   chan []byte
    userID UserID
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            h.userClients[client.userID] = append(h.userClients[client.userID], client)
            
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                h.removeUserClient(client.userID, client)
            }
            
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}
```

**Week 2: 通知サービス実装**
```go
// internal/notification/service.go
type NotificationService struct {
    hub      *Hub
    userRepo UserRepository
    notifRepo NotificationRepository
}

func (s *NotificationService) SendNotification(ctx context.Context, userID UserID, notification Notification) error {
    // 1. 通知保存
    if err := s.notifRepo.Save(ctx, &notification); err != nil {
        return err
    }
    
    // 2. WebSocketでリアルタイム送信
    message := NotificationMessage{
        Type: notification.Type,
        Title: notification.Title,
        Body: notification.Body,
        Data: notification.Data,
        Timestamp: notification.CreatedAt,
    }
    
    s.hub.SendToUser(userID, message)
    
    return nil
}

func (s *NotificationService) SendMatchStartNotification(ctx context.Context, matchID MatchID) error {
    // 1. 予想したユーザーを取得
    predictions, err := s.predRepo.FindByMatch(ctx, matchID)
    if err != nil {
        return err
    }
    
    // 2. 各ユーザーに通知送信
    for _, pred := range predictions {
        notification := Notification{
            ID:        NewNotificationID(),
            UserID:    pred.userID,
            Type:      NotificationTypeMatchStart,
            Title:     "試合開始",
            Body:      fmt.Sprintf("予想した試合が開始しました"),
            Data:      map[string]interface{}{"match_id": matchID.String()},
            CreatedAt: time.Now(),
        }
        
        if err := s.SendNotification(ctx, pred.userID, notification); err != nil {
            // ログ出力
            continue
        }
    }
    
    return nil
}
```

**Week 3: 通知API実装**
```graphql
# graph/schema.graphqls
type Notification {
    id: ID!
    type: NotificationType!
    title: String!
    body: String!
    data: JSON
    read: Boolean!
    createdAt: DateTime!
}

enum NotificationType {
    MATCH_START
    MATCH_RESULT
    RANKING_UPDATE
    SYSTEM_ANNOUNCEMENT
}

type Subscription {
    notifications: Notification!
}

type Mutation {
    markNotificationAsRead(id: ID!): Boolean!
    markAllNotificationsAsRead: Boolean!
}
```

#### 4.2.2 マイクロサービス化（3週間）

**Week 1: サービス分割設計**
```go
// cmd/user-service/main.go
func main() {
    // ユーザーサービス専用の設定
    config := loadConfig()
    
    // 依存性注入
    container := fx.New(
        fx.Provide(
            NewUserRepository,
            NewUserService,
            NewAuthService,
            NewUserController,
        ),
        fx.Invoke(RegisterUserRoutes),
    )
    
    // サーバー起動
    server := NewServer(config)
    server.Run()
}

// cmd/prediction-service/main.go
func main() {
    // 予想サービス専用の設定
    config := loadConfig()
    
    // 依存性注入
    container := fx.New(
        fx.Provide(
            NewPredictionRepository,
            NewPredictionService,
            NewPointService,
            NewPredictionController,
        ),
        fx.Invoke(RegisterPredictionRoutes),
    )
    
    // サーバー起動
    server := NewServer(config)
    server.Run()
}
```

**Week 2: サービス間通信実装**
```go
// internal/communication/grpc.go
type UserServiceClient struct {
    client userpb.UserServiceClient
}

func NewUserServiceClient(conn *grpc.ClientConn) *UserServiceClient {
    return &UserServiceClient{
        client: userpb.NewUserServiceClient(conn),
    }
}

func (c *UserServiceClient) GetUser(ctx context.Context, userID UserID) (*User, error) {
    req := &userpb.GetUserRequest{
        UserId: userID.String(),
    }
    
    resp, err := c.client.GetUser(ctx, req)
    if err != nil {
        return nil, err
    }
    
    return &User{
        id:    UserID(resp.User.Id),
        email: Email(resp.User.Email),
        // ... その他のフィールド
    }, nil
}

// internal/communication/eventbus.go
type EventBus struct {
    publisher EventPublisher
    subscribers map[string][]EventSubscriber
}

func (b *EventBus) Publish(ctx context.Context, event DomainEvent) error {
    return b.publisher.Publish(ctx, event)
}

func (b *EventBus) Subscribe(eventType string, subscriber EventSubscriber) {
    b.subscribers[eventType] = append(b.subscribers[eventType], subscriber)
}
```

**Week 3: API Gateway実装**
```go
// cmd/api-gateway/main.go
func main() {
    // API Gateway設定
    config := loadConfig()
    
    // サービス接続設定
    userServiceConn := connectToUserService(config.UserServiceAddr)
    predictionServiceConn := connectToPredictionService(config.PredictionServiceAddr)
    
    // ルーティング設定
    router := gin.Default()
    
    // 認証ミドルウェア
    authMiddleware := NewAuthMiddleware(userServiceConn)
    router.Use(authMiddleware)
    
    // ルート設定
    router.POST("/users/*path", proxyToUserService(userServiceConn))
    router.POST("/predictions/*path", proxyToPredictionService(predictionServiceConn))
    
    // サーバー起動
    server := &http.Server{
        Addr:    config.Addr,
        Handler: router,
    }
    
    server.ListenAndServe()
}
```

#### 4.2.3 監視・運用機能（2週間）

**Week 1: メトリクス収集**
```go
// internal/monitoring/metrics.go
type MetricsCollector struct {
    requestDuration prometheus.HistogramVec
    requestTotal    prometheus.CounterVec
    activeUsers     prometheus.Gauge
    errorRate       prometheus.CounterVec
}

func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "http_request_duration_seconds",
                Help: "HTTP request duration in seconds",
            },
            []string{"method", "endpoint", "status"},
        ),
        requestTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "endpoint", "status"},
        ),
        activeUsers: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "active_users_total",
                Help: "Total number of active users",
            },
        ),
        errorRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "errors_total",
                Help: "Total number of errors",
            },
            []string{"service", "error_type"},
        ),
    }
}

// ミドルウェア
func MetricsMiddleware(collector *MetricsCollector) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        
        collector.requestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            status,
        ).Observe(duration)
        
        collector.requestTotal.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            status,
        ).Inc()
    }
}
```

**Week 2: ログ管理・アラート**
```go
// internal/logging/logger.go
type Logger struct {
    logger *zap.Logger
}

func NewLogger(config LogConfig) (*Logger, error) {
    var zapConfig zap.Config
    if config.Environment == "production" {
        zapConfig = zap.NewProductionConfig()
    } else {
        zapConfig = zap.NewDevelopmentConfig()
    }
    
    logger, err := zapConfig.Build()
    if err != nil {
        return nil, err
    }
    
    return &Logger{logger: logger}, nil
}

func (l *Logger) Info(ctx context.Context, message string, fields ...zap.Field) {
    l.logger.Info(message, append(fields, zap.String("trace_id", getTraceID(ctx)))...)
}

func (l *Logger) Error(ctx context.Context, message string, err error, fields ...zap.Field) {
    l.logger.Error(message, append(fields, zap.Error(err), zap.String("trace_id", getTraceID(ctx)))...)
}

// アラート設定
type AlertManager struct {
    slackClient *slack.Client
    emailClient *email.Client
}

func (a *AlertManager) SendAlert(ctx context.Context, alert Alert) error {
    switch alert.Severity {
    case AlertSeverityCritical:
        return a.sendSlackAlert(ctx, alert)
    case AlertSeverityWarning:
        return a.sendEmailAlert(ctx, alert)
    default:
        return a.sendLogAlert(ctx, alert)
    }
}
```

### 4.3 成果物
- 通知機能の完全実装
- マイクロサービス化の実装
- 監視・運用機能の実装
- セキュリティ強化の実装

## 5. 品質保証

### 5.1 テスト戦略
- **単体テスト**: 各機能の単体テスト（カバレッジ80%以上）
- **統合テスト**: サービス間の統合テスト
- **E2Eテスト**: エンドツーエンドテスト
- **負荷テスト**: パフォーマンステスト

### 5.2 コード品質
- **golangci-lint**: コード品質チェック
- **gosec**: セキュリティチェック
- **コードレビュー**: プルリクエストレビュー
- **リファクタリング**: 定期的なリファクタリング

### 5.3 セキュリティ
- **脆弱性スキャン**: 定期的な脆弱性スキャン
- **セキュリティテスト**: ペネトレーションテスト
- **セキュリティ監査**: 外部監査の実施

## 6. リスク管理

### 6.1 技術的リスク
- **GraphQLスキーマの変更**: 段階的な移行
- **データベーススキーマの変更**: マイグレーション戦略
- **パフォーマンスの劣化**: プロファイリングと最適化

### 6.2 運用リスク
- **開発期間の延長**: 段階的な開発
- **チームの学習コスト**: 教育・研修の実施
- **システムの複雑性**: ドキュメント整備

### 6.3 対策
- **段階的移行**: リスクを分散した移行計画
- **監視体制**: 24/7監視体制の構築
- **バックアップ**: 複数のバックアップ戦略
- **テスト**: 十分なテストの実施

## 7. 成功指標

### 7.1 機能指標
- ユーザー登録・認証機能の実装完了
- 予想機能の実装完了
- ポイント・ランキング機能の実装完了
- 検索機能の実装完了
- 通知機能の実装完了

### 7.2 非機能指標
- API レスポンス時間 200ms 以内
- 同時接続数 1,000 ユーザー対応
- 稼働率 99.9% 以上
- テストカバレッジ 80% 以上

### 7.3 品質指標
- コード品質スコア 8.0 以上
- セキュリティ脆弱性 0件
- 技術的負債の解消率 90% 以上

---

**文書作成日**: 2024年12月19日  
**バージョン**: 1.0  
**作成者**: AI Assistant  
**承認者**: [承認者名]
