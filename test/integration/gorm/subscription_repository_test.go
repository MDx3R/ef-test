package gorm_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/MDx3R/ef-test/internal/config"
	"github.com/MDx3R/ef-test/internal/domain/entity"
	gormdb "github.com/MDx3R/ef-test/internal/infra/database/gorm"
	"github.com/MDx3R/ef-test/internal/usecase"
	"github.com/MDx3R/ef-test/internal/usecase/dto"
	"github.com/MDx3R/ef-test/internal/usecase/model"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDB *gorm.DB
	repo   usecase.SubscriptionRepository
	pgC    testcontainers.Container
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var (
		username = "test"
		password = "test"
		database = "testdb"
	)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		Env:          map[string]string{"POSTGRES_USER": username, "POSTGRES_PASSWORD": password, "POSTGRES_DB": database},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}
	var err error
	pgC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}

	host, err := pgC.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %v", err)
	}
	port, err := pgC.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Failed to get mapped port: %v", err)
	}

	cfg := config.DatabaseConfig{
		Driver:   "postgres",
		Host:     host,
		Port:     port.Port(),
		Username: username,
		Password: password,
		Database: database,
	}

	gormDB, err := gormdb.NewGormDatabase(&cfg)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	err = gormDB.Migrate()
	if err != nil {
		log.Fatalf("Failed to migrate DB: %v", err)
	}

	testDB = gormDB.GetDB()
	repo = gormdb.NewGormSubscriptionRepository(testDB)

	code := m.Run()

	if err := gormDB.Dispose(); err != nil {
		log.Fatalf("Failed to dispose database: %v", err)
	}

	if err := pgC.Terminate(ctx); err != nil {
		log.Fatalf("Failed to terminate container: %v", err)
	}

	os.Exit(code)
}

func clearTable(t *testing.T) {
	err := testDB.Exec("TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE").Error
	if err != nil {
		t.Fatalf("Failed to clear table: %v", err)
	}
}

func makeTestSubscription(t *testing.T) *entity.Subscription {
	id := uuid.New()
	sub, _ := entity.NewSubscriptionWithID(
		id,
		"test_service",
		uuid.New(),
		100,
		time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		nil,
	)
	return sub
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func getIDs(subs []*entity.Subscription) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(subs))
	for _, s := range subs {
		ids = append(ids, s.ID())
	}
	return ids
}

func TestGormSubscriptionRepository_AddAndGet(t *testing.T) {
	clearTable(t)

	// Arrange
	sub := makeTestSubscription(t)

	// Act
	errAdd := repo.Add(sub)
	got, errGet := repo.Get(sub.ID())

	// Assert
	assert.NoError(t, errAdd)
	assert.NoError(t, errGet)
	assert.Equal(t, sub.ID(), got.ID())
	assert.Equal(t, sub.UserID(), got.UserID())
	assert.Equal(t, sub.ServiceName(), got.ServiceName())
	assert.Equal(t, sub.Price(), got.Price())
}

func TestGormSubscriptionRepository_Get_NotFound(t *testing.T) {
	clearTable(t)

	// Arrange
	randomID := uuid.New()

	// Act
	_, err := repo.Get(randomID)

	// Assert
	assert.ErrorIs(t, err, usecase.ErrNotFound)
}

func TestGormSubscriptionRepository_List(t *testing.T) {
	clearTable(t)

	// Arrange
	sub1 := makeTestSubscription(t)
	sub2 := makeTestSubscription(t)
	err := repo.Add(sub1)
	assert.NoError(t, err)
	err = repo.Add(sub2)
	assert.NoError(t, err)

	filter := dto.SubscriptionFilter{
		Page:     1,
		PageSize: 10,
	}

	// Act
	list, err := repo.List(filter)

	// Assert
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
}

