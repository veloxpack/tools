package thumbnail_test

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

const thumbnailImage = "ghcr.io/veloxpack/ffmpeg:8.0-thumbnail"

func TestThumbnail_PNG_Generation(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Generate PNG thumbnail at 5 seconds
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-ss", "5",
		"-vframes", "1",
		"/output/thumbnail.png",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: thumbnailImage,
			Cmd:   containerCmd,
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
	defer c.Terminate(ctx)

	printContainerLogs(t, ctx, c)

	// Then: Verify PNG file exists
	thumbnailPath := filepath.Join(outputPath, "thumbnail.png")
	verifyFileExists(t, thumbnailPath)
	verifyFileSize(t, thumbnailPath, 1024) // At least 1KB
}

func TestThumbnail_JPEG_Generation(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Generate JPEG thumbnail at 10 seconds
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-ss", "10",
		"-vframes", "1",
		"/output/thumbnail.jpg",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: thumbnailImage,
			Cmd:   containerCmd,
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
	defer c.Terminate(ctx)

	printContainerLogs(t, ctx, c)

	// Then: Verify JPEG file exists
	thumbnailPath := filepath.Join(outputPath, "thumbnail.jpg")
	verifyFileExists(t, thumbnailPath)
	verifyFileSize(t, thumbnailPath, 1024) // At least 1KB
}

func TestThumbnail_Storyboard_5x5(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Create 5x5 storyboard
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "fps=1/10,scale=160:90,tile=5x5",
		"/output/storyboard.jpg",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: thumbnailImage,
			Cmd:   containerCmd,
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
	defer c.Terminate(ctx)

	printContainerLogs(t, ctx, c)

	// Then: Verify storyboard exists
	storyboardPath := filepath.Join(outputPath, "storyboard.jpg")
	verifyFileExists(t, storyboardPath)
	verifyFileSize(t, storyboardPath, 10*1024) // At least 10KB for grid
}

func TestThumbnail_BestFrame_Using_ThumbnailFilter(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Use thumbnail filter to find best frame
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "thumbnail",
		"-frames:v", "1",
		"/output/best-frame.jpg",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: thumbnailImage,
			Cmd:   containerCmd,
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
	defer c.Terminate(ctx)

	printContainerLogs(t, ctx, c)

	// Then: Verify best frame exists
	bestFramePath := filepath.Join(outputPath, "best-frame.jpg")
	verifyFileExists(t, bestFramePath)
	verifyFileSize(t, bestFramePath, 1024)
}

func TestThumbnail_Multiple_At_Intervals(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Generate thumbnails every 60 seconds (limited to 10s for testing)
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "fps=1/60",
		"/output/thumb-%04d.jpg",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: thumbnailImage,
			Cmd:   containerCmd,
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
	defer c.Terminate(ctx)

	printContainerLogs(t, ctx, c)

	// Then: Verify at least one thumbnail exists
	files, err := filepath.Glob(filepath.Join(outputPath, "thumb-*.jpg"))
	require.NoError(t, err)
	assert.NotEmpty(t, files, "Should have generated at least one thumbnail")
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
