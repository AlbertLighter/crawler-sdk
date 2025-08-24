import hmac
import hashlib

def sign_with_hmac_sha256(key: bytes, message: bytes) -> str:
    """
    Signs a message with the given key using HMAC-SHA256.

    Args:
        key (bytes): The secret key.
        message (bytes): The message to sign.

    Returns:
        str: The hexadecimal representation of the signature.
    """
    signature = hmac.new(key, message, hashlib.sha256)
    return signature.hexdigest()

if __name__ == "__main__":
    # --- Example Usage ---

    # Example 1: Using bytes
    secret_key_bytes = b'mysecretkey'
    message_to_sign_bytes = b'This is the message to sign'
    
    hex_signature_bytes = sign_with_hmac_sha256(secret_key_bytes, message_to_sign_bytes)
    
    print(f"Secret Key (bytes): {secret_key_bytes}")
    print(f"Message (bytes): {message_to_sign_bytes}")
    print(f"HMAC-SHA256 Signature (hex): {hex_signature_bytes}")
    print("-" * 20)

    # Example 2: Using strings (must be encoded to bytes first)
    secret_key_str = ""
    message_to_sign_str = "sign"

    # Encode strings to bytes using UTF-8
    key_as_bytes = secret_key_str.encode('utf-8')
    message_as_bytes = message_to_sign_str.encode('utf-8')

    hex_signature_str = sign_with_hmac_sha256(key_as_bytes, message_as_bytes)

    print(f"Secret Key (str): '{secret_key_str}'")
    print(f"Message (str): '{message_to_sign_str}'")
    print(f"HMAC-SHA256 Signature (hex): {hex_signature_str}")
    bebeca83c9ad71a4ea4e625d45af67666e5113fd1e486abe20e445b6add69635