# OAuth 2.0 Client Credentials 認証

## 概要

SP APIはOAuth 2.0 Client Credentials フローを使用してAIエージェントの認証を行います。この認証方式は、マシン間通信に適しており、AIエージェントがAPIにアクセスする際の標準的な認証方法です。

## 認証フロー

```mermaid
sequenceDiagram
    participant Agent as AI Agent
    participant API as SP API
    participant Auth as Auth Service
    
    Agent->>Auth: 1. Client Credentials Request
    Note over Agent,Auth: POST /oauth/token<br/>client_id, client_secret
    Auth->>Auth: 2. Validate Credentials
    Auth->>Agent: 3. Access Token (JWT)
    Agent->>API: 4. API Request with Token
    Note over Agent,API: Authorization: Bearer {token}
    API->>API: 5. Validate Token
    API->>Agent: 6. API Response
```

## トークン取得

### エンドポイント

**URL**: `POST /oauth/token`  
**Content-Type**: `application/x-www-form-urlencoded`

### リクエストパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|----|----|----- |
| grant_type | string | 必須 | `client_credentials` を指定 |
| client_id | string | 必須 | エージェントのクライアントID |
| client_secret | string | 必須 | エージェントのクライアントシークレット |
| scope | string | 任意 | アクセスしたいスコープ（スペース区切り） |

### リクエスト例

```bash
curl -X POST https://api.example.com/oauth/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=your-client-id&client_secret=your-client-secret&scope=events:read teams:read"
```

### レスポンス

**成功時 (200 OK)**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "scope": "events:read teams:read"
}
```

**エラー時 (400 Bad Request)**:
```json
{
  "error": "invalid_client",
  "error_description": "Client authentication failed"
}
```

## トークン使用

### 認証ヘッダー

取得したアクセストークンは、APIリクエストのAuthorizationヘッダーに含めて送信します。

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### リクエスト例

```bash
curl -X POST https://api.example.com/query \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { events { id name } }"
  }'
```

## スコープ

### 利用可能なスコープ

| スコープ | 説明 |
|---------|----- |
| `events:read` | イベント情報の読み取り |
| `events:write` | イベント情報の書き込み |
| `teams:read` | チーム情報の読み取り |
| `teams:write` | チーム情報の書き込み |
| `matches:read` | 試合情報の読み取り |
| `matches:write` | 試合情報の書き込み |
| `players:read` | 選手情報の読み取り |
| `players:write` | 選手情報の書き込み |
| `predictions:read` | 予想情報の読み取り |
| `predictions:write` | 予想情報の書き込み |
| `admin` | 管理者権限 |

### スコープの組み合わせ

複数のスコープを要求する場合は、スペースで区切って指定します。

```
scope=events:read teams:read matches:read
```

## エラーレスポンス

### 認証エラー

| エラーコード | HTTPステータス | 説明 |
|-------------|---------------|----- |
| `invalid_client` | 400 | クライアント認証に失敗 |
| `invalid_grant` | 400 | 認証情報が無効 |
| `invalid_scope` | 400 | 要求されたスコープが無効 |
| `unauthorized_client` | 400 | クライアントが認可されていない |

### トークンエラー

| エラーコード | HTTPステータス | 説明 |
|-------------|---------------|----- |
| `invalid_token` | 401 | トークンが無効または期限切れ |
| `insufficient_scope` | 403 | スコープが不足 |

## セキュリティ考慮事項

### 1. クライアントシークレットの管理

- クライアントシークレットは安全に保管してください
- 環境変数や設定ファイルに保存する場合は、適切な権限設定を行ってください
- ログにクライアントシークレットが出力されないよう注意してください

### 2. トークンの管理

- アクセストークンは安全に保管してください
- トークンの有効期限を適切に管理してください
- 不要になったトークンは破棄してください

### 3. HTTPSの使用

- 本番環境では必ずHTTPSを使用してください
- トークンの送信時は暗号化された通信を使用してください

## 実装例

### Go言語での実装例

```go
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type TokenResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int    `json:"expires_in"`
    Scope       string `json:"scope"`
}

type APIClient struct {
    baseURL      string
    clientID     string
    clientSecret string
    accessToken  string
    expiresAt    time.Time
    httpClient   *http.Client
}

