package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestBatchUpload_Success(t *testing.T) {
	// Create mock readers
	readers := []io.Reader{
		strings.NewReader("file1 content"),
		strings.NewReader("file2 content"),
		strings.NewReader("file3 content"),
	}

	opts := []UploadOptions{
		{Name: "file1.txt"},
		{Name: "file2.txt"},
		{Name: "file3.txt"},
	}

	// Note: This test requires a mock client
	// For now, we're testing the validation logic
	client := &Client{restClient: nil}

	results, err := client.BatchUpload(context.Background(), readers, opts, 3, nil)

	// Should succeed but all uploads will fail due to nil client
	if err != nil {
		t.Errorf("BatchUpload should not return error for parameter validation, got: %v", err)
	}

	// Check that we got results for all files
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// All should have errors due to nil client
	for i, result := range results {
		if result.Error == nil {
			t.Errorf("Expected error for result %d", i)
		}
	}
}

func TestBatchUpload_EmptyReaders(t *testing.T) {
	client := &Client{}

	_, err := client.BatchUpload(context.Background(), []io.Reader{}, []UploadOptions{}, 10, nil)

	if err == nil {
		t.Error("Expected error for empty readers")
	}

	expectedMsg := "readers cannot be empty"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message to contain %q, got %q", expectedMsg, err.Error())
	}
}

func TestBatchUpload_MismatchedLengths(t *testing.T) {
	client := &Client{}

	readers := []io.Reader{
		strings.NewReader("content1"),
		strings.NewReader("content2"),
	}

	opts := []UploadOptions{
		{Name: "file1.txt"},
	}

	_, err := client.BatchUpload(context.Background(), readers, opts, 10, nil)

	if err == nil {
		t.Error("Expected error for mismatched lengths")
	}

	expectedMsg := "opts length (1) must match readers length (2)"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestBatchUpload_DefaultBatchSize(t *testing.T) {
	// Test that batch size defaults to 50 when <= 0
	client := &Client{}

	readers := []io.Reader{strings.NewReader("content")}
	opts := []UploadOptions{{Name: "file.txt"}}

	// Test with 0 - should default to 50
	results, err := client.BatchUpload(context.Background(), readers, opts, 0, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// Test with negative - should default to 50
	results, err = client.BatchUpload(context.Background(), readers, opts, -1, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

func TestBatchUploadResult(t *testing.T) {
	// Test BatchUploadResult structure
	result := &BatchUploadResult{
		Index: 0,
		File: &File{
			ID:       "file-123",
			Filename: "test.txt",
			Size:     100,
		},
		Error: nil,
	}

	if result.Index != 0 {
		t.Errorf("Expected Index 0, got %d", result.Index)
	}

	if result.File == nil {
		t.Error("Expected File to be set")
	}

	if result.Error != nil {
		t.Errorf("Expected no error, got %v", result.Error)
	}
}

func TestBatchUploadResult_WithError(t *testing.T) {
	// Test BatchUploadResult with error
	result := &BatchUploadResult{
		Index: 1,
		File:  nil,
		Error: fmt.Errorf("upload failed"),
	}

	if result.Index != 1 {
		t.Errorf("Expected Index 1, got %d", result.Index)
	}

	if result.File != nil {
		t.Error("Expected File to be nil")
	}

	if result.Error == nil {
		t.Error("Expected error to be set")
	}
}

func TestBatchUploadCallback(t *testing.T) {
	// Test that callback type is correct
	var callback BatchUploadCallback

	callbackCalled := false
	callback = func(index int, reader io.Reader, result interface{}) {
		callbackCalled = true
		if index != 0 {
			t.Errorf("Expected index 0, got %d", index)
		}
		if reader == nil {
			t.Error("Expected reader to be set")
		}
		if result == nil {
			t.Error("Expected result to be set")
		}
	}

	// Simulate callback
	mockReader := strings.NewReader("test")
	mockResult := &BatchUploadResult{Index: 0}
	callback(0, mockReader, mockResult)

	if !callbackCalled {
		t.Error("Callback was not called")
	}
}

func TestBatchUpload_ParameterValidation(t *testing.T) {
	tests := []struct {
		name        string
		readers     []io.Reader
		opts        []UploadOptions
		batchSize   int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty_readers",
			readers:     []io.Reader{},
			opts:        []UploadOptions{},
			batchSize:   10,
			expectError: true,
			errorMsg:    "readers cannot be empty",
		},
		{
			name: "mismatched_lengths_more_readers",
			readers: []io.Reader{
				strings.NewReader("1"),
				strings.NewReader("2"),
				strings.NewReader("3"),
			},
			opts: []UploadOptions{
				{Name: "file1.txt"},
				{Name: "file2.txt"},
			},
			batchSize:   10,
			expectError: true,
			errorMsg:    "opts length (2) must match readers length (3)",
		},
		{
			name: "mismatched_lengths_more_opts",
			readers: []io.Reader{
				strings.NewReader("1"),
			},
			opts: []UploadOptions{
				{Name: "file1.txt"},
				{Name: "file2.txt"},
			},
			batchSize:   10,
			expectError: true,
			errorMsg:    "opts length (2) must match readers length (1)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}
			_, err := client.BatchUpload(context.Background(), tt.readers, tt.opts, tt.batchSize, nil)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got %v", err)
				}
			}
		})
	}
}

