package worker

import (
	"log"

	service "github.com/gabutlabs/godevopin/internal/services"
)

type MySQLActivityWorker struct{ service service.MySQLActivityService }

func NewMySQLActivityWorker(svc service.MySQLActivityService) *MySQLActivityWorker {
	return &MySQLActivityWorker{service: svc}
}

func (w *MySQLActivityWorker) CollectAndPersist() {
	if err := w.service.CollectAndPersistMySQLActivity(); err != nil {
		log.Printf("MySQL activity collection failed: %v", err)
	}
}

func (w *MySQLActivityWorker) CleanupRetention() {
	deleted, err := w.service.DeleteExpiredMySQLActivity()
	if err != nil {
		log.Printf("MySQL activity retention cleanup failed: %v", err)
		return
	}
	if deleted > 0 {
		log.Printf("deleted %d MySQL activity records older than 30 days", deleted)
	}
}
