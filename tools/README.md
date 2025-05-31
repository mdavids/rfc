# Tools to play and test the _for-sale draft

## check_dns_for_sale.py
~~~
(venv) mdavids@iMac-van-Marco forsale % python3 ./check_dns_for_sale.py 

Checking for TXT records for: _for-sale.testdns.nl

  [DEBUG] Failed to decode TXT record content as ASCII: 'ascii' codec can't decode byte 0x80 in position 79: ordinal not in range(128). Content might contain non-ASCII characters.
  Raw record content: b'v=FORSALE1;fcod=SNAG-an fcod with characters below %x20:[\x15][\x00], above %x7E:[\x7f][\x80]'
  Skipping this raw record for further processing due to decoding error.
  [DEBUG] Failed to decode TXT record content as ASCII: 'ascii' codec can't decode byte 0x80 in position 77: ordinal not in range(128). Content might contain non-ASCII characters.
  Raw record content: b'v=FORSALE1;ftxt=not recommended characters below %x20:[\x15][\x00], above %x7E:[\x7f][\x80]'
  Skipping this raw record for further processing due to decoding error.

Found TXT record: "I'm not even a for sale TXT record, so ignore me!"
  [ERROR] Record does not start with the required 'v=FORSALE1;' version tag. Ignoring for validation.

Found TXT record: "idcode=NGYyYjEyZWYtZTUzYi00M2U0LTliNmYtNTcxZjBhMzA2NWQy"
  [ERROR] Record does not start with the required 'v=FORSALE1;' version tag. Ignoring for validation.

Found TXT record: "v=FORSALE0;fcod=FLAW-b3BlbiBzZXNhbWUK"
  [ERROR] Record does not start with the required 'v=FORSALE1;' version tag. Ignoring for validation.

Found TXT record: "v=FORSALE1"
  [ERROR] Record does not start with the required 'v=FORSALE1;' version tag. Ignoring for validation.

Found TXT record: "v=FORSALE1;"
  [INFO] Record contains only the version tag (after stripping whitespace). Processors MAY assume the domain is for sale.

Found TXT record: "v=FORSALE1;"
  [INFO] Record contains only the version tag (after stripping whitespace). Processors MAY assume the domain is for sale.

Found TXT record: "v=FORSALE1;"
  [INFO] Record contains only the version tag (after stripping whitespace). Processors MAY assume the domain is for sale.

Found TXT record: "v=FORSALE1;   fcod=HAZY-NGYyYjEyZWYtZTUzYi00M2U0LTliNmYtNTcxZjBhMzA2NWQy"
  Found fcod-tag with value: "HAZY-NGYyYjEyZWYtZTUzYi00M2U0LTliNmYtNTcxZjBhMzA2NWQy"
  [OK] fcod-value length is valid.

Found TXT record: "v=FORSALE1;   fcod=NLFS-FcodIsAlsoUnknownToUsAndNeedsRobustnessPrinciple"
  Found fcod-tag with value: "NLFS-FcodIsAlsoUnknownToUsAndNeedsRobustnessPrinciple"
  [OK] fcod-value length is valid.

Found TXT record: "v=FORSALE1;fcod="
  Found fcod-tag with value: ""
  [ERROR] fcod-value length is 0 octets, but MUST be between 1 and 239 octets.
  This record SHOULD be treated as if the tag-value pair were absent. Processors MAY assume the domain is for sale.

Found TXT record: "v=FORSALE1;fcod=ACME-b3BlbiBzZXNhbWUK"
  Found fcod-tag with value: "ACME-b3BlbiBzZXNhbWUK"
  [OK] fcod-value length is valid.

Found TXT record: "v=FORSALE1;fcod=LONG-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  Found fcod-tag with value: "LONG-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  [OK] fcod-value length is valid.

Found TXT record: "v=FORSALE1;fcod=NLFS-NGYyYjEyZWYtZTUzYi00M2U0LTliNmYtNTcxZjBhMzA2NWQy"
  Found fcod-tag with value: "NLFS-NGYyYjEyZWYtZTUzYi00M2U0LTliNmYtNTcxZjBhMzA2NWQy"
  [OK] fcod-value length is valid.

