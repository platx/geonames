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
	ErrInvalidType           = errors.New("invalid type")
)

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient httpDoer
	baseURL    string
}

type Option func(*Client)

func WithHTTPClient(httpClient httpDoer) Option {
	return func(client *Client) {
		client.httpClient = httpClient
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
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func (c *Client) geoNames(ctx context.Context, fileName string) (Iterator[GeoName], error) {
	res, err := c.downloadAndParseZIPFile(ctx, fileName)
	if err != nil {
		return nil, err
	}

	return withUnmarshalRows[GeoName](res), nil
}

func (c *Client) downloadAndParseFile(ctx context.Context, fileName string) (Iterator[[]string], error) {
	file, err := c.downloadFile(ctx, fileName)
	if err != nil {
		return nil, fmt.Errorf("download file => %w", err)
	}

	defer func() {
		_ = os.Remove(file.Name())
	}()

	fileReader, err := os.Open(file.Name())
	if err != nil {
		return nil, fmt.Errorf("open file => %w", err)
	}

	return c.parseTSV(ctx, fileReader), nil
}

func (c *Client) downloadAndParseZIPFile(ctx context.Context, fileName string) (Iterator[[]string], error) {
	file, err := c.downloadFile(ctx, fileName)
	if err != nil {
		return nil, fmt.Errorf("download file => %w", err)
	}

	defer func() {
		_ = os.Remove(file.Name())
	}()

	zipFile := strings.Replace(fileName, ".zip", ".txt", 1)

	return c.parseZIPFile(ctx, file, zipFile)
}

func (c *Client) downloadFile(ctx context.Context, fileName string) (*os.File, error) {
	req, err := c.createHTTPRequest(ctx, fileName)
	if err != nil {
		return nil, fmt.Errorf("create http request => %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client do => %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, res.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "*"+fileName)
	if err != nil {
		return nil, fmt.Errorf("create temp file => %w", err)
	}

	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return nil, fmt.Errorf("copy file content => %w", err)
	}

	_ = tmpFile.Close()

	return tmpFile, nil
}

func (c *Client) parseZIPFile(ctx context.Context, zipArchive *os.File, fileName string) (Iterator[[]string], error) {
	file, err := zip.OpenReader(zipArchive.Name())
	if err != nil {
		return nil, fmt.Errorf("open zip archive => %w", err)
	}

	var targetFile *zip.File

	for _, f := range file.File {
		if f.Name == fileName {
			targetFile = f

			break
		}
	}

	if targetFile == nil {
		return nil, ErrFileNotFoundInArchive
	}

	fileReader, err := targetFile.Open()
	if err != nil {
		return nil, fmt.Errorf("open file from archive => %w", err)
	}

	return withClose(c.parseTSV(ctx, fileReader), file), nil
}

func (c *Client) parseTSV(ctx context.Context, file io.ReadCloser) Iterator[[]string] {
	return func(yield func([]string, error) bool) {
		defer func() {
			_ = file.Close()
		}()

		scanner := bufio.NewScanner(file)

		for {
			select {
			case <-ctx.Done():
				yield(nil, ctx.Err())

				return
			default:
			}

			if !scanner.Scan() {
				break
			}

			text := scanner.Text()
			if len(text) == 0 {
				continue
			}

			if strings.HasPrefix(text, commentPrefix) {
				continue
			}

			if !yield(parseLine(text, columnSeparator), nil) {
				return
			}
		}

		if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
			yield(nil, err)
		}
	}
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

func withClose(rows Iterator[[]string], closer io.Closer) Iterator[[]string] {
	return func(yield func([]string, error) bool) {
		defer func() {
			_ = closer.Close()
		}()

		for rowRes, rowErr := range rows {
			if !yield(rowRes, rowErr) {
				return
			}
		}
	}
}

func withSkipHeader(rows Iterator[[]string]) Iterator[[]string] {
	return func(yield func([]string, error) bool) {
		header := true

		for res, err := range rows {
			if header {
				header = false

				continue
			}

			if !yield(res, err) {
				return
			}
		}
	}
}

func withUnmarshalRows[T any](rows Iterator[[]string]) Iterator[T] {
	return func(yield func(T, error) bool) {
		for res, err := range rows {
			ptr := new(T)

			if err != nil {
				yield(*ptr, err)

				return
			}

			casted, ok := any(ptr).(interface{ UnmarshalRow(row []string) error })
			if !ok {
				yield(*ptr, fmt.Errorf("%w => type %T does not implement UnmarshalRow", ErrInvalidType, ptr))

				return
			}

			err = casted.UnmarshalRow(res)
			if !yield(*ptr, err) {
				return
			}
		}
	}
}
