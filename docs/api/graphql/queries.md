# GraphQL クエリ API

## 概要

SP APIはGraphQLを使用してサッカー関連のデータを提供します。GraphQLの強力な型システムと効率的なデータ取得機能により、AIエージェントが必要なデータを1回のリクエストで取得できます。

## エンドポイント

**URL**: `POST /query`  
**Content-Type**: `application/json`  
**認証**: Bearer Token (OAuth 2.0)

## 基本的な使用方法

### リクエスト形式

```json
{
  "query": "GraphQLクエリ文字列",
  "variables": {
    "変数名": "変数値"
  },
  "operationName": "操作名（オプション）"
}
```

### レスポンス形式

```json
{
  "data": {
    // クエリの結果
  },
  "errors": [
    {
      "message": "エラーメッセージ",
      "locations": [
        {
          "line": 2,
          "column": 3
        }
      ],
      "path": ["field", "subfield"],
      "extensions": {
        "code": "ERROR_CODE",
        "details": "詳細情報"
      }
    }
  ]
}
```

## クエリ一覧

### イベント関連

#### イベント一覧の取得

```graphql
query GetEvents {
  events {
    id
    name
    nameJp
    participants {
      id
      name
      fifaTrigramme
      confederation
    }
  }
}
```

**レスポンス例**:
```json
{
  "data": {
    "events": [
      {
        "id": "1",
        "name": "World Cup 2024",
        "nameJp": "ワールドカップ 2024",
        "participants": [
          {
            "id": "1",
            "name": "Japan",
            "fifaTrigramme": "JPN",
            "confederation": "AFC"
          },
          {
            "id": "2",
            "name": "Brazil",
            "fifaTrigramme": "BRA",
            "confederation": "CONMEBOL"
          }
        ]
      }
    ]
  }
}
```

#### 特定イベントの取得

```graphql
query GetEvent($id: ID!) {
  event(id: $id) {
    id
    name
    nameJp
    matches {
      id
      kickoff
      category
      matchups {
        id
        name
        fifaTrigramme
      }
    }
  }
}
```

**変数**:
```json
{
  "id": "1"
}
```

### チーム関連

#### チーム一覧の取得

```graphql
query GetTeams {
  teams {
    id
    name
    nameJp
    fifaTrigramme
    confederation
    members {
      id
      name
      number
      position
    }
  }
}
```

#### 特定チームの取得

```graphql
query GetTeam($fifa: String!) {
  team(fifa: $fifa) {
    id
    name
    nameJp
    fifaTrigramme
    confederation
    coach {
      id
      name
      birthday
    }
    members {
      id
      name
      number
      position
      birthday
    }
  }
}
```

**変数**:
```json
{
  "fifa": "JPN"
}
```

### 試合関連

#### 試合一覧の取得

```graphql
query GetMatches($eventId: ID!) {
  event(id: $eventId) {
    matches {
      id
      kickoff
      category
      stadium {
        id
        name
        venue
      }
      matchups {
        id
        name
        fifaTrigramme
      }
      timelines {
        id
        type
        time
        player {
          id
          name
          number
        }
      }
    }
  }
}
```

### 選手関連

#### 選手情報の取得

```graphql
query GetPlayer($id: ID!) {
  player(id: $id) {
    id
    name
    number
    position
    birthday
    team {
      id
      name
      fifaTrigramme
    }
  }
}
```

## 変数の使用

### 基本的な変数

```graphql
query GetEventsWithLimit($limit: Int, $offset: Int) {
  events(limit: $limit, offset: $offset) {
    id
    name
  }
}
```

**変数**:
```json
{
  "limit": 10,
  "offset": 0
}
```

### 複雑な変数

```graphql
query GetMatchesByDateRange($startDate: DateTime!, $endDate: DateTime!) {
  matches(startDate: $startDate, endDate: $endDate) {
    id
    kickoff
    matchups {
      name
    }
  }
}
```

**変数**:
```json
{
  "startDate": "2024-06-01T00:00:00Z",
  "endDate": "2024-06-30T23:59:59Z"
}
```

## フラグメントの使用

### フラグメントの定義

```graphql
fragment TeamInfo on Team {
  id
  name
  nameJp
  fifaTrigramme
  confederation
}

fragment PlayerInfo on Player {
  id
  name
  number
  position
}
```

### フラグメントの使用

```graphql
query GetEventWithTeams {
  events {
    id
    name
    participants {
      ...TeamInfo
      members {
        ...PlayerInfo
      }
    }
  }
}
```

## インライン・フラグメント

### 型による分岐

```graphql
query GetEventDetails($id: ID!) {
  event(id: $id) {
    id
    name
    participants {
      ... on Team {
        id
        name
        fifaTrigramme
      }
    }
  }
}
```

## エイリアスの使用

### 同じフィールドの複数取得

```graphql
query GetTeamComparison($team1: String!, $team2: String!) {
  team1: team(fifa: $team1) {
    id
    name
    fifaTrigramme
  }
  team2: team(fifa: $team2) {
    id
    name
    fifaTrigramme
  }
}
```

**変数**:
```json
{
  "team1": "JPN",
  "team2": "BRA"
}
```

## ディレクティブ

### @include ディレクティブ

```graphql
query GetEvent($id: ID!, $includeMatches: Boolean!) {
  event(id: $id) {
    id
    name
    matches @include(if: $includeMatches) {
      id
      kickoff
    }
  }
}
```

### @skip ディレクティブ

```graphql
query GetTeam($fifa: String!, $skipPlayers: Boolean!) {
  team(fifa: $fifa) {
    id
    name
    members @skip(if: $skipPlayers) {
      id
      name
    }
  }
}
```

## パフォーマンス最適化

### 必要なフィールドのみ取得

