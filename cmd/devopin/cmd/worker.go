// cmd/nama-aplikasi/cmd/worker.go
package cmd

import (
	"log"

	"github.com/gabutlabs/devopin/internal/config"
	"github.com/gabutlabs/devopin/internal/database"
	"github.com/gabutlabs/devopin/internal/worker"
	"github.com/spf13/cobra"
)

// Variabel untuk menyimpan nilai dari flag --task
var taskName string

// workerCmd merepresentasikan command 'worker'
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Menjalankan sebuah tugas worker",
	Long:  `Mensimulasikan eksekusi sebuah tugas di background, seperti mengirim email atau memproses data.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("could not load config: %v", err)
		}
		db, err := database.ConnectDB(cfg.Database) // Sekarang db bertipe *gorm.DB
		if err != nil {
			log.Fatalf("could not connect to database: %v", err)
		}
		// if taskName == "" {
		// 	fmt.Println("Error: Nama tugas tidak boleh kosong. Gunakan flag --task.")
		// 	return
		// }
		worker.StartWorkers(&cfg, db)
	},
}

func init() {
	// Menambahkan flag --task atau -t ke subcommand 'worker'
	workerCmd.Flags().StringVarP(&taskName, "task", "t", "", "Nama tugas yang akan dijalankan (wajib diisi)")
}
