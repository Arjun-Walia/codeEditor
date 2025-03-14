package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// executeCode runs the given code in a specified language and returns output/error
func executeCode(language string, code string, input string) (string, string) {
	langConfig := map[string]struct {
		ext     string
		compile string
		run     string
	}{
		"python": {".py", "", "python {filename}"},
		"cpp":    {".cpp", "g++ {filename} -o {exe_name}", "{exe_name}.exe"},
	}

	config, exists := langConfig[language]
	if !exists {
		return "", "Unsupported language"
	}

	tempFile, err := os.CreateTemp("", "code_*"+config.ext)
	if err != nil {
		return "", "Failed to create temp file"
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(code); err != nil {
		tempFile.Close()
		return "", "Failed to write code to file"
	}
	tempFile.Close()

	var exeName string
	if config.compile != "" {
		exeName = strings.TrimSuffix(tempFile.Name(), config.ext)
		compileCmd := strings.Replace(config.compile, "{filename}", tempFile.Name(), -1)
		compileCmd = strings.Replace(compileCmd, "{exe_name}", exeName, -1)

		cmd := exec.Command("cmd", "/C", compileCmd)
		if err := cmd.Run(); err != nil {
			return "", "Compilation failed: " + err.Error()
		}
	}

	runCmd := strings.Replace(config.run, "{filename}", tempFile.Name(), -1)
	runCmd = strings.Replace(runCmd, "{exe_name}", exeName, -1)

	// ✅ FIX: Ensure Proper Timeout Handling
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // Timeout in 2s
	defer cancel()

	cmd := exec.CommandContext(ctx, "cmd", "/C", runCmd)
	cmd.Stdin = strings.NewReader(input)

	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer

	err = cmd.Run()

	// 🚨 Properly check if execution timed out
	if ctx.Err() == context.DeadlineExceeded {
		return "", "Execution timed out! Possible infinite loop detected."
	}

	if err != nil {
		return "", fmt.Sprintf("Execution error: %s\n%s", err.Error(), outputBuffer.String())
	}

	return outputBuffer.String(), ""
}