func TestGormSubscriptionRepository_List_EmptyResult(t *testing.T) {
	clearTable(t)

	// Arrange
	serviceName := "non-existent-service"
	filter := dto.SubscriptionFilter{
		Page:        1,
		PageSize:    10,
		ServiceName: &serviceName,
	}

	// Act
	list, err := repo.List(filter)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}

func TestGormSubscriptionRepository_List_Pagination(t *testing.T) {
	clearTable(t)

	// Arrange
	for range 15 {
		sub := makeTestSubscription(t)
		err := repo.Add(sub)
		assert.NoError(t, err)
	}

	filterPage1 := dto.SubscriptionFilter{
		Page:     1,
		PageSize: 10,
	}
	filterPage2 := dto.SubscriptionFilter{
		Page:     2,
		PageSize: 10,
	}

	// Act
	listPage1, err1 := repo.List(filterPage1)
	listPage2, err2 := repo.List(filterPage2)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Len(t, listPage1, 10)
	assert.Len(t, listPage2, 5)
}

func TestGormSubscriptionRepository_List_FilterByUserID(t *testing.T) {
	clearTable(t)

	// Arrange
	userID1 := uuid.New()
	userID2 := uuid.New()

	sub1, _ := entity.NewSubscriptionWithID(uuid.New(), "service1", userID1, 50, time.Now(), nil)
	sub2, _ := entity.NewSubscriptionWithID(uuid.New(), "service2", userID2, 100, time.Now(), nil)

	assert.NoError(t, repo.Add(sub1))
	assert.NoError(t, repo.Add(sub2))

	userID1Str := userID1.String()
	filter := dto.SubscriptionFilter{
		UserID:   &userID1Str,
		Page:     1,
		PageSize: 10,
	}

	// Act
	list, err := repo.List(filter)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, userID1, list[0].UserID())
}

func TestGormSubscriptionRepository_List_FilterByServiceName(t *testing.T) {
	clearTable(t)

	// Arrange
	sub1, _ := entity.NewSubscriptionWithID(uuid.New(), "serviceA", uuid.New(), 50, time.Now(), nil)
	sub2, _ := entity.NewSubscriptionWithID(uuid.New(), "serviceB", uuid.New(), 100, time.Now(), nil)

	assert.NoError(t, repo.Add(sub1))
	assert.NoError(t, repo.Add(sub2))

	serviceName := "serviceB"
	filter := dto.SubscriptionFilter{
		ServiceName: &serviceName,
		Page:        1,
		PageSize:    10,
	}

	// Act
	list, err := repo.List(filter)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "serviceB", list[0].ServiceName())
}

