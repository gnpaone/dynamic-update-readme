package main

import (
	"fmt"
	"strings"
	"log"
	"os"
	"os/exec"

	"github.com/gnpaone/dynamic-update-readme"
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

	if err := dynreadme.UpdateContent(readmePath, markerText, mdText, isTable, tableOptions); err != nil {
		log.Fatalf("Failed to update README: %s", err)
	}

	if err := updateGitRepo(readmePath, commitMessage, gitUsername, gitEmail, confirmAndPush); err != nil {
		log.Fatalf("Failed to update Git repository: %s", err)
	}
}

func updateGitRepo(readmePath, commitMessage, gitUsername, gitEmail, confirmAndPush string) error {
	safeCmd := exec.Command("git", "config", "--global", "--add", "safe.directory", "/github/workspace")
	if err := safeCmd.Run(); err != nil {
		log.Fatalf("Error setting safe directory: %q", err)
	}

	nameCmd := exec.Command("git", "config", "user.name", gitUsername)
	if err := nameCmd.Run(); err != nil {
		log.Fatalf("Error setting git user: %q", err)
	}

	emailCmd := exec.Command("git", "config", "user.email", gitEmail)
	if err := emailCmd.Run(); err != nil {
		log.Fatalf("Error setting git email: %q", err)
	}

	statusCmd, err := exec.Command("git", "status").Output()
	if err != nil {
		log.Fatalf("Error getting status: %q", err)
	}

	statusOutput := string(statusCmd)
	if !strings.Contains(statusOutput, "nothing to commit") {
		if err := exec.Command("git", "add", readmePath).Run(); err != nil {
			log.Fatalf("Error adding to staging area: %q", err)
		}
		if err := exec.Command("git", "commit", "-m", commitMessage).Run(); err != nil {
			log.Fatalf("Error committing to repo: %q", err)
		}
		if confirmAndPush == "true" {
			if err := exec.Command("git", "push").Run(); err != nil {
				log.Fatalf("Error pushing to repo: %q", err)
			}
		} else if confirmAndPush == "false" {
			output := fmt.Sprintf("git_username=%s\ngit_email=%s\ncommit_message=%s\n", gitUsername, gitEmail, commitMessage)
			if err := appendToFile(os.Getenv("GITHUB_OUTPUT"), output); err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

func appendToFile(filename, text string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %q", err)
	}
	defer file.Close()

	if _, err := file.WriteString(text); err != nil {
		log.Fatalf("Error writing to file: %q", err)
	}
	return nil
}
