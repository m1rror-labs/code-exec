package anchor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func (r *runtime) BuildProgram(code string, deleteArtifacts bool) ([]byte, string, error) {
	id := randString(20)
	newCmd := exec.Command("anchor", "new", id)
	newCmd.Dir = "./pkg/dependencies/anchor/code-exec"
	if output, err := newCmd.CombinedOutput(); err != nil {
		log.Println(err)
		return nil, id, fmt.Errorf("error running Anchor new: %s", string(output))
	}
	if deleteArtifacts {
		defer os.RemoveAll("./pkg/dependencies/anchor/code-exec/programs/" + id)
		defer os.RemoveAll("./pkg/dependencies/anchor/code-exec/target/deploy")
	}

	err := os.WriteFile(fmt.Sprintf("./pkg/dependencies/anchor/code-exec/programs/%s/src/lib.rs", id), []byte(code), 0644)
	if err != nil {
		return nil, id, err
	}

	cmd := exec.Command("anchor", "build", "--program-name", id)
	cmd.Dir = "./pkg/dependencies/anchor/code-exec"
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println(err)
		return nil, id, fmt.Errorf("error building Anchor: %s", string(output))
	}
	outputFile := fmt.Sprintf("./pkg/dependencies/anchor/code-exec/target/deploy/%s.so", id)
	output, err := os.ReadFile(outputFile)
	if err != nil {
		log.Println("Error reading output file:", err)
		return nil, id, fmt.Errorf("error reading output file: %s", err)
	}
	return output, id, nil
}
