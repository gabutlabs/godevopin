package worker

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabutlabs/devopin/internal/logparser"
	"github.com/gabutlabs/devopin/internal/model"
	service "github.com/gabutlabs/devopin/internal/services"
)

type LogParserWorker struct {
	projectService service.ProjectService
	logService     service.LogHistoryService
	parser         *logparser.Parser
	fileOffsets    map[string]int64
}

func NewLogParserWorker(ps service.ProjectService, ls service.LogHistoryService) *LogParserWorker {
	return &LogParserWorker{
		projectService: ps,
		logService:     ls,
		parser:         logparser.NewParser(),
		fileOffsets:    make(map[string]int64),
	}
}

func (w *LogParserWorker) Run() {
	projects, err := w.projectService.ListProjects()
	if err != nil {
		log.Printf("Worker [LogParser]: failed to fetch projects: %v", err)
		return
	}

	for _, project := range projects {
		if project.Path == "" {
			continue
		}

		// Find .log files in project path
		files, err := w.findLogFiles(project.Path)
		if err != nil {
			log.Printf("Worker [LogParser]: failed to find log files in %s: %v", project.Path, err)
			continue
		}

		var allLogs []model.LogHistory
		for _, file := range files {
			w.parser.ClearEntries()

			pattern := project.LogFormat
			if project.ProjectType == "laravel" && pattern == "" {
				// Default Laravel format: [timestamp] source.level: message
				pattern = "[{timestamp}] {source}.{level}: {message}"
			}

			startOffset := w.fileOffsets[file]
			entries, newOffset, err := w.parser.ParseFile(file, pattern, startOffset)
			if err != nil {
				log.Printf("Worker [LogParser]: failed to parse file %s: %v", file, err)
				continue
			}

			w.fileOffsets[file] = newOffset

			for _, entry := range entries {
				allLogs = append(allLogs, model.LogHistory{
					ProjectID: project.ID,
					Timestamp: entry.Timestamp,
					Level:     entry.Level,
					Message:   entry.Message,
					Source:    entry.Source,
				})
			}
		}

		if len(allLogs) > 0 {
			if err := w.logService.BatchInsert(allLogs); err != nil {
				log.Printf("Worker [LogParser]: failed to batch insert logs for project %s: %v", project.Name, err)
			} else {
				log.Printf("Worker [LogParser]: successfully parsed and saved %d logs for project %s", len(allLogs), project.Name)
			}
		}
	}
}

func (w *LogParserWorker) findLogFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files with errors (e.g. permission)
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".log") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
