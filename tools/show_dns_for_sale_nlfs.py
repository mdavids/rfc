import dns.resolver
import re
import sys
import urllib.parse

def show_for_sale_records(domain):
    """
    Retrieves and displays ONLY specific '_for-sale' TXT records with an fcod
    starting with "NLFS-" followed by 48 characters, based on
    draft-davids-forsalereg-08.txt.
    All other record types or malformed fcod values are silently skipped.
    """
    for_sale_subdomain = f"_for-sale.{domain}"
    print(f"Attempting to retrieve and display specific '_for-sale' records for: {for_sale_subdomain}\n")

    try:
        answers = dns.resolver.resolve(for_sale_subdomain, 'TXT')
    except dns.resolver.NoAnswer:
        print(f"No '_for-sale' TXT record found for {domain}. Domain is likely not declared for sale via this method.")
        return
    except dns.resolver.NXDOMAIN:
        print(f"The subdomain '{for_sale_subdomain}' does not exist. Domain is likely not declared for sale via this method.")
        return
    except Exception as e:
        print(f"An error occurred while querying DNS records: {e}")
        return

    decoded_records = []
    for rdata in answers:
        try:
            txt_record_content = rdata.strings[0].decode('ascii')
            decoded_records.append(txt_record_content)
        except UnicodeDecodeError:
            continue # Silently skip undecodable records
    
    decoded_records.sort()

    found_specific_nlfs_fcod = False
    print("--- Specific For Sale Information (NLFS- fcod) ---")
    record_count = 0

    # Specific fcod pattern: "NLFS-" followed by exactly 48 characters
    # We use '.' to match any character for robustness as requested,
    # but still expect base64 characters for the special string match.
    nlfs_fcod_pattern = re.compile(r"^NLFS-(.{48})$")
    
    # The specific 48-character string that triggers the special action
    specific_identifier_string = "NGYyYjEyZWYtZTUzYi00M2U0LTliNmYtNTcxZjBhMzA2NWQy"

    for txt_record_content in decoded_records:
        version_tag = "v=FORSALE1;"
        
        # Apply robustness principle: allow optional whitespace after the version tag.
        if not txt_record_content.startswith(version_tag):
            continue # Skip records not starting with the correct version tag

        content_part = txt_record_content[len(version_tag):].strip()

        # ONLY consider records that are exactly an 'fcod=' tag
        # We enforce that there's only one tag-value pair and it's an fcod.
        if not content_part.startswith("fcod="):
            continue # Skip if not an fcod tag
        
        # Extract the fcod value
        fcod_value = content_part[len("fcod="):]

        # Now, specifically check if this fcod value matches the NLFS- pattern
        nlfs_match = nlfs_fcod_pattern.match(fcod_value)
        
        if nlfs_match:
            record_count += 1
            print(f"\nRecord {record_count}:")
            found_specific_nlfs_fcod = True
            
            identifier_part = nlfs_match.group(1) # The 48 characters after "NLFS-"

            print(f"  For Sale Code (NLFS- format found): {fcod_value}")
            
            if identifier_part == specific_identifier_string:
                print(f"  Action: Visit https://www.sidn.nl/en/landing-page-buying-and-selling-example?domain={domain} for more information.")
            else:
                print(f"  Warning: Specific fcod not found in database.")
        # If it's an fcod but doesn't match the NLFS- pattern, it's silently skipped as per requirements.

    if not found_specific_nlfs_fcod:
        print("\nNo '_for-sale' TXT records with a specific 'NLFS-' fcod pattern were found.")
        print("Other types of '_for-sale' records (ftxt, furi, or fcod not matching 'NLFS-{48}')")
        print("are ignored by this script. Use 'check_dns_for_sale.py' for full validation.")

# --- Command-line argument handling ---
if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python show_dns_for_sale.py <domain_name>")
        print("Example: python show_dns_for_sale.py example.com")
        sys.exit(1)

    domain_to_show = sys.argv[1]
    show_for_sale_records(domain_to_show)
