import base58
import json
import os

# Base58 encoded string
base58_string = "588FU4PktJWfGfxtzpAAXywSNt74AvtroVzGfKkVN1LwRuvHwKGr851uH8czM5qm4iqLbs1kKoMKtMJG4ATR7Ld2"

# Decode Base58 string
decoded_bytes = base58.b58decode(base58_string)

# Convert bytes to a list of integers
byte_list = list(decoded_bytes)

# Define output file path
output_path = os.path.expanduser("~/.config/solana/id.json")

# Write to file
with open(output_path, "w") as f:
    json.dump(byte_list, f)

print(f"Decoded bytes written to {output_path}")