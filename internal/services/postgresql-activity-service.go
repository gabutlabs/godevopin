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
	PostgreSQLActivityRetention  = 30 * 24 * time.Hour
	postgresActivityLiveTimeout  = 10 * time.Second
	postgresActivityDefaultLimit = 500
)

type PostgreSQLTargetRequest struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
	Enabled  *bool  `json:"enabled"`
}

type PostgreSQLTargetView struct {
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

type PostgreSQLActivityLiveResult struct {
	Activities []model.PostgreSQLActivity `json:"activities"`
	Targets    []PostgreSQLTargetView     `json:"targets"`
}

type PostgreSQLActivityService interface {
	ListPostgreSQLTargets() ([]PostgreSQLTargetView, error)
	CreatePostgreSQLTarget(request PostgreSQLTargetRequest) (PostgreSQLTargetView, error)
	UpdatePostgreSQLTarget(id uint, request PostgreSQLTargetRequest) (PostgreSQLTargetView, error)
	DeletePostgreSQLTarget(id uint) error
	TestPostgreSQLTarget(request PostgreSQLTargetRequest) (string, error)
	GetLivePostgreSQLActivity(targetID uint, search string, limit int) (PostgreSQLActivityLiveResult, error)
	GetPostgreSQLActivityHistory(activityKey string, targetID uint, filter, search, state string, blockedOnly bool, minimumDurationMS int64, limit int) ([]model.PostgreSQLActivity, error)
	CollectAndPersistPostgreSQLActivity() error
	DeleteExpiredPostgreSQLActivity() (int64, error)
}

type postgresqlActivityService struct {
	targetRepository   repository.PostgreSQLTargetRepository
	activityRepository repository.PostgreSQLActivityRepository
	collector          *monitoring.PostgreSQLCollector
	secret             string
}

func NewPostgreSQLActivityService(
	targetRepository repository.PostgreSQLTargetRepository,
	activityRepository repository.PostgreSQLActivityRepository,
	secret string,
) PostgreSQLActivityService {
	return &postgresqlActivityService{
		targetRepository:   targetRepository,
		activityRepository: activityRepository,
		collector:          monitoring.NewPostgreSQLCollector(),
		secret:             secret,
	}
}

func (s *postgresqlActivityService) ListPostgreSQLTargets() ([]PostgreSQLTargetView, error) {
	targets, err := s.targetRepository.ListPostgreSQLTargets()
	if err != nil {
		return nil, err
	}
	result := make([]PostgreSQLTargetView, 0, len(targets))
	for _, target := range targets {
		result = append(result, targetView(target))
	}
	return result, nil
}

func (s *postgresqlActivityService) CreatePostgreSQLTarget(request PostgreSQLTargetRequest) (PostgreSQLTargetView, error) {
	target, err := normalizeTargetRequest(request, 0)
	if err != nil {
		return PostgreSQLTargetView{}, err
	}
	if strings.TrimSpace(request.Password) == "" {
		return PostgreSQLTargetView{}, errors.New("password is required")
	}
	target.PasswordCiphertext, err = s.encrypt(request.Password)
	if err != nil {
		return PostgreSQLTargetView{}, err
	}
	if err := s.targetRepository.CreatePostgreSQLTarget(&target); err != nil {
		return PostgreSQLTargetView{}, err
	}
	return targetView(target), nil
}

func (s *postgresqlActivityService) UpdatePostgreSQLTarget(id uint, request PostgreSQLTargetRequest) (PostgreSQLTargetView, error) {
	target, err := s.targetRepository.GetPostgreSQLTargetByID(id)
	if err != nil {
		return PostgreSQLTargetView{}, err
	}
	updated, err := normalizeTargetRequest(request, id)
	if err != nil {
		return PostgreSQLTargetView{}, err
	}
	updated.PasswordCiphertext = target.PasswordCiphertext
	updated.CreatedAt = target.CreatedAt
	updated.LastCheckedAt = target.LastCheckedAt
	updated.LastError = target.LastError
	if strings.TrimSpace(request.Password) != "" {
		updated.PasswordCiphertext, err = s.encrypt(request.Password)
		if err != nil {
			return PostgreSQLTargetView{}, err
		}
	}
	if err := s.targetRepository.UpdatePostgreSQLTarget(&updated); err != nil {
		return PostgreSQLTargetView{}, err
	}
	return targetView(updated), nil
}

func (s *postgresqlActivityService) DeletePostgreSQLTarget(id uint) error {
	return s.targetRepository.DeletePostgreSQLTarget(id)
}

func (s *postgresqlActivityService) TestPostgreSQLTarget(request PostgreSQLTargetRequest) (string, error) {
	target, err := normalizeTargetRequest(request, 0)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(request.Password) == "" {
		return "", errors.New("password is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), postgresActivityLiveTimeout)
	defer cancel()
	return s.collector.TestConnection(ctx, target, request.Password)
}

func (s *postgresqlActivityService) GetLivePostgreSQLActivity(targetID uint, search string, limit int) (PostgreSQLActivityLiveResult, error) {
	targets, err := s.targetRepository.ListPostgreSQLTargets()
	if err != nil {
		return PostgreSQLActivityLiveResult{}, err
	}
	if targetID > 0 {
		filtered := targets[:0]
		for _, target := range targets {
			if target.ID == targetID {
				filtered = append(filtered, target)
			}
		}
		if len(filtered) == 0 {
			return PostgreSQLActivityLiveResult{}, gorm.ErrRecordNotFound
		}
		targets = filtered
	}

	activities := make([]model.PostgreSQLActivity, 0)
	for _, target := range targets {
		if !target.Enabled {
			continue
		}
		password, decryptErr := s.decrypt(target.PasswordCiphertext)
		if decryptErr != nil {
			s.recordTargetHealth(target, decryptErr)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), postgresActivityLiveTimeout)
		collection, collectErr := s.collector.Collect(ctx, target, password)
		cancel()
		if collectErr != nil {
			s.recordTargetHealth(target, collectErr)
			continue
		}
		s.recordTargetHealth(target, nil)
		activities = append(activities, collection.Activities...)
	}

	search = strings.ToLower(strings.TrimSpace(search))
	if search != "" {
		filtered := activities[:0]
		for _, activity := range activities {
			if activityContains(activity, search) {
				filtered = append(filtered, activity)
			}
		}
		activities = filtered
	}
	sort.SliceStable(activities, func(i, j int) bool {
		if activities[i].QueryDurationMS == activities[j].QueryDurationMS {
			return activities[i].TargetName < activities[j].TargetName
		}
		return activities[i].QueryDurationMS > activities[j].QueryDurationMS
	})
	if limit <= 0 {
		limit = postgresActivityDefaultLimit
	}
	if limit > 2000 {
		limit = 2000
	}
	if len(activities) > limit {
		activities = activities[:limit]
	}

	views, err := s.ListPostgreSQLTargets()
	if err != nil {
		return PostgreSQLActivityLiveResult{}, err
	}
	return PostgreSQLActivityLiveResult{Activities: activities, Targets: views}, nil
}

func (s *postgresqlActivityService) GetPostgreSQLActivityHistory(activityKey string, targetID uint, filter, search, state string, blockedOnly bool, minimumDurationMS int64, limit int) ([]model.PostgreSQLActivity, error) {
	duration, err := postgresActivityDuration(filter)
	if err != nil {
		return nil, err
	}
	return s.activityRepository.ListPostgreSQLActivityHistory(repository.PostgreSQLActivityHistoryFilter{
		ActivityKey:     activityKey,
		TargetID:        targetID,
		Cutoff:          time.Now().UTC().Add(-duration),
		Search:          search,
		State:           state,
		BlockedOnly:     blockedOnly,
		MinimumDuration: minimumDurationMS,
		Limit:           limit,
	})
}

func (s *postgresqlActivityService) CollectAndPersistPostgreSQLActivity() error {
	targets, err := s.targetRepository.ListPostgreSQLTargets()
	if err != nil {
		return err
	}
	activities := make([]model.PostgreSQLActivity, 0)
	for _, target := range targets {
		if !target.Enabled {
			continue
		}
		password, decryptErr := s.decrypt(target.PasswordCiphertext)
		if decryptErr != nil {
			s.recordTargetHealth(target, decryptErr)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), postgresActivityLiveTimeout)
		collection, collectErr := s.collector.Collect(ctx, target, password)
		cancel()
		if collectErr != nil {
			s.recordTargetHealth(target, collectErr)
			continue
		}
		s.recordTargetHealth(target, nil)
		for i := range collection.Activities {
			activity := &collection.Activities[i]
			activity.ID = fmt.Sprintf("%s:%d", activity.ActivityKey, activity.ObservedAt.UnixNano())
			activities = append(activities, *activity)
		}
	}
	return s.activityRepository.CreatePostgreSQLActivities(activities)
}

func (s *postgresqlActivityService) DeleteExpiredPostgreSQLActivity() (int64, error) {
	return s.activityRepository.DeleteOlderPostgreSQLActivities(time.Now().UTC().Add(-PostgreSQLActivityRetention))
}

func (s *postgresqlActivityService) recordTargetHealth(target model.PostgreSQLTarget, collectErr error) {
	checkedAt := time.Now().UTC()
	lastError := ""
	if collectErr != nil {
		lastError = collectErr.Error()
		if password, err := s.decrypt(target.PasswordCiphertext); err == nil && password != "" {
			lastError = strings.ReplaceAll(lastError, password, "[redacted]")
		}
	}
	_ = s.targetRepository.UpdatePostgreSQLTargetHealth(target.ID, &checkedAt, lastError)
}

func (s *postgresqlActivityService) encrypt(value string) (string, error) {
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

func (s *postgresqlActivityService) decrypt(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", errors.New("encrypted PostgreSQL password is missing")
	}
	block, err := s.cipherBlock()
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.RawStdEncoding.DecodeString(value)
	if err != nil {
		return "", errors.New("encrypted PostgreSQL password is invalid")
	}
	if len(ciphertext) < block.NonceSize() {
		return "", errors.New("encrypted PostgreSQL password is invalid")
	}
	nonce, ciphertext := ciphertext[:block.NonceSize()], ciphertext[block.NonceSize():]
	plaintext, err := block.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("unable to decrypt PostgreSQL password")
	}
	return string(plaintext), nil
}

