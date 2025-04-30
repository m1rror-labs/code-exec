package anchor

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

func (r *runtime) TestCode(code string, blockchainID uuid.UUID, id string) (string, error) {
	defer os.RemoveAll("./pkg/dependencies/anchor/code-exec/programs/" + id)
	defer os.RemoveAll("./pkg/dependencies/anchor/code-exec/target/deploy")
	defer os.RemoveAll("./pkg/dependencies/anchor/code-exec/tests/" + id + ".ts")

	err := os.WriteFile(fmt.Sprintf("./pkg/dependencies/anchor/code-exec/tests/%s.ts", id), []byte(code), 0644)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("anchor", "build", "--program-name", id, "--provider.cluster", fmt.Sprintf("https://engine.mirror.ad/rpc/%s", blockchainID.String()))
	cmd.Dir = "./pkg/dependencies/anchor/code-exec"
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("error testing Anchor: %s", string(output))
	}

	return string(output), nil
}
