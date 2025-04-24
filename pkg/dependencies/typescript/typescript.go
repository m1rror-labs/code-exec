package typescript

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
	// mu *multisync.Mutex
}

func NewRuntime() pkg.CodeExecutor {
	return &runtime{}
}

func (r *runtime) ExecuteCode(code string) (string, error) {
	// awaiting := r.mu.Acquire()
	// defer r.mu.Release()
	// <-awaiting

	polyfill := `if (typeof CustomEvent !== 'function') {
		class CustomEvent extends Event {
			constructor(event: string, params: { detail?: any; bubbles?: boolean; cancelable?: boolean } = {}) {
				super(event, params);
				this.detail = params.detail || null;
			}
			detail: any;
		}
		(global as any).CustomEvent = CustomEvent;
	}`
	code = polyfill + code

	now := time.Now()
	id := uuid.NewString()
	fullFilename := "./pkg/dependencies/typescript/dist/" + id + ".ts"
	err := os.WriteFile(fullFilename, []byte(code), 0644)
	if err != nil {
		log.Println("Error writing TypeScript file:", err)
		return "", err
	}
	defer os.Remove(fullFilename)

	shortMjsFilename := "./dist/" + id + ".ts"
	cmd := exec.Command("bun", shortMjsFilename)
	cmd.Dir = "./pkg/dependencies/typescript"
	cmd.Env = append(os.Environ(), "NODE_OPTIONS=--no-warnings")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error running JavaScript:", err)
		return "", fmt.Errorf("error running JavaScript: %s", string(output))
	}
	fmt.Println("time taken to run JavaScript:", time.Since(now)) // Log the time taken to run the JavaScript
	return string(output), nil
}
