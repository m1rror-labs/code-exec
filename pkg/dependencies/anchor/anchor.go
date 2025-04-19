package anchor

import (
	"code-exec/pkg"
	"fmt"
	"log"
	"math/rand"

	"os"
	"os/exec"
	"time"
)

type runtime struct {
}

func NewRuntime() pkg.ProgramBuilder {
	return &runtime{}
}

func (r *runtime) BuildProgram(code string) ([]byte, error) {
	now := time.Now()
	id := randString(20)
	newCmd := exec.Command("anchor", "new", id)
	newCmd.Dir = "./pkg/dependencies/anchor/code-exec"
	if output, err := newCmd.CombinedOutput(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error running Anchor new: %s", string(output))
	}
	defer os.RemoveAll("./pkg/dependencies/anchor/code-exec/programs/" + id)

	err := os.WriteFile(fmt.Sprintf("./pkg/dependencies/anchor/code-exec/programs/%s/src/lib.rs", id), []byte(code), 0644)
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll("./pkg/dependencies/anchor/code-exec/target/deploy")
	cmd := exec.Command("anchor", "build", "--program-name", id)
	cmd.Dir = "./pkg/dependencies/anchor/code-exec"
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error building Anchor: %s", string(output))
	}
	outputFile := fmt.Sprintf("./pkg/dependencies/anchor/code-exec/target/deploy/%s.so", id)
	output, err := os.ReadFile(outputFile)
	if err != nil {
		log.Println("Error reading output file:", err)
		return nil, fmt.Errorf("error reading output file: %s", err)
	}
	fmt.Println("time taken to run Anchor:", time.Since(now)) // Log the time taken to run the JavaScript
	return output, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