Found TXT record: "v=FORSALE1;fcod=NLFS-ThisFcodIsUnknownToUsAndWeCannotRedirectToAnURL."
  Found fcod-tag with value: "NLFS-ThisFcodIsUnknownToUsAndWeCannotRedirectToAnURL."
  [OK] fcod-value length is valid.

Found TXT record: "v=FORSALE1;fcod=NLFS-Too_Short_No_Base64_Invalid_NLFS_Synax"
  Found fcod-tag with value: "NLFS-Too_Short_No_Base64_Invalid_NLFS_Synax"
  [OK] fcod-value length is valid.

Found TXT record: "v=FORSALE1;fcod=SNAG-an fcod with escaped characters:["][\]"
  Found fcod-tag with value: "SNAG-an fcod with escaped characters:["][\]"
  [OK] fcod-value length is valid.

Found TXT record: "v=FORSALE1;fcod=TRIP-not-a-syntax-error-but-very-confusing;ftxt=so-dont-do-this!"
  [WARNING] Record suggests there is more than one tag-value pair. This is not conformant to the specification.
  Found fcod-tag with value: "TRIP-not-a-syntax-error-but-very-confusing;ftxt=so-dont-do-this!"
  [OK] fcod-value length is valid.

Found TXT record: "v=FORSALE1;foo=bar"
  [ERROR] No valid content-tag (fcod, ftxt, furi) found or the tag-value structure is invalid.
  This record SHOULD be treated as if the tag-value pair were absent. Processors MAY assume the domain is for sale.

Found TXT record: "v=FORSALE1;ftxt="
  Found ftxt-tag with value: ""
  [ERROR] ftxt-value length is 0 octets, but MUST be between 1 and 239 octets.
  This record SHOULD be treated as if the tag-value pair were absent. Processors MAY assume the domain is for sale.
  [INFO] ftxt-value does not contain disallowed characters.

Found TXT record: "v=FORSALE1;ftxt=$$[:free - format - string:]$$; test,test."
  Found ftxt-tag with value: "$$[:free - format - string:]$$; test,test."
  [INFO] ftxt-value does not contain disallowed characters.

Found TXT record: "v=FORSALE1;ftxt=I am "
  Found ftxt-tag with value: "I am"
  [INFO] ftxt-value does not contain disallowed characters.

Found TXT record: "v=FORSALE1;ftxt=I am here twice"
  Found ftxt-tag with value: "I am here twice"
  [INFO] ftxt-value does not contain disallowed characters.

Found TXT record: "v=FORSALE1;ftxt=https://:a.uri.in.ftxt.is.not.like.a.furi.but.in.essence.just.text/"
  Found ftxt-tag with value: "https://:a.uri.in.ftxt.is.not.like.a.furi.but.in.essence.just.text/"
  [INFO] ftxt-value does not contain disallowed characters.

Found TXT record: "v=FORSALE1;ftxt=not recommended escaped characters:["][\]"
  Found ftxt-tag with value: "not recommended escaped characters:["][\]"
  [ERROR] ftxt-value contains disallowed characters (double quote or backslash): ", \.

Found TXT record: "v=FORSALE1;furi="
  Found furi-tag with URI: ""
  [ERROR] furi-value is empty. MUST contain exactly one URI.

Found TXT record: "v=FORSALE1;furi=data:text/html,<script>alert('hi');</script>"
  Found furi-tag with URI: "data:text/html,<script>alert('hi');</script>"
  [ERROR] furi-value does not have a valid URI structure or contains unencoded characters.
  [WARNING] URI scheme 'data' is not recommended. Recommended schemes are: http, https, mailto, tel.

Found TXT record: "v=FORSALE1;furi=example:foo"
  Found furi-tag with URI: "example:foo"
  [INFO] furi-value has a valid URI structure.
  [WARNING] URI scheme 'example' is not recommended. Recommended schemes are: http, https, mailto, tel.

