package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spitsynv2/yt-audio-cutter/config"
	"github.com/spitsynv2/yt-audio-cutter/internal/model"
	"github.com/spitsynv2/yt-audio-cutter/internal/store"
)

func downloadAndTrim(url, output, start, end string) error {
	downloadCmd := exec.Command("yt-dlp",
		"-f", "bestaudio",
		"-x", "--audio-format", "mp3",
		"-o", "temp.%(ext)s",
		url,
	)
	downloadCmd.Stdout = log.Writer()
	downloadCmd.Stderr = log.Writer()

	if err := downloadCmd.Run(); err != nil {
		return fmt.Errorf("error downloading: %w", err)
	}

	log.Printf("Running ffmpeg: output=%s start=%s end=%s", output, start, end)

	trimCmd := exec.Command("ffmpeg",
		"-y",
		"-i", "temp.mp3",
		"-ss", start,
		"-to", end,
		"-c", "copy",
		output,
	)
	trimCmd.Stdout = log.Writer()
	trimCmd.Stderr = log.Writer()

	if err := trimCmd.Run(); err != nil {
		return fmt.Errorf("error trimming: %w", err)
	}

	return nil
}

func uploadToDropbox(localFile, dropboxPath, accessToken string) error {
	data, err := os.ReadFile(localFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	url := "https://content.dropboxapi.com/2/files/upload"
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg",
		fmt.Sprintf(`{"path": "%s", "mode": "overwrite", "autorename": true}`, dropboxPath))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("upload error: %s", string(body))
	}

	log.Println("Uploaded successfully to Dropbox:", dropboxPath)
	return nil
}

func createSharedLink(dropboxPath, accessToken string) (string, error) {
	url := "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings"

	payload := map[string]string{
		"path": dropboxPath,
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("share link request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("share link error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if url, ok := result["url"].(string); ok {
		// Convert dl=0 to dl=1 for direct download
		sharedURL := url
		if len(sharedURL) > 0 {
			sharedURL = sharedURL[:len(sharedURL)-1] + "1"
		}
		return sharedURL, nil
	}

	return "", fmt.Errorf("no url found in response: %s", string(body))
}

func ProcessJob(job model.Job) string {
	log.Printf("Processing job: %s", job.Id)
	job.Status = model.StatusRunning

	ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()
	store.UpdateJob(ctx1, job)

	safeName := strings.NewReplacer(" ", "_", "/", "_").Replace(job.Name)
	localFile := fmt.Sprintf("./%s-%s.mp3", safeName, job.Id)
	dropboxPath := fmt.Sprintf("/%s-%s.mp3", safeName, job.Id)

	url := job.YoutubeURL
	start := model.FormatDuration(job.StartTime.Duration)
	end := model.FormatDuration(job.EndTime.Duration)

	if err := downloadAndTrim(url, localFile, start, end); err != nil {
		log.Fatal(err)
	}
	log.Println("Download and trim completed successfully:", localFile)

	accessToken := config.Conf.DropboxToken
	if err := uploadToDropbox(localFile, dropboxPath, accessToken); err != nil {
		log.Fatal(err)
	}
	log.Println("Uploaded to Dropbox:", dropboxPath)

	sharedURL, err := createSharedLink(dropboxPath, accessToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created shared link:", sharedURL)

	// Update job
	job.FileUrl = sharedURL
	job.Status = model.StatusDone

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	store.UpdateJob(ctx2, job)

	return sharedURL
}
