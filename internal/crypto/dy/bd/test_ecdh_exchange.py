from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.kdf.hkdf import HKDF
from cryptography.hazmat.primitives import serialization
from cryptography.x509 import load_pem_x509_certificate

def derive_ecdh_key(private_key_pem_path, peer_cert_pem_path):
    """
    Derives a shared key using ECDH and HKDF, similar to the JavaScript implementation.

    Args:
        private_key_pem_path (str): Path to the local private key PEM file.
        peer_cert_pem_path (str): Path to the peer's certificate PEM file.

    Returns:
        bytes: The derived 32-byte key.
    """
    # 1. Load local private key
    with open(private_key_pem_path, "rb") as f:
        private_key = serialization.load_pem_private_key(
            f.read(),
            password=None
        )

    # 2. Load peer's certificate and extract public key
    with open(peer_cert_pem_path, "rb") as f:
        peer_cert = load_pem_x509_certificate(f.read())
        peer_public_key = peer_cert.public_key()

    # 3. Perform ECDH to get the shared secret
    shared_secret = private_key.exchange(ec.ECDH(), peer_public_key)

    # 4. Apply HKDF to the shared secret
    derived_key = HKDF(
        algorithm=hashes.SHA256(),
        length=32,
        salt=None,
        info=None
    ).derive(shared_secret)

    return derived_key

if __name__ == "__main__":
    try:
        # These files should exist in the same directory
        my_private_key_file = "my_private_key.pem"
        peer_certificate_file = "peer_certificate.crt"

        derived_key_bytes = derive_ecdh_key(my_private_key_file, peer_certificate_file)

        print(f"Derived Key (hex): {derived_key_bytes.hex()}")

    except FileNotFoundError as e:
        print(f"Error: {e}. Please make sure '{my_private_key_file}' and '{peer_certificate_file}' exist.")
    except Exception as e:
        print(f"An error occurred: {e}")
