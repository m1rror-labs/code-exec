/**
 * Program IDL in camelCase format in order to be used in JS/TS.
 *
 * Note that this is only a type helper and is not the actual IDL. The original
 * IDL can be found at `target/idl/cpi.json`.
 */
export type Cpi = {
  "address": "9AvUNHjxscdkiKQ8tUn12QCMXtcnbR9BVGq3ULNzFMRi",
  "metadata": {
    "name": "cpi",
    "version": "0.1.0",
    "spec": "0.1.0",
    "description": "Created with Anchor"
  },
  "instructions": [
    {
      "name": "solTransfer",
      "discriminator": [
        135,
        254,
        247,
        202,
        217,
        48,
        184,
        165
      ],
      "accounts": [
        {
          "name": "sender",
          "writable": true,
          "signer": true
        },
        {
          "name": "recipient",
          "writable": true
        },
        {
          "name": "systemProgram",
          "address": "11111111111111111111111111111111"
        }
      ],
      "args": [
        {
          "name": "amount",
          "type": "u64"
        }
      ]
    }
  ]
};
