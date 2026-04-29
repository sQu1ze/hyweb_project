# Golang RESTful API 專案 SPEC（修正版）

本文件用於指導 AI coding agent 與開發者完成本專案，範圍僅包含題目中的必做項目，不包含加分題的天氣 API 與排程功能。

## 專案目標

使用 Golang 搭配 Gin 建立一個 RESTful API 後端服務，串接 MySQL，完成使用者註冊、登入、修改密碼三個核心功能，並提供統一格式 API 回應、Docker 啟動方式、README、資料表建立方式，以及可在本機呼叫本機服務的 Swagger API 文件。

## 範圍界定

### 必做功能

- 統一 API 回應格式。
- 使用者註冊，密碼需使用 bcrypt 雜湊後儲存。
- 使用者登入，登入成功需回傳 JWT，payload 必須包含 `email` 與 `updated` 欄位。
- 修改密碼，需驗證 JWT，並驗證舊密碼後更新新密碼。
- User 資料表設計。
- README、Dockerfile、docker-compose、資料表建立語法或 migrate 指令、API 文件。
- 提供 GitHub Public repo。

### 明確不做

- 不實作加分項目的中央氣象署 API。
- 不建立 Weather 資料表。
- 不實作天氣查詢 API。
- 不實作 JWT 黑名單或登出失效機制，因題目已明示修改密碼後原 Token 仍有效。

## 技術限制

- Language: Go 1.22+，且實際 Docker builder image 必須與 `go.mod` 要求版本一致或更高。
- Framework: Gin。
- Database: MySQL 8。
- Password Hash: bcrypt。
- Auth: JWT。
- Container: Docker + Docker Compose。
- API doc: 必須使用 Swagger UI，不採用 Postman collection 作為主要交付。

## 系統需求

### 功能清單

1. 註冊 API
2. 登入 API
3. 修改密碼 API
4. Health Check API（自定義補充，方便 Docker 與驗收）
5. Swagger 文件路由，可直接在本機瀏覽器測試

### 統一回應格式

所有 API 都必須回傳以下結構：

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

錯誤時：

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

此格式為題目強制要求，不可自行更改欄位名稱。

## API 設計

### 1. 註冊

- Method: `POST`
- Path: `/api/v1/users/register`
- Auth: 不需要 JWT。

Request body:

```json
{
  "email": "user@example.com",
  "password": "P@ssw0rd123"
}
```

成功行為：
- 驗證 email 與 password 格式。
- 檢查 email 是否已存在。
- 使用 bcrypt 進行密碼雜湊。
- 新增使用者資料。
- 回傳成功訊息與必要欄位，不回傳密碼。

失敗行為：
- email 已存在時，回傳明確錯誤訊息。
- 參數格式錯誤時，回傳 400。

### 2. 登入

- Method: `POST`
- Path: `/api/v1/users/login`
- Auth: 不需要 JWT。

Request body:

```json
{
  "email": "user@example.com",
  "password": "P@ssw0rd123"
}
```

成功行為：
- 驗證 email 是否存在。
- 使用 bcrypt 驗證密碼。
- 成功後簽發 JWT。
- JWT payload 至少包含：`email`、`updated`。

失敗行為：
- 帳號或密碼錯誤時，統一回傳 `invalid email or password` 類似訊息，不可揭露到底哪一個錯誤。

### 3. 修改密碼

- Method: `PUT`
- Path: `/api/v1/users/password`
- Auth: 必須帶 `Authorization: Bearer <token>`。

Request body:

```json
{
  "old_password": "OldP@ssw0rd123",
  "new_password": "NewP@ssw0rd123"
}
```

成功行為：
- 驗證 JWT。
- 從 JWT 取得 email。
- 查詢使用者。
- 驗證舊密碼。
- bcrypt 雜湊新密碼後更新資料。
- 同步更新 `updated` 欄位。
- 回傳成功訊息。

注意：
- 舊 JWT 在修改密碼後仍然有效，不需要實作 token revoke 或 blacklist。

### 4. Health Check

- Method: `GET`
- Path: `/api/v1/health`
- Auth: 不需要。
- 用途：確認服務可用，並方便 Docker 與驗收。

## 資料表設計

### User Table

```sql
CREATE TABLE IF NOT EXISTS users (
  email VARCHAR(50) NOT NULL,
  password VARCHAR(255) NOT NULL,
  created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (email)
);
```

此結構依題目給定欄位設計，`email` 是唯一識別鍵，`password` 儲存 bcrypt hash，並保留 `created`、`updated` 時間欄位。

## 專案結構（強制統一）

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

重要：
- `go.mod` 必須位於專案根目錄。
- `main.go` 統一位於 `cmd/server/main.go`。
- Dockerfile、Swagger 產生指令、README、匯入路徑都必須以此結構為準，不可混用根目錄 `main.go` 與 `cmd/server/main.go`。

