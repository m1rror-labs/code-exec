package typescript

import (
	"testing"
)

func TestExecute(t *testing.T) {
	t.Skip()
	runtime := NewRuntime()

	code := `
	import {
		airdropFactory,
		appendTransactionMessageInstructions,
		createSolanaRpc,
		createSolanaRpcSubscriptions,
		createTransactionMessage,
		generateKeyPairSigner,
		getSignatureFromTransaction,
		lamports,
		pipe,
		sendAndConfirmTransactionFactory,
		setTransactionMessageFeePayerSigner,
		setTransactionMessageLifetimeUsingBlockhash,
		signTransactionMessageWithSigners,
	  } from "@solana/kit";
	  import { getCreateAccountInstruction } from "@solana-program/system";
	  import {
		getInitializeMintInstruction,
		getMintSize,
		TOKEN_2022_PROGRAM_ADDRESS,
	  } from "@solana-program/token-2022";
	  
	  // Create Connection, local validator in this example
	  const rpc = createSolanaRpc("https://engine.mirror.ad/rpc/85e01483-a7f7-4cf7-9192-c1db0690764f");
	  const rpcSubscriptions = createSolanaRpcSubscriptions("wss://engine.mirror.ad/rpc/85e01483-a7f7-4cf7-9192-c1db0690764f");
	  
	  // Generate keypairs for fee payer
	  const feePayer = await generateKeyPairSigner();
	  
	  // Fund fee payer
	  await airdropFactory({ rpc, rpcSubscriptions })({
		recipientAddress: feePayer.address,
		lamports: lamports(1_000_000_000n),
		commitment: "confirmed",
	  });
	  
	  // Generate keypair to use as address of mint
	  const mint = await generateKeyPairSigner();
	  
	  // Get default mint account size (in bytes), no extensions enabled
	  const space = BigInt(getMintSize());
	  
	  // Get minimum balance for rent exemption
	  const rent = await rpc.getMinimumBalanceForRentExemption(space).send();
	  
	  // Instruction to create new account for mint (token 2022 program)
	  // Invokes the system program
	  const createAccountInstruction = getCreateAccountInstruction({
		payer: feePayer,
		newAccount: mint,
		lamports: rent,
		space,
		programAddress: TOKEN_2022_PROGRAM_ADDRESS,
	  });
	  
	  // Instruction to initialize mint account data
	  // Invokes the token 2022 program
	  const initializeMintInstruction = getInitializeMintInstruction({
		mint: mint.address,
		decimals: 9,
		mintAuthority: feePayer.address,
	  });
	  
	  const instructions = [createAccountInstruction, initializeMintInstruction];
	  
	  // Get latest blockhash to include in transaction
	  const { value: latestBlockhash } = await rpc.getLatestBlockhash().send();
	  
	  // Create transaction message
	  const transactionMessage = pipe(
		createTransactionMessage({ version: 0 }), // Create transaction message
		(tx) => setTransactionMessageFeePayerSigner(feePayer, tx), // Set fee payer
		(tx) => setTransactionMessageLifetimeUsingBlockhash(latestBlockhash, tx), // Set transaction blockhash
		(tx) => appendTransactionMessageInstructions(instructions, tx), // Append instructions
	  );

	  console.log("1")
	  
	  // Sign transaction message with required signers (fee payer and mint keypair)
	  const signedTransaction =
		await signTransactionMessageWithSigners(transactionMessage);
		console.log("2")
	  // Send and confirm transaction
	  await sendAndConfirmTransactionFactory({ rpc, rpcSubscriptions })(
		signedTransaction,
		{ commitment: "confirmed" },
	  );
	  console.log("")
	  // Get transaction signature
	  const transactionSignature = getSignatureFromTransaction(signedTransaction);
	  
	  console.log("Mint Address:", mint.address);
	  console.log("Transaction Signature:", transactionSignature);
	`
	output, err := runtime.ExecuteCode(code)
	if err != nil {
		t.Error(err)
	}
	t.Fatal(output)
}
