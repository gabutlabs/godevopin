package worker

import (
	"log"

	service "github.com/gabutlabs/godevopin/internal/services"
)

type PostgreSQLActivityWorker struct {
	service service.PostgreSQLActivityService
}

func NewPostgreSQLActivityWorker(svc service.PostgreSQLActivityService) *PostgreSQLActivityWorker {
	return &PostgreSQLActivityWorker{service: svc}
}

func (w *PostgreSQLActivityWorker) CollectAndPersist() {
	if err := w.service.CollectAndPersistPostgreSQLActivity(); err != nil {
		log.Printf("PostgreSQL activity collection failed: %v", err)
	}
}

func (w *PostgreSQLActivityWorker) CleanupRetention() {
	deleted, err := w.service.DeleteExpiredPostgreSQLActivity()
	if err != nil {
		log.Printf("PostgreSQL activity retention cleanup failed: %v", err)
		return
	}
	if deleted > 0 {
		log.Printf("deleted %d PostgreSQL activity records older than 30 days", deleted)
	}
}
