# fohnhab
Fohnhab is a poc implementation of some of the crypto features within the go standard library

## Getting Started
Make sure you have go installed on your computer > 1.8
You can visit the go offical documentation in order to determine your version or installation: [golang](https://golang.org/doc/install)

Fohnhab uses dep as a package manager. After forking and pulling down the repo into your go path
make sure you have dep installed: [dep](https://github.com/golang/dep)

After dep is confirmed to be installed you can run `dep ensure` at the fohnhab root directory to install the dependencies

Fohnhab also uses ginkgo for its testing suite of choice. You can get started with ginkgo here: [ginkgo](https://github.com/onsi/ginkgo)

## Design Notes
Generate Keys
Encipher
Decipher
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
- Symmetric key encipher commands use 256-bit keys (analogous calls to decrypt use the inverse function) 
- In addition to enciphering plaintext to produce ciphertext, it computes an authentication tag over the cipher text and and any additionald ata opver which authentication is required
- Authentication tags help ensure that the data is from the purported source and that the ciphertext/AAD have not been modified.

There are two digital signature schemes utilized in AWS KMS

- The elliptic curve digital signature algorithm(ECDSA)
- RSA

All service host entities have an elliptic curve digital signature algorithm key pair
They perform ECDSA as defined here <https://tools.ietf.org/html/rfc5753> and <http://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.180-4.pdf> SHA384

Keys are generated on the curve <http://www.secg.org/sec2-v2.pdf>

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
AES-GCM is only constant time if the AES Instructions are installed on your machine's processor

## Precaution
AEADs should not be used to encrypt large amounts of data in one go. The API is designed to discourage this.
Encrypting large amounts of data in a single operation means that 
    a) either all the data has to be held in memory or 
    b) the API has to operate in a streaming fashion, by returning unauthenticated plaintext.

Returning unauthenticated data is dangerous it's not hard to find people on the internet suggesting things like gpg -d your_archive.tgz.gpg | tar xzbecause the gpg command also provides a streaming interface.
With constructions like AES-GCM it's, of course, very easy to manipulate the plaintext at will if the application doesn't authenticate it before processing. 
Even if the application is careful not to "release" plaintext to the UI until the authenticity has been established, a streaming design exposes more program attack surface.

By normalising large ciphertexts and thus streaming APIs, the next protocol that comes along is more likely to use them without realising the issues and thus the problem persists.

Preferably, plaintext inputs would be chunked into reasonably large parts (say 16KiB) and encrypted separately. 
The chunks only need to be large enough that the overhead from the additional authenticators is negligible. 
With such a design, large messages can be incrementally processed without having to deal with unauthenticated plaintext, and AEAD APIs can be safer. 
(Not to mention that larger messages can be processed since AES-GCM, for one, has a 64GiB limit for a single plaintext.)

Some thought is needed to ensure that the chunks are in the correct order, 
i.e. by counting nonces, that the first chunk should be first, 
i.e. by starting the nonce at zero, and that the last chunk should be last, 
i.e. by appending an empty, terminator chunk with special additional data. But that's not hard.

For an example, see the chunking used in miniLock.

Even with such a design it's still the case that an attacker can cause the message to be detectably truncated. 
If you want to aim higher, an all-or-nothing transform can be used, although that requires two passes over the input and isn't always viable.
