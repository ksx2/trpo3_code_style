# Руководство по стилю кода на Go

## 1. Форматирование кода с использованием отступов и переносов строк

### Общие принципы

В Go форматирование кода стандартизировано инструментом `gofmt`. Всегда используйте `gofmt` или `go fmt` для автоматического форматирования кода перед коммитом.

### Отступы

- Используйте табы (TAB), а не пробелы для отступов
- Всегда используйте `gofmt`, который автоматически применяет правильные отступы

### Переносы строк

- Переносите длинные строки логически
- При переносе параметров функции выравнивайте их по первому параметру


**Пример правильного переноса длинной функции:**

```go
func SendEmail(
	to string,
	subject string,
	body string,
	attachments []string,
) error {
	// реализация
	return nil
}
```

**Пример правильного переноса длинного условия:**

```go
if user != nil &&
	user.IsActive() &&
	user.HasPermission("read") {
	// обработка
}
```

**Пример неправильного форматирования (слишком длинная строка):**

```go
func SendEmail(to string, subject string, body string, attachments []string) error {
	return nil
}
```

### Пустые строки

- Используйте пустые строки для разделения логических блоков
- Одна пустая строка между функциями
- Пустые строки внутри функций для разделения логических секций

**Пример использования пустых строк:**

```go
package main

import (
	"fmt"
	"log"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func NewUser(name, email string) *User {
	return &User{
		Name:  name,
		Email: email,
	}
}

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}

	if u.Email == "" {
		return fmt.Errorf("email is required")
	}

	return nil
}

func main() {
	user := NewUser("John Doe", "john@example.com")
	
	if err := user.Validate(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User created: %s\n", user.Name)
}
```

## 2. Разделение функционала на независимые модули или компоненты

### Принципы модульности

Разделение кода на независимые модули улучшает читаемость, тестируемость и переиспользование кода.

### Пакеты (Packages)

- Каждый пакет должен иметь одну четкую ответственность
- Именуйте пакеты коротко и понятно
- Избегайте общих имен типа `util`, `common`, `helpers`

**Пример структуры проекта:**

```
project/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── user/
│   │   ├── user.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── order/
│   │   ├── order.go
│   │   ├── repository.go
│   │   └── service.go
│   └── payment/
│       ├── payment.go
│       └── processor.go
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   └── validator/
│       └── validator.go
└── go.mod
```

**Пример определения типа пользователя:**

```go
package user

import "time"

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}
```

**Пример репозитория пользователя:**

```go
package user

type Repository interface {
	FindByID(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	Save(user *User) error
	Delete(id int) error
}

type repository struct {
	// зависимости для работы с БД
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) FindByID(id int) (*User, error) {
	// реализация
	return nil, nil
}

func (r *repository) FindByEmail(email string) (*User, error) {
	// реализация
	return nil, nil
}

func (r *repository) Save(user *User) error {
	// реализация
	return nil
}

func (r *repository) Delete(id int) error {
	// реализация
	return nil
}
```

**Пример сервисного слоя:**

```go
package user

import "errors"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}
	
	return s.repo.FindByID(id)
}

func (s *Service) CreateUser(name, email string) (*User, error) {
	user := &User{
		Name:  name,
		Email: email,
	}
	
	if err := s.repo.Save(user); err != nil {
		return nil, err
	}
	
	return user, nil
}
```

### Интерфейсы

- Определяйте интерфейсы в месте использования, а не в месте реализации
- Используйте маленькие, сфокусированные интерфейсы

**Пример правильного определения интерфейса (там, где используется):**

```go
package handler

import "project/internal/user"

type UserGetter interface {
	GetUser(id int) (*user.User, error)
}

func GetUserHandler(service UserGetter, id int) (*user.User, error) {
	return service.GetUser(id)
}
```

**Пример неправильного определения (слишком большой интерфейс):**

```go
type UserService interface {
	GetUser(id int) (*user.User, error)
	CreateUser(name, email string) (*user.User, error)
	UpdateUser(user *user.User) error
	DeleteUser(id int) error
	ValidateUser(user *user.User) error
	// ... еще 10 методов
}
```

### Зависимости

- Используйте dependency injection
- Передавайте зависимости через конструкторы
- Избегайте глобальных переменных

**Пример правильного использования dependency injection:**

```go
package main

import (
	"project/internal/user"
	"project/pkg/logger"
)

func main() {
	log := logger.New()
	repo := user.NewRepository()
	service := user.NewService(repo)
	
	// использование service
	_ = service
}
```

**Пример неправильного использования (глобальные переменные):**

```go
var globalRepo user.Repository

func init() {
	globalRepo = user.NewRepository()
}
```

---

## 3. Управление зависимостями и версиями библиотек

