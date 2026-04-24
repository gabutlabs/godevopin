package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd adalah command dasar saat aplikasi dijalankan tanpa subcommand
var rootCmd = &cobra.Command{
	Use:   "nama-aplikasi", // Ganti 'nama-aplikasi' sesuai keinginan
	Short: "Aplikasi serbaguna dengan web server dan worker",
	Long: `Ini adalah aplikasi contoh yang menunjukkan cara mengintegrasikan
Fiber web framework dengan Cobra untuk membuat aplikasi dengan beberapa fungsi.`,
}

// Execute adalah fungsi yang dipanggil oleh main.go.
// Fungsi ini akan menjalankan command yang sesuai.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init() akan dipanggil saat package ini di-load.
// Di sini kita mendaftarkan semua subcommand.
func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(workerCmd)
}