func TestBatchUpload_ConcurrencyControl(t *testing.T) {
	// This test verifies that the batch size parameter is respected
	// We can't fully test concurrency without a real client, but we can
	// verify the parameter validation works

	client := &Client{}

	// Create 100 readers
	readers := make([]io.Reader, 100)
	opts := make([]UploadOptions, 100)
	for i := 0; i < 100; i++ {
		readers[i] = strings.NewReader(fmt.Sprintf("content%d", i))
		opts[i] = UploadOptions{Name: fmt.Sprintf("file%d.txt", i)}
	}

	// Test with various batch sizes
	batchSizes := []int{1, 5, 10, 50, 100}

	for _, batchSize := range batchSizes {
		t.Run(fmt.Sprintf("batch_size_%d", batchSize), func(t *testing.T) {
			results, err := client.BatchUpload(context.Background(), readers, opts, batchSize, nil)
			// Should succeed with results containing errors
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(results) != 100 {
				t.Errorf("Expected 100 results, got %d", len(results))
			}
		})
	}
}

func TestBatchUpload_CallbackInvocation(t *testing.T) {
	// Test that callback signature is correct and can be invoked
	callCount := 0
	successCount := 0
	failureCount := 0

	callback := func(_ int, _ io.Reader, result interface{}) {
		callCount++

		if res, ok := result.(*BatchUploadResult); ok {
			if res.Error != nil {
				failureCount++
			} else {
				successCount++
			}
		}
	}

	// Simulate callbacks
	callback(0, strings.NewReader("test"), &BatchUploadResult{Index: 0, File: &File{ID: "1"}})
	callback(1, strings.NewReader("test"), &BatchUploadResult{Index: 1, Error: fmt.Errorf("failed")})
	callback(2, strings.NewReader("test"), &BatchUploadResult{Index: 2, File: &File{ID: "3"}})

	if callCount != 3 {
		t.Errorf("Expected 3 callback invocations, got %d", callCount)
	}

	if successCount != 2 {
		t.Errorf("Expected 2 successes, got %d", successCount)
	}

	if failureCount != 1 {
		t.Errorf("Expected 1 failure, got %d", failureCount)
	}
}

func TestBatchUpload_ResultsMap(t *testing.T) {
	// Test that results map structure is correct
	results := make(map[int]*BatchUploadResult)

	// Add some results
	results[0] = &BatchUploadResult{Index: 0, File: &File{ID: "file-1"}}
	results[1] = &BatchUploadResult{Index: 1, Error: fmt.Errorf("failed")}
	results[2] = &BatchUploadResult{Index: 2, File: &File{ID: "file-3"}}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Check successful upload
	if results[0].File == nil {
		t.Error("Expected File at index 0")
	}
	if results[0].Error != nil {
		t.Error("Expected no error at index 0")
	}

	// Check failed upload
	if results[1].File != nil {
		t.Error("Expected no File at index 1")
	}
	if results[1].Error == nil {
		t.Error("Expected error at index 1")
	}

	// Check another successful upload
	if results[2].File == nil {
		t.Error("Expected File at index 2")
	}
	if results[2].Error != nil {
		t.Error("Expected no error at index 2")
	}
}

func TestBatchUpload_Documentation(t *testing.T) {
	// This test verifies that the API matches the documented example
	client := &Client{}

	readers := []io.Reader{
		bytes.NewReader([]byte("content1")),
		bytes.NewReader([]byte("content2")),
	}

	opts := []UploadOptions{
		{Name: "file1.txt"},
		{Name: "file2.txt"},
	}

	callback := func(idx int, r io.Reader, result interface{}) {
		if res, ok := result.(*BatchUploadResult); ok {
			if res.Error != nil {
				t.Logf("File %d failed: %v\n", idx, res.Error)
			} else if res.File != nil {
				t.Logf("File %d uploaded: %s\n", idx, res.File.ID)
			}
		}
	}

	results, err := client.BatchUpload(context.Background(), readers, opts, 10, callback)

	// Should succeed with results containing errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}