func NewAPIClient(baseURL, clientID, clientSecret string) *APIClient {
    return &APIClient{
        baseURL:      baseURL,
        clientID:     clientID,
        clientSecret: clientSecret,
        httpClient:   &http.Client{Timeout: 30 * time.Second},
    }
}

func (c *APIClient) GetAccessToken(ctx context.Context) error {
    data := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s&scope=events:read teams:read",
        c.clientID, c.clientSecret)
    
    req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/oauth/token", bytes.NewBufferString(data))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("token request failed: %s", string(body))
    }
    
    var tokenResp TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
        return err
    }
    
    c.accessToken = tokenResp.AccessToken
    c.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
    
    return nil
}

func (c *APIClient) ensureValidToken(ctx context.Context) error {
    if c.accessToken == "" || time.Now().After(c.expiresAt.Add(-5*time.Minute)) {
        return c.GetAccessToken(ctx)
    }
    return nil
}

func (c *APIClient) MakeRequest(ctx context.Context, query string) (*http.Response, error) {
    if err := c.ensureValidToken(ctx); err != nil {
        return nil, err
    }
    
    reqBody := map[string]string{"query": query}
    jsonBody, _ := json.Marshal(reqBody)
    
    req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/query", bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Authorization", "Bearer "+c.accessToken)
    req.Header.Set("Content-Type", "application/json")
    
    return c.httpClient.Do(req)
}

// 使用例
func main() {
    client := NewAPIClient("https://api.example.com", "your-client-id", "your-client-secret")
    
    ctx := context.Background()
    resp, err := client.MakeRequest(ctx, "query { events { id name } }")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    // レスポンスの処理
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Pythonでの実装例

```python
import requests
import time
from typing import Optional

class APIClient:
    def __init__(self, base_url: str, client_id: str, client_secret: str):
        self.base_url = base_url
        self.client_id = client_id
        self.client_secret = client_secret
        self.access_token: Optional[str] = None
        self.expires_at: Optional[float] = None
        self.session = requests.Session()
    
    def get_access_token(self) -> None:
        """アクセストークンを取得"""
        data = {
            'grant_type': 'client_credentials',
            'client_id': self.client_id,
            'client_secret': self.client_secret,
            'scope': 'events:read teams:read'
        }
        
        response = self.session.post(
            f"{self.base_url}/oauth/token",
            data=data,
            headers={'Content-Type': 'application/x-www-form-urlencoded'}
        )
        
        if response.status_code != 200:
            raise Exception(f"Token request failed: {response.text}")
        
        token_data = response.json()
        self.access_token = token_data['access_token']
        self.expires_at = time.time() + token_data['expires_in']
    
    def ensure_valid_token(self) -> None:
        """有効なトークンを確保"""
        if not self.access_token or time.time() > (self.expires_at - 300):
            self.get_access_token()
    
    def make_request(self, query: str) -> requests.Response:
        """APIリクエストを実行"""
        self.ensure_valid_token()
        
        headers = {
            'Authorization': f'Bearer {self.access_token}',
            'Content-Type': 'application/json'
        }
        
        data = {'query': query}
        
        return self.session.post(
            f"{self.base_url}/query",
            json=data,
            headers=headers
        )

# 使用例
if __name__ == "__main__":
    client = APIClient(
        base_url="https://api.example.com",
        client_id="your-client-id",
        client_secret="your-client-secret"
    )
    
    response = client.make_request("query { events { id name } }")
    print(response.json())
```

## トラブルシューティング

### よくある問題

**問題**: `invalid_client` エラーが発生する  
**原因**: クライアントIDまたはクライアントシークレットが間違っている  
**解決法**: 正しい認証情報を確認してください

**問題**: `invalid_scope` エラーが発生する  
**原因**: 要求されたスコープが無効または許可されていない  
**解決法**: 利用可能なスコープを確認し、適切なスコープを要求してください

**問題**: トークンがすぐに期限切れになる  
**原因**: システム時刻のずれまたはトークンの有効期限が短い  
**解決法**: システム時刻を確認し、トークンの自動更新を実装してください

## 更新履歴

| バージョン | 日付 | 変更内容 |
|-----------|-----|---------|
| 1.0.0 | 2024-12-19 | 初版 |

---

**作成日**: 2024年12月19日  
**バージョン**: 1.0  
**承認者**: [承認者名]
