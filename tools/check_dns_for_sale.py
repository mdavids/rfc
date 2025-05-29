import dns.resolver
import re

def check_for_sale_record(domain):
    """
    Checks the '_for-sale' TXT records for a domain name according to
    draft-davids-forsalereg-08.txt, applying the robustness principle.
    """
    for_sale_subdomain = f"_for-sale.{domain}"
    print(f"\nChecking for TXT records for: {for_sale_subdomain}\n")

    try:
        answers = dns.resolver.resolve(for_sale_subdomain, 'TXT')
    except dns.resolver.NoAnswer:
        print(f"No '_for-sale' TXT record found for {domain}. Domain is likely not for sale according to this specification.")
        return
    except dns.resolver.NXDOMAIN:
        print(f"The domain {domain} or subdomain '{for_sale_subdomain}' does not exist. Domain is likely not declared for sale via this method, if it exists at all")
        return
    except Exception as e:
        print(f"Error querying DNS records: {e}")
        return

    # Decode and collect all TXT record strings, handling potential decoding errors
    decoded_records = []
    for rdata in answers:
        try:
            # Decode using 'ascii' and handle potential errors
            # TODO: is this fcod= compatible ?
            txt_record_raw_content = rdata.strings[0]
            # Strip trailing/leading spaces from the raw byte string before decoding,
            # as these might cause issues with 'ascii' decoding if they are not standard space characters.
            # However, the RFC talks about characters after the version tag, so we'll strip *after* the version check.
            txt_record_content = txt_record_raw_content.decode('ascii')
            decoded_records.append(txt_record_content)
        except UnicodeDecodeError as e:
            print(f"  [DEBUG] Failed to decode TXT record content as ASCII: {e}. Content might contain non-ASCII characters.")
            print(f"  Raw record content: {rdata.strings[0]}")
            print("  Skipping this raw record for further processing due to decoding error.")
            continue # Move to the next TXT record
    
    # Sort the decoded records for consistent processing order
    decoded_records.sort()

    if not decoded_records:
        print("\nNo valid (decodeable) '_for-sale' TXT record found that conforms to the specification after initial decoding attempt.")
        return

    found_valid_record = False
    for txt_record_content in decoded_records:
        print(f"\nFound TXT record: \"{txt_record_content}\"")

        # 3.1. General Record Format - version tag MUST begin with "v=FORSALE1;"
        # Apply robustness principle: allow optional whitespace after the version tag.
        # We check for the presence of the tag at the start and then strip.
        version_tag = "v=FORSALE1;"
        if not txt_record_content.startswith(version_tag):
            print("  [ERROR] Record does not start with the required 'v=FORSALE1;' version tag. Ignoring for validation.")
            continue

        found_valid_record = True
        
        # Extract content part AFTER the version tag and strip leading/trailing whitespace
        # This is where the robustness principle for parsing extra spaces comes in.
        content_part = txt_record_content[len(version_tag):].strip()

        # 3.4. RRset Limitations - The RDATA [RFC9499] of each TXT record MUST consist of a single character-string [RFC1035]
        # This is implicitly handled by dnspython returning a list of strings, each being one character-string.
        # The RFC also states: "Each '_for-sale' TXT record MUST NOT contain more than one tag-value pair."
        
        # Check for multiple tag-value pairs within the *single* character-string
        # This regex looks for any of the defined tags followed by an equals sign
        tags_found_in_content = re.findall(r'(fcod=|ftxt=|furi=)', content_part)
        if len(tags_found_in_content) > 1:
            print("  [WARNING] Record suggests there is more than one tag-value pair. This is not conformant to the specification.") 
           
            # Continue processing to show individual tag issues, but mark as invalid overall for this specific record.
            
        if not content_part:
            print("  [INFO] Record contains only the version tag (after stripping whitespace). Processors MAY assume the domain is for sale.")
            continue

        # Regex for content tags
        # Using re.match to ensure it starts with the tag
        fcod_match = re.match(r'fcod=(.*)', content_part)
        ftxt_match = re.match(r'ftxt=(.*)', content_part)
        furi_match = re.match(r'furi=(.*)', content_part)

        if fcod_match:
            value = fcod_match.group(1)
            print(f"  Found fcod-tag with value: \"{value}\"")
            # 3.1. General Record Format - fcod-value = 1*239OCTET
            # Since we're decoding to ASCII, len(value) will be the octet length for ASCII chars.
            if not (1 <= len(value) <= 239):
                print(f"  [ERROR] fcod-value length is {len(value)} octets, but MUST be between 1 and 239 octets.")
                print("  This record SHOULD be treated as if the tag-value pair were absent. Processors MAY assume the domain is for sale.")
            else:
                print("  [OK] fcod-value length is valid.")
        elif ftxt_match:
            value = ftxt_match.group(1)
            print(f"  Found ftxt-tag with value: \"{value}\"")
            # 3.1. General Record Format - ftxt-char = %x20-21 / %x23-5B / %x5D-7E (excluding " and \)
            # 5. Operational Guidelines - recommended-char = %x20-21 / %x23-5B / %x5D-7E
            if not (1 <= len(value) <= 239): # Octet length check
                print(f"  [ERROR] ftxt-value length is {len(value)} octets, but MUST be between 1 and 239 octets.")
                print("  This record SHOULD be treated as if the tag-value pair were absent. Processors MAY assume the domain is for sale.")
            
            # Check for excluded characters: double quote (") and backslash (\)
            invalid_chars = re.findall(r'["\\]', value)
            if invalid_chars:
                print(f"  [ERROR] ftxt-value contains disallowed characters (double quote or backslash): {', '.join(invalid_chars)}.")
            else:
                print("  [INFO] ftxt-value does not contain disallowed characters.")
            
        elif furi_match:
            uri = furi_match.group(1)
            print(f"  Found furi-tag with URI: \"{uri}\"")
            # 3.1. General Record Format - furi-value = URI; exactly one URI
            # 3.2.3. furi - Only http, https, mailto and tel schemes RECOMMENDED
            if not uri:
                print("  [ERROR] furi-value is empty. MUST contain exactly one URI.")
            else:
                # Basic URI validation (not exhaustive, as RFC3986 Appendix A is complex)
                # This regex is a simplified check for URI structure
                # It also checks for percent-encoding requirements
                uri_pattern = re.compile(r"^[a-zA-Z][a-zA-Z0-9+.-]*:/{0,2}[%a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=.]*$")
                if not uri_pattern.match(uri):
                    print("  [ERROR] furi-value does not have a valid URI structure or contains unencoded characters.")
                else:
                    print("  [INFO] furi-value has a valid URI structure.")
                    
                # Check recommended schemes
                parsed_scheme = uri.split(':')[0]
                if parsed_scheme not in ['http', 'https', 'mailto', 'tel']:
                    print(f"  [WARNING] URI scheme '{parsed_scheme}' is not recommended. Recommended schemes are: http, https, mailto, tel.")
                else:
                    print(f"  [INFO] URI scheme '{parsed_scheme}' is recommended.")
        else:
            # "If a tag-value pair is present but invalid, this constitutes a syntax error and SHOULD be treated as if it were absent."
            # "In such cases, if the version tag itself is valid, processors MAY assume that the domain is for sale."
            print("  [ERROR] No valid content-tag (fcod, ftxt, furi) found or the tag-value structure is invalid.")
            print("  This record SHOULD be treated as if the tag-value pair were absent. Processors MAY assume the domain is for sale.")

    if not found_valid_record:
        print("\nNo valid '_for-sale' TXT record found that conforms to the specification.")
    else:
        print("\nValidation of '_for-sale' TXT records completed.")

# Use the script for testdns.nl
domain_to_check = "example.nl"
check_for_sale_record(domain_to_check)
