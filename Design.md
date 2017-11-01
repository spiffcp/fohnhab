# Design Notes

MVP
Generate Key
Encrypt
Decrypt
Implementation of a few fips approved algorythms
TDD

## Vocabulary

HSA - Hardened Security Appliance
AES-GCM - Authenticated encryption scheme
AAD - additionally authenticated data
ECDSA - (Elliptic Curve Digital signature algorithm)
RSA - (Rivest–Shamir–Adleman)
AES - Advanced Encryption Standard

## KMS White paper notes

Key generation is performed on dedicated HSA

- Phyisical Devices without a virtualization layer
- Hybrid random number generator
- Random Number generator is seeded with system entropy and then updated periodically
- Symmetric key encrypt commands use 256-bit keys (analogous calls to decrypt use the inverse function) 
- In addition to encrypting plaintext to produce ciphertext, it computes an authentication tag over the cipher text and and any additionald ata opver which authentication is required
- Authentication tags help ensure that the data is from the purported source and that the ciphertext/AAD have not been modified.

There are two digital signature schemes utilized in AWS KMS

- The elliptic curve digital signature algorithm(ECDSA)
- RSA

All service host entities have an elliptic curve digital signature algorithm key pair
They perform ECDSA as defined here <https://tools.ietf.org/html/rfc5753> and <http://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.180-4.pdf> SHA384

Keys are generated on the curve <http://www.secg.org/sec2-v2.pdf>

## Generate Key Implementation Rijndael Algorithm (AES) in Galois Counter Mode (AES-GCM for short)

256bit cipher key
AES-256

## Other Resources

<https://csrc.nist.gov/csrc/media/publications/fips/140/2/final/documents/fips1402annexa.pdf>
<https://csrc.nist.gov/csrc/media/publications/fips/197/final/documents/fips-197.pdf>
<https://csrc.nist.gov/publications/detail/sp/800-38d/final>
<https://tools.ietf.org/html/rfc5753>
<http://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-56Ar2.pdf>


## golang cryto/rand implementation
On Linux, Reader uses getrandom(2) if available, /dev/urandom otherwise.
On OpenBSD, Reader uses getentropy(2).
On other Unix-like systems, Reader reads from /dev/urandom.
On Windows systems, Reader uses the CryptGenRandom API.
