//go:build integration
// +build integration

package files_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/files"
)

func TestFilesIntegration(t *testing.T) {
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		t.Skip("XAI_API_KEY not set, skipping integration test")
	}

	client, err := xai.NewClient(&xai.Config{APIKey: apiKey})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	var uploadedFileID string

	t.Run("UploadFile", func(t *testing.T) {
		content := strings.NewReader("This is a test file for integration testing.")

		file, err := client.Files().Upload(ctx, content, files.UploadOptions{
			Name:    "test-integration.txt",
			Purpose: "assistants",
		})
		if err != nil {
			t.Fatalf("Upload failed: %v", err)
		}

		if file.ID == "" {
			t.Fatal("Expected file ID")
		}

		uploadedFileID = file.ID
		t.Logf("Uploaded file: %s (%d bytes)", file.ID, file.Size)
	})

	t.Run("ListFiles", func(t *testing.T) {
		result, err := client.Files().List(ctx, &files.ListOptions{
			Limit: 10,
		})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(result.Files) == 0 {
			t.Log("No files found (this is okay)")
		} else {
			t.Logf("Found %d files", len(result.Files))
		}
	})

	t.Run("GetFile", func(t *testing.T) {
		if uploadedFileID == "" {
			t.Skip("No file uploaded")
		}

		file, err := client.Files().Get(ctx, uploadedFileID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if file.ID != uploadedFileID {
			t.Fatalf("Expected file ID %s, got %s", uploadedFileID, file.ID)
		}

		t.Logf("Retrieved file: %s", file.Filename)
	})

	t.Run("GetFileURL", func(t *testing.T) {
		if uploadedFileID == "" {
			t.Skip("No file uploaded")
		}

		url, err := client.Files().GetURL(ctx, uploadedFileID)
		if err != nil {
			t.Fatalf("GetURL failed: %v", err)
		}

		if url == "" {
			t.Fatal("Expected non-empty URL")
		}

		t.Logf("Download URL: %s", url)
	})

	t.Run("DownloadFile", func(t *testing.T) {
		if uploadedFileID == "" {
			t.Skip("No file uploaded")
		}

		reader, err := client.Files().Download(ctx, uploadedFileID)
		if err != nil {
			t.Fatalf("Download failed: %v", err)
		}
		defer reader.Close()

		t.Log("File download successful")
	})

	// Cleanup
	t.Run("DeleteFile", func(t *testing.T) {
		if uploadedFileID == "" {
			t.Skip("No file uploaded")
		}

		err := client.Files().Delete(ctx, uploadedFileID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		t.Logf("Deleted file: %s", uploadedFileID)
	})
}