Found TXT record: "v=FORSALE1;furi=https://example.nl/for-sale.txt"
  Found furi-tag with URI: "https://example.nl/for-sale.txt"
  [INFO] furi-value has a valid URI structure.
  [INFO] URI scheme 'https' is recommended.

Found TXT record: "v=FORSALE1;furi=mailto:demo.doesnotwork@example.nl"
  Found furi-tag with URI: "mailto:demo.doesnotwork@example.nl"
  [INFO] furi-value has a valid URI structure.
  [INFO] URI scheme 'mailto' is recommended.

Found TXT record: "v=FORSALE1;furi=tel:+1-201-555-0123"
  Found furi-tag with URI: "tel:+1-201-555-0123"
  [INFO] furi-value has a valid URI structure.
  [INFO] URI scheme 'tel' is recommended.

Found TXT record: "v=FORSALE1;lorumipsum"
  [ERROR] No valid content-tag (fcod, ftxt, furi) found or the tag-value structure is invalid.
  This record SHOULD be treated as if the tag-value pair were absent. Processors MAY assume the domain is for sale.

Validation of '_for-sale' TXT records completed.
~~~

## generate_for_sale_txt.py
~~~
(venv) mdavids@iMac-van-Marco forsale % python3 ./generate_for_sale_txt.py

--- Generate _for-sale TXT Record ---
This tool will help you create a DNS TXT record for a domain that is for sale,
following the specifications in draft-davids-forsalereg-08.txt.

Choose a content tag type (only one allowed per record):
1. fcod (For Sale Code: Short alphanumeric/text code)
2. ftxt (For Sale Text: Short text description)
3. furi (For Sale URI: Link to a sales page or contact info)
4. No additional content tag (Just the version tag)
Enter your choice (1-4): 2
Enter the ftxt value (1-239 visible ASCII characters, excluding " and \): Hello World - Buy Me !

--- Generated TXT Record ---
To be placed in your DNS zone file for the domain 'example.com':
_for-sale IN TXT "v=FORSALE1;ftxt=Hello World - Buy Me !"

Replace 'example.com' with your actual domain name.

Important Notes:
- Some DNS providers require the TXT record content to be wrapped in double quotes.
- Ensure no extra spaces are added inside the quotes when copying.
- After adding, it may take some time for DNS changes to propagate.
- You can verify the record using the 'check_dns_for_sale.py' script.
~~~

## show_dns_for_sale.py
~~~
(venv) mdavids@iMac-van-Marco forsale % python3 ./show_dns_for_sale.py example.nl

Attempting to retrieve and display valid '_for-sale' records for: _for-sale.example.nl

--- For Sale Information ---

Usable _for-sale record found:
  For Sale Code: NLFS-NGYyYjEyZWYtZTUzYi00M2U0LTliNmYtNTcxZjBhMzA2NWQy

Usable _for-sale record found:
  For Sale Text: See the URL for important information!

Usable _for-sale record found:
  Action: Visit this URL for more information or to make an offer:
  https://example.nl/for-sale.txt

Usable _for-sale record found:
  This domain is declared for sale.
  Cannot display specific content details for this record (malformed tag or empty).
  Please refer to check_dns_for_sale.py for full validation details if unexpected.
~~~

## show_dns_for_sale_nlfs.go
(`show_dns_for_sale_nlfs.py` does the same)

~~~
./show_dns_for_sale_nlfs example.nl
Attempting to retrieve and display specific '_for-sale' records for: _for-sale.example.nl

--- Specific For Sale Information (NLFS- fcod) ---

Record 1:
  For Sale Code (NLFS- format found): NLFS-NGYyYjEyZWYtZTUzYi00M2U0LTliNmYtNTcxZjBhMzA2NWQy
  Action: Visit https://www.sidn.nl/en/landing-page-buying-and-selling-example?domain=example.nl for more information.
~~~

