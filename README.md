# 이슈 관리 API (백엔드 개발자 채용 과제)

Go 언어로 구현된 RESTful API로, 이슈 관리 시스템을 제공합니다.

## 🚀 기능 요약

- 이슈 생성, 조회, 수정 기능
- 상태별 이슈 필터링
- 담당자 할당/변경 시 자동 상태 관리
- RESTful API 설계

## 📋 요구사항

- Go 1.22+
- 포트 8080 사용

## 🛠️ 설치 및 실행

```bash
# 저장소 클론 후
cd backed-assignment-aro

# 의존성 설치
go mod tidy

# 서버 실행
go run ./cmd/server
# 서버가 8080 포트에서 시작됩니다
```

## 🧑‍💻 사전 등록된 사용자

```json
[
  { "id": 1, "name": "김개발" },
  { "id": 2, "name": "이디자인" },
  { "id": 3, "name": "박기획" }
]
```

## 📚 API 명세

### 1. 이슈 생성 [POST] /issue

**요청 예시:**
```http
POST /issue
Content-Type: application/json

{
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생",
    "userId": 1
}
```

**성공 응답 (201):**
```json
{
    "id": 1,
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생",
    "status": "IN_PROGRESS",
    "user": { "id": 1, "name": "김개발" },
    "createdAt": "2025-07-10T11:30:00Z",
    "updatedAt": "2025-07-10T11:30:00Z"
}
```

### 2. 이슈 목록 조회 [GET] /issues

**쿼리 파라미터:**
- `status`: 필터링할 상태값 (선택사항)

**예시:**
```http
GET /issues?status=PENDING
```

**성공 응답 (200):**
```json
{
    "issues": [
        {
            "id": 1,
            "title": "버그 수정 필요",
            "description": "로그인 페이지에서 오류 발생",
            "status": "PENDING",
            "createdAt": "2025-07-10T11:30:00Z",
            "updatedAt": "2025-07-10T11:30:00Z"
        }
    ]
}
```

### 3. 이슈 상세 조회 [GET] /issue/{id}

**예시:**
```http
GET /issue/1
```

**성공 응답 (200):**
```json
{
    "id": 1,
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생",
    "status": "PENDING",
    "user": { "id": 1, "name": "김개발" },
    "createdAt": "2025-07-10T11:30:00Z",
    "updatedAt": "2025-07-10T11:30:00Z"
}
```

### 4. 이슈 수정 [PATCH] /issue/{id}

**요청 예시:**
```http
PATCH /issue/1
Content-Type: application/json

{
    "title": "로그인 버그 수정",
    "status": "IN_PROGRESS",
    "userId": 2
}
```

**성공 응답 (200):**
```json
{
    "id": 1,
    "title": "로그인 버그 수정",
    "description": "로그인 페이지에서 오류 발생",
    "status": "IN_PROGRESS",
    "user": { "id": 2, "name": "이디자인" },
    "createdAt": "2025-07-10T11:30:00Z",
    "updatedAt": "2025-07-10T11:35:00Z"
}
```

## 🚨 에러 응답

**형식:**
```json
{
    "error": "에러 메시지",
    "code": 400
}
```

**주요 에러 케이스:**
- 400: 잘못된 요청 (필수 파라미터 누락, 유효하지 않은 값 등)
- 404: 리소스를 찾을 수 없음
- 409: 충돌 (예: 완료된 이슈 수정 시도)

## 🧪 테스트 방법

### cURL 예시

```bash
# 1. 이슈 생성
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{"title":"버그 수정","description":"로그인 오류","userId":1}'

# 2. 이슈 목록 조회
curl "http://localhost:8080/issues?status=IN_PROGRESS"

# 3. 이슈 상세 조회
curl http://localhost:8080/issue/1

# 4. 이슈 수정
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{"status":"COMPLETED"}'
```

## 📂 프로젝트 구조

```
.
├── cmd
│   └── server          # 애플리케이션 진입점
├── internal
│   ├── handler         # HTTP 요청 처리
│   ├── model           # 도메인 모델
│   ├── repository      # 데이터 접근 계층
│   └── service         # 비즈니스 로직
├── go.mod
└── README.md
```

## 📝 구현 내용

- **RESTful API** 설계
- **비동기 안전**한 메모리 저장소 구현
- **상태 관리** 자동화
  - 담당자 할당 시 자동으로 `IN_PROGRESS` 상태 변경
  - 담당자 제거 시 자동으로 `PENDING` 상태 변경
  - `COMPLETED`/`CANCELLED` 상태에서는 수정 불가
- **에러 처리** 및 유효성 검증

## 🛠 기술 스택

- **언어**: Go 1.22+
- **웹 프레임워크**: Chi Router
- **의존성 관리**: Go Modules
- **코드 포맷팅**: gofmt

## 📜 라이선스

MIT
