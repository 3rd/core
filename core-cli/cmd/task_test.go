package cmd

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestTaskCurrentCommandFindsFirstInProgressTask(t *testing.T) {
	root := t.TempDir()
	writeTaskNode(t, root, "a-zulu", "Zulu", "zulu task")
	writeTaskNode(t, root, "z-alpha", "Alpha", "alpha task")

	previousEnv := env
	env.WIKI_ROOT = root
	t.Cleanup(func() {
		env = previousEnv
	})

	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	previousStdout := os.Stdout
	os.Stdout = write
	t.Cleanup(func() {
		os.Stdout = previousStdout
	})

	taskCurrentCommand.Run(taskCurrentCommand, nil)

	if err := write.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = previousStdout
	output, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}
	if err := read.Close(); err != nil {
		t.Fatal(err)
	}

	if !regexp.MustCompile(`^Alpha - alpha task \(.+\)\n$`).Match(output) {
		t.Fatalf("unexpected current task output: %q", output)
	}
}

func writeTaskNode(t *testing.T, root string, fileName string, title string, task string) {
	t.Helper()
	content := "@meta\n" +
		"  type: project\n" +
		"  title: " + title + "\n" +
		"@end\n\n" +
		"* Tasks\n" +
		"  [-] " + task + "\n" +
		"    Session: 2024.01.01 01:00\n"
	if err := os.WriteFile(filepath.Join(root, fileName), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
