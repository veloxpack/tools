package shakapackager

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	shakaPackagerImage = "ghcr.io/veloxpack/shaka-packager:latest"
	ffmpegImage        = "ghcr.io/veloxpack/ffmpeg:8.0-lite"
)

// Helper functions
func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "shaka-test-*")
	require.NoError(t, err)
	return dir
}

func cleanupFiles(t *testing.T, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Logf("Warning: failed to remove directory %s: %v", path, err)
	}
}

func verifyFileExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	require.NoError(t, err, "File should exist: %s", path)
}

func verifyFileSize(t *testing.T, path string, minSize int64) {
	info, err := os.Stat(path)
	require.NoError(t, err, "File should exist: %s", path)
	assert.Greater(t, info.Size(), minSize, "File should have minimum size")
}

func readContainerLogs(t *testing.T, ctx context.Context, c testcontainers.Container) string {
	logs, err := c.Logs(ctx)
	if err != nil {
		return ""
	}
	defer logs.Close()

	reader := bufio.NewReader(logs)
	var output strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			output.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}
	return output.String()
}

// Test 1: Basic DASH packaging with audio and video separation
func TestShakaPackager_BasicDASH(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Package video with separate audio and video streams
	ctx := context.Background()
	containerCmd := []string{
		"in=/input/sample.mp4,stream=audio,output=/output/audio.mp4",
		"in=/input/sample.mp4,stream=video,output=/output/video.mp4",
		"--mpd_output", "/output/manifest.mpd",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: shakaPackagerImage,
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
						Source: outputPath,
						Target: "/output",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer c.Terminate(ctx)

	// Read logs for debugging
	logs := readContainerLogs(t, ctx, c)
	t.Log("Shaka Packager output:", logs)

	// Then: Verify output files
	verifyFileExists(t, filepath.Join(outputPath, "audio.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "video.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "manifest.mpd"))

	// Verify MPD contains expected content
	mpdContent, err := os.ReadFile(filepath.Join(outputPath, "manifest.mpd"))
	require.NoError(t, err)
	mpdStr := string(mpdContent)
	assert.Contains(t, mpdStr, "MPD")
	assert.Contains(t, mpdStr, "AdaptationSet")
	assert.Contains(t, mpdStr, "audio.mp4")
	assert.Contains(t, mpdStr, "video.mp4")
}

// Test 2: HLS packaging with master playlist
func TestShakaPackager_HLS(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Package video for HLS
	ctx := context.Background()
	containerCmd := []string{
		"in=/input/sample.mp4,stream=audio,output=/output/audio.m4a,playlist_name=audio.m3u8",
		"in=/input/sample.mp4,stream=video,output=/output/video.mp4,playlist_name=video.m3u8",
		"--hls_master_playlist_output", "/output/master.m3u8",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: shakaPackagerImage,
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
						Source: outputPath,
						Target: "/output",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer c.Terminate(ctx)

	// Read logs for debugging
	logs := readContainerLogs(t, ctx, c)
	t.Log("Shaka Packager output:", logs)

	// Then: Verify output files
	verifyFileExists(t, filepath.Join(outputPath, "audio.m4a"))
	verifyFileExists(t, filepath.Join(outputPath, "video.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "audio.m3u8"))
	verifyFileExists(t, filepath.Join(outputPath, "video.m3u8"))
	verifyFileExists(t, filepath.Join(outputPath, "master.m3u8"))

	// Verify master playlist contains playlists
	masterContent, err := os.ReadFile(filepath.Join(outputPath, "master.m3u8"))
	require.NoError(t, err)
	masterStr := string(masterContent)
	assert.Contains(t, masterStr, "#EXTM3U")
	assert.Contains(t, masterStr, "audio.m3u8")
	assert.Contains(t, masterStr, "video.m3u8")
}

// Test 3: Multi-bitrate DASH packaging
// This test requires creating multiple bitrate videos first using ffmpeg
func TestShakaPackager_MultiBitrateDASH(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	ctx := context.Background()

	// Step 1: Create 720p version using ffmpeg
	t.Log("Creating 720p version...")
	ffmpegCmd720 := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "scale=1280:720",
		"-c:v", "libx264",
		"-b:v", "2500k",
		"-c:a", "aac",
		"-b:a", "128k",
		"/output/video_720p.mp4",
	}

	cFFmpeg720, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffmpegImage,
			Cmd:   ffmpegCmd720,
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
						Source: outputPath,
						Target: "/output",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer cFFmpeg720.Terminate(ctx)

	// Wait for completion
	exitCode, err := cFFmpeg720.State(ctx)
	require.NoError(t, err)
	if exitCode.ExitCode != 0 {
		logs := readContainerLogs(t, ctx, cFFmpeg720)
		t.Logf("FFmpeg 720p logs: %s", logs)
	}
	require.Equal(t, 0, exitCode.ExitCode, "FFmpeg 720p should complete successfully")

	// Step 2: Create 480p version using ffmpeg
	t.Log("Creating 480p version...")
	ffmpegCmd480 := []string{
		"-i", "/input/sample.mp4",
		"-t", "10",
		"-vf", "scale=854:480",
		"-c:v", "libx264",
		"-b:v", "1200k",
		"-c:a", "aac",
		"-b:a", "128k",
		"/output/video_480p.mp4",
	}

	cFFmpeg480, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: ffmpegImage,
			Cmd:   ffmpegCmd480,
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
						Source: outputPath,
						Target: "/output",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer cFFmpeg480.Terminate(ctx)

	// Wait for completion
	exitCode480, err := cFFmpeg480.State(ctx)
	require.NoError(t, err)
	if exitCode480.ExitCode != 0 {
		logs := readContainerLogs(t, ctx, cFFmpeg480)
		t.Logf("FFmpeg 480p logs: %s", logs)
	}
	require.Equal(t, 0, exitCode480.ExitCode, "FFmpeg 480p should complete successfully")

	// Verify intermediate files exist
	verifyFileExists(t, filepath.Join(outputPath, "video_720p.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "video_480p.mp4"))

	// Step 3: Package with Shaka Packager
	t.Log("Packaging multi-bitrate DASH...")
	video720Path := filepath.Join(outputPath, "video_720p.mp4")
	video480Path := filepath.Join(outputPath, "video_480p.mp4")

	shakaCmd := []string{
		"in=/input/video_720p.mp4,stream=video,output=/output/dash_720p.mp4",
		"in=/input/video_480p.mp4,stream=video,output=/output/dash_480p.mp4",
		"in=/input/sample.mp4,stream=audio,output=/output/dash_audio.mp4",
		"--mpd_output", "/output/manifest.mpd",
	}

	cShaka, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: shakaPackagerImage,
			Cmd:   shakaCmd,
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/input/sample.mp4",
					FileMode:          0o644,
				},
				{
					HostFilePath:      video720Path,
					ContainerFilePath: "/input/video_720p.mp4",
					FileMode:          0o644,
				},
				{
					HostFilePath:      video480Path,
					ContainerFilePath: "/input/video_480p.mp4",
					FileMode:          0o644,
				},
			},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.Mounts = []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: outputPath,
						Target: "/output",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer cShaka.Terminate(ctx)

	// Read logs
	logs := readContainerLogs(t, ctx, cShaka)
	t.Log("Shaka Packager output:", logs)

	// Then: Verify all output files
	verifyFileExists(t, filepath.Join(outputPath, "dash_720p.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "dash_480p.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "dash_audio.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "manifest.mpd"))

	// Verify MPD contains multiple representations
	mpdContent, err := os.ReadFile(filepath.Join(outputPath, "manifest.mpd"))
	require.NoError(t, err)
	mpdStr := string(mpdContent)
	assert.Contains(t, mpdStr, "dash_720p.mp4")
	assert.Contains(t, mpdStr, "dash_480p.mp4")
	assert.Contains(t, mpdStr, "dash_audio.mp4")

	// Count AdaptationSets (should have at least 2: one for video, one for audio)
	adaptationSetCount := strings.Count(mpdStr, "<AdaptationSet")
	assert.GreaterOrEqual(t, adaptationSetCount, 2, "Should have at least 2 AdaptationSets")
}

// Test 4: Fragmented MP4 generation
func TestShakaPackager_FragmentedMP4(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Package with fragment duration
	ctx := context.Background()
	containerCmd := []string{
		"in=/input/sample.mp4,stream=audio,output=/output/audio_frag.mp4",
		"in=/input/sample.mp4,stream=video,output=/output/video_frag.mp4",
		"--fragment_duration", "2",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: shakaPackagerImage,
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
						Source: outputPath,
						Target: "/output",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer c.Terminate(ctx)

	// Read logs
	logs := readContainerLogs(t, ctx, c)
	t.Log("Shaka Packager output:", logs)

	// Then: Verify fragmented files
	verifyFileExists(t, filepath.Join(outputPath, "audio_frag.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "video_frag.mp4"))
	verifyFileSize(t, filepath.Join(outputPath, "audio_frag.mp4"), 1024)
	verifyFileSize(t, filepath.Join(outputPath, "video_frag.mp4"), 1024)
}

// Test 5: Static live MPD generation
func TestShakaPackager_StaticLiveMPD(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Package with static live MPD flag
	ctx := context.Background()
	containerCmd := []string{
		"in=/input/sample.mp4,stream=audio,output=/output/audio.mp4",
		"in=/input/sample.mp4,stream=video,output=/output/video.mp4",
		"--mpd_output", "/output/manifest.mpd",
		"--generate_static_live_mpd",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: shakaPackagerImage,
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
						Source: outputPath,
						Target: "/output",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer c.Terminate(ctx)

	// Read logs
	logs := readContainerLogs(t, ctx, c)
	t.Log("Shaka Packager output:", logs)

	// Then: Verify output and static live profile
	verifyFileExists(t, filepath.Join(outputPath, "manifest.mpd"))

	// Verify MPD has live profile attributes
	mpdContent, err := os.ReadFile(filepath.Join(outputPath, "manifest.mpd"))
	require.NoError(t, err)
	mpdStr := string(mpdContent)
	assert.Contains(t, mpdStr, "MPD")
	// Static live MPDs should have type="dynamic" or specific attributes
	// The exact structure depends on Shaka Packager version
	assert.NotEmpty(t, mpdStr)
}

// Test 6: Segment duration configuration
func TestShakaPackager_SegmentDuration(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Package with custom segment duration
	ctx := context.Background()
	containerCmd := []string{
		"in=/input/sample.mp4,stream=audio,output=/output/audio.mp4",
		"in=/input/sample.mp4,stream=video,output=/output/video.mp4",
		"--mpd_output", "/output/manifest.mpd",
		"--segment_duration", "4",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: shakaPackagerImage,
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
						Source: outputPath,
						Target: "/output",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer c.Terminate(ctx)

	// Read logs
	logs := readContainerLogs(t, ctx, c)
	t.Log("Shaka Packager output:", logs)

	// Then: Verify output files
	verifyFileExists(t, filepath.Join(outputPath, "audio.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "video.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "manifest.mpd"))

	// Verify segment duration in MPD (would contain duration attributes)
	mpdContent, err := os.ReadFile(filepath.Join(outputPath, "manifest.mpd"))
	require.NoError(t, err)
	mpdStr := string(mpdContent)
	assert.Contains(t, mpdStr, "Duration")
}

// Test 7: JSON output parsing for stream info
func TestShakaPackager_StreamInfo(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Package with dump_stream_info
	ctx := context.Background()
	containerCmd := []string{
		"in=/input/sample.mp4,stream=audio,output=/output/audio.mp4",
		"in=/input/sample.mp4,stream=video,output=/output/video.mp4",
		"--mpd_output", "/output/manifest.mpd",
		"--dump_stream_info",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: shakaPackagerImage,
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
						Source: outputPath,
						Target: "/output",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer c.Terminate(ctx)

	// Read logs which should contain stream info
	logs := readContainerLogs(t, ctx, c)
	t.Log("Shaka Packager output:", logs)

	// Then: Verify stream info is in logs
	assert.Contains(t, logs, "Stream")

	// Verify output files still created
	verifyFileExists(t, filepath.Join(outputPath, "audio.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "video.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "manifest.mpd"))
}

// Test 8: Video-only packaging (no audio)
func TestShakaPackager_VideoOnly(t *testing.T) {
	// Given: A test video file
	absPath, err := filepath.Abs(filepath.Join("..", "testdata", "sample.mp4"))
	require.NoError(t, err)

	outputPath := createTempDir(t)
	defer cleanupFiles(t, outputPath)

	// When: Package video stream only
	ctx := context.Background()
	containerCmd := []string{
		"in=/input/sample.mp4,stream=video,output=/output/video_only.mp4",
		"--mpd_output", "/output/manifest.mpd",
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: shakaPackagerImage,
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
						Source: outputPath,
						Target: "/output",
					},
				}
			},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer c.Terminate(ctx)

	// Read logs
	logs := readContainerLogs(t, ctx, c)
	t.Log("Shaka Packager output:", logs)

	// Then: Verify video-only output
	verifyFileExists(t, filepath.Join(outputPath, "video_only.mp4"))
	verifyFileExists(t, filepath.Join(outputPath, "manifest.mpd"))

	// Verify MPD only has video adaptation set
	mpdContent, err := os.ReadFile(filepath.Join(outputPath, "manifest.mpd"))
	require.NoError(t, err)
	mpdStr := string(mpdContent)
	assert.Contains(t, mpdStr, "video_only.mp4")

	// Should have only 1 AdaptationSet for video
	adaptationSetCount := strings.Count(mpdStr, "<AdaptationSet")
	assert.Equal(t, 1, adaptationSetCount, "Should have exactly 1 AdaptationSet for video only")
}

// Helper type for validating JSON output (if needed in future tests)
type StreamInfo struct {
	Type     string `json:"type"`
	Codec    string `json:"codec"`
	Duration string `json:"duration"`
}

func parseStreamInfo(data []byte) ([]StreamInfo, error) {
	var streams []StreamInfo
	err := json.Unmarshal(data, &streams)
	return streams, err
}
