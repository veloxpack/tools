package ffmpeg_test

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

const ffmpegImage = "ghcr.io/veloxpack/ffmpeg:8.0-lite"

func TestFFmpeg_Transcode_1080p_H264(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Transcode to 1080p H.264
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "scale=1920:1080",
		"-c:v", "libx264",
		"-preset", "medium",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		"/output/output_1080p.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffmpegImage,
			Cmd:   containerCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
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

	// Then: Verify 1080p output exists
	outputFile := filepath.Join(outputPath, "output_1080p.mp4")
	verifyFileExists(t, outputFile)
	verifyFileSize(t, outputFile, 100*1024)
}

func TestFFmpeg_Transcode_720p_H264(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Transcode to 720p H.264
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "scale=1280:720",
		"-c:v", "libx264",
		"-preset", "medium",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		"/output/output_720p.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffmpegImage,
			Cmd:   containerCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
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

	// Then: Verify 720p output exists
	outputFile := filepath.Join(outputPath, "output_720p.mp4")
	verifyFileExists(t, outputFile)
	verifyFileSize(t, outputFile, 50*1024)
}

func TestFFmpeg_Transcode_480p_H264(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Transcode to 480p H.264
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "scale=854:480",
		"-c:v", "libx264",
		"-preset", "medium",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "96k",
		"/output/output_480p.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffmpegImage,
			Cmd:   containerCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
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

	// Then: Verify 480p output exists
	outputFile := filepath.Join(outputPath, "output_480p.mp4")
	verifyFileExists(t, outputFile)
	verifyFileSize(t, outputFile, 30*1024)
}

func TestFFmpeg_Transcode_VP9_WebM(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Transcode to VP9/WebM
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-c:v", "libvpx-vp9",
		"-crf", "30",
		"-b:v", "0",
		"-c:a", "libopus",
		"-b:a", "128k",
		"/output/output.webm",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffmpegImage,
			Cmd:   containerCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
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

	// Then: Verify WebM output exists
	outputFile := filepath.Join(outputPath, "output.webm")
	verifyFileExists(t, outputFile)
	verifyFileSize(t, outputFile, 50*1024)
}

func TestFFmpeg_Scale_CustomResolution(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Scale to custom resolution (640x360)
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "scale=640:360",
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"-c:a", "copy",
		"/output/scaled_640x360.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffmpegImage,
			Cmd:   containerCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
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

	// Then: Verify scaled output exists
	outputFile := filepath.Join(outputPath, "scaled_640x360.mp4")
	verifyFileExists(t, outputFile)
	verifyFileSize(t, outputFile, 20*1024)
}

func TestFFmpeg_Audio_AAC_Transcode(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Extract and transcode audio to AAC
	ctx := context.Background()
	containerCmd := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vn",
		"-c:a", "aac",
		"-b:a", "192k",
		"/output/audio.mp4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffmpegImage,
			Cmd:   containerCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
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

	// Then: Verify audio file exists
	audioFile := filepath.Join(outputPath, "audio.mp4")
	verifyFileExists(t, audioFile)
	verifyFileSize(t, audioFile, 10*1024)
}

func TestFFmpeg_MultiBitrate_ABR(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	ctx := context.Background()

	// When: Create 720p version
	containerCmd720p := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "scale=1280:720",
		"-c:v", "libx264",
		"-b:v", "2800k",
		"-preset", "medium",
		"-an",
		"/output/video_720p.mp4",
	}

	c720p, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffmpegImage,
			Cmd:   containerCmd720p,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
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
	defer c720p.Terminate(ctx)

	// When: Create 480p version
	containerCmd480p := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "scale=854:480",
		"-c:v", "libx264",
		"-b:v", "1400k",
		"-preset", "medium",
		"-an",
		"/output/video_480p.mp4",
	}

	c480p, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffmpegImage,
			Cmd:   containerCmd480p,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
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
	defer c480p.Terminate(ctx)

	printContainerLogs(t, ctx, c720p)
	printContainerLogs(t, ctx, c480p)

	// Then: Verify both outputs exist
	video720p := filepath.Join(outputPath, "video_720p.mp4")
	video480p := filepath.Join(outputPath, "video_480p.mp4")

	verifyFileExists(t, video720p)
	verifyFileExists(t, video480p)

	// Verify 720p is larger than 480p
	info720p, _ := os.Stat(video720p)
	info480p, _ := os.Stat(video480p)
	assert.Greater(t, info720p.Size(), info480p.Size(), "720p should be larger than 480p")
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
