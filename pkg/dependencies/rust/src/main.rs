use solana_client::rpc_client::RpcClient;
use solana_sdk::commitment_config::CommitmentConfig;

fn main() -> anyhow::Result<()> {
    let client = RpcClient::new_with_commitment(
        String::from("https://engine.mirror.ad/rpc/3bcc19a3-9e24-44b8-a956-bc7fa250f4d0"),
        CommitmentConfig::confirmed(),
    );

    let data_len = 1500;
    let rent_exemption_amount = client.get_minimum_balance_for_rent_exemption(data_len)?;

    println!("{rent_exemption_amount}");

    Ok(())
}
