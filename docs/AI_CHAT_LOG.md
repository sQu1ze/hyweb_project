# AI 對話紀錄

以下紀錄用來整理這個專案要如何與 AI 協作，先完成 SPEC，再依 SPEC 讓 Gemini CLI 產出專案，最後用 CLI 完成後的專案進行測試。

## 協作流程

### 1. 先按照需求完成 SPEC.md

先把題目需求整理成可執行的 SPEC.md，內容要明確定義：
- 功能範圍。
- 專案結構。
- API 路由。
- 資料表設計。
- 環境變數。
- Docker / Swagger / README 規範。
- 驗收條件。

這一步的目標是先讓規格固定下來，避免 AI 在開發時自由發揮，導致結構、路由或文件不一致。

### 2. Gemini CLI 按照 SPEC.md 實作專案

接著把 SPEC.md 丟給 Gemini CLI，要求它完全依照規格建立專案。

建議明確告訴 Gemini CLI：
- 請嚴格遵守 SPEC.md。
- 不要自行新增超出範圍的功能。
- 入口檔固定在 `cmd/server/main.go`。
- Swagger 入口固定使用 `swag init -g cmd/server/main.go`。
- README、Dockerfile、migration、env 命名都要與 SPEC 一致。

這一步的重點是讓 CLI 產出的專案可以直接啟動、文件可以直接產生、Swagger 可以正常打開。

### 3. 按照 CLI 完成的專案進行測試

專案產出後，依照 SPEC 的驗收條件做測試，確認實作是否真的符合需求。

測試時建議依序檢查：
- `docker compose up --build` 是否能正常啟動。
- MySQL 是否成功建立 users table。
- migration 是否成功載入測試資料。
- 註冊、登入、修改密碼 API 是否正常。
- JWT payload 是否包含 `email` 與 `updated`。
- Swagger UI 是否可開啟，且有正確 operations。
- README 指令是否與實際專案一致。

如果有不符合的地方，再回頭修改 SPEC 或要求 Gemini CLI 重做對應部分。

## 最終檢查清單

- SPEC.md 已完成並凍結。
- Gemini CLI 已依 SPEC 實作。
- 專案可成功啟動。
- Swagger 可正常使用。
- README 與實際專案一致。
- API 驗證與測試都符合 SPEC。
