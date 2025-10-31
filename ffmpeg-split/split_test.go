package split_test

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

const splitImage = "ghcr.io/veloxpack/ffmpeg:8.0-split"

func TestSplit_TimeBased_StreamCopy(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Split first 10 seconds
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-c", "copy",
		"/output/first-10s.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: splitImage,
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

	// Then: Verify split file exists
	splitPath := filepath.Join(outputPath, "first-10s.mp4")
	verifyFileExists(t, splitPath)
	verifyFileSize(t, splitPath, 100*1024) // At least 100KB
}

func TestSplit_Segments_ByDuration(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Split into 5-second segments (limited to 10s for testing)
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-c", "copy",
		"-f", "segment",
		"-segment_time", "5",
		"-reset_timestamps", "1",
		"/output/part-%03d.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: splitImage,
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

	// Then: Verify at least one segment exists
	files, err := filepath.Glob(filepath.Join(outputPath, "part-*.mp4"))
	require.NoError(t, err)
	assert.NotEmpty(t, files, "Should have generated at least one segment")
}

func TestSplit_SceneDetection(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Split on scene changes
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "select='gt(scene,0.4)'",
		"-vsync", "vfr",
		"/output/scene_%03d.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: splitImage,
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

	// Then: Verify scene files exist
	files, err := filepath.Glob(filepath.Join(outputPath, "scene_*.mp4"))
	require.NoError(t, err)
	assert.NotEmpty(t, files, "Should have generated at least one scene file")
}

func TestSplit_SceneDetection_WithMetadata(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Split on scene changes with metadata export
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "select='gt(scene,0.4)',metadata=print:file=/output/scenes.txt",
		"-vsync", "vfr",
		"/output/scene_%03d.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: splitImage,
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

	// Then: Verify metadata file exists
	metadataPath := filepath.Join(outputPath, "scenes.txt")
	verifyFileExists(t, metadataPath)

	// Verify scene files exist
	files, err := filepath.Glob(filepath.Join(outputPath, "scene_*.mp4"))
	require.NoError(t, err)
	assert.NotEmpty(t, files, "Should have generated at least one scene file")
}

func TestSplit_SceneDetection_CustomThreshold(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Split with lower threshold (more sensitive)
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "select='gt(scene,0.3)'",
		"-vsync", "vfr",
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"/output/scene_%03d.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: splitImage,
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

	// Then: Verify scene files exist
	files, err := filepath.Glob(filepath.Join(outputPath, "scene_*.mp4"))
	require.NoError(t, err)
	assert.NotEmpty(t, files, "Should have generated at least one scene file with custom threshold")
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
