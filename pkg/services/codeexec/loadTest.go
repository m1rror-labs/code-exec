package codeexec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func LoadTestCodeExec(
	ctx context.Context,
	numRustThreads int,
	numTsThreads int,
) error {
	for range numRustThreads {
		go func() {
			if err := RunRustCodeExec(ctx); err != nil {
				log.Println("error running rust code exec", err)
			}
		}()
	}

	for range numTsThreads {
		go func() {
			if err := RunTsCodeExec(ctx); err != nil {
				log.Println("error running ts code exec", err)
			}
		}()
	}

	time.Sleep(time.Second)

	for range numRustThreads {
		go func() {
			if err := RunRustCodeExec(ctx); err != nil {
				log.Println("error running rust code exec", err)
			}
		}()
	}

	for range numTsThreads {
		go func() {
			if err := RunTsCodeExec(ctx); err != nil {
				log.Println("error running ts code exec", err)
			}
		}()
	}

	return nil
}

func RunRustCodeExec(
	ctx context.Context,
) error {
	start := time.Now()
	rustCode := `use solana_client::nonblocking::rpc_client::RpcClient;
	use solana_sdk::{
		commitment_config::CommitmentConfig, native_token::LAMPORTS_PER_SOL, signature::Keypair,
		signer::Signer, system_instruction::transfer, transaction::Transaction,
	};
	
	#[tokio::main]
	async fn main() -> anyhow::Result<()> {
		let client = RpcClient::new_with_commitment(
			String::from("https://engine.mirror.ad/rpc/<mirror-id>"),
			CommitmentConfig::confirmed(),
		);
	
		let from_keypair = Keypair::new();
		let to_keypair = Keypair::new();
	
		let transfer_ix = transfer(
			&from_keypair.pubkey(),
			&to_keypair.pubkey(),
			LAMPORTS_PER_SOL,
		);
	
	
		let transaction_signature = client
			.request_airdrop(&from_keypair.pubkey(), 5 * LAMPORTS_PER_SOL)
			.await?;
		loop {
			if client.confirm_transaction(&transaction_signature).await? {
				break;
			}
		}
	
		let mut transaction = Transaction::new_with_payer(&[transfer_ix], Some(&from_keypair.pubkey()));
		transaction.sign(&[&from_keypair], client.get_latest_blockhash().await?);
	
		match client.send_and_confirm_transaction(&transaction).await {
			Ok(signature) => println!("Transaction Signature: {}", signature),
			Err(err) => eprintln!("Error sending transaction: {}", err),
		}
	
		Ok(())
	}`

	sessionUrl, _, err := GetSession(ctx)
	if err != nil {
		return fmt.Errorf("error getting session: %v", err)
	}
	rustCode = strings.Replace(rustCode, "https://engine.mirror.ad/rpc/<mirror-id>", sessionUrl, -1)
	if err := runCode(ctx, rustCode, "rust"); err != nil {
		return fmt.Errorf("error running code: %v", err)
	}
	fmt.Println("Time taken to run Rust:", time.Since(start)) // Log the time taken to run the TypeScript

	return nil
}

func RunTsCodeExec(
	ctx context.Context,
) error {
	start := time.Now()
	tsCode := `import {
		address,
		lamports,
		createTransaction,
		createSolanaClient,
		signTransactionMessageWithSigners,
		generateKeyPairSigner,
		airdropFactory,
	  } from "gill";
	  // @ts-ignore
	  import { getTransferSolInstruction } from "gill/programs";
	  
	  const { rpc, rpcSubscriptions, sendAndConfirmTransaction } = createSolanaClient({
		urlOrMoniker: "https://engine.mirror.ad/rpc/<mirror-id>",
	  });
	  
	  const signer = await generateKeyPairSigner();
	  await airdropFactory({ rpc, rpcSubscriptions })({
		commitment: "confirmed",
		lamports: lamports(5_000_000n),
		recipientAddress: signer.address,
	  });
	  
	  const destination = address("nick6zJc6HpW3kfBm4xS2dmbuVRyb5F3AnUvj5ymzR5");
	  
	  const { value: latestBlockhash } = await rpc.getLatestBlockhash().send();
	  
	  const tx = createTransaction({
		version: "legacy",
		feePayer: signer,
		instructions: [
		  getTransferSolInstruction({
			source: signer,
			destination,
			amount: lamports(1_000_000n),
		  }),
		],
		latestBlockhash,
	  });
	  
	  const signedTransaction = await signTransactionMessageWithSigners(tx);
	  await sendAndConfirmTransaction(signedTransaction);`

	sessionUrl, _, err := GetSession(ctx)
	if err != nil {
		return fmt.Errorf("error getting session: %v", err)
	}
	tsCode = strings.Replace(tsCode, "https://engine.mirror.ad/rpc/<mirror-id>", sessionUrl, -1)
	if err := runCode(ctx, tsCode, "typescript"); err != nil {
		return fmt.Errorf("error running code: %v", err)
	}
	fmt.Println("Time taken to run TypeScript:", time.Since(start)) // Log the time taken to run the TypeScript

	return nil
}

func GetSession(ctx context.Context) (string, string, error) {
	apiKey := ""
	url := "https://api.mirror.ad/blockchains/sessions"

	// Request payload (if needed, otherwise use nil for empty body)
	payload := []byte(`{}`)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "http://localhost:8899", "ws://localhost:8900", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_key", apiKey)
	req.Header.Set("user_id", uuid.NewString())

	// Execute HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "http://localhost:8899", "ws://localhost:8900", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "http://localhost:8899", "ws://localhost:8900", fmt.Errorf(
			"error fetching mirror instance: %d %s", resp.StatusCode, resp.Status)
	}

	// Parse JSON response
	var session struct {
		URL   string `json:"url"`
		WsURL string `json:"wsUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return "http://localhost:8899", "ws://localhost:8900", fmt.Errorf("error decoding response: %v", err)
	}

	// Validate session data
	if session.URL == "" || session.WsURL == "" {
		return "http://localhost:8899", "ws://localhost:8900", nil
	}

	return session.URL, session.WsURL, nil
}

func runCode(ctx context.Context, code string, lang string) error {
	apiKey := ""
	url := "https://api.mirror.ad/code-exec/" + lang

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]string{
		"code": code,
	})
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_key", apiKey)

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errorDetails map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorDetails); err != nil {
			return fmt.Errorf("error decoding error response: %v", err)
		}
		return fmt.Errorf("error running code: %v", errorDetails)
	}
	return nil
}
