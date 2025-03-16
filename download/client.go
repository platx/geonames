package download

import (
	"archive/zip"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultBaseURL        = "https://download.geonames.org/export/dump"
	defaultRequestTimeout = 10 * time.Minute
	columnSeparator       = "\t"
	commentPrefix         = "#"
)

var (
	ErrFileNotFoundInArchive = errors.New("file not found in archive")
	ErrUnexpectedStatusCode  = errors.New("unexpected status code")
)

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}

type Client struct {
	httpClient httpDoer
	baseURL    string
	logger     logger
}

type Option func(*Client)

func WithHTTPClient(httpClient httpDoer) Option {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

func WithLogger(logger logger) Option {
	return func(client *Client) {
		client.logger = logger
	}
}

func WithBaseURL(baseURL string) Option {
	return func(client *Client) {
		baseURL = strings.TrimSpace(baseURL)
		if idx := strings.Index(baseURL, "?"); idx != -1 {
			baseURL = baseURL[:idx]
		}

		client.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

func NewClient(opts ...Option) *Client {
	res := &Client{
		httpClient: &http.Client{
			Transport:     http.DefaultTransport,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       defaultRequestTimeout,
		},
		baseURL: defaultBaseURL,
		logger:  nopLogger{},
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func (c *Client) geoNames(ctx context.Context, fileName string, callback func(parsed GeoName) error) error {
	return c.downloadAndParseZIPFile(ctx, fileName, func(row []string) error {
		var parsed GeoName

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}

func (c *Client) downloadAndParseFile(ctx context.Context, fileName string, callback func(row []string) error) error {
	file, err := c.downloadFile(ctx, fileName)
	if err != nil {
		return fmt.Errorf("download file => %w", err)
	}

	defer func() {
		_ = os.Remove(file.Name())

		c.logger.Debug("removed temp file", "file", file.Name())
	}()

	fileReader, err := os.Open(file.Name())
	if err != nil {
		return fmt.Errorf("open file => %w", err)
	}

	defer func() {
		_ = fileReader.Close()
	}()

	if err = c.parseFile(ctx, fileReader, callback); err != nil {
		return fmt.Errorf("parse file => %w", err)
	}

	return nil
}

func (c *Client) downloadAndParseZIPFile(
	ctx context.Context,
	fileName string,
	callback func(row []string) error,
) error {
	file, err := c.downloadFile(ctx, fileName)
	if err != nil {
		return fmt.Errorf("download file => %w", err)
	}

	defer func() {
		_ = os.Remove(file.Name())

		c.logger.Debug("removed temp file", "file", file.Name())
	}()

	zipFile := strings.Replace(fileName, ".zip", ".txt", 1)
	if err = c.parseZIPFile(ctx, file, zipFile, callback); err != nil {
		return fmt.Errorf("parse file %q in archive => %w", zipFile, err)
	}

	return nil
}

func (c *Client) downloadFile(ctx context.Context, fileName string) (*os.File, error) {
	req, err := c.createHTTPRequest(ctx, fileName)
	if err != nil {
		return nil, fmt.Errorf("create http request => %w", err)
	}

	c.logger.Debug("try download remote file", "url", req.URL.String())

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client do => %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, res.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "*"+fileName)
	if err != nil {
		return nil, fmt.Errorf("create temp file => %w", err)
	}

	c.logger.Debug("created temp file", "file", tmpFile.Name())

	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return nil, fmt.Errorf("copy file content => %w", err)
	}

	_ = tmpFile.Close()

	c.logger.Info("remote file downloaded", "url", req.URL.String())

	return tmpFile, nil
}

func (c *Client) parseZIPFile(
	ctx context.Context,
	zipArchive *os.File,
	fileName string,
	callback func(row []string) error,
) error {
	file, err := zip.OpenReader(zipArchive.Name())
	if err != nil {
		return fmt.Errorf("open zip archive => %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	var targetFile *zip.File

	for _, f := range file.File {
		if f.Name == fileName {
			c.logger.Debug("found target file in archive", "file", fileName)

			targetFile = f

			break
		}
	}

	if targetFile == nil {
		return ErrFileNotFoundInArchive
	}

	fileReader, err := targetFile.Open()
	if err != nil {
		return fmt.Errorf("open file from archive => %w", err)
	}

	defer func() {
		_ = fileReader.Close()
	}()

	if err = c.parseFile(ctx, fileReader, callback); err != nil {
		return fmt.Errorf("parse file => %w", err)
	}

	return nil
}

func (c *Client) parseFile(ctx context.Context, file io.Reader, callback func(row []string) error) error {
	var (
		text string
		line int
	)

	c.logger.Debug("try parse file")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line++

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			text = scanner.Text()
			if len(text) == 0 {
				c.logger.Debug("skip empty line", "line", line)

				continue
			}

			if strings.HasPrefix(text, commentPrefix) {
				c.logger.Debug("skip comment line", "line", line, "text", text)

				continue
			}

			if err := callback(parseLine(text, columnSeparator)); err != nil {
				c.logger.Warn("failed parse line", "line", line, "text", text, "error", err)

				continue
			}

			c.logger.Debug("parsed line", "line", line, "text", text)
		}
	}

	c.logger.Info("parsed file", "lines", line)

	return scanner.Err()
}

func (c *Client) createHTTPRequest(ctx context.Context, fileName string) (*http.Request, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(fileName), nil)
	if err != nil {
		return nil, err
	}

	return httpReq, nil
}

func (c *Client) url(path string) string {
	return c.baseURL + "/" + path
}

func parseLine(line string, separator string) []string {
	return strings.Split(line, separator)
}

type nopLogger struct{}

func (nopLogger) Debug(string, ...any) {}
func (nopLogger) Info(string, ...any)  {}
func (nopLogger) Warn(string, ...any)  {}