```graphql
# 良い例：必要なフィールドのみ
query GetEventNames {
  events {
    id
    name
  }
}

# 悪い例：不要なフィールドも取得
query GetEventNames {
  events {
    id
    name
    nameJp
    participants {
      id
      name
      members {
        id
        name
        number
        position
        birthday
      }
    }
  }
}
```

### 適切な深さの制限

```graphql
# 適切な深さ
query GetEventSummary {
  events {
    id
    name
    participants {
      id
      name
    }
  }
}

# 深すぎるクエリ（避ける）
query GetDeepEvent {
  events {
    participants {
      members {
        team {
          events {
            participants {
              members {
                # さらに深く...
              }
            }
          }
        }
      }
    }
  }
}
```

## エラーハンドリング

### 一般的なエラー

```json
{
  "errors": [
    {
      "message": "Field 'invalidField' doesn't exist on type 'Event'",
      "locations": [
        {
          "line": 3,
          "column": 5
        }
      ],
      "path": ["events", "invalidField"],
      "extensions": {
        "code": "GRAPHQL_VALIDATION_FAILED"
      }
    }
  ]
}
```

### 認証エラー

```json
{
  "errors": [
    {
      "message": "Authentication required",
      "extensions": {
        "code": "UNAUTHORIZED"
      }
    }
  ]
}
```

### 権限エラー

```json
{
  "errors": [
    {
      "message": "Insufficient permissions",
      "extensions": {
        "code": "INSUFFICIENT_SCOPES",
        "details": {
          "required_scopes": ["events:read"]
        }
      }
    }
  ]
}
```

## レート制限

### クエリコスト制限

各プランにはクエリコスト制限があります：

| プラン | 制限 |
|-------|-----|
| Free | 50 |
| Standard | 200 |
| Premium | 1000 |
| Enterprise | 5000 |

### コスト計算例

```graphql
# コスト: 5 (events) + 10 (participants) + 5 (members) = 20
query GetEventWithTeams {
  events {
    id
    name
    participants {
      id
      name
      members {
        id
        name
      }
    }
  }
}
```

## 実装例

### Go言語での実装

```go
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type GraphQLRequest struct {
    Query     string                 `json:"query"`
    Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
    Data   interface{} `json:"data"`
    Errors []struct {
        Message   string `json:"message"`
        Locations []struct {
            Line   int `json:"line"`
            Column int `json:"column"`
        } `json:"locations"`
        Extensions map[string]interface{} `json:"extensions"`
    } `json:"errors"`
}

func (c *APIClient) ExecuteQuery(ctx context.Context, query string, variables map[string]interface{}) (*GraphQLResponse, error) {
    reqBody := GraphQLRequest{
        Query:     query,
        Variables: variables,
    }
    
    jsonBody, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/query", bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Authorization", "Bearer "+c.accessToken)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var graphqlResp GraphQLResponse
    if err := json.Unmarshal(body, &graphqlResp); err != nil {
        return nil, err
    }
    
    if len(graphqlResp.Errors) > 0 {
        return &graphqlResp, fmt.Errorf("GraphQL errors: %v", graphqlResp.Errors)
    }
    
    return &graphqlResp, nil
}

// 使用例
func main() {
    client := NewAPIClient("https://api.example.com", "client-id", "client-secret")
    
    query := `
        query GetEvents {
            events {
                id
                name
                participants {
                    id
                    name
                    fifaTrigramme
                }
            }
        }
    `
    
    resp, err := client.ExecuteQuery(context.Background(), query, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Response: %+v\n", resp.Data)
}
```

### Pythonでの実装

```python
import requests
import json
from typing import Dict, Any, Optional

class GraphQLClient:
    def __init__(self, base_url: str, access_token: str):
        self.base_url = base_url
        self.access_token = access_token
        self.session = requests.Session()
    
    def execute_query(self, query: str, variables: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """GraphQLクエリを実行"""
        headers = {
            'Authorization': f'Bearer {self.access_token}',
            'Content-Type': 'application/json'
        }
        
        payload = {
            'query': query,
            'variables': variables or {}
        }
        
        response = self.session.post(
            f"{self.base_url}/query",
            json=payload,
            headers=headers
        )
        
        if response.status_code != 200:
            raise Exception(f"HTTP error: {response.status_code}")
        
        result = response.json()
        
        if 'errors' in result and result['errors']:
            raise Exception(f"GraphQL errors: {result['errors']}")
        
        return result

# 使用例
if __name__ == "__main__":
    client = GraphQLClient("https://api.example.com", "your-access-token")
    
    query = """
        query GetEvents {
            events {
                id
                name
                participants {
                    id
                    name
                    fifaTrigramme
                }
            }
        }
    """
    
    try:
        result = client.execute_query(query)
        print(json.dumps(result['data'], indent=2))
    except Exception as e:
        print(f"Error: {e}")
```

## トラブルシューティング

### よくある問題

**問題**: `Field doesn't exist` エラー  
**原因**: 存在しないフィールドを指定している  
**解決法**: GraphQLスキーマを確認し、正しいフィールド名を使用してください

**問題**: `Query complexity exceeds limit` エラー  
**原因**: クエリの複雑さがプランの制限を超えている  
**解決法**: クエリを簡素化するか、より高いプランにアップグレードしてください

**問題**: レスポンスが遅い  
**原因**: 不要なフィールドを取得している、または深すぎるクエリ  
**解決法**: 必要なフィールドのみを取得し、クエリの深さを制限してください

## 更新履歴

| バージョン | 日付 | 変更内容 |
|-----------|-----|---------|
| 1.0.0 | 2024-12-19 | 初版 |

---

**作成日**: 2024年12月19日  
**バージョン**: 1.0  
**承認者**: [承認者名]

