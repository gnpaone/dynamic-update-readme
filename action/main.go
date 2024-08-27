package main

import (
	"fmt"
	"strings"
	"log"
	"os"
	"os/exec"
	"errors"

	dynreadme "github.com/gnpaone/dynamic-update-readme"
)

func main() {
	readmePath := os.Getenv("INPUT_README_PATH")
	markerText := os.Getenv("INPUT_MARKER_TEXT")
	mdText := os.Getenv("INPUT_MARKDOWN_TEXT")
	isTable := os.Getenv("INPUT_TABLE")
	tableOptions := os.Getenv("INPUT_TABLE_OPTIONS")
	gitUsername := os.Getenv("INPUT_COMMIT_USER")
	gitEmail := os.Getenv("INPUT_COMMIT_EMAIL")
	commitMessage := os.Getenv("INPUT_COMMIT_MESSAGE")
	if commitMessage == "" {
		commitMessage = "Update README.md"
	}
	confirmAndPush := os.Getenv("INPUT_CONFIRM_AND_PUSH")

	// updater := dynreadme.Update{}

	if err := dynreadme.UpdateContent(readmePath, markerText, mdText, isTable, tableOptions); err != nil {
		log.Fatalf("Failed to update README: %s", err)
	}

	if err := updateGitRepo(readmePath, commitMessage, gitUsername, gitEmail, confirmAndPush); err != nil {
		log.Fatalf(err)
	}
}

func updateGitRepo(readmePath, commitMessage, gitUsername, gitEmail, confirmAndPush string) error {
	safeCmd := exec.Command("git", "config", "--global", "--add", "safe.directory", "/github/workspace")
	if err := safeCmd.Run(); err != nil {
		return errors.New("Error setting safe directory: %w", err)
	}

	nameCmd := exec.Command("git", "config", "user.name", gitUsername)
	if err := nameCmd.Run(); err != nil {
		return errors.New("Error setting git user: %w", err)
	}

	emailCmd := exec.Command("git", "config", "user.email", gitEmail)
	if err := emailCmd.Run(); err != nil {
		return errors.New("Error setting git email: %w", err)
	}

	statusCmd, err := exec.Command("git", "status").Output()
	if err != nil {
		return errors.New("Error getting status: %w", err)
	}

	statusOutput := string(statusCmd)
	if !strings.Contains(statusOutput, "nothing to commit") {
		if err := exec.Command("git", "add", readmePath).Run(); err != nil {
			return errors.New("Error adding to staging area: %w", err)
		}
		if err := exec.Command("git", "commit", "-m", commitMessage).Run(); err != nil {
			return errors.New("Error committing to repo: %w", err)
		}
		if confirmAndPush == "true" {
			if err := exec.Command("git", "push").Run(); err != nil {
				return errors.New("Error pushing to repo: %w", err)
			}
		} else if confirmAndPush == "false" {
			output := fmt.Sprintf("git_username=%s\ngit_email=%s\ncommit_message=%s\n", gitUsername, gitEmail, commitMessage)
			return appendToFile(os.Getenv("GITHUB_OUTPUT"), output)
		}
	}

	return nil
}

func appendToFile(filename, text string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("Error opening file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(text); err != nil {
		return errors.New("Error writing to file: %w", err)
	}
	return nil
}
