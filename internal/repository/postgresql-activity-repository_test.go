package repository

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestPostgreSQLActivityRepositoryFiltersHistoryAndRetention(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "metrics.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.PostgreSQLActivity{}); err != nil {
		t.Fatalf("migrate activity: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	activities := []model.PostgreSQLActivity{
		{ID: "active-old", ActivityKey: "1:42:100", TargetID: 1, TargetName: "prod", DatabaseName: "app", Username: "alice", State: "active", BlockingPID: 7, BlockingQuery: "UPDATE devopin_test_lock", Query: "SELECT 1", QueryDurationMS: 9000, ObservedAt: now.Add(-2 * time.Minute)},
		{ID: "active-new", ActivityKey: "1:42:100", TargetID: 1, TargetName: "prod", DatabaseName: "app", Username: "alice", State: "active", Query: "SELECT 2", QueryDurationMS: 1000, ObservedAt: now.Add(-time.Minute)},
		{ID: "other-target", ActivityKey: "2:7:200", TargetID: 2, TargetName: "staging", DatabaseName: "app", Username: "bob", State: "active", Query: "SELECT 3", QueryDurationMS: 12000, ObservedAt: now.Add(-time.Minute)},
		{ID: "expired", ActivityKey: "1:42:100", TargetID: 1, TargetName: "prod", DatabaseName: "app", Username: "alice", State: "active", ObservedAt: now.Add(-31 * 24 * time.Hour)},
	}
	if err := db.Create(&activities).Error; err != nil {
		t.Fatalf("insert activities: %v", err)
	}

	repo := NewPostgreSQLActivityRepository(db)
	history, err := repo.ListPostgreSQLActivityHistory(PostgreSQLActivityHistoryFilter{
		TargetID:        1,
		ActivityKey:     "1:42:100",
		Cutoff:          now.Add(-time.Hour),
		MinimumDuration: 5000,
		Limit:           20,
	})
	if err != nil {
		t.Fatalf("ListPostgreSQLActivityHistory() error = %v", err)
	}
	if len(history) != 1 || history[0].ID != "active-old" {
		t.Fatalf("unexpected history: %+v", history)
	}
	blockedHistory, err := repo.ListPostgreSQLActivityHistory(PostgreSQLActivityHistoryFilter{
		TargetID:    1,
		Cutoff:      now.Add(-time.Hour),
		BlockedOnly: true,
		Limit:       20,
	})
	if err != nil {
		t.Fatalf("ListPostgreSQLActivityHistory() blocked filter error = %v", err)
	}
	if len(blockedHistory) != 1 || blockedHistory[0].BlockingPID != 7 {
		t.Fatalf("unexpected blocked history: %+v", blockedHistory)
	}

	deleted, err := repo.DeleteOlderPostgreSQLActivities(now.Add(-30 * 24 * time.Hour))
	if err != nil {
		t.Fatalf("DeleteOlderPostgreSQLActivities() error = %v", err)
	}
	if deleted != 1 {
		t.Fatalf("deleted = %d, want 1", deleted)
	}
}
