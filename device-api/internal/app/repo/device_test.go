package repo

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/model"
)

const (
	dropDB            = `DROP DATABASE IF EXISTS test_devices;`
	createDB          = `CREATE DATABASE test_devices;`
	createDeviceTable = `CREATE TABLE IF NOT EXISTS devices (
    id          serial      PRIMARY KEY,
    platform    varchar(32) NOT NULL,
    user_id     bigint      NOT NULL,
    entered_at  timestamp   NOT NULL
);`
)

func TestMain(m *testing.M) {
	//Подготовка данных
	//Создали соединение с контейнером
	psql, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=test password=test sslmode=disable",
	)
	if err != nil {
		panic(fmt.Errorf("sql.Open() err: %v", err))
	}
	defer psql.Close()
	// Создание и удаление базы данных
	_, err = psql.Exec(dropDB)
	if err != nil {
		panic(fmt.Errorf("drop db err: %v", err))
	}

	_, err = psql.Exec(createDB)
	if err != nil {
		panic(fmt.Errorf("create db err: %v", err))
	}

	defer func() {
		_, err = psql.Exec(dropDB)
		if err != nil {
			panic(fmt.Errorf("drop db err: %v", err))
		}
	}()
	//Соединяемся с базой данных
	db, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=test password=test sslmode=disable dbname=test_devices",
	)
	if err != nil {
		panic(fmt.Errorf("sql.Open() err: %v", err))
	}
	defer db.Close()
	//Создаем таблицу в базе данных
	_, err = db.Exec(createDeviceTable)
	if err != nil {
		panic(fmt.Errorf("create device table err: %v", err))
	}

	m.Run()
	//Очистка данных
}

func TestCreateDevice(t *testing.T) {
	db, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=test password=test sslmode=disable dbname=test_devices",
	)
	require.NoError(t, err)
	defer db.Close()

	database := repo{db: &sqlx.DB{DB: db}}

	t.Run("simple device", testCreateDevice(database))
}

func testCreateDevice(database repo) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		var err error
		enteredAt := time.Now()
		device := model.Device{
			Platform:  "iOS",
			UserID:    1,
			EnteredAt: &enteredAt,
		}

		id, err := database.CreateDevice(ctx, &device)

		require.NoError(t, err)
		assert.NotZero(t, id)
	}
}
