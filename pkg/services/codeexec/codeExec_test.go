package codeexec

import (
	"testing"

	"github.com/google/uuid"
)

func TestParseEngineUrl(t *testing.T) {
	t.Skip()
	uuid := uuid.NewString()
	localhost := "http://localhost:8899/rpc/" + uuid
	live := "https://engine.mirror.ad/rpc/" + uuid
	localhostID, ok := parseEngineUrl(localhost)
	if !ok || localhostID.String() != uuid {
		t.Errorf("Expected %s to be parsed as %s", localhostID.String(), uuid)
	}

	liveID, ok := parseEngineUrl(live)
	if !ok || liveID.String() != uuid {
		t.Errorf("Expected %s to be parsed as %s", liveID.String(), uuid)
	}
	t.Fatal(localhostID)
}
