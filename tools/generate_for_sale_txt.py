import re

def validate_fcod(value):
    """
    Validates an fcod value according to RFC draft-davids-forsalereg-08.txt.
    fcod-value = 1*239OCTET
    """
    if not value:
        return False, "fcod value cannot be empty."
    
    # Since we're dealing with US-ASCII, len(value) directly gives octet length.
    if not (1 <= len(value) <= 239):
        return False, f"fcod value length is {len(value)} octets, but MUST be between 1 and 239 octets."
    
    # Check if all characters are within the visible US-ASCII range (excluding " and \ if applicable, though fcod doesn't explicitly exclude them)
    # The RFC 3.1 ABNF for TXT-char excludes DQUOTE and BACKSLASH for general TXT, but fcod isn't explicitly defined that way.
    # However, 5. Operational Guidelines recommends avoiding non-ASCII.
    # For simplicity and consistency with ftxt, we'll recommend visible US-ASCII.
    for char in value:
        if not (0x20 <= ord(char) <= 0x7E): # Visible ASCII range
            return False, f"fcod value contains non-ASCII or control characters: '{char}' (U+{ord(char):04x}). Recommended to use visible US-ASCII."
    
    return True, "fcod value is valid."

def validate_ftxt(value):
    """
    Validates an ftxt value according to RFC draft-davids-forsalereg-08.txt.
    ftxt-value = 1*239ftxt-char
    ftxt-char = %x20-21 / %x23-5B / %x5D-7E (visible US-ASCII excluding " and \)
    """
    if not value:
        return False, "ftxt value cannot be empty."

    if not (1 <= len(value) <= 239):
        return False, f"ftxt value length is {len(value)} octets, but MUST be between 1 and 239 octets."
    
    invalid_chars = re.findall(r'["\\]', value)
    if invalid_chars:
        return False, f"ftxt value contains disallowed characters (double quote '\"' or backslash '\\'): {', '.join(invalid_chars)}."
    
    for char in value:
        # Check against ftxt-char ABNF: %x20-21 (space, !), %x23-5B (#-Z), %x5D-7E (]-~)
        # This covers visible US-ASCII excluding " (0x22) and \ (0x5C).
        if not ((0x20 <= ord(char) <= 0x21) or \
                (0x23 <= ord(char) <= 0x5B) or \
                (0x5D <= ord(char) <= 0x7E)):
            return False, f"ftxt value contains an invalid character: '{char}' (U+{ord(char):04x}). Only visible US-ASCII excluding double quotes and backslashes are allowed."
            
    return True, "ftxt value is valid."

def validate_furi(uri):
    """
    Validates an furi value according to RFC draft-davids-forsalereg-08.txt.
    furi-value = URI; exactly one URI
    Only http, https, mailto and tel schemes RECOMMENDED
    """
    if not uri:
        return False, "furi value cannot be empty. It MUST contain exactly one URI."

    # RFC3986 Appendix A (URI Generic Syntax) is complex.
    # A simplified regex for basic URI structure, allowing percent-encoding.
    uri_pattern = re.compile(r"^[a-zA-Z][a-zA-Z0-9+.-]*:/{0,2}[%a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=.]*$")
    
    if not uri_pattern.match(uri):
        return False, "furi value does not have a valid URI structure or contains unencoded characters. Ensure proper percent-encoding for reserved characters if needed."

    # Check recommended schemes
    parsed_scheme = uri.split(':')[0]
    if parsed_scheme.lower() not in ['http', 'https', 'mailto', 'tel']:
        # This is a warning, not a hard error for generation, but we'll flag it.
        return True, f"furi value is valid, but URI scheme '{parsed_scheme}' is not recommended. Recommended schemes are: http, https, mailto, tel."

    return True, "furi value is valid."


def generate_for_sale_txt_record():
    """
    Interactively generates a '_for-sale' TXT record for DNS zone files.
    """
    print("--- Generate _for-sale TXT Record ---")
    print("This tool will help you create a DNS TXT record for a domain that is for sale,")
    print("following the specifications in draft-davids-forsalereg-08.txt.\n")

    version_tag = "v=FORSALE1;"
    tag_content = ""
    selected_tag_type = None

    while True:
        print("Choose a content tag type (only one allowed per record):")
        print("1. fcod (For Sale Code: Short alphanumeric/text code)")
        print("2. ftxt (For Sale Text: Short text description)")
        print("3. furi (For Sale URI: Link to a sales page or contact info)")
        print("4. No additional content tag (Just the version tag)")
        choice = input("Enter your choice (1-4): ")

        if choice == '1':
            selected_tag_type = "fcod"
            while True:
                value = input("Enter the fcod value (1-239 visible ASCII characters): ")
                is_valid, message = validate_fcod(value)
                if is_valid:
                    tag_content = f"fcod={value}"
                    break
                else:
                    print(f"[ERROR] {message}")
                    print("Please try again.")
            break
        elif choice == '2':
            selected_tag_type = "ftxt"
            while True:
                value = input("Enter the ftxt value (1-239 visible ASCII characters, excluding \" and \\): ")
                is_valid, message = validate_ftxt(value)
                if is_valid:
                    tag_content = f"ftxt={value}"
                    break
                else:
                    print(f"[ERROR] {message}")
                    print("Please try again.")
            break
        elif choice == '3':
            selected_tag_type = "furi"
            while True:
                value = input("Enter the furi value (a valid URI, e.g., 'https://www.example.com/buy-domain', 'mailto:sales@example.com'): ")
                is_valid, message = validate_furi(value)
                if is_valid:
                    tag_content = f"furi={value}"
                    if "URI scheme" in message: # This indicates it's a warning, not a hard error
                        print(f"[WARNING] {message}")
                    break
                else:
                    print(f"[ERROR] {message}")
                    print("Please try again.")
            break
        elif choice == '4':
            print("No additional content tag will be added.")
            break
        else:
            print("Invalid choice. Please enter a number between 1 and 4.")

    full_txt_record_value = f"{version_tag}{tag_content}"

    print("\n--- Generated TXT Record ---")
    print(f"To be placed in your DNS zone file for the domain 'example.com':")
    print(f"_for-sale IN TXT \"{full_txt_record_value}\"")
    print("\nReplace 'example.com' with your actual domain name.")
    print("\nImportant Notes:")
    print("- Some DNS providers require the TXT record content to be wrapped in double quotes.")
    print("- Ensure no extra spaces are added inside the quotes when copying.")
    print("- After adding, it may take some time for DNS changes to propagate.")
    print("- You can verify the record using the 'check_dns_for_sale.py' script.")

# Run the generator
generate_for_sale_txt_record()
