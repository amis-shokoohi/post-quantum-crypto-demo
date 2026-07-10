# Post-Quantum Cryptography Secure Chat Demo
This project is a simple demonstration of building an end-to-end encrypted chat session between two parties over TCP using **ML-KEM-768** for key establishment and **AES-256-GCM** for encryption.

<div align="center">
  <img
    src="https://github.com/user-attachments/assets/46dacab2-340a-49ad-8bb5-bd7a3dce64fd"
    alt="Usage demo"
  >
</div>

The handshake works as follows:

1. Bob generates an ML-KEM-768 key pair.
2. Bob sends his public key to Alice over TCP.
3. Alice encapsulates a shared secret using Bob's public key and sends the resulting ciphertext back.
4. Bob decapsulates the ciphertext to recover the same shared secret.
5. Both parties derive a symmetric AES-256 key from the shared secret using HKDF.
6. All subsequent chat messages are encrypted with AES-GCM before being sent over the TCP connection.

<div align="center">
  <img
    src="https://github.com/user-attachments/assets/621e868c-aa30-4797-8e3c-aceaa0442254"
    alt="Handshake"
    width="550"
  >
</div>

The goal of this project is to demonstrate how a post-quantum key encapsulation mechanism can be used to establish a shared secret and bootstrap an encrypted communication channel.

## Disclaimer

This project is intended **for educational purposes only** and should **not** be considered a production-ready secure messaging implementation.

To keep the code focused on the ML-KEM handshake and symmetric encryption, several aspects of a real-world secure protocol are intentionally omitted or simplified, including:

* **Authentication / Digital Signatures:** No digital signature algorithm (DSA) is used. As a result, the handshake is **not authenticated** and is vulnerable to man-in-the-middle attacks.
* **Hybrid Key Exchange:** Modern deployments of ML-KEM typically use a **hybrid key exchange**, combining ML-KEM with a classical algorithm such as **X25519**. This provides security against both classical and quantum adversaries during the transition to post-quantum cryptography. This demo intentionally uses ML-KEM-768 by itself to keep the protocol easy to understand.
* **Forward Secrecy and Session Management:** The project establishes a single session key for the lifetime of the connection and does not implement features such as key rotation, rekeying, session resumption, or ratcheting algorithms.
* **Production Hardening:** The project omits many features expected from production protocols, such as peer authentication, replay protection, protocol negotiation, versioning, and comprehensive error handling.