func TestSubscriptionRepository_List_FilterDates(t *testing.T) {
	clearTable(t)

	// Arrange
	startDate := model.NewMonthYear(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC))
	endDate := model.NewMonthYear(time.Date(2025, 8, 31, 23, 59, 59, 0, time.UTC))

	// 1. start_date: 2025-07-01, end_date: NULL
	sub1, _ := entity.NewSubscriptionWithID(uuid.New(), "serviceA", uuid.New(), 100, time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC), nil)

	// 2. start_date: 2025-08-05, end_date: 2025-08-20
	sub2, _ := entity.NewSubscriptionWithID(uuid.New(), "serviceA", uuid.New(), 150, time.Date(2025, 8, 5, 0, 0, 0, 0, time.UTC), timePtr(time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC)))

	// 3. start_date: 2025-08-15, end_date: 2025-09-01
	sub3, _ := entity.NewSubscriptionWithID(uuid.New(), "serviceB", uuid.New(), 200, time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), timePtr(time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)))

	for _, s := range []*entity.Subscription{sub1, sub2, sub3} {
		assert.NoError(t, repo.Add(s))
	}

	// Act & Assert
	t.Run("StartDate set, EndDate nil", func(t *testing.T) {
		filter := dto.SubscriptionFilter{
			StartDate: &startDate,
			Page:      1,
			PageSize:  10,
		}
		subs, err := repo.List(filter)
		require.NoError(t, err)

		// start_date >= 2025-08-01 => sub2 и sub3
		assert.ElementsMatch(t, []uuid.UUID{sub2.ID(), sub3.ID()}, getIDs(subs))
	})

	t.Run("EndDate set, StartDate nil", func(t *testing.T) {
		filter := dto.SubscriptionFilter{
			EndDate:  &endDate,
			Page:     1,
			PageSize: 10,
		}
		subs, err := repo.List(filter)
		require.NoError(t, err)

		// end_date <= 2025-08-31 OR end_date IS NULL
		// sub1: end_date=NUL
		// sub2: end_date=2025-08-20 <= endDate
		// sub3: end_date=2025-09-01 > endDate -> exclude
		assert.ElementsMatch(t, []uuid.UUID{sub1.ID(), sub2.ID()}, getIDs(subs))
	})

	t.Run("StartDate и EndDate установлены", func(t *testing.T) {
		filter := dto.SubscriptionFilter{
			StartDate: &startDate,
			EndDate:   &endDate,
			Page:      1,
			PageSize:  10,
		}
		subs, err := repo.List(filter)
		require.NoError(t, err)

		// start_date >= 2025-08-01 AND (end_date <= 2025-08-31 OR end_date IS NULL)
		// sub1: start_date=2025-07-01 < startDate -> exclude
		// sub2: start_date=2025-08-05 >= startDate AND end_date=2025-08-20 <= endDate
		// sub3: start_date=2025-08-15 >= startDate AND end_date=2025-09-01 > endDate - exclude
		assert.ElementsMatch(t, []uuid.UUID{sub2.ID()}, getIDs(subs))
	})
}

func TestGormSubscriptionRepository_Update(t *testing.T) {
	clearTable(t)

	// Arrange
	sub := makeTestSubscription(t)
	err := repo.Add(sub)
	assert.NoError(t, err)

	sub.SetServiceName("updated_service")
	sub.SetPrice(200)

	// Act
	errUpdate := repo.Update(sub)
	got, errGet := repo.Get(sub.ID())

	// Assert
	assert.NoError(t, errUpdate)
	assert.NoError(t, errGet)
	assert.Equal(t, "updated_service", got.ServiceName())
	assert.Equal(t, 200, got.Price())
}

func TestGormSubscriptionRepository_Delete(t *testing.T) {
	clearTable(t)

	// Arrange
	sub := makeTestSubscription(t)
	err := repo.Add(sub)
	assert.NoError(t, err)

	// Act
	errDelete := repo.Delete(sub.ID())
	_, errGet := repo.Get(sub.ID())

	// Assert
	assert.NoError(t, errDelete)
	assert.ErrorIs(t, errGet, usecase.ErrNotFound)
}

func TestGormSubscriptionRepository_CalculateTotalCost(t *testing.T) {
	clearTable(t)

	// Arrange
	userID := uuid.New()
	sub1, _ := entity.NewSubscriptionWithID(uuid.New(), "serviceA", userID, 100, time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC), nil)
	sub2, _ := entity.NewSubscriptionWithID(uuid.New(), "serviceA", userID, 150, time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC), nil)
	sub3, _ := entity.NewSubscriptionWithID(uuid.New(), "serviceB", userID, 200, time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC), nil)

	assert.NoError(t, repo.Add(sub1))
	assert.NoError(t, repo.Add(sub2))
	assert.NoError(t, repo.Add(sub3))

	startDate := model.NewMonthYear(time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC))
	endDate := model.NewMonthYear(time.Date(2025, 8, 31, 23, 59, 59, 0, time.UTC))
	filter := dto.TotalCostFilter{
		UserID:      userID.String(),
		ServiceName: "serviceA",
		PeriodStart: &startDate,
		PeriodEnd:   &endDate,
	}

	// Act
	total, err := repo.CalculateTotalCost(filter)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 250, total) // sub1 + sub2
}
