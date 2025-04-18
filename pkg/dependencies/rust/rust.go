package rust

import (
	"code-exec/pkg"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

type runtime struct {
}

func NewRuntime() pkg.CodeExecutor {
	return &runtime{}
}

func (r *runtime) ExecuteCode(code string) (string, error) {
	// awaiting := r.mu.Acquire()
	// defer r.mu.Release()
	// <-awaiting

	now := time.Now()
	id := uuid.NewString()
	fullFilename := "./pkg/dependencies/runtimes/rust/src/bin/" + id + ".rs"
	err := os.WriteFile(fullFilename, []byte(code), 0644)
	if err != nil {
		return "", err
	}
	defer os.Remove(fullFilename)
	defer os.Remove("./pkg/dependencies/runtimes/rust/target/debug/" + id)
	defer os.Remove("./pkg/dependencies/runtimes/rust/target/debug/" + id + ".d")

	cmd := exec.Command("cargo", "run", "--locked", "--bin", id)
	cmd.Dir = "./pkg/dependencies/runtimes/rust"
	// cmd.Env = append(os.Environ(), "CARGO_TARGET_DIR=./pkg/dependencies/runtimes/rust/target")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("error running Rust: %s", string(output))
	}
	fmt.Println("time taken to run Rust:", time.Since(now)) // Log the time taken to run the JavaScript
	return string(output), nil
}