func (s *postgresqlActivityService) cipherBlock() (cipher.AEAD, error) {
	if strings.TrimSpace(s.secret) == "" {
		return nil, errors.New("app.jwt_secret is required to protect PostgreSQL credentials")
	}
	key := sha256.Sum256([]byte(s.secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func normalizeTargetRequest(request PostgreSQLTargetRequest, id uint) (model.PostgreSQLTarget, error) {
	name := strings.TrimSpace(request.Name)
	host := strings.TrimSpace(request.Host)
	username := strings.TrimSpace(request.Username)
	databaseName := strings.TrimSpace(request.Database)
	if name == "" || host == "" || username == "" {
		return model.PostgreSQLTarget{}, errors.New("name, host, and username are required")
	}
	port := request.Port
	if port == 0 {
		port = 5432
	}
	if port < 1 || port > 65535 {
		return model.PostgreSQLTarget{}, errors.New("port must be between 1 and 65535")
	}
	if databaseName == "" {
		databaseName = "postgres"
	}
	sslMode := strings.ToLower(strings.TrimSpace(request.SSLMode))
	if sslMode == "" {
		sslMode = "prefer"
	}
	switch sslMode {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
	default:
		return model.PostgreSQLTarget{}, errors.New("ssl_mode must be disable, allow, prefer, require, verify-ca, or verify-full")
	}
	enabled := true
	if request.Enabled != nil {
		enabled = *request.Enabled
	}
	return model.PostgreSQLTarget{
		ID:       id,
		Name:     name,
		Host:     host,
		Port:     port,
		Database: databaseName,
		Username: username,
		SSLMode:  sslMode,
		Enabled:  enabled,
	}, nil
}

func targetView(target model.PostgreSQLTarget) PostgreSQLTargetView {
	return PostgreSQLTargetView{
		ID:                 target.ID,
		Name:               target.Name,
		Host:               target.Host,
		Port:               target.Port,
		Database:           target.Database,
		Username:           target.Username,
		SSLMode:            target.SSLMode,
		Enabled:            target.Enabled,
		PasswordConfigured: target.PasswordCiphertext != "",
		LastCheckedAt:      target.LastCheckedAt,
		LastError:          target.LastError,
		CreatedAt:          target.CreatedAt,
		UpdatedAt:          target.UpdatedAt,
	}
}

func activityContains(activity model.PostgreSQLActivity, search string) bool {
	for _, value := range []string{
		activity.TargetName,
		activity.DatabaseName,
		activity.Username,
		activity.ApplicationName,
		activity.State,
		activity.Query,
		fmt.Sprint(activity.PID),
	} {
		if strings.Contains(strings.ToLower(value), search) {
			return true
		}
	}
	return false
}

func postgresActivityDuration(filter string) (time.Duration, error) {
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
		return PostgreSQLActivityRetention, nil
	default:
		return 0, errors.New("invalid filter format")
	}
}