## 分層責任

### handler
- 負責接收 HTTP request。
- 驗證 request body。
- 呼叫 service。
- 將結果轉為統一回應格式。
- 必須包含 Swagger operation 註解。

### service
- 實作核心商業邏輯。
- 註冊、登入、改密碼流程都集中在這裡。
- 決定錯誤訊息與流程控制。

### repository
- 負責與 MySQL 溝通。
- 提供查詢使用者、建立使用者、更新密碼等方法。

### middleware
- JWT 驗證。
- 將 token payload 寫入 context 供 handler 或 service 使用。

### response
- 封裝成功與失敗回應格式，避免每個 handler 重複寫 response JSON。

### router
- 註冊所有 API 路由。
- 註冊 `/swagger/*any` 路由。
- 不可將 Swagger operation 註解錯放在 router 註冊處。

## JWT 設計

### Payload

JWT payload 至少包含：

```json
{
  "email": "user@example.com",
  "updated": "2025-09-12T07:58:43Z",
  "exp": 1757663923
}
```

### 規則

- Secret 透過環境變數提供。
- Token 有有效期限，例如 24 小時。
- 驗證失敗時回傳 401。
- 修改密碼時只需驗證 token 可用，不需比對 token 內 `updated` 與 DB 是否一致，因題目已明示原 token 仍有效。

## 驗證規則

### email
- 必填。
- 最大長度 50。
- 必須符合基本 email 格式。

### password
- 必填。
- 建議最少 8 碼。
- 註冊與新密碼都套用相同規則。

### old_password / new_password
- 必填。
- `new_password` 不可與 `old_password` 完全相同。

## 錯誤處理規則

### 建議錯誤碼
- 200：成功
- 201：建立成功
- 400：參數錯誤
- 401：未授權或 JWT 無效
- 409：email 已存在
- 500：伺服器內部錯誤

### 錯誤訊息原則

- 註冊重複 email：可明確提示 email 已存在。
- 登入失敗：只能回覆帳號或密碼錯誤，不可透露具體哪個欄位錯。
- 修改密碼失敗：若舊密碼錯誤，可直接提示 `old password incorrect`。

## 環境變數（強制一致）

建立 `.env.example`：

```env
APP_PORT=8080
APP_ENV=local
MYSQL_HOST=mysql-db
MYSQL_PORT=3306
MYSQL_DATABASE=hyweb_db
MYSQL_USER=api_user
MYSQL_PASSWORD=api_password
JWT_SECRET=change_me_secret_key
JWT_EXPIRE_HOURS=24
```

強制規則：
- `.env.example` 與 `docker-compose.yml` 必須使用同一套命名，統一使用 `MYSQL_*`。
- 不可混用 `DB_HOST` / `DB_NAME` 與 `MYSQL_HOST` / `MYSQL_DATABASE`。
- 若 app 執行於 Docker Compose 網路中，`MYSQL_HOST` 必須為 `mysql-db`，不可為 `127.0.0.1` 或 `localhost`。
- 若 app 在宿主機本機直接執行，才可另外準備本機專用 `.env`。

## Docker 要求（強制規範）

### Dockerfile

- 必須提供可執行的 multi-stage Dockerfile。
- builder image 的 Go 版本不得低於 `go.mod` 要求。
- 若使用 Alpine-based Go image，方可使用 `apk add --no-cache git`。
- `go build` 必須指定正確入口：

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server
```

- 不可假設專案根目錄有可直接 build 的 Go 檔。

### docker-compose

必須包含：
- app container
- mysql container

app service 必須明確指定：

```yaml
build:
  context: .
  dockerfile: Dockerfile
```

MySQL service 必須包含：

```yaml
volumes:
  - mysql_data:/var/lib/mysql
  - ./migrations:/docker-entrypoint-initdb.d:ro
```

MySQL healthcheck 建議：

```yaml
healthcheck:
  test: ["CMD-SHELL", "mysqladmin ping -h 127.0.0.1 -uapi_user -papi_password --silent"]
  interval: 10s
  timeout: 5s
  retries: 10
  start_period: 20s
```

app service 必須包含：

```yaml
depends_on:
  mysql-db:
    condition: service_healthy
```

並可使用：

```yaml
env_file:
  - .env
```

注意：
- 若 `environment:` 與 `env_file:` 使用相同 key，需清楚說明覆蓋關係。
- 不可讓 `.env` 中的 `MYSQL_HOST=127.0.0.1` 與 Compose 中的 Docker 網路配置衝突。

### 啟動指令

```bash
docker compose up --build
```

## Migration 強制規範

- migration 檔案必須放在 `migrations/` 目錄。
- 必須至少提供一個初始化 SQL，例如：

```text
migrations/000001_init_users.sql
```

- migration 除了建立 `users` table 外，還必須插入測試資料。
- 至少需包含以下測試帳號：
  - `admin@example.com`
  - `user1@example.com`
  - `user2@example.com`
- 測試資料密碼必須使用 bcrypt hash，不可存明文。

建議 migration 內容：

```sql
CREATE TABLE IF NOT EXISTS users (
  email VARCHAR(50) NOT NULL,
  password VARCHAR(255) NOT NULL,
  created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (email)
);

