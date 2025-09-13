package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAgent(t *testing.T) {
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		agentName    string
		plan         Plan
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "正常ケース",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			agentName:    "Test Agent",
			plan:         PlanFree,
			wantErr:      false,
		},
		{
			name:         "ClientIDが空",
			clientID:     "",
			clientSecret: "test-secret",
			agentName:    "Test Agent",
			plan:         PlanFree,
			wantErr:      true,
			errMsg:       "client ID cannot be empty",
		},
		{
			name:         "ClientSecretが空",
			clientID:     "test-client-id",
			clientSecret: "",
			agentName:    "Test Agent",
			plan:         PlanFree,
			wantErr:      true,
			errMsg:       "client secret cannot be empty",
		},
		{
			name:         "AgentNameが空",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			agentName:    "",
			plan:         PlanFree,
			wantErr:      true,
			errMsg:       "agent name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewAgent(tt.clientID, tt.clientSecret, tt.agentName, tt.plan)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, agent)
			} else {
				require.NoError(t, err)
				require.NotNil(t, agent)
				assert.Equal(t, tt.clientID, agent.ClientID())
				assert.Equal(t, tt.clientSecret, agent.ClientSecret())
				assert.Equal(t, tt.agentName, agent.Name())
				assert.Equal(t, tt.plan, agent.Plan())
				assert.True(t, agent.IsActive())
				assert.NotZero(t, agent.GetCreatedAt())
				assert.NotZero(t, agent.GetUpdatedAt())
			}
		})
	}
}

func TestAgent_AddScope(t *testing.T) {
	agent, err := NewAgent("test-id", "test-secret", "Test Agent", PlanFree)
	require.NoError(t, err)

	tests := []struct {
		name    string
		scope   Scope
		wantErr bool
	}{
		{
			name:    "正常ケース",
			scope:   ScopeReadEvents,
			wantErr: false,
		},
		{
			name:    "重複スコープ",
			scope:   ScopeReadEvents,
			wantErr: false, // 重複はエラーにならない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := agent.AddScope(tt.scope)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, agent.HasScope(tt.scope))
			}
		})
	}
}

func TestAgent_AddScope_InactiveAgent(t *testing.T) {
	agent, err := NewAgent("test-id", "test-secret", "Test Agent", PlanFree)
	require.NoError(t, err)

	// エージェントを非アクティブにする
	err = agent.Deactivate()
	require.NoError(t, err)

	// 非アクティブなエージェントにスコープを追加しようとする
	err = agent.AddScope(ScopeReadEvents)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot add scope to inactive agent")
}

func TestAgent_RemoveScope(t *testing.T) {
	agent, err := NewAgent("test-id", "test-secret", "Test Agent", PlanFree)
	require.NoError(t, err)

	// スコープを追加
	err = agent.AddScope(ScopeReadEvents)
	require.NoError(t, err)
	assert.True(t, agent.HasScope(ScopeReadEvents))

	// スコープを削除
	err = agent.RemoveScope(ScopeReadEvents)
	assert.NoError(t, err)
	assert.False(t, agent.HasScope(ScopeReadEvents))
}

func TestAgent_RemoveScope_NonExistentScope(t *testing.T) {
	agent, err := NewAgent("test-id", "test-secret", "Test Agent", PlanFree)
	require.NoError(t, err)

	// 存在しないスコープを削除しようとする
	err = agent.RemoveScope(ScopeReadEvents)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scope not found")
}

func TestAgent_Deactivate(t *testing.T) {
	agent, err := NewAgent("test-id", "test-secret", "Test Agent", PlanFree)
	require.NoError(t, err)
	assert.True(t, agent.IsActive())

	// エージェントを非アクティブにする
	err = agent.Deactivate()
	assert.NoError(t, err)
	assert.False(t, agent.IsActive())

	// 既に非アクティブなエージェントを非アクティブにしようとする
	err = agent.Deactivate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent is already inactive")
}

func TestAgent_Activate(t *testing.T) {
	agent, err := NewAgent("test-id", "test-secret", "Test Agent", PlanFree)
	require.NoError(t, err)

	// エージェントを非アクティブにする
	err = agent.Deactivate()
	require.NoError(t, err)
	assert.False(t, agent.IsActive())

	// エージェントをアクティブにする
	err = agent.Activate()
	assert.NoError(t, err)
	assert.True(t, agent.IsActive())

	// 既にアクティブなエージェントをアクティブにしようとする
	err = agent.Activate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent is already active")
}

func TestAgent_UpdateLastUsed(t *testing.T) {
	agent, err := NewAgent("test-id", "test-secret", "Test Agent", PlanFree)
	require.NoError(t, err)

	// 初期状態ではLastUsedAtはnil
	assert.Nil(t, agent.LastUsedAt())

	// UpdateLastUsedを呼び出す
	agent.UpdateLastUsed()

	// LastUsedAtが設定される
	assert.NotNil(t, agent.LastUsedAt())
	assert.True(t, time.Since(*agent.LastUsedAt()) < time.Second)
}

func TestAgent_CanAccess(t *testing.T) {
	agent, err := NewAgent("test-id", "test-secret", "Test Agent", PlanFree)
	require.NoError(t, err)

	// スコープを追加
	err = agent.AddScope(ScopeReadEvents)
	require.NoError(t, err)

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
			result := agent.CanAccess(tt.resource)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAgent_CanAccess_InactiveAgent(t *testing.T) {
	agent, err := NewAgent("test-id", "test-secret", "Test Agent", PlanFree)
	require.NoError(t, err)

	// スコープを追加
	err = agent.AddScope(ScopeReadEvents)
	require.NoError(t, err)

	// エージェントを非アクティブにする
	err = agent.Deactivate()
	require.NoError(t, err)

	// 非アクティブなエージェントはアクセスできない
	result := agent.CanAccess("events")
	assert.False(t, result)
}

func TestAgent_DomainEvents(t *testing.T) {
	agent, err := NewAgent("test-id", "test-secret", "Test Agent", PlanFree)
	require.NoError(t, err)

	// 初期状態ではAgentCreatedイベントが1つある
	events := agent.GetDomainEvents()
	assert.Len(t, events, 1)
	assert.Equal(t, "agent.created", events[0].EventType())

	// ドメインイベントをクリア
	agent.ClearDomainEvents()
	events = agent.GetDomainEvents()
	assert.Empty(t, events)

	// スコープを追加してドメインイベントを生成
	err = agent.AddScope(ScopeReadEvents)
	require.NoError(t, err)

	// ドメインイベントが追加される
	events = agent.GetDomainEvents()
	assert.Len(t, events, 1)
	assert.Equal(t, "agent.scope_added", events[0].EventType())

	// ドメインイベントをクリア
	agent.ClearDomainEvents()
	events = agent.GetDomainEvents()
	assert.Empty(t, events)
}
