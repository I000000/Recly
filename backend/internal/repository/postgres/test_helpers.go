//go:build integration
// +build integration

package postgres

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	// Проверяем доступность Docker
	if _, err := testcontainers.NewDockerClient(); err != nil {
		t.Skip("Docker not available, skipping integration test")
	}

	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithPollInterval(200*time.Millisecond).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	// Даем время на полную инициализацию
	time.Sleep(2 * time.Second)

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	// Используем localhost
	connStr := fmt.Sprintf("postgres://testuser:testpass@localhost:%s/testdb?sslmode=disable", port.Port())

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	// Применяем миграции
	if err := applyMigrations(ctx, pool); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	cleanup := func() {
		pool.Close()
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return pool, cleanup
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// Получаем текущий путь к файлу тестов
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}
	// Поднимаемся наверх до корня проекта: internal/repository/postgres -> ../../
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(currentFile))))
	migDir := filepath.Join(projectRoot, "internal", "repository", "postgres", "migrations")

	// Проверяем, существует ли директория
	if _, err := os.Stat(migDir); os.IsNotExist(err) {
		// Если не существует, попробуем альтернативные пути
		possiblePaths := []string{
			"internal/repository/postgres/migrations",
			"./internal/repository/postgres/migrations",
			"../internal/repository/postgres/migrations",
			"../../internal/repository/postgres/migrations",
			migDir,
		}
		var found bool
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				migDir = path
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("migrations directory not found in any of the paths: %v", possiblePaths)
		}
	}

	files, err := os.ReadDir(migDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".up.sql") && !strings.HasSuffix(name, ".up") {
			continue
		}
		if !strings.Contains(name, ".up") {
			continue
		}

		path := filepath.Join(migDir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", name, err)
		}

		lines := strings.Split(string(content), "\n")
		cleaned := ""
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "--") || trimmed == "" {
				continue
			}
			cleaned += line + "\n"
		}

		queries := strings.Split(cleaned, ";")
		for _, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}
			timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			_, err := pool.Exec(timeoutCtx, q)
			if err != nil {
				return fmt.Errorf("exec migration %s: %w\nSQL: %s", name, err, q)
			}
		}
	}
	return nil
}
