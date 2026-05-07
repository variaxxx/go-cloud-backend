package app_load

import (
	"bytes"
	core_logger "cloud/internal/core/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

const sampleTextFileContent = "load test"

type App struct {
	logger *core_logger.Logger
	config Config
	client *http.Client
}

type authResponse struct {
	AccessToken string `json:"access_token"`
}

func New() (*App, error) {
	loggerConfig, err := core_logger.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger config: %w", err)
	}

	config, err := NewConfig()
	if err != nil {
		return nil, err
	}

	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	return &App{
		logger: logger,
		config: config,
		client: &http.Client{
			Timeout: config.RequestTimeout,
		},
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info(
		"Load test started",
		zap.String("api_base_url", a.config.APIBaseURL),
		zap.Int("upload_count", a.config.UploadCount),
		zap.Int("concurrency", a.config.Concurrency),
	)

	token, err := a.ensureAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("ensure access token: %w", err)
	}

	startedAt := time.Now()
	var succeeded atomic.Int64
	var failed atomic.Int64

	jobs := make(chan int)
	errCh := make(chan error, a.config.UploadCount)
	var wg sync.WaitGroup

	for workerID := 0; workerID < a.config.Concurrency; workerID++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			for jobID := range jobs {
				filename := fmt.Sprintf(
					"load-%d-%d-%d.txt",
					workerID,
					jobID,
					time.Now().UnixNano(),
				)

				if err := a.uploadFile(ctx, token, filename, []byte(sampleTextFileContent)); err != nil {
					failed.Add(1)
					errCh <- err
					continue
				}

				succeeded.Add(1)
			}
		}(workerID)
	}

	for i := 0; i < a.config.UploadCount; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	close(errCh)

	for err := range errCh {
		a.logger.Error("Upload request failed", zap.Error(err))
	}

	a.logger.Info(
		"Load test finished",
		zap.Int64("succeeded", succeeded.Load()),
		zap.Int64("failed", failed.Load()),
		zap.Duration("duration", time.Since(startedAt)),
	)

	if succeeded.Load() == 0 && failed.Load() > 0 {
		return fmt.Errorf("all uploads failed")
	}

	return nil
}

func (a *App) Close() {
	a.logger.Close()
}

func (a *App) ensureAccessToken(ctx context.Context) (string, error) {
	token, statusCode, err := a.register(ctx)
	if err == nil {
		return token, nil
	}

	if statusCode != http.StatusConflict {
		return "", fmt.Errorf("register load user: %w", err)
	}

	token, _, err = a.login(ctx)
	if err != nil {
		return "", fmt.Errorf("login load user: %w", err)
	}

	return token, nil
}

func (a *App) register(ctx context.Context) (string, int, error) {
	return a.authenticate(
		ctx,
		http.MethodPost,
		"/auth/register",
	)
}

func (a *App) login(ctx context.Context) (string, int, error) {
	return a.authenticate(
		ctx,
		http.MethodPost,
		"/auth/login",
	)
}

func (a *App) authenticate(
	ctx context.Context,
	method string,
	path string,
) (string, int, error) {
	payload, err := json.Marshal(map[string]string{
		"username": a.config.Username,
		"password": a.config.Password,
	})
	if err != nil {
		return "", 0, fmt.Errorf("marshal auth payload: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		method,
		strings.TrimRight(a.config.APIBaseURL, "/")+path,
		bytes.NewReader(payload),
	)
	if err != nil {
		return "", 0, fmt.Errorf("build auth request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := a.client.Do(request)
	if err != nil {
		return "", 0, fmt.Errorf("send auth request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", response.StatusCode, fmt.Errorf(
			"unexpected auth status %d: %s",
			response.StatusCode,
			strings.TrimSpace(string(body)),
		)
	}

	var auth authResponse
	if err := json.NewDecoder(response.Body).Decode(&auth); err != nil {
		return "", response.StatusCode, fmt.Errorf("decode auth response: %w", err)
	}

	if auth.AccessToken == "" {
		return "", response.StatusCode, fmt.Errorf("empty access token in auth response")
	}

	return auth.AccessToken, response.StatusCode, nil
}

func (a *App) uploadFile(
	ctx context.Context,
	token string,
	filename string,
	content []byte,
) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("create multipart file: %w", err)
	}

	if _, err := part.Write(content); err != nil {
		return fmt.Errorf("write multipart file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		strings.TrimRight(a.config.APIBaseURL, "/")+"/files/upload",
		&body,
	)
	if err != nil {
		return fmt.Errorf("build upload request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := a.client.Do(request)
	if err != nil {
		return fmt.Errorf("send upload request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf(
			"unexpected upload status %d for %s: %s",
			response.StatusCode,
			filename,
			strings.TrimSpace(string(body)),
		)
	}

	return nil
}
