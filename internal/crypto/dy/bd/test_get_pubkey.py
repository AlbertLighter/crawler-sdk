import base64
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec

def get_pub_key_base64(private_key_pem: str = None, public_key_pem: str = None) -> str:
    """
    Extracts the raw, uncompressed public key and returns it as a base64 string.
    This mimics the behavior of the getPubKeyBase64 function in the provided JS.

    Args:
        private_key_pem (str, optional): The private key in PEM format.
        public_key_pem (str, optional): The public key in PEM format.

    Returns:
        str: The base64 encoded uncompressed public key, or an empty string if unable to extract.
    """
    public_key = None
    try:
        if public_key_pem:
            # Load public key from PEM string
            public_key = serialization.load_pem_public_key(public_key_pem.encode('utf-8'))
        elif private_key_pem:
            # Load private key from PEM string
            private_key = serialization.load_pem_private_key(
                private_key_pem.encode('utf-8'),
                password=None
            )
            # Derive the public key from the private key
            public_key = private_key.public_key()

        if public_key and isinstance(public_key, ec.EllipticCurvePublicKey):
            # Get the raw uncompressed public key point (starts with 0x04)
            # This corresponds to the 'rawHex' in the JavaScript code.
            raw_public_key_bytes = public_key.public_bytes(
                encoding=serialization.Encoding.X962,
                format=serialization.PublicFormat.UncompressedPoint
            )
            # Base64 encode the raw bytes and return as a string
            return base64.b64encode(raw_public_key_bytes).decode('utf-8')

    except Exception:
        # If any error occurs (e.g., invalid key format), return an empty string
        return ""

    return ""

if __name__ == "__main__":
    private_key_file = "my_private_key.pem"
    try:
        with open(private_key_file, "r") as f:
            private_key_pem_content = f.read()

        # --- Test Case 1: Extract public key from a Private Key PEM ---
        print(f"--- Testing with '{private_key_file}' ---")
        b64_pub_key = get_pub_key_base64(private_key_pem=private_key_pem_content)
        print(f"Extracted Base64 Public Key: {b64_pub_key}")

        # For P-256, the decoded length should be 65 bytes:
        # 1 byte for the 0x04 prefix (uncompressed) + 32 bytes for X + 32 bytes for Y
        if b64_pub_key:
            decoded_key = base64.b64decode(b64_pub_key)
            print(f"Decoded length: {len(decoded_key)} bytes")
            # The hex of the decoded key should start with '04'
            print(f"Decoded key starts with '04': {decoded_key.hex().startswith('04')}")

    except FileNotFoundError:
        print(f"Error: Test file '{private_key_file}' not found. Please make sure it exists.")
    except Exception as e:
        print(f"An error occurred during testing: {e}")
