package concat_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const concatImage = "ghcr.io/veloxpack/ffmpeg:8.0-concat"
const splitImage = "ghcr.io/veloxpack/ffmpeg:8.0-split"

func TestConcat_MP4_Files(t *testing.T) {
	// Given: Split a video into segments first
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	ctx := context.Background()

	// Step 1: Split video into 3 segments using ffmpeg-split
	splitCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-c", "copy",
		"-f", "segment",
		"-segment_time", "5",
		"-reset_timestamps", "1",
		"/output/part-%03d.mp4",
	}

	splitContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: splitImage,
			Cmd:   splitCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.AutoRemove = true
				hc.Mounts = []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: fmt.Sprintf("%s/", outputPath),
						Target: "/output/",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer splitContainer.Terminate(ctx)

	printContainerLogs(t, ctx, splitContainer)

	// Verify segments exist
	segments, err := filepath.Glob(filepath.Join(outputPath, "part-*.mp4"))
	require.NoError(t, err)
	require.NotEmpty(t, segments, "Should have created segments")
	t.Logf("Created %d segments", len(segments))

	// Step 2: Create concat list file
	listPath := filepath.Join(outputPath, "list.txt")
	listFile, err := os.Create(listPath)
	require.NoError(t, err)

	for _, segment := range segments {
		_, err = fmt.Fprintf(listFile, "file '%s'\n", filepath.Base(segment))
		require.NoError(t, err)
	}
	listFile.Close()

	// Step 3: Concatenate using ffmpeg-concat
	concatCmd := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", "/workspace/list.txt",
		"-c", "copy",
		"/workspace/concatenated.mp4",
	}

	concatContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: concatImage,
			Cmd:   concatCmd,
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.AutoRemove = true
				hc.Mounts = []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: fmt.Sprintf("%s/", outputPath),
						Target: "/workspace/",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer concatContainer.Terminate(ctx)

	printContainerLogs(t, ctx, concatContainer)

	// Then: Verify concatenated file exists
	concatPath := filepath.Join(outputPath, "concatenated.mp4")
	verifyFileExists(t, concatPath)
	verifyFileSize(t, concatPath, 100*1024) // At least 100KB
}

func TestConcat_WithDurationMetadata(t *testing.T) {
	// Given: Split video into segments
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	ctx := context.Background()

	// Step 1: Split video
	splitCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-c", "copy",
		"-f", "segment",
		"-segment_time", "3",
		"-reset_timestamps", "1",
		"/output/clip-%03d.mp4",
	}

	splitContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: splitImage,
			Cmd:   splitCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.AutoRemove = true
				hc.Mounts = []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: fmt.Sprintf("%s/", outputPath),
						Target: "/output/",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer splitContainer.Terminate(ctx)

	// Step 2: Create list with duration metadata
	listPath := filepath.Join(outputPath, "list.txt")
	listFile, err := os.Create(listPath)
	require.NoError(t, err)

	segments, _ := filepath.Glob(filepath.Join(outputPath, "clip-*.mp4"))
	for _, segment := range segments {
		fmt.Fprintf(listFile, "file '%s'\n", filepath.Base(segment))
		fmt.Fprintf(listFile, "duration 15.0\n")
	}
	listFile.Close()

	// Step 3: Concatenate
	concatCmd := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", "/workspace/list.txt",
		"-c", "copy",
		"/workspace/output.mp4",
	}

	concatContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: concatImage,
			Cmd:   concatCmd,
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.AutoRemove = true
				hc.Mounts = []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: fmt.Sprintf("%s/", outputPath),
						Target: "/workspace/",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer concatContainer.Terminate(ctx)

	printContainerLogs(t, ctx, concatContainer)

	// Then: Verify output
	concatPath := filepath.Join(outputPath, "output.mp4")
	verifyFileExists(t, concatPath)
}

func TestConcat_WithTrimPoints(t *testing.T) {
	// Given: Split video into segments
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	ctx := context.Background()

	// Step 1: Split video
	splitCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-c", "copy",
		"-f", "segment",
		"-segment_time", "5",
		"-reset_timestamps", "1",
		"/output/segment-%03d.mp4",
	}

	splitContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: splitImage,
			Cmd:   splitCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.AutoRemove = true
				hc.Mounts = []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: fmt.Sprintf("%s/", outputPath),
						Target: "/output/",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer splitContainer.Terminate(ctx)

	// Step 2: Create list with trim points (inpoint/outpoint)
	listPath := filepath.Join(outputPath, "trimlist.txt")
	listFile, err := os.Create(listPath)
	require.NoError(t, err)

	segments, _ := filepath.Glob(filepath.Join(outputPath, "segment-*.mp4"))
	if len(segments) > 0 {
		// Trim first segment from 5s to 15s
		fmt.Fprintf(listFile, "file '%s'\n", filepath.Base(segments[0]))
		fmt.Fprintf(listFile, "inpoint 5.0\n")
		fmt.Fprintf(listFile, "outpoint 15.0\n")
	}
	if len(segments) > 1 {
		// Use second segment from 0s to 10s
		fmt.Fprintf(listFile, "file '%s'\n", filepath.Base(segments[1]))
		fmt.Fprintf(listFile, "inpoint 0.0\n")
		fmt.Fprintf(listFile, "outpoint 10.0\n")
	}
	listFile.Close()

	// Step 3: Concatenate with trim points
	concatCmd := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", "/workspace/trimlist.txt",
		"-c", "copy",
		"/workspace/trimmed.mp4",
	}

	concatContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: concatImage,
			Cmd:   concatCmd,
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.AutoRemove = true
				hc.Mounts = []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: fmt.Sprintf("%s/", outputPath),
						Target: "/workspace/",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer concatContainer.Terminate(ctx)

	printContainerLogs(t, ctx, concatContainer)

	// Then: Verify trimmed output
	trimmedPath := filepath.Join(outputPath, "trimmed.mp4")
	verifyFileExists(t, trimmedPath)
}

func TestConcat_WebM_Files(t *testing.T) {
	// Given: We need WebM files - convert MP4 to WebM first
	// Note: This test assumes sample.mp4 exists and we can convert it
	// For a real test, you'd need ffmpeg image or pre-existing WebM files

	t.Skip("WebM test requires pre-existing WebM files or ffmpeg image for conversion")

	// Implementation would be similar to MP4 test but with WebM format
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

func verifyFileExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	require.NoError(t, err, "File should exist: %s", path)
}

func verifyFileSize(t *testing.T, path string, minSize int64) {
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, info.Size(), minSize, "File should be at least %d bytes", minSize)
}

func printContainerLogs(t *testing.T, ctx context.Context, c testcontainers.Container) {
	if log, err := c.Logs(ctx); err == nil {
		reader := bufio.NewReader(log)
		data, _ := io.ReadAll(reader)
		if len(data) > 0 {
			t.Log("Container logs:", string(data))
		}
	}
}
