package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	cmd "todo/cmd"
	"todo/config"
	td "todo/td"
)

func Pick() {
	tasks, err := cmd.SelectTasksWithFzf(true, config.DisableFzf)
	if err != nil || len(tasks) == 0 {
		fail("No task selected: %v", err)
		return
	}

	// --json
	if contains("--json") {
		out, _ := json.MarshalIndent(tasks, "", "  ")
		fmt.Println(string(out))
		return
	}

	// --edit
	if contains("--edit") {
		for _, t := range tasks {
			newText := td.PromptInput("✏️ Edit "+t.Text, t.Text)
			if newText != "" {
				_ = td.UpdateTaskByID(t.ID, func(n *td.Task) { n.Text = newText })
				success("Edited #%d", t.ID)
			}
		}
		return
	}

	// --note
	if contains("--note") {
		for _, t := range tasks {
			note := td.PromptMultiline("📝 Note for "+t.Text, t.Notes)
			if note != "" {
				_ = td.UpdateTaskByID(t.ID, func(n *td.Task) { n.Notes = note })
				success("Note saved for #%d", t.ID)
			}
		}
		return
	}

	// Default → output IDs
	for _, t := range tasks {
		fmt.Println(t.ID)
	}
}

func contains(flag string) bool {
	for _, a := range os.Args {
		if strings.ToLower(a) == flag {
			return true
		}
	}
	return false
}
