package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/gabutlabs/godevopin/internal/monitoring"
	"github.com/gabutlabs/godevopin/internal/repository"
	"gorm.io/gorm"
)

const (
	MySQLActivityRetention    = 30 * 24 * time.Hour
	mysqlActivityTimeout      = 10 * time.Second
	mysqlActivityDefaultLimit = 500
)

type MySQLTargetRequest struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
	Enabled  *bool  `json:"enabled"`
}

type MySQLTargetView struct {
	ID                 uint       `json:"id"`
	Name               string     `json:"name"`
	Host               string     `json:"host"`
	Port               int        `json:"port"`
	Database           string     `json:"database"`
	Username           string     `json:"username"`
	SSLMode            string     `json:"ssl_mode"`
	Enabled            bool       `json:"enabled"`
	PasswordConfigured bool       `json:"password_configured"`
	LastCheckedAt      *time.Time `json:"last_checked_at,omitempty"`
	LastError          string     `json:"last_error,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type MySQLActivityLiveResult struct {
	Activities []model.MySQLActivity `json:"activities"`
	Targets    []MySQLTargetView     `json:"targets"`
}

type MySQLActivityService interface {
	ListMySQLTargets() ([]MySQLTargetView, error)
	CreateMySQLTarget(request MySQLTargetRequest) (MySQLTargetView, error)
	UpdateMySQLTarget(id uint, request MySQLTargetRequest) (MySQLTargetView, error)
	DeleteMySQLTarget(id uint) error
	TestMySQLTarget(request MySQLTargetRequest) (string, error)
	GetLiveMySQLActivity(targetID uint, search string, limit int) (MySQLActivityLiveResult, error)
	GetMySQLActivityHistory(activityKey string, targetID uint, filter, search, state string, minimumDurationMS int64, limit int) ([]model.MySQLActivity, error)
	CollectAndPersistMySQLActivity() error
	DeleteExpiredMySQLActivity() (int64, error)
}

type mysqlActivityService struct {
	targetRepository   repository.MySQLTargetRepository
	activityRepository repository.MySQLActivityRepository
	collector          *monitoring.MySQLCollector
	secret             string
}

func NewMySQLActivityService(targetRepo repository.MySQLTargetRepository, activityRepo repository.MySQLActivityRepository, secret string) MySQLActivityService {
	return &mysqlActivityService{
		targetRepository:   targetRepo,
		activityRepository: activityRepo,
		collector:          monitoring.NewMySQLCollector(),
		secret:             secret,
	}
}

func (s *mysqlActivityService) ListMySQLTargets() ([]MySQLTargetView, error) {
	targets, err := s.targetRepository.ListMySQLTargets()
	if err != nil {
		return nil, err
	}
	result := make([]MySQLTargetView, 0, len(targets))
	for _, target := range targets {
		result = append(result, mysqlTargetView(target))
	}
	return result, nil
}

func (s *mysqlActivityService) CreateMySQLTarget(request MySQLTargetRequest) (MySQLTargetView, error) {
	target, err := normalizeMySQLTargetRequest(request, 0)
	if err != nil {
		return MySQLTargetView{}, err
	}
	if strings.TrimSpace(request.Password) == "" {
		return MySQLTargetView{}, errors.New("password is required")
	}
	target.PasswordCiphertext, err = s.encrypt(request.Password)
	if err != nil {
		return MySQLTargetView{}, err
	}
	if err := s.targetRepository.CreateMySQLTarget(&target); err != nil {
		return MySQLTargetView{}, err
	}
	return mysqlTargetView(target), nil
}

func (s *mysqlActivityService) UpdateMySQLTarget(id uint, request MySQLTargetRequest) (MySQLTargetView, error) {
	target, err := s.targetRepository.GetMySQLTargetByID(id)
	if err != nil {
		return MySQLTargetView{}, err
	}
	updated, err := normalizeMySQLTargetRequest(request, id)
	if err != nil {
		return MySQLTargetView{}, err
	}
	updated.PasswordCiphertext = target.PasswordCiphertext
	updated.CreatedAt = target.CreatedAt
	updated.LastCheckedAt = target.LastCheckedAt
	updated.LastError = target.LastError
	if strings.TrimSpace(request.Password) != "" {
		updated.PasswordCiphertext, err = s.encrypt(request.Password)
		if err != nil {
			return MySQLTargetView{}, err
		}
	}
	if err := s.targetRepository.UpdateMySQLTarget(&updated); err != nil {
		return MySQLTargetView{}, err
	}
	return mysqlTargetView(updated), nil
}

func (s *mysqlActivityService) DeleteMySQLTarget(id uint) error {
	return s.targetRepository.DeleteMySQLTarget(id)
}

func (s *mysqlActivityService) TestMySQLTarget(request MySQLTargetRequest) (string, error) {
	target, err := normalizeMySQLTargetRequest(request, 0)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(request.Password) == "" {
		return "", errors.New("password is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), mysqlActivityTimeout)
	defer cancel()
	return s.collector.TestConnection(ctx, target, request.Password)
}

func (s *mysqlActivityService) GetLiveMySQLActivity(targetID uint, search string, limit int) (MySQLActivityLiveResult, error) {
	targets, err := s.targetRepository.ListMySQLTargets()
	if err != nil {
		return MySQLActivityLiveResult{}, err
	}
	if targetID > 0 {
		filtered := targets[:0]
		for _, target := range targets {
			if target.ID == targetID {
				filtered = append(filtered, target)
			}
		}
		if len(filtered) == 0 {
			return MySQLActivityLiveResult{}, gorm.ErrRecordNotFound
		}
		targets = filtered
	}

	activities := make([]model.MySQLActivity, 0)
	for _, target := range targets {
		if !target.Enabled {
			continue
		}
		password, err := s.decrypt(target.PasswordCiphertext)
		if err != nil {
			s.recordHealth(target, err)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), mysqlActivityTimeout)
		collection, collectErr := s.collector.Collect(ctx, target, password)
		cancel()
		if collectErr != nil {
			s.recordHealth(target, collectErr)
			continue
		}
		s.recordHealth(target, nil)
		activities = append(activities, collection.Activities...)
	}

	search = strings.ToLower(strings.TrimSpace(search))
	if search != "" {
		filtered := activities[:0]
		for _, activity := range activities {
			if mysqlActivityContains(activity, search) {
				filtered = append(filtered, activity)
			}
		}
		activities = filtered
	}
	sort.SliceStable(activities, func(i, j int) bool {
		return activities[i].QueryDurationMS > activities[j].QueryDurationMS
	})
	if limit <= 0 {
		limit = mysqlActivityDefaultLimit
	}
	if limit > 2000 {
		limit = 2000
	}
	if len(activities) > limit {
		activities = activities[:limit]
	}
	views, err := s.ListMySQLTargets()
	if err != nil {
		return MySQLActivityLiveResult{}, err
	}
	return MySQLActivityLiveResult{Activities: activities, Targets: views}, nil
}

func (s *mysqlActivityService) GetMySQLActivityHistory(activityKey string, targetID uint, filter, search, state string, minimumDurationMS int64, limit int) ([]model.MySQLActivity, error) {
	duration, err := mysqlActivityDuration(filter)
	if err != nil {
		return nil, err
	}
	return s.activityRepository.ListMySQLActivityHistory(repository.MySQLActivityHistoryFilter{
		ActivityKey:     activityKey,
		TargetID:        targetID,
		Cutoff:          time.Now().UTC().Add(-duration),
		Search:          search,
		State:           state,
		MinimumDuration: minimumDurationMS,
		Limit:           limit,
	})
}

func (s *mysqlActivityService) CollectAndPersistMySQLActivity() error {
	targets, err := s.targetRepository.ListMySQLTargets()
	if err != nil {
		return err
	}
	activities := make([]model.MySQLActivity, 0)
	for _, target := range targets {
		if !target.Enabled {
			continue
		}
		password, decryptErr := s.decrypt(target.PasswordCiphertext)
		if decryptErr != nil {
			s.recordHealth(target, decryptErr)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), mysqlActivityTimeout)
		collection, collectErr := s.collector.Collect(ctx, target, password)
		cancel()
		if collectErr != nil {
			s.recordHealth(target, collectErr)
			continue
		}
		s.recordHealth(target, nil)
		for i := range collection.Activities {
			activity := &collection.Activities[i]
			activity.ID = fmt.Sprintf("%s:%d", activity.ActivityKey, activity.ObservedAt.UnixNano())
			activities = append(activities, *activity)
		}
	}
	return s.activityRepository.CreateMySQLActivities(activities)
}

func (s *mysqlActivityService) DeleteExpiredMySQLActivity() (int64, error) {
	return s.activityRepository.DeleteOlderMySQLActivities(time.Now().UTC().Add(-MySQLActivityRetention))
}

func (s *mysqlActivityService) recordHealth(target model.MySQLTarget, collectErr error) {
	checkedAt := time.Now().UTC()
	lastError := ""
	if collectErr != nil {
		lastError = collectErr.Error()
		if password, err := s.decrypt(target.PasswordCiphertext); err == nil && password != "" {
			lastError = strings.ReplaceAll(lastError, password, "[redacted]")
		}
	}
	_ = s.targetRepository.UpdateMySQLTargetHealth(target.ID, &checkedAt, lastError)
}

func (s *mysqlActivityService) encrypt(value string) (string, error) {
	block, err := s.cipherBlock()
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, block.NonceSize()+len(value))
	nonce := ciphertext[:block.NonceSize()]
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := block.Seal(nonce, nonce, []byte(value), nil)
	return base64.RawStdEncoding.EncodeToString(sealed), nil
}

func (s *mysqlActivityService) decrypt(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", errors.New("encrypted MySQL password is missing")
	}
	block, err := s.cipherBlock()
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.RawStdEncoding.DecodeString(value)
	if err != nil || len(ciphertext) < block.NonceSize() {
		return "", errors.New("encrypted MySQL password is invalid")
	}
	nonce, ciphertext := ciphertext[:block.NonceSize()], ciphertext[block.NonceSize():]
	plaintext, err := block.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("unable to decrypt MySQL password")
	}
	return string(plaintext), nil
}

func (s *mysqlActivityService) cipherBlock() (cipher.AEAD, error) {
	if strings.TrimSpace(s.secret) == "" {
		return nil, errors.New("app.jwt_secret is required to protect MySQL credentials")
	}
	key := sha256.Sum256([]byte(s.secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func normalizeMySQLTargetRequest(request MySQLTargetRequest, id uint) (model.MySQLTarget, error) {
	name, host, username := strings.TrimSpace(request.Name), strings.TrimSpace(request.Host), strings.TrimSpace(request.Username)
	if name == "" || host == "" || username == "" {
		return model.MySQLTarget{}, errors.New("name, host, and username are required")
	}
	port := request.Port
	if port == 0 {
		port = 3306
	}
	if port < 1 || port > 65535 {
		return model.MySQLTarget{}, errors.New("port must be between 1 and 65535")
	}
	databaseName := strings.TrimSpace(request.Database)
	if databaseName == "" {
		databaseName = "mysql"
	}
	sslMode := strings.ToLower(strings.TrimSpace(request.SSLMode))
	if sslMode == "" {
		sslMode = "preferred"
	}
	if sslMode != "disable" && sslMode != "preferred" && sslMode != "required" {
		return model.MySQLTarget{}, errors.New("ssl_mode must be disable, preferred, or required")
	}
	enabled := true
	if request.Enabled != nil {
		enabled = *request.Enabled
	}
	return model.MySQLTarget{ID: id, Name: name, Host: host, Port: port, Database: databaseName, Username: username, SSLMode: sslMode, Enabled: enabled}, nil
}

func mysqlTargetView(target model.MySQLTarget) MySQLTargetView {
	return MySQLTargetView{ID: target.ID, Name: target.Name, Host: target.Host, Port: target.Port, Database: target.Database, Username: target.Username, SSLMode: target.SSLMode, Enabled: target.Enabled, PasswordConfigured: target.PasswordCiphertext != "", LastCheckedAt: target.LastCheckedAt, LastError: target.LastError, CreatedAt: target.CreatedAt, UpdatedAt: target.UpdatedAt}
}

func mysqlActivityContains(activity model.MySQLActivity, search string) bool {
	for _, value := range []string{activity.TargetName, activity.DatabaseName, activity.Username, activity.ClientAddress, activity.Command, activity.State, activity.Query, fmt.Sprint(activity.PID)} {
		if strings.Contains(strings.ToLower(value), search) {
			return true
		}
	}
	return false
}

func mysqlActivityDuration(filter string) (time.Duration, error) {
	switch strings.ToLower(strings.TrimSpace(filter)) {
	case "1h", "1 hour":
		return time.Hour, nil
	case "6h", "6 hours":
		return 6 * time.Hour, nil
	case "12h", "12 hours":
		return 12 * time.Hour, nil
	case "1d", "1 day":
		return 24 * time.Hour, nil
	case "7d", "7 days":
		return 7 * 24 * time.Hour, nil
	case "30d", "30 days":
		return MySQLActivityRetention, nil
	default:
		return 0, errors.New("invalid filter format")
	}
}