INSERT INTO users (email, password, created, updated) VALUES
('admin@example.com', '$2a$10$fy9n/DpSqGx8MLoH7WWeUOYVMrgkVDiYHUg9pp3WSi5hxdyiaVJuq', NOW(), NOW()),
('user1@example.com', '$2a$10$fy9n/DpSqGx8MLoH7WWeUOYVMrgkVDiYHUg9pp3WSi5hxdyiaVJuq', NOW(), NOW()),
('user2@example.com', '$2a$10$fy9n/DpSqGx8MLoH7WWeUOYVMrgkVDiYHUg9pp3WSi5hxdyiaVJuq', NOW(), NOW());
```

重要說明：
- `/docker-entrypoint-initdb.d` 內的 SQL 只會在 MySQL data directory 第一次初始化時執行。
- 若已存在 named volume，需使用以下指令清除後重新建立：

```bash
docker compose down -v
docker compose up --build
```

## Swagger 強制規範

- 必須使用 Swagger UI。
- Swagger UI 必須可透過以下網址開啟：

```text
http://localhost:8080/swagger/index.html
```

- Gin router 必須註冊：

```go
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

- `cmd/server/main.go` 必須包含 Swagger 全域註解：
  - `@title`
  - `@version`
  - `@description`
  - `@host`
  - `@BasePath`
  - `@schemes`
  - `@securityDefinitions.apikey`

- Swagger 全域註解必須放在 `cmd/server/main.go`，不可放在根目錄 `main.go`、`router.go` 或其他檔案冒充入口。
- `swag init` 的 `-g` 參數必須使用 `cmd/server/main.go`。
- README、Dockerfile、Makefile（若有）、CI 指令（若有）中的 Swagger 指令都必須一致。
- 若實作者搬動目錄結構，必須同步更新 README 與 Swagger 產生指令，不能只改程式不改文件。

- 每個 handler function 都必須包含完整 Swagger operation 註解，至少包含：
  - `@Summary`
  - `@Description`
  - `@Tags`
  - `@Accept`
  - `@Produce`
  - `@Param`
  - `@Success`
  - `@Failure`
  - `@Router`

- 不可只在 `router.go` 上方填寫 Swagger 全域資訊後就結束。
- 產生出的 `swagger.json` 不可為空 spec，`paths` 不可為 `{}`。
- 必須能在 Swagger UI 中看到以下 operations：
  - `/health`
  - `/users/register`
  - `/users/login`
  - `/users/password`

Swagger 產生指令統一為：

```bash
swag init -g cmd/server/main.go -o docs/swagger --parseInternal
```

若使用跨 package model，也可視需要加上：

```bash
swag init -g cmd/server/main.go -o docs/swagger --parseInternal --parseDependency
```

## go.mod / import path 強制規範

- `go.mod` 必須存在於專案根目錄。
- module path 必須與程式中的 import path 完全一致。
- 所有 import 必須以 `go.mod` module 名稱為準，例如：

```go
import "hyweb-api/internal/router"
```

- 不可遺漏 `go.mod`、`go.sum`。
- 不可讓 SPEC、Dockerfile、README、Swagger 指令、實際程式入口位置彼此不一致。

## README 必須包含

- 專案簡介
- 技術棧
- 專案結構
- 環境變數說明
- 本機啟動方式
- Docker 啟動方式
- DB migration 使用方式
- `docker compose down -v` 何時需要使用
- Swagger 使用方式
- 測試 API 的 curl 範例
- 題目對應功能清單

## README 與 Swagger 一致性規範

### README 必須明確反映實際專案結構

README 內容必須與實際程式入口、Swagger 產生方式、Docker 啟動方式完全一致，特別是本專案入口已固定為 `cmd/server/main.go`，不可再使用根目錄 `main.go` 的描述。

README 至少必須明確寫出：

- 專案入口為：

```text
cmd/server/main.go
```

- 本機執行指令範例：

```bash
go run ./cmd/server
```

- 編譯指令範例：

```bash
go build -o bin/server ./cmd/server
```

- Swagger 文件產生指令：

```bash
swag init -g cmd/server/main.go -o docs/swagger --parseInternal
```

- 若有跨 package model 需要補充：

```bash
swag init -g cmd/server/main.go -o docs/swagger --parseInternal --parseDependency
```

- Dockerfile 的 build target 必須對應：

```bash
go build -o main ./cmd/server
```

