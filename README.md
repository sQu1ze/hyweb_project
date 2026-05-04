# Golang RESTful API

此專案使用 Go、Gin、MySQL 實作 RESTful API，提供使用者註冊、登入、修改密碼與 health check，並整合 Swagger UI 供本機直接測試 API。[file:1]

## 技術棧

- Go 1.22+[file:1]
- Gin[file:1]
- MySQL 8[file:1]
- JWT[file:1]
- bcrypt[file:1]
- Docker / Docker Compose[file:1]
- Swagger UI[file:1]

## 題目對應功能

- 統一 API 回應格式。[file:1]
- 使用者註冊，密碼以 bcrypt hash 儲存。[file:1]
- 使用者登入，成功後回傳 JWT，payload 至少包含 `email` 與 `updated`。[file:1]
- 使用者修改密碼，需帶 Bearer Token，並驗證舊密碼。[file:1]
- Health Check API。[file:1]
- Docker 啟動方式、migration、Swagger 文件。[file:1]

## 專案結構

```text
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── response/
│   ├── router/
│   └── service/
├── docs/
│   ├── swagger/
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   ├── SPEC.md
│   └── AI_CHAT_LOG.md
├── migrations/
│   └── 000001_init_users.sql
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

重要：本專案唯一入口為 `cmd/server/main.go`，README、Dockerfile、Swagger 指令都必須以此結構為準，不可混用根目錄 `main.go`。[file:1]

## 環境變數

請先建立 `.env`，可直接由 `.env.example` 複製：[file:1]

```bash
cp .env.example .env
```

`.env.example` 內容如下：[file:1]

```env
APP_PORT=8080
APP_ENV=local
MYSQL_HOST=mysql-db
MYSQL_PORT=3306
MYSQL_DATABASE=hyweb_db
MYSQL_USER=api_user
MYSQL_PASSWORD=api_password
JWT_SECRET=QmFhU1NlY3JldEtleUluTWFuZUEzMjhCaXRzX1NhZmVGb3JfSFMjMjU2
JWT_EXPIRE_HOURS=24
```

說明：若透過 Docker Compose 啟動，`MYSQL_HOST` 必須使用 `mysql-db`；若是在宿主機本機直接執行，才可改成本機資料庫位址。[file:1]

## 本機啟動

先確認本機已準備 Go 1.22+、MySQL 8，並建立對應資料庫。[file:1]

1. 安裝依賴

```bash
go mod tidy
```

2. 啟動程式

```bash
go run ./cmd/server
```

3. 或先編譯再執行

```bash
go build -o bin/server ./cmd/server
./bin/server
```

本機啟動後，預期服務網址為 `http://localhost:8080`。[file:1]

## Docker 啟動

本專案使用 app + mysql 兩個服務，並透過 Docker Compose 啟動。[file:1]

```bash
docker compose up --build
```

若要背景執行：

```bash
docker compose up --build -d
```

停止服務：

```bash
docker compose down
```

## DB migration 使用方式

migration 檔案必須放在 `migrations/`，其中至少包含 `000001_init_users.sql`，用於建立 `users` table 與插入測試資料。[file:1]

若使用 Docker Compose，MySQL 容器會自動讀取 `./migrations` 掛載到 `/docker-entrypoint-initdb.d` 的 SQL 檔進行初始化。[file:1]

建議至少插入以下測試帳號，且密碼必須為 bcrypt hash，不可存明文：[file:1]

- `admin@example.com`[file:1]
- `user1@example.com`[file:1]
- `user2@example.com`[file:1]

## 何時需要 docker compose down -v

MySQL 只會在資料目錄第一次初始化時執行 `/docker-entrypoint-initdb.d` 內的 SQL。[file:1]

當你遇到以下情況時，需要清除 volume 後重新建立：[file:1]
- 修改了 migration SQL，但容器資料已存在。[file:1]
- 想重新初始化資料表與測試資料。[file:1]
- 發現舊資料導致驗證結果不一致。[file:1]

指令如下：

```bash
docker compose down -v
docker compose up --build
```

## Swagger 使用方式

Swagger UI 必須可透過下列網址開啟：[file:1]

```text
http://localhost:8080/swagger/index.html
```

Swagger 文件產生指令固定如下：[file:1]

```bash
swag init -g cmd/server/main.go -o docs/swagger --parseInternal
```

若有跨 package model 解析需求，可改用：[file:1]

```bash
swag init -g cmd/server/main.go -o docs/swagger --parseInternal --parseDependency
```

注意事項：[file:1]
- Swagger 全域註解必須放在 `cmd/server/main.go`。[file:1]
- 不可使用 `swag init -g main.go`。[file:1]
- Router 必須註冊 `/swagger/*any`。[file:1]
- 若 Swagger UI 顯示 `No operations defined in spec!`，請先確認 handler 註解與 `swag init` 入口是否正確。[file:1]

## API 路由

| Method | Path | 說明 |
|---|---|---|
| GET | `/api/v1/health` | Health Check。[file:1] |
| POST | `/api/v1/users/register` | 使用者註冊。[file:1] |
| POST | `/api/v1/users/login` | 使用者登入並取得 JWT。[file:1] |
| PUT | `/api/v1/users/password` | 修改密碼，需帶 Bearer Token。[file:1] |

## curl 測試範例

### 1. Health Check

```bash
curl -X GET http://localhost:8080/api/v1/health
```

### 2. Register

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "P@ssw0rd123"
  }'
```

### 3. Login

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

登入成功後，請從回應中取得 JWT token。[file:1]

### 4. Change Password

```bash
curl -X PUT http://localhost:8080/api/v1/users/password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "old_password": "password123",
    "new_password": "NewP@ssw0rd123"
  }'
```

## 統一回應格式

成功時範例：[file:1]

```json
{
  "success": true,
  "message": "成功",
  "data": {},
  "error": null,
  "code": 200,
  "timestamp": "2025-09-12T07:58:43Z"
}
```

失敗時範例：[file:1]

```json
{
  "success": false,
  "message": "失敗原因",
  "data": null,
  "error": {},
  "code": 400,
  "timestamp": "2025-09-12T07:58:43Z"
}
```

## 開發與驗收重點

- `go.mod` 必須位於專案根目錄，且 import path 必須與 module path 完全一致。[file:1]
- Docker、Swagger、migration、env、README、實際入口位置必須彼此一致。[file:1]
- `cmd/server/main.go` 為唯一入口，不可殘留 `go run main.go` 或 `swag init -g main.go` 之類的舊寫法。[file:1]
- Swagger UI 中必須看得到 `/health`、`/users/register`、`/users/login`、`/users/password` 這些 operations。[file:1]
