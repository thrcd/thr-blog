package main

import (
	"github.com/thrcd/thr-blog/internal/testkit"
	"slices"
	"testing"
)

func TestGetSubDirs(t *testing.T) {
	t.Log("Test getting subdirectories from root folder")
	{
		root := "content/test"
		want := "content/test/tech"

		dirs := getSubDirs(root)
		args := []any{want, dirs}
		testkit.Check(t, slices.Contains(dirs, want), "Should return dir %s. Received %v", args...)
	}
}

func TestGetFilePaths(t *testing.T) {
	dir := "content/test"

	t.Log("Test getting files in specific dir")
	{
		t.Log("Test 0: When passing a valid dir")
		{
			want := "content/test/about.md"

			files, err := getFilePaths(dir)
			if err != nil {
				errorArgs := []any{dir, err}
				testkit.ErrorT(t, "Should return files from %s. Got Error: %v", errorArgs)
			}

			args := []any{want, files}
			testkit.Check(t, slices.Contains(files, want), "Should return file %s. Received %v", args...)
		}

		t.Log("Test 1: When passing an invalid dir")
		{
			files, err := getFilePaths("folder")
			testkit.Check(t, err != nil, "Should return an err: file does not exist. Received %s", err)
			testkit.Check(t, len(files) == 0, "Should return an empty slice for files. Received %v", files)
		}
	}
}
