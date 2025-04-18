package rpcengine

import (
	"bytes"
	"code-exec/pkg"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"

	"github.com/google/uuid"
)

type rpcEngine struct {
	url string
}

func New(url string) pkg.RpcEngine {
	return &rpcEngine{url}
}

type createBlockchainRequest struct {
	Config                 *uuid.UUID `json:"config"`
	DeferAccountInitiation bool       `json:"defer_account_initiation"`
}

func (e *rpcEngine) CreateBlockchain(ctx context.Context, apiKey uuid.UUID, user_id *string, config *uuid.UUID) (uuid.UUID, error) {
	reqBody := createBlockchainRequest{Config: config, DeferAccountInitiation: false}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("error marshalling request", err)
		return uuid.Nil, pkg.ErrHttpRequest
	}

	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/blockchains", e.url), bytes.NewReader(reqBytes))
	if err != nil {
		log.Println("error creating request", err)
		return uuid.Nil, pkg.ErrHttpRequest
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("api_key", apiKey.String())
	if user_id != nil {
		r.Header.Set("user_id", *user_id)
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println("error sending request", err)
		return uuid.Nil, pkg.ErrHttpRequest
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body", err)
		return uuid.Nil, pkg.ErrHttpRequest
	}

	type Error struct {
		Message string `json:"message"`
	}

	if resp.StatusCode != http.StatusOK {
		var res Error
		if err := json.Unmarshal(body, &res); err != nil {
			log.Println("error unmarshalling error", resp.StatusCode, err, string(body))
			return uuid.Nil, pkg.ErrHttpRequest
		}
		return uuid.Nil, errors.New(res.Message)
	}

	type Response struct {
		Url string `json:"url"`
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		log.Println("error unmarshalling response", err, string(body))
		return uuid.Nil, pkg.ErrHttpRequest
	}

	id := removeMirrorRPC(res.Url)
	log.Println(id)
	return uuid.Parse(id)
}

func removeMirrorRPC(url string) string {
	re := regexp.MustCompile(`https?://(rpc\.mirror\.ad/rpc/|localhost:8899/rpc/)`)
	return re.ReplaceAllString(url, "")
}

func (e *rpcEngine) DeleteBlockchain(ctx context.Context, apiKey uuid.UUID, blockchainID uuid.UUID) error {
	r, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/rpc/%s", e.url, blockchainID.String()), nil)
	if err != nil {
		log.Println("error creating request", err)
		return pkg.ErrHttpRequest
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("api_key", apiKey.String())

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println("error sending request", err)
		return pkg.ErrHttpRequest
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body", err)
		return pkg.ErrHttpRequest
	}

	type Error struct {
		Message string `json:"message"`
	}

	if resp.StatusCode != http.StatusOK {
		var res Error
		if err := json.Unmarshal(body, &res); err != nil {
			log.Println("error unmarshalling error", err, string(body))
			return pkg.ErrHttpRequest
		}

		return errors.New(res.Message)
	}

	return nil
}

func (e *rpcEngine) ExpireBlockchains(ctx context.Context) error {
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/blockchains/expire", e.url), nil)
	if err != nil {
		log.Println("error creating request", err)
		return pkg.ErrHttpRequest
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println("error sending request", err)
		return pkg.ErrHttpRequest
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body", err)
		return pkg.ErrHttpRequest
	}

	type Error struct {
		Message string `json:"message"`
	}

	if resp.StatusCode != http.StatusOK {
		var res Error
		if err := json.Unmarshal(body, &res); err != nil {
			log.Println("error unmarshalling error", err, string(body))
			return pkg.ErrHttpRequest
		}

		return errors.New(res.Message)
	}

	return nil
}

func (e *rpcEngine) LoadProgram(ctx context.Context, blockchainID uuid.UUID, programID string, programBinary []byte) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the program binary as a file part
	part, err := writer.CreateFormFile("program", "program.so")
	if err != nil {
		log.Println("error creating form file", err)
		return pkg.ErrHttpRequest
	}
	if _, err := part.Write(programBinary); err != nil {
		log.Println("error writing program binary", err)
		return pkg.ErrHttpRequest
	}

	// Add other form fields if necessary
	if err := writer.WriteField("program_id", programID); err != nil {
		log.Println("error writing form field", err)
		return pkg.ErrHttpRequest
	}

	// Close the writer to finalize the multipart message
	if err := writer.Close(); err != nil {
		log.Println("error closing writer", err)
		return pkg.ErrHttpRequest
	}

	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/programs/%s", e.url, blockchainID.String()), &buf)
	if err != nil {
		log.Println("error creating request", err)
		return pkg.ErrHttpRequest
	}
	r.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	r.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println("error sending request", err)
		return pkg.ErrHttpRequest
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body", err)
		return pkg.ErrHttpRequest
	}

	type Error struct {
		Message string `json:"message"`
	}

	if resp.StatusCode != http.StatusOK {
		var res Error
		if err := json.Unmarshal(body, &res); err != nil {
			log.Println("error unmarshalling error", err, string(body))
			return pkg.ErrHttpRequest
		}
		log.Println("error unmarshalling error", err, string(body))

		return errors.New(res.Message)
	}

	return nil
}

type setBlockchainRequest struct {
	Address       string  `json:"address"`
	Lamports      uint    `json:"lamports"`
	Data          string  `json:"data"`
	Owner         string  `json:"owner"`
	RentEpoch     uint    `json:"rent_epoch"`
	Executable    bool    `json:"executable"`
	Label         *string `json:"label"`
	TokenMintAuth *string `json:"token_mint_auth"`
}

func (e *rpcEngine) SetAccounts(
	ctx context.Context,
	blockchainID uuid.UUID,
	accounts []pkg.SolanaAccount,
	label *string,
	tokenMintAuth *string,
) error {
	var reqBody []setBlockchainRequest
	for _, account := range accounts {
		encodedData := base64.StdEncoding.EncodeToString([]byte(account.Data))
		reqBody = append(reqBody, setBlockchainRequest{
			Address:       account.Address,
			Lamports:      account.Lamports,
			Data:          encodedData,
			Owner:         account.Owner,
			RentEpoch:     account.RentEpoch,
			Executable:    account.Executable,
			Label:         label,
			TokenMintAuth: tokenMintAuth,
		})
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("error marshalling request", err)
		return pkg.ErrHttpRequest
	}
	log.Println("request body", string(reqBytes))

	r, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/accounts/%s", e.url, blockchainID.String()), bytes.NewReader(reqBytes))
	if err != nil {
		log.Println("error creating request", err)
		return pkg.ErrHttpRequest
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println("error sending request", err)
		return pkg.ErrHttpRequest
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("error response", resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("error reading response body", err)
			return pkg.ErrHttpRequest
		}
		log.Println("error unmarshalling error", err, string(body))
		return pkg.ErrHttpRequest
	}

	return nil
}
