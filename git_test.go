package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseArgs_Defaults(t *testing.T) {
	cfg, err := parseArgs(nil)
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if cfg.theme != "default" {
		t.Fatalf("expected default theme, got %q", cfg.theme)
	}
	if cfg.dir != "" {
		t.Fatalf("expected empty dir, got %q", cfg.dir)
	}
}

func TestParseArgs_WithDirectoryAndFlags(t *testing.T) {
	cfg, err := parseArgs([]string{"--theme", "matrix", "--output", "out.gif", "/tmp/repo"})
	if err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}
	if cfg.theme != "matrix" || cfg.output != "out.gif" || cfg.dir != "/tmp/repo" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestParseArgs_RejectsMultipleDirectories(t *testing.T) {
	_, err := parseArgs([]string{"repo1", "repo2"})
	if err == nil || !strings.Contains(err.Error(), "only one target directory") {
		t.Fatalf("expected multiple directory error, got %v", err)
	}
}

func TestGetRepoInfo_CurrentDirectory(t *testing.T) {
	repoDir := setupTestRepo(t)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(cwd)
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	info, err := getRepoInfo("")
	if err != nil {
		t.Fatalf("getRepoInfo returned error: %v", err)
	}
	if info.name != filepath.Base(repoDir) {
		t.Fatalf("expected repo name %q, got %q", filepath.Base(repoDir), info.name)
	}
	if info.totalCommits < 1 {
		t.Fatalf("expected commits, got %d", info.totalCommits)
	}
}

func TestGetRepoInfo_TargetDirectory(t *testing.T) {
	repoDir := setupTestRepo(t)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	info, err := getRepoInfo(repoDir)
	if err != nil {
		t.Fatalf("getRepoInfo returned error: %v", err)
	}
	if info.name != filepath.Base(repoDir) {
		t.Fatalf("expected repo name %q, got %q", filepath.Base(repoDir), info.name)
	}
	if after, _ := os.Getwd(); after != cwd {
		t.Fatalf("working directory changed from %q to %q", cwd, after)
	}
}

func TestGetRepoInfo_InvalidDirectory(t *testing.T) {
	_, err := getRepoInfo("/definitely/not/a/repo")
	if err == nil {
		t.Fatal("expected error for invalid directory")
	}
}

func setupTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runInDir(t, dir, "git", "init")
	runInDir(t, dir, "git", "config", "user.name", "Test User")
	runInDir(t, dir, "git", "config", "user.email", "test@example.com")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	runInDir(t, dir, "git", "add", "README.md")
	runInDir(t, dir, "git", "commit", "-m", "feat: initial commit")
	return dir
}

func runInDir(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, string(out))
	}
}
