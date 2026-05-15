package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	service "github.com/gabutlabs/devopin/internal/services"
	"github.com/gabutlabs/devopin/internal/services/ai"
	"gopkg.in/telebot.v3"
)

type TelegramWorker struct {
	settingService service.SettingService
	sysService     service.SystemMetricService
	wrkService     service.WorkerServiceService
	docService     *service.DockerService
	logService     service.LogHistoryService
	prjService     service.ProjectService
}

func NewTelegramWorker(
	settingSvc service.SettingService,
	sysSvc service.SystemMetricService,
	wrkSvc service.WorkerServiceService,
	docSvc *service.DockerService,
	logSvc service.LogHistoryService,
	prjSvc service.ProjectService,
) *TelegramWorker {
	return &TelegramWorker{
		settingService: settingSvc,
		sysService:     sysSvc,
		wrkService:     wrkSvc,
		docService:     docSvc,
		logService:     logSvc,
		prjService:     prjSvc,
	}
}

func (w *TelegramWorker) Start() {
	var token string
	for {
		settings, err := w.settingService.GetSettings()
		if err != nil {
			log.Printf("TelegramWorker: failed to get settings: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if settings.TelegramBotToken == "" {
			log.Println("TelegramWorker: TelegramBotToken is empty, waiting 10s...")
			time.Sleep(10 * time.Second)
			continue
		}

		token = settings.TelegramBotToken
		break
	}

	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Printf("TelegramWorker: failed to create bot: %v", err)
		return
	}

	log.Println("TelegramWorker: Bot started successfully")

	toolProvider := ai.NewAIToolProvider(
		w.sysService,
		w.wrkService,
		w.docService,
		w.logService,
		w.prjService,
	)

	b.Handle(telebot.OnText, func(c telebot.Context) error {
		userMsg := c.Text()
		log.Printf("TelegramWorker: received message: %s", userMsg)

		// Create a dynamic LLM client based on current settings
		// Note: In production, you might want to cache this or handle reloads better
		currentSettings, _ := w.settingService.GetSettings()

		llm, err := ai.NewLLMClient(
			context.Background(),
			currentSettings.AIProvider,
			currentSettings.AIApiKey,
			currentSettings.AIModelName,
			currentSettings.AIBaseURL,
		)

		if err != nil {
			return c.Send(fmt.Sprintf("AI Error: %v", err))
		}

		agent := ai.NewAIAgent(llm, toolProvider)

		// Send "typing..." action
		c.Notify(telebot.Typing)
		fmt.Println("TelegramWorker: received message: %s", userMsg)
		resp, err := agent.Chat(context.Background(), userMsg)
		if err != nil {
			log.Printf("TelegramWorker: agent error: %v", err)
			return c.Send(fmt.Sprintf("Agent Error: %v", err))
		}

		// Sanitize response (ensure it doesn't break Telegram markdown if any, though we use plain text here)
		return c.Send(resp)
	})

	b.Start()
}