- README 中的專案結構圖，必須顯示 `cmd/server/main.go`，不可再出現根目錄 `main.go`。

### README 必須包含的操作範例（補充）

README 除了列出章節名稱，還必須提供可直接操作的命令範例，至少包含：

1. 本機啟動
2. Docker 啟動
3. Swagger 產生
4. migration 初始化方式
5. `docker compose down -v` 重新初始化資料庫時機
6. curl 測試 register / login / change password / health

## Git / Repo 強制規範

- 必須提供 GitHub Public repo。
- 專案根目錄必須包含 `.gitignore`。
- `.gitignore` 至少需忽略：
  - `.env`
  - `bin/`
  - `dist/`
  - `tmp/`
  - `*.log`
  - `coverage.out`
  - IDE 與系統暫存檔
- README 必須說明如何從乾淨環境 clone、安裝、啟動與測試。
- 若有需要，應提供初始化 git repository 的標準流程，但不得影響專案目錄結構。
- 不可將敏感資訊提交到 repo，例如 `.env`、DB 密碼、JWT secret。

## API 文件方案

本專案必須使用 Swagger UI，因為需可直接在本機透過瀏覽器與本機服務互動，符合題目要求 API 文件必須可在本機發送 request 到本機服務。

## 驗收清單

### 功能驗收

- 可以成功註冊新使用者。
- 重複 email 註冊會失敗。
- 可以成功登入並拿到 JWT。
- 錯誤登入不會洩漏是帳號錯還是密碼錯。
- 帶 JWT 可成功修改密碼。
- 舊密碼錯誤會被拒絕。
- 修改後可用新密碼登入。
- 舊 token 仍可通過需要 JWT 的 API 驗證。

### 交付驗收

- 有 Public GitHub repo。
- 有 README。
- 有 Dockerfile。
- 有 docker-compose。
- 有 schema SQL 或 migration。
- 有 Swagger API 文件。
- docs 中包含 SPEC 與 AI 對話紀錄。

### 啟動與配置驗收

- `docker compose up --build` 可成功啟動。
- app 可成功連線到 mysql。
- 容器內 DB host 不可使用 `127.0.0.1`。
- MySQL 可成功建立 `users` table。
- migration 可插入測試資料。
- Swagger UI 可成功開啟。
- Swagger UI 不可出現 `No operations defined in spec!`。
- `swagger.json` 的 `paths` 必須包含 register、login、change password、health。
- README 不可殘留 `main.go` 舊指令。
- Swagger 指令、Docker build path、實際程式入口必須一致。

## 開發順序

1. 建立專案骨架與 Gin router。
2. 建立 `go.mod` 並確認 module path。
3. 完成 config 與 env 載入。
4. 建立 DB 連線。
5. 建立 users table migration 與測試資料。
6. 完成統一 response helper。
7. 完成 register API。
8. 完成 login API 與 JWT。
9. 完成 auth middleware。
10. 完成 change password API。
11. 補 Swagger 註解與文件產生。
12. 補 Dockerfile 與 docker-compose。
13. 撰寫 README。
14. 整理 docs/SPEC.md 與 docs/AI_CHAT_LOG.md。

## 給 AI coding agent 的執行指令

請嚴格遵守以下規則：

- 只實作必做項目，不要實作加分題。
- 使用 Gin + MySQL。
- 所有 API 回應都要套用統一格式。
- 使用 bcrypt 處理密碼。
- login 成功時簽發 JWT，payload 必含 `email` 與 `updated`。
- login 失敗訊息不得洩漏帳號或密碼何者錯誤。
- change password 需要 Bearer Token。
- 不需要 token blacklist。
- 需產出 Dockerfile、docker-compose、README、migration、Swagger API 文件。
- Docker、Swagger、migration、env、go.mod、專案結構必須彼此一致，不可只產出看似完整但無法啟動的模板。
- 若新增額外功能，只能是為了驗收與可維護性，例如 health check、graceful shutdown、structured logging，不得偏離題目主軸。

## 建議依賴

- `github.com/gin-gonic/gin`
- `github.com/golang-jwt/jwt/v5`
- `golang.org/x/crypto/bcrypt`
- `github.com/go-sql-driver/mysql`
- `github.com/joho/godotenv`
- `github.com/swaggo/gin-swagger`
- `github.com/swaggo/files`
- `github.com/swaggo/swag`

## Done 定義

當以下條件全部成立，視為完成：

- `docker compose up --build` 可成功啟動。
- MySQL 成功建立 users table 並插入測試資料。
- 註冊、登入、修改密碼 API 均可用。
- 統一回應格式正確。
- `http://localhost:8080/swagger/index.html` 可打開且有可測試 operations。
- README 可讓審查者在乾淨環境重現專案。
- docs 中包含 SPEC 與 AI 對話紀錄。
