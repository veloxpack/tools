package ffprobe_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const ffprobeImage = "ghcr.io/veloxpack/ffprobe:latest"

type FFProbeOutput struct {
	Format  FFProbeFormat   `json:"format"`
	Streams []FFProbeStream `json:"streams"`
}

type FFProbeFormat struct {
	Filename   string            `json:"filename"`
	FormatName string            `json:"format_name"`
	Duration   string            `json:"duration"`
	Size       string            `json:"size"`
	BitRate    string            `json:"bit_rate"`
	Tags       map[string]string `json:"tags"`
}

type FFProbeStream struct {
	Index      int               `json:"index"`
	CodecName  string            `json:"codec_name"`
	CodecType  string            `json:"codec_type"`
	Width      int               `json:"width,omitempty"`
	Height     int               `json:"height,omitempty"`
	SampleRate string            `json:"sample_rate,omitempty"`
	Channels   int               `json:"channels,omitempty"`
	Duration   string            `json:"duration"`
	BitRate    string            `json:"bit_rate"`
	Tags       map[string]string `json:"tags"`
}

func TestFFProbe_BasicInfo(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Run ffprobe
	ctx := context.Background()
	containerCmd := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"/input/sample.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffprobeImage,
			Cmd:   containerCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer c.Terminate(ctx)

	// Then: Verify JSON output
	logs, err := c.Logs(ctx)
	require.NoError(t, err)

	logData, err := io.ReadAll(logs)
	require.NoError(t, err)

	// Clean up any potential Docker stream headers or binary data
	logData = cleanLogData(logData)

	var output FFProbeOutput
	err = json.Unmarshal(logData, &output)
	require.NoError(t, err)

	// Verify format info
	assert.NotEmpty(t, output.Format.Filename)
	assert.NotEmpty(t, output.Format.FormatName)
	assert.NotEmpty(t, output.Format.Duration)

	// Verify streams
	assert.NotEmpty(t, output.Streams)
}

func TestFFProbe_JSONOutput(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	// When: Run ffprobe with JSON output
	ctx := context.Background()
	containerCmd := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"/input/sample.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffprobeImage,
			Cmd:   containerCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer c.Terminate(ctx)

	// Then: Parse and verify JSON
	logs, err := c.Logs(ctx)
	require.NoError(t, err)

	logData, err := io.ReadAll(logs)
	require.NoError(t, err)

	// Clean up any potential Docker stream headers or binary data
	logData = cleanLogData(logData)

	var output FFProbeOutput
	err = json.Unmarshal(logData, &output)
	require.NoError(t, err)

	assert.NotEmpty(t, output.Format.FormatName)
	assert.Contains(t, output.Format.FormatName, "mov")
}

func TestFFProbe_StreamInfo(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	// When: Run ffprobe for stream info
	ctx := context.Background()
	containerCmd := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		"/input/sample.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffprobeImage,
			Cmd:   containerCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer c.Terminate(ctx)

	// Then: Verify stream information
	logs, err := c.Logs(ctx)
	require.NoError(t, err)

	logData, err := io.ReadAll(logs)
	require.NoError(t, err)

	// Clean up any potential Docker stream headers or binary data
	logData = cleanLogData(logData)

	var output FFProbeOutput
	err = json.Unmarshal(logData, &output)
	require.NoError(t, err)

	assert.NotEmpty(t, output.Streams)

	// Verify we have video stream
	hasVideo := false
	for _, stream := range output.Streams {
		if stream.CodecType == "video" {
			hasVideo = true
			assert.NotEmpty(t, stream.CodecName)
			assert.Greater(t, stream.Width, 0)
			assert.Greater(t, stream.Height, 0)
		}
	}
	assert.True(t, hasVideo, "Should have at least one video stream")
}

// Helper functions
func createTempDir(t *testing.T) string {
	outputPath, err := filepath.Abs(filepath.Join("..", "testdata", strings.ReplaceAll(uuid.NewString(), "-", "")))
	require.NoError(t, err)

	err = os.MkdirAll(outputPath, 0755)
	require.NoError(t, err)

	return outputPath
}

func cleanupFiles(t *testing.T, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Logf("failed to remove output directory: %s", path)
	}
}

// cleanLogData strips Docker stream headers from log data
func cleanLogData(data []byte) []byte {
	// Docker stream format: 8-byte header (1 byte stream type + 3 bytes padding + 4 bytes size)
	// Scan through the data and remove all Docker headers
	var cleaned bytes.Buffer
	i := 0

	for i < len(data) {
		// Look for potential Docker header pattern: 01 00 00 00 00 00 00 XX
		// Stream type 01 (stdout), three zero bytes, then 4-byte size
		if i+8 <= len(data) &&
			data[i] == 0x01 &&
			data[i+1] == 0x00 &&
			data[i+2] == 0x00 &&
			data[i+3] == 0x00 {

			// This looks like a Docker stdout header, skip these 8 bytes
			i += 8
			continue
		}

		// Also check for stderr (02) or stdin (00) headers
		if i+8 <= len(data) &&
			(data[i] == 0x00 || data[i] == 0x02) &&
			data[i+1] == 0x00 &&
			data[i+2] == 0x00 &&
			data[i+3] == 0x00 {

			// This looks like a Docker stream header, skip these 8 bytes
			i += 8
			continue
		}

		// Not a header, keep this byte
		cleaned.WriteByte(data[i])
		i++
	}

	return cleaned.Bytes()
}