### Go Modules

Начиная с Go 1.11, используется система модулей (Go Modules) для управления зависимостями.

### Инициализация модуля

**Пример создания нового модуля:**

```bash
go mod init github.com/username/projectname
```

**Пример файла go.mod:**

```go
module github.com/username/projectname

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.14.0
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	// ...
)
```

### Версионирование

Go Modules использует семантическое версионирование (SemVer):

- **v1.2.3** - основная.минорная.патч
- **v0.x.x** - нестабильная версия
- **v1.x.x** - стабильная версия

**Пример использования версий в go.mod:**

```go
module myproject

go 1.21

require (
	// Точная версия
	github.com/pkg/errors v0.9.1
	
	// Минимальная версия (>=)
	github.com/gin-gonic/gin v1.9.0
	
	// Версия из ветки
	github.com/some/repo v0.0.0-20231015123456-abcdef123456
)
```

### Управление версиями в коде

**Пример использования зависимостей в коде:**

```go
package config

import (
	"database/sql"
	_ "github.com/lib/pq" // драйвер PostgreSQL
	
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Dependencies struct {
	DB     *sql.DB
	Router *gin.Engine
}

func NewDependencies() (*Dependencies, error) {
	db, err := sql.Open("postgres", "connection string")
	if err != nil {
		return nil, err
	}
	
	router := gin.Default()
	
	return &Dependencies{
		DB:     db,
		Router: router,
	}, nil
}
```

### Вендоринг зависимостей

**Пример использования vendor:**

```bash
# Создание vendor директории
go mod vendor

# Использование vendor при сборке
go build -mod=vendor
```

### Best Practices

**Пример регулярной проверки обновлений:**

```bash
# Еженедельно проверяйте обновления
go list -u -m all
```

**Пример очистки зависимостей:**

```bash
# После каждого изменения зависимостей
go mod tidy
```

**Пример проверки уязвимостей:**

```bash
# Использование govulncheck (Go 1.18+)
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

---

## 4. Обновление и поддержка программного обеспечения после выпуска

### Стратегия версионирования

Используйте семантическое версионирование для вашего проекта:

- **MAJOR.MINOR.PATCH** (например, v1.2.3)
- **MAJOR** - несовместимые изменения API
- **MINOR** - новая функциональность с обратной совместимостью
- **PATCH** - исправления багов с обратной совместимостью

**Пример версионирования в коде:**

```go
package main

import "fmt"

var (
	Version   = "1.2.3"
	BuildTime = "2024-01-15T10:30:00Z"
	GitCommit = "abc123def"
)

func PrintVersion() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Git Commit: %s\n", GitCommit)
}
```

### Чейнджлоги

Ведите файл CHANGELOG.md для отслеживания изменений:

```markdown
# Changelog

## [1.2.3] - 2024-01-15

### Fixed
- Исправлена утечка памяти в обработчике запросов
- Исправлена ошибка валидации email

## [1.2.2] - 2024-01-10

### Added
- Добавлена поддержка Redis кэширования
- Новый endpoint для статистики

### Changed
- Улучшена производительность запросов к БД

## [1.2.1] - 2024-01-05

### Security
- Обновлена зависимость для устранения уязвимости CVE-2024-XXXX
```

### Health checks

**Пример health check endpoint:**

```go
package health

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type HealthStatus struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

type HealthChecker struct {
	db      *sql.DB
	version string
}

func NewHealthChecker(db *sql.DB, version string) *HealthChecker {
	return &HealthChecker{
		db:      db,
		version: version,
	}
}

func (h *HealthChecker) Check(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Version:   h.version,
		Timestamp: time.Now(),
		Checks:    make(map[string]string),
	}

	// Проверка БД
	if err := h.db.Ping(); err != nil {
		status.Status = "unhealthy"
		status.Checks["database"] = "failed"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		status.Checks["database"] = "ok"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
```

### Best Practices для поддержки

**Пример версионирования API:**

```go
// api/v1/users.go
package v1

// api/v2/users.go  
package v2

// Поддержка обеих версий одновременно
```

**Пример тестирования перед релизом:**

```bash
# Обязательно тестируйте перед релизом
go test ./...
go test -race ./...
go test -cover ./...
```

**Пример CI/CD конфигурации для автоматических обновлений:**

```yaml
# .github/workflows/release.yml
name: Release
on:
  push:
    tags:
      - 'v*'
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: go test ./...
      - run: go build -o app
```

---

## Заключение

Следование этим принципам обеспечивает:
- Читаемость и поддерживаемость кода
- Легкость тестирования и отладки
- Простое управление зависимостями
- Плавные обновления без простоев

Всегда используйте `gofmt`, `go vet`, и `golint` для поддержания качества кода.
