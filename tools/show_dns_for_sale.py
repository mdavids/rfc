import dns.resolver
import re
import sys # Import the sys module for command-line arguments
import urllib.parse # For potential URL parsing if needed, though split will suffice for scheme

def show_for_sale_records(domain):
    """
    Retrieves and displays valid '_for-sale' TXT records for a domain name
    in a human-readable format, based on draft-davids-forsalereg-08.txt.
    Invalid or malformed records are silently skipped (use check_dns_for_sale.py for full validation).
    """
    for_sale_subdomain = f"_for-sale.{domain}"
    print(f"\nAttempting to retrieve and display valid '_for-sale' records for: {for_sale_subdomain}\n")

    try:
        answers = dns.resolver.resolve(for_sale_subdomain, 'TXT')
    except dns.resolver.NoAnswer:
        print(f"No '_for-sale' TXT record found for {domain}. Domain is likely not declared for sale via this method.")
        return
    except dns.resolver.NXDOMAIN:
        print(f"The domain '{domain}' or subdomain '{for_sale_subdomain}' does not exist. Domain is likely not declared for sale via this method, if it exists at all")
        return
    except Exception as e:
        print(f"An error occurred while querying DNS records: {e}")
        return

    # Decode and collect all TXT record strings that can be ASCII decoded
    decoded_records = []
    for rdata in answers:
        try:
            txt_record_content = rdata.strings[0].decode('ascii')
            decoded_records.append(txt_record_content)
        except UnicodeDecodeError:
            # Silently skip records that cannot be decoded as ASCII for display purposes.
            # check_dns_for_sale.py will report these errors.
            continue
    
    # Sort the decoded records for consistent display order
    decoded_records.sort()

    if not decoded_records:
        print("No valid (decodeable) '_for-sale' TXT records found.")
        print("Please refer to check_dns_for_sale.py for full validation details if unexpected.")
        return

    found_displayable_record = False
    print("--- For Sale Information ---")
    record_count = 0

    for txt_record_content in decoded_records:
        version_tag = "v=FORSALE1;"
        
        # Apply robustness principle: allow optional whitespace after the version tag.
        if not txt_record_content.startswith(version_tag):
            # Skip records not starting with the correct version tag
            continue

        # Extract content part AFTER the version tag and strip leading/trailing whitespace
        content_part = txt_record_content[len(version_tag):].strip()

        # Check for multiple tag-value pairs - for display, we only show the first valid one if multiple exist.
        # check_dns_for_sale.py already flags this as an error.
        tags_found_in_content = re.findall(r'(fcod=|ftxt=|furi=)', content_part)
        if len(tags_found_in_content) > 1:
            # For display purposes, we might just try to parse the first one or skip.
            # Given the RFC states "MUST NOT contain more than one", we'll skip for clear display.
            continue # Skip this record, as its format is ambiguous for friendly display

        record_count += 1
        #print(f"\nRecord {record_count}:")
        print("\nUsable _for-sale record found:")

        if not content_part:
            print("  This domain is declared for sale.")
            print("  No specific contact or sale information provided in this record.")
            found_displayable_record = True
            continue # Move to next record

        fcod_match = re.match(r'fcod=(.*)', content_part)
        ftxt_match = re.match(r'ftxt=(.*)', content_part)
        furi_match = re.match(r'furi=(.*)', content_part)

        if fcod_match:
            value = fcod_match.group(1)
            # Basic validation check for display - if it looks wrong, don't show
            if 1 <= len(value) <= 239 and all(0x20 <= ord(c) <= 0x7E for c in value):
                # TODO: the ABNF has it a little different
                print(f"  For Sale Code: {value}")
                found_displayable_record = True
            else:
                print("  This domain is declared for sale.")
                print("  Specific content details are to be considered absent.")
                print("  Please refer to check_dns_for_sale.py for full validation details if unexpected.")
                print(f"  [Skipped] Invalid fcod value found: \"{value}\" (check_dns_for_sale.py for details).")
        elif ftxt_match:
            value = ftxt_match.group(1)
            # Basic validation for display - check length and disallowed chars
            if 1 <= len(value) <= 239 and not re.search(r'["\\]', value):
                print(f"  For Sale Text: {value}")
                found_displayable_record = True
            else:
                print("  This domain is declared for sale.")
                print("  Specific content details are to be considered absent.")
                print("  Please refer to check_dns_for_sale.py for full validation details if unexpected.")            
                print(f"  [Skipped] Invalid ftxt value found: \"{value}\" (check_dns_for_sale.py for details).")
        elif furi_match:
            uri = furi_match.group(1)
            # Basic URI validation for display
            uri_pattern = re.compile(r"^[a-zA-Z][a-zA-Z0-9+.-]*:/{0,2}[%a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=.]*$")
            if uri_pattern.match(uri):
                parsed_scheme = uri.split(':')[0].lower()
                #print(f"  For Sale URI: {uri}")
                if parsed_scheme == 'http' or parsed_scheme == 'https':
                    print("  Action: Visit this URL for more information or to make an offer:")
                    print(f"  {uri}")
                elif parsed_scheme == 'mailto':
                    # Optional: extract email address if desired
                    email_address = uri[len('mailto:'):]
                    print(f"  Action: Email us at {email_address} for more information.")
                elif parsed_scheme == 'tel':
                    # Optional: extract phone number if desired
                    phone_number = uri[len('tel:'):]
                    print(f"  Action: Call us at {phone_number} for more information.")
                else:
                    print("  Action: Use this URI for more information:")
                    print(f"  {uri}")
                    
                found_displayable_record = True
            else:
                print("  This domain is declared for sale.")
                print("  Specific content details are to be considered absent.")
                print("  Please refer to check_dns_for_sale.py for full validation details if unexpected.")
                print(f"  [Skipped] Invalid furi value found: \"{uri}\" (check_dns_for_sale.py for details).")
        else:
            # If a record has a valid version tag but an unparseable content tag,
            # it's considered valid for "for sale" but content cannot be displayed.
            print("  This domain is declared for sale.")
            print("  Cannot display specific content details for this record (malformed tag or empty).")
            print("  Please refer to check_dns_for_sale.py for full validation details if unexpected.")
            found_displayable_record = True # Still considered a "for sale" record in principle

    if not found_displayable_record and record_count > 0:
        print("\nAll found '_for-sale' records were malformed or contained multiple tags and could not be displayed clearly.")
        print("Please run 'check_dns_for_sale.py' for detailed validation errors.")
    elif not found_displayable_record and record_count == 0:
        print("\nNo valid '_for-sale' TXT records were found or able to be displayed.")
        print("Please refer to check_dns_for_sale.py for full validation details if unexpected.")


# --- Command-line argument handling ---
if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python show_dns_for_sale.py <domain_name>")
        print("Example: python show_dns_for_sale.py example.com")
        sys.exit(1) # Exit with an error code

    domain_to_show = sys.argv[1]
    show_for_sale_records(domain_to_show)
