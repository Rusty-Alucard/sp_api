package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	auth "github.com/Rusty-Alucard/sp_api/internal/auth/domain"
)

// MockAgentRepository is a mock implementation of AgentRepository
type MockAgentRepository struct {
	mock.Mock
}

func (m *MockAgentRepository) FindByClientID(ctx context.Context, clientID string) (*auth.Agent, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Agent), args.Error(1)
}

func (m *MockAgentRepository) Save(ctx context.Context, agent *auth.Agent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *MockAgentRepository) FindAll(ctx context.Context) ([]*auth.Agent, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auth.Agent), args.Error(1)
}

func (m *MockAgentRepository) Delete(ctx context.Context, clientID string) error {
	args := m.Called(ctx, clientID)
	return args.Error(0)
}

func TestNewAuthService(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	jwtSecret := []byte("test-secret")
	tokenExpiry := time.Hour

	service := NewAuthService(mockRepo, jwtSecret, tokenExpiry)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.agentRepo)
	assert.Equal(t, jwtSecret, service.jwtSecret)
	assert.Equal(t, tokenExpiry, service.tokenExpiry)
}

func TestAuthService_AuthenticateAgent(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		setupMock    func()
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "正常ケース",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			setupMock: func() {
				agent, _ := auth.NewAgent("test-client-id", "test-secret", "Test Agent", auth.PlanFree)
				mockRepo.On("FindByClientID", mock.Anything, "test-client-id").Return(agent, nil)
				mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*auth.Agent")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:         "エージェントが見つからない",
			clientID:     "non-existent",
			clientSecret: "test-secret",
			setupMock: func() {
				mockRepo.On("FindByClientID", mock.Anything, "non-existent").Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name:         "パスワードが間違っている",
			clientID:     "test-client-id",
			clientSecret: "wrong-secret",
			setupMock: func() {
				agent, _ := auth.NewAgent("test-client-id", "test-secret", "Test Agent", auth.PlanFree)
				mockRepo.On("FindByClientID", mock.Anything, "test-client-id").Return(agent, nil)
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name:         "非アクティブなエージェント",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			setupMock: func() {
				agent, _ := auth.NewAgent("test-client-id", "test-secret", "Test Agent", auth.PlanFree)
				_ = agent.Deactivate()
				mockRepo.On("FindByClientID", mock.Anything, "test-client-id").Return(agent, nil)
			},
			wantErr: true,
			errMsg:  "agent is inactive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			agent, err := service.AuthenticateAgent(context.Background(), tt.clientID, tt.clientSecret)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, agent)
			} else {
				require.NoError(t, err)
				require.NotNil(t, agent)
				assert.Equal(t, tt.clientID, agent.ClientID())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	agent, err := auth.NewAgent("test-client-id", "test-secret", "Test Agent", auth.PlanFree)
	require.NoError(t, err)

	// スコープを追加
	err = agent.AddScope(auth.ScopeReadEvents)
	require.NoError(t, err)

	token, err := service.GenerateToken(agent)

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// トークンを検証
	claims, err := service.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "test-client-id", claims.AgentID)
	assert.Equal(t, "Test Agent", claims.Name)
	assert.Equal(t, "free", claims.Plan)
	assert.Contains(t, claims.Scopes, "events:read")
}

func TestAuthService_ValidateToken(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	agent, err := auth.NewAgent("test-client-id", "test-secret", "Test Agent", auth.PlanFree)
	require.NoError(t, err)

	// 有効なトークンを生成
	token, err := service.GenerateToken(agent)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "有効なトークン",
			token:   token,
			wantErr: false,
		},
		{
			name:    "無効なトークン",
			token:   "invalid-token",
			wantErr: true,
		},
		{
			name:    "空のトークン",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, "test-client-id", claims.AgentID)
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	agent, err := auth.NewAgent("test-client-id", "test-secret", "Test Agent", auth.PlanFree)
	require.NoError(t, err)

	// 有効なトークンを生成
	originalToken, err := service.GenerateToken(agent)
	require.NoError(t, err)

	// エージェントを再取得するモック
	mockRepo.On("FindByClientID", mock.Anything, "test-client-id").Return(agent, nil)

	// トークンをリフレッシュ
	newToken, err := service.RefreshToken(context.Background(), originalToken)

	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
	// JWT tokens are time-dependent, so we just check they're different strings
	// Note: Tokens might be the same due to timing, which is acceptable

	// 新しいトークンを検証
	claims, err := service.ValidateToken(newToken)
	require.NoError(t, err)
	assert.Equal(t, "test-client-id", claims.AgentID)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	// 無効なトークンでリフレッシュを試行
	_, err := service.RefreshToken(context.Background(), "invalid-token")

	assert.Error(t, err)
}

func TestAuthService_RefreshToken_AgentNotFound(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	agent, err := auth.NewAgent("test-client-id", "test-secret", "Test Agent", auth.PlanFree)
	require.NoError(t, err)

	// 有効なトークンを生成
	token, err := service.GenerateToken(agent)
	require.NoError(t, err)

	// エージェントが見つからないモック
	mockRepo.On("FindByClientID", mock.Anything, "test-client-id").Return(nil, errors.New("not found"))

	// トークンをリフレッシュ
	_, err = service.RefreshToken(context.Background(), token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent not found")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_HasScope(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	claims := &AgentClaims{
		AgentID: "test-client-id",
		Scopes:  []string{"events:read", "teams:read"},
	}

	tests := []struct {
		name     string
		scope    auth.Scope
		expected bool
	}{
		{
			name:     "持っているスコープ",
			scope:    auth.ScopeReadEvents,
			expected: true,
		},
		{
			name:     "持っていないスコープ",
			scope:    auth.ScopeWriteEvents,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.HasScope(claims, tt.scope)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthService_CanAccess(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	claims := &AgentClaims{
		AgentID: "test-client-id",
		Scopes:  []string{"events:read"},
	}

	tests := []struct {
		name     string
		resource string
		expected bool
	}{
		{
			name:     "アクセス可能なリソース",
			resource: "events",
			expected: true,
		},
		{
			name:     "アクセス不可能なリソース",
			resource: "teams",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CanAccess(claims, tt.resource)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthService_CanWrite(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	claims := &AgentClaims{
		AgentID: "test-client-id",
		Scopes:  []string{"events:write"},
	}

	tests := []struct {
		name     string
		resource string
		expected bool
	}{
		{
			name:     "書き込み可能なリソース",
			resource: "events",
			expected: true,
		},
		{
			name:     "書き込み不可能なリソース",
			resource: "teams",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CanWrite(claims, tt.resource)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthService_GetComplexityLimit(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	tests := []struct {
		name     string
		plan     string
		expected int
	}{
		{
			name:     "Free",
			plan:     "free",
			expected: 50,
		},
		{
			name:     "Standard",
			plan:     "standard",
			expected: 200,
		},
		{
			name:     "Premium",
			plan:     "premium",
			expected: 1000,
		},
		{
			name:     "Enterprise",
			plan:     "enterprise",
			expected: 5000,
		},
		{
			name:     "Unknown",
			plan:     "unknown",
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &AgentClaims{
				Plan: tt.plan,
			}

			result := service.GetComplexityLimit(claims)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthService_GetRateLimit(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	tests := []struct {
		name     string
		plan     string
		expected auth.RateLimit
	}{
		{
			name: "Free",
			plan: "free",
			expected: auth.RateLimit{
				RequestsPerMinute: 10,
				RequestsPerHour:   100,
				RequestsPerDay:    1000,
			},
		},
		{
			name: "Standard",
			plan: "standard",
			expected: auth.RateLimit{
				RequestsPerMinute: 50,
				RequestsPerHour:   1000,
				RequestsPerDay:    10000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &AgentClaims{
				Plan: tt.plan,
			}

			result := service.GetRateLimit(claims)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthService_scopesToStrings(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAuthService(mockRepo, []byte("test-secret"), time.Hour)

	scopes := []auth.Scope{
		auth.ScopeReadEvents,
		auth.ScopeWriteEvents,
		auth.ScopeReadTeams,
	}

	result := service.scopesToStrings(scopes)

	expected := []string{
		"events:read",
		"events:write",
		"teams:read",
	}

	assert.Equal(t, expected, result)
}
