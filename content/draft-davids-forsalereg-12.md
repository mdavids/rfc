%%%
# This is a comment - but only in this block
title = "The \"_for-sale\" Underscored and Globally Scoped DNS Node Name"
abbrev = "forsalereg"
ipr = "trust200902" 
# area = "Internet"
workgroup = "Internet Engineering Task Force (IETF)"
submissiontype = "IETF"
keyword = [""]
# https://www.rfc-editor.org/rfc/rfc7991#section-2.45.14
tocdepth = 3
# date = 2022-12-22T00:00:00Z

# See FAQ: "How Do I Create an Independent IETF Document?"
# https://mmark.miek.nl/post/faq/
[seriesInfo]
name = "Internet-Draft"
value = "draft-davids-forsalereg-12"
stream = "IETF"
status = "informational"

[[author]]
initials="M."
surname="Davids"
fullname="Marco Davids"
abbrev = "SIDN Labs"
organization = "SIDN Labs"
  [author.address]
  email = "marco.davids@sidn.nl"
  phone = "+31 26 352 5500"
  [author.address.postal]
  street = "Meander 501"
  city = "Arnhem"
  code = "6825 MD"
  country = "Netherlands"
%%%

<!-- TODO: alle comments nalopen, want het zijn niet alleen TODO en kijken of het opgeschoond kan/moet -->

<!-- hint: use Title Case everywhere -->

.# Abstract

This document defines an operational convention for using the reserved
underscored DNS leaf node name "\_for-sale" to indicate that the 
parent domain name is available for purchase. This approach offers 
the advantage of easy deployment without affecting ongoing operations. 
As such, the method can be applied to a domain name that is still in full use.

{removeInRFC="true"}
.# About This Document

This note is to be removed before publishing as an RFC.

This document contains several "Notes to the RFC Editor", including this section. 
These should be reviewed and resolved prior to publication.

{mainmatter}

# Introduction {#introsect}

Well-established services [@RFC3912; @RFC9083] exist to determine whether a domain name is registered. However, the fact that a domain name exists does not necessarily mean it
is unavailable; it may still be for sale.

Some registrars and other parties offer brokerage services between domain name holders and interested buyers.
Such services are of limited value when the domain name is not for sale, but they may be beneficial 
for domain names that are clearly being offered for sale.

This specification defines a lightweight method to ascertain whether a domain name, although registered, is available for purchase. It enables a domain name holder to add a reserved underscored
leaf node name [@!RFC8552] in the zone, indicating that the domain name is for sale.

The TXT RR type [@!RFC1035] created for this purpose **MUST** follow the formal definition of
(#conventions). Its content **MAY** contain a pointer, such as a Uniform Resource Identifier (URI) 
[@!RFC3986], or another string, allowing interested parties to obtain information or 
contact the domain name holder for further negotiations.

With due caution, such information can also be incorporated into automated availability services. When checking a domain name for availability, the service may indicate whether it is for sale and provide a pointer to the seller's information.

Note: In this document, the term "for sale" is used in a broad sense and 
**MAY** also refer to cases where the domain name is available for lease, 
or where the contractual right to use the domain name is offered to another party.

## Terminology

The key words "**MUST**", "**MUST NOT**", "**REQUIRED**", "**SHALL**", "**SHALL NOT**",
"**SHOULD**", "**SHOULD NOT**", "**RECOMMENDED**", "**NOT RECOMMENDED**", "**MAY**", and
"**OPTIONAL**" in this document are to be interpreted as described in BCP 14 [@!RFC2119] [@!RFC8174]
when, and only when, they appear in all capitals, as shown here.

# Rationale

There are undoubtedly more ways to address this problem space. The reasons for the approach defined in this document are primarily accessibility and simplicity. The indicator can be easily turned on and off at will and moreover, it is immediately deployable and does not require significant changes in existing services. This allows for a smooth introduction of the concept.

Furthermore, the chosen approach aligns with ethical considerations by promoting a more equitable domain aftermarket and minimizing potential for unintended commercial entanglements by registries, as detailed in (#ethicalconsids).

# Conventions {#conventions}

## General Record Format {#abnf}

Each "\_for-sale" TXT record **MUST** begin with a version tag, optionally followed by a string containing content that follows a simple "tag=value" syntax.

The formal definition of the record format, using ABNF [@!RFC5234; @!RFC7405], is as follows:

~~~
forsale-record  = forsale-version [forsale-content]
                  ; referred to as content or RDATA
                  ; in a single character-string

forsale-version = %s"v=FORSALE1;"
                  ; %x76.3D.46.4F.52.53.41.4C.45.31.3B
                  ; version tag, case sensitive, no spaces

forsale-content = fcod-pair / ftxt-pair / furi-pair / fval-pair
                  ; referred to as tag-value pairs
                  ; only one tag-value pair per record

fcod-pair       = fcod-tag fcod-value
ftxt-pair       = ftxt-tag ftxt-value
furi-pair       = furi-tag furi-value
fval-pair       = fval-tag fval-value
                  ; the tags are referred to as content tags
                  ; the values are referred to as content values

fcod-tag        = %s"fcod="
ftxt-tag        = %s"ftxt="
furi-tag        = %s"furi="
fval-tag        = %s"fval="
                  ; all content tags case sensitive lowercase

fcod-value      = 1*239OCTET
                  ; must be at least 1 OCTET

ftxt-value      = 1*239ftxt-char
ftxt-char       = %x20-21 / %x23-5B / %x5D-7E
                  ; excluding " and \ to avoid escape issues

furi-value      = URI
                  ; http, https, mailto and tel URI schemes
                  ; exactly one URI

URI             = <as defined in RFC3986, Appendix A>

fval-value      = fval-currency fval-amount
                  ; total length: 2 to 239 characters 
fval-currency   = 1*ALPHA
                  ; one or more uppercase letters (A-Z)
                  ; indicating (crypto)currency
                  ; e.g., USD, EUR, BTC, ETH
fval-amount     = int-part [ %x2E frac-part ]
                  ; integer part with optional fractional part
                  ; e.g., 0.00010
int-part        = 1*DIGIT
frac-part       = 1*DIGIT
~~~
<!-- hint: make sure [@!RFC3986 remains somewhere in the document-->
<!-- hint: double check on https://author-tools.ietf.org/abnf -->

See (#tagdefs) for more detailed format definitions per content tag type. 

Each "\_for-sale" TXT record **MUST NOT** contain more than one tag-value
pair, but multiple TXT records **MAY** be present in a single RRset.

Every tag-value pair in the RRset **MUST** be unique, but multiple 
instances of the same content tag **MAY** occur within a single RRset 
(e.g., two "fcod=" content tags, each with a different content value).


See (#rrsetlimits) for additional RRset limitations.

The **OPTIONAL** forsale-content provides information to interested parties as explained
in (#introsect). 

If the forsale-content is absent or invalid, but a valid version tag
is present, processors **SHOULD** assume that the domain is for sale. For
example:

```
_for-sale.example.com. IN TXT "v=FORSALE1;"
_for-sale.example.com. IN TXT "v=FORSALE1;fcod="
_for-sale.example.com. IN TXT "v=FORSALE1;foo=bar"
```

In such cases, processors **SHOULD** determine how to proceed. 
An approach might be to signal that the domain is for sale and 
to rely on traditional mechanisms such as WHOIS or RDAP to retrieve and present contact
information.

TXT records in the same RRset, but without a version tag, **MUST NOT** be interpreted or processed as a valid "\_for-sale" indicator. 
However, they may still offer some additional information for humans when considered alongside a valid
record. For example:

```
_for-sale.example.com. IN TXT "I am for sale"
_for-sale.example.com. IN TXT "v=FORSALE1;fcod=XX-NGYyYjEyZWY"
```

If no TXT records at a leaf node contain a valid version tag, processors **MUST** consider the node name invalid and discard it.

See (#contentlimits) for additional content limitations.

## Content Tag Type Definitions {#tagdefs}

A new IANA registry for known content tags is created in (#ianaconsid), with 
this document registering the initial set. Implementations **SHOULD** 
process only registered tags they support, and **MAY** ignore any others.

The following content tags are defined as the initial valid content tags.

<!-- author tip: there are two spaces behind the content tag, to enforce a new line -->
### fcod {#fcoddef}  
This content tag is intended to contain a code that is meaningful only to processors 
that understand its semantics. The content value MUST consist of at least one octet. 

The manner in which the "fcod=" content tag is used is determined by agreement
between cooperating parties.

For example, a domain name registry may allow registrars to enter a "for sale" URL into their system. 
From that URL, a unique code is generated. This code is inserted as the value of
the "fcod=" content tag of the "\_for-sale" TXT record of a domain name, as shown in the example below.

When a user checks the availability of the domain name using a registry-provided tool 
(e.g., a web interface), the domain name registry may use the code to redirect the user to the 
appropriate "for sale" URL, which may include a query component containing the domain name, for example:

~~~
https://forsale-url.example.com/acme?d=example.org
~~~

The rationale for this approach is that controlling parties retain 
authority over redirection URLs and any other information derived 
from the content tag, thereby preventing users from being sent 
to unintended or malicious destinations or from being presented 
with unintended content.

The following example shows a string encoded using base 64 [@?RFC4648] 
preceded by the prefix "ACME-" as the value of the content tag:

~~~
_for-sale IN TXT "v=FORSALE1;fcod=ACME-S2lscm95IHdhcyBoZXJl"
~~~

See the (#examples, use title) section for other possible uses of this
content tag.

Note: As an implementation consideration, when multiple parties are involved in 
the domain sale process and use the same mechanism, it may be difficult to identify 
the relevant content in an RRset. Adding a recognizable prefix to the content (e.g.,
"ACME-") is one possible approach. However, this is left to the implementor, 
as it is not enforced in this document. In this case, ACME would recognize its 
content tag and interpret it as intended. This example uses base 64 encoding 
to avoid escaping and ensure printable characters, though this is also not required.

### ftxt  
This content tag is intended to contain human-readable text that conveys information to interested parties. For example:

~~~
_for-sale IN TXT "v=FORSALE1;ftxt=Call for info."
~~~

While a single visible character is the minimum, it is **RECOMMENDED** to provide more context.

While a URI in this field is not syntactically prohibited, its 
interpretation as a URI is not guaranteed. Use of URIs in this 
field **SHOULD** be avoided in favor of the "furi=" content tag.

See (#fvalpar) for a way to explicitly indicate an asking price for easier machine parsing.

<!-- TODO https://www.rfc-editor.org/rfc/rfc7553.html noemen, of zelfs opnemen? -->
### furi  
This content tag is intended to contain a human-readable and machine-parseable URI that conveys information to interested parties.

While the syntax allows any URI scheme, only the following schemes are **RECOMMENDED** 
for use: `http` and `https` [@RFC9110], `mailto` [@RFC6068], and `tel` [@RFC3966].

The content value **MUST** contain exactly one URI. For example:

~~~
_for-sale IN TXT "v=FORSALE1;furi=https://example.com/foo%20bar"
~~~

URIs **MUST** conform to the syntax and encoding requirements specified in 
[@!RFC3986, section 2.1], including the percent-encoding of characters 
not allowed unencoded (e.g., spaces **MUST** be encoded as `%20` in a URI).

See the (#security, use title) section for possible risks.

### fval {#fvalpar}
This content tag is intended to contain human-readable and machine-parseable 
text that explicitly indicates an asking price in a certain currency, as opposed to 
the price being loosly incorporated in an "ftxt=" content tag. For example:

~~~
_for-sale IN TXT "v=FORSALE1;fval=EUR999"
~~~

## Content Limitations {#contentlimits}

The "\_for-sale" TXT record [@RFC8553, (see) section 2.1] **MUST** contain content deemed valid under this specification.

Any text suggesting that a domain is not for sale is invalid content. If a domain name is not or no longer for sale, 
a "\_for-sale" indicator **MUST NOT** exist. The presence of a valid "_for-sale" TXT record
**SHOULD** therefore be regarded as an indication that the domain name is for sale.

The existence of a "\_for-sale" leaf node does not obligate the holder to sell the domain name; 
it may have been published in error, or withdrawn later for other reasons.

This specification does not dictate the exact use of any content values in the "\_for-sale" TXT record.
Parties **MAY** use it in their tools, perhaps even by defining specific requirements that the content
value must meet. Content values can also be represented in a human-readable format for individuals to
interpret. See the (#examples, use title) section for clarification.

See (#guidelines) for additional guidelines.

## RRset Limitations {#rrsetlimits}

This specification does not define restrictions on the number of TXT records in the
RRset.

When multiple content TXT records are present, the processor **MAY** select one or more of them.

For example, a domain name registry might extract content from an RRset that includes 
a recognizable "fcod=" content tag and use it to direct visitors to a sales page as 
part of its services. An individual, on the other hand, might extract a 
phone number (if present) from a "furi=" tag in the same RRset and use it to contact a potential seller.

An example of such a combined record is provided in (#combiexample).

The RDATA [@RFC9499] of each TXT record **MUST** consist of a single character-string
[@!RFC1035] with a maximum length of 255 octets, in order to avoid the need to concatenate multiple
character-strings during processing. 

The following example illustrates an invalid TXT record due to the presence of multiple
character-strings:

~~~
_for-sale IN TXT "v=FORSALE1;" "ftxt=foo" "bar" "invalid"
~~~

## Wildcard Limitation

Wildcards are only interpreted as leaf names, so "\_for-sale.*.example." is not a valid wildcard and is non-conformant.
Hence, it is not possible to put all domains under a TLD for sale with just one TXT record.

The example below, however, shows a common use case where a "\_for-sale" leaf node exists alongside a
wildcard:

~~~
*         IN A    198.51.100.80
          IN AAAA 2001:DB8::80
_for-sale IN TXT  "v=FORSALE1;ftxt=Only $99 at ACME"
~~~

<!-- TODO dit is nieuwe tekst, goed checken en over nadenken nog! -->


## Placement of the Leaf Node Name

The "\_for-sale" leaf node name can essentially be placed at any level of
the DNS except in the in-addr.arpa. infrastructure TLD.

(#placements) illustrates this:

Name | Situation | Verdict
-----|-----------|--------
\_for-sale.example. | root zone | For sale
\_for-sale.aaa.example. | second level | For sale
\_for-sale.acme.bbb.example. | third level with public registry | For sale
\_for-sale.www.ccc.example. | third level without public registry | See note 1
\_for-sale.51.198.in-addr.arpa. | infrastructure TLD | See note 2
xyz.\_for-sale.example. | Invalid placement, not a leaf | non-conformant
Table: Placements of TXT record {#placements}

Note 1: 
When the "\_for-sale" leaf node is applied to a label under a subdomain, 
there may not be a public domain name registry [@?RFC8499] capable of properly recording the rights associated with that label. 
Nevertheless, this does not constitute a violation of this document. 
One possible approach is for the involved parties to establish a mutual agreement to formalize these rights.

Note 2:
If a "\_for-sale" leaf node were to appear under the .arpa infrastructure top-level 
domain, it might be interpreted as an offer to sell IP address space. 
However, such use is explicitly out of scope for this document, and processors
**MUST** ignore any such records.

# Additional Examples {#examples}

## Example 1: Code Format

A proprietary format, defined and used by agreement between parties - for example, 
a domain name registry and its registrars - without a clearly specified meaning for third parties.
For example, it may be used to automatically redirect visitors to a web page, as described in
(#fcoddef):

~~~
_for-sale IN TXT "v=FORSALE1;fcod=XX-aHR0cHM...wbGUuY29t"
~~~

The use of the "fcod=" content tag is, in principle, unrestricted, allowing implementers to define additional 
uses as needed. For example, it may convey arbitrary formatting or conditional display 
instructions, such as adding an extra banner (e.g., "eligibility criteria apply") or 
specifying a style, including color, font, emojis, or logos.

## Example 2: Free Text Format

Free format text, with some additional unstructured information, aimed at
being human-readable:

~~~
_for-sale IN TXT "v=FORSALE1;ftxt=Eligibility criteria apply."
~~~

The content in the following example could be malicious, but it is not in violation of this specification (see
the (#security, use title)):

~~~
_for-sale IN TXT "v=FORSALE1;ftxt=<script>...</script>"
~~~

## Example 3: URI Format

The holder of "example.com" wishes to signal that the domain is for sale and adds this record to the "example.com" zone:

~~~
_for-sale IN TXT "v=FORSALE1;furi=https://example.com/fs?d=eHl6"
~~~

An interested party notices this signal and can visit the URI mentioned for further information. The TXT record
may also be processed by automated tools, but see the (#security, use title) section for possible risks. 

As an alternative, a mailto: URI could also be used:

~~~
_for-sale IN TXT "v=FORSALE1;furi=mailto:hq@example.com?subject=foo"
~~~

Or a telephone URI:

~~~
_for-sale IN TXT "v=FORSALE1;furi=tel:+1-201-555-0123"
~~~

There can be a use case for these URIs, especially since WHOIS (or RDAP) often has privacy restrictions.
But see the (#privacy, use title) section for possible downsides.

## Example 4: Asking Price Format

Consists of an uppercase currency code (e.g., USD, EUR), followed by a
numeric amount. See (#guidelines) for additional guidelines.

In Bitcoins:

~~~
_for-sale IN TXT "v=FORSALE1;fval=BTC0.000010"
~~~

In US dollars:

~~~
_for-sale IN TXT "v=FORSALE1;fval=USD750"
~~~

## Example 5: Combinations {#combiexample}

An example of multiple valid TXT records from which a processor can choose:

~~~
_for-sale IN TXT "v=FORSALE1;furi=https://fs.example.com/"
          IN TXT "v=FORSALE1;ftxt=This domain name is for sale"
          IN TXT "v=FORSALE1;fval=EUR500"
          IN TXT "v=FORSALE1;fcod=ACME-ZGVhZGJlZWYx"
          IN TXT "v=FORSALE1;fcod=XYZ1-MTExLTIyMi0zMzMtNDQ0"
~~~

# Operational Guidelines {#guidelines}
## DNS Wildcards

DNS wildcards interact poorly with underscored names [@RFC8552, (see) section 1.4],
but they may still be encountered in practice, especially with operators who 
are not implementing this mechanism. This is why the version 
tag is a **REQUIRED** element: it allows processors to distinguish 
valid "\_for-sale" records from unrelated TXT records.

Nonetheless, any assumptions about the content of "\_for-sale" TXT 
records **SHOULD** be made with caution, particularly in edge 
cases where wildcard expansion - possibly combined with DNS aliases 
(e.g., CNAMEs) or redirections (e.g., DNAMEs [@?RFC6672]) - might 
result in misleading listings or unintended references to third-party domains.

## Character Set

For the "ftxt=" content tag, the content value **MUST** be limited to visible US-ASCII characters, 
excluding the double quote (") and backslash (\\).

In ABNF syntax, this would be:

~~~
forsale-content  = 0*244recommended-char
recommended-char = %x20-21 / %x23-5B / %x5D-7E
~~~

For the content value of the "fcod=" content tag, this is **RECOMMENDED**. 

For example, base 64 uses only characters within this range, and therefore conforms to 
this recommendation.

## Currency

While the ABNF for the "fval=" content value in (#abnf) allows flexibility
regarding the currency indication, it is **RECOMMENDED** to use a three-letter uppercase 
currency code, such as those listed in [@?ISO4217], followed by a numeric amount, 

~~~
fval-value    = fval-currency fval-amount
                ; total length: 4 to 239 characters 
fval-currency = 3ALPHA
                ; 3-letter uppercase currency code (A-Z)
                ; e.g., USD, EUR, BTC, ETH
fval-amount   = int-part [ "." frac-part ]
                ; integer part with optional fractional part
                ; e.g., 0.00010
int-part      = 1*DIGIT
                ; at least one digit before the decimal point
frac-part     = 1*DIGIT
                ; at least one digit after the decimal point
~~~

## TTLs

Long TTLs [@!RFC1035, (see) section 3.2.1] increase the risk of outdated data misleading buyers into thinking the domain is still
available. 

## Ambiguous Constructs

Ambiguous constructs in content values **SHOULD** be avoided, as illustrated by the following
example:

~~~
_for-sale IN TXT "v=FORSALE1;fcod=TRIP-confusing;ftxt=dont-do-this"
~~~

The above example is a valid "fcod=" content tag that includes the 
string ";ftxt=" in the content value, which may be confusing, 
as it does not actually represent an "ftxt=" content tag.

## Robustness

Because the format of the content part is not strictly defined in this
document, processors **MAY** apply the robustness principle of being 
liberal in what they accept. This also applies to space 
characters (`%x20`) immediately following the version tag.
Alternatively, parties may agree on a more strictly defined proprietary format
for the content value to reduce ambiguity.

## Scope of Application

Note that this mechanism relies on the domain name being resolvable in the DNS.
This is not guaranteed, for example during a redemption period, in pending delete status [@?STD69],
or when the domain is DNSSEC-signed but fails validation (i.e., has a bogus state).

# IANA Considerations {#ianaconsid} <!-- See RFC8126 -->

IANA has established the "Underscored and Globally Scoped DNS Node Names" registry [@!RFC8552; @IANA]. The underscored
leaf node name defined in this specification should be added as follows:

RR Type | _NODE NAME | Reference
--------|------------|-----------
TXT | \_for-sale | <this memo>
Table: Entry for the "Underscored and Globally Scoped DNS Node Names" registry

<NOTE TO RFC EDITOR: Adjust the text in the table above before publication with a citation for the (this) document making the addition as per RFC8552.>

<!-- INFO zie https://www.rfc-editor.org/rfc/rfc8726.html#name-creating-new-iana-registrie -->
<!-- INFO zie ook https://www.iana.org/help/protocol-registration -->
<!-- INFO en zie ook https://www.rfc-editor.org/rfc/rfc8126.html -->
<!-- INFO of deze: https://www.ietf.org/id/draft-baber-ianabis-rfc8126bis-00.html -->
<!-- TODO niet vergeten reference anchor op te ruimen indien alsnog niet nodig -->

A registry group called "The '_for-sale' Underscored and Globally Scoped DNS Node Name" [@?FORSALEREG] is to be created, 
along with a registry called "Content Tags" within it. This registry group will be
maintained by IANA.

<NOTE TO RFC EDITOR: Remove the text about the example registry below, prior to publication.>

An early example of such IANA registry is publicly accessible at:

~~~
https://forsalereg.sidnlabs.nl/
~~~

The registry entries consist of content tags as defined in
(#tagdefs).

The initial set of entries in this registry is as follows:

Tag Name | Reference | Status | Description
---------|-----------|--------|-------------
fcod | RFCXXXX | active | For Sale Proprietary Code
ftxt | RFCXXXX | active | For Sale Free Format Text
furi | RFCXXXX | active | For Sale URI
fval | RFCXXXX | active | For Sale Asking Price
Table: Initial set of entries in the "Content Tags" registry

<NOTE TO RFC EDITOR: Adjust the text in the table above before publication with a citation for the (this) document making the addition as per RFC8552.>

Future updates will be managed by the Change Controller.

Entries are assigned only for values that have been documented in 
a manner consistent with the "RFC Required" registration 
policy defined in [@!RFC8126].

Newly defined content tags MUST NOT alter the semantics of existing content tags.

The addition of a new content tag to the registered list does not require the 
definition of a new version tag. However, any modification to existing content tags does.

The "status" column can have one of the following values:

* active - the tag is in use in current implementations.
* historic - the tag is deprecated and not expected to be used in current implementations.

This registry group is maintained by IANA as per [@?RFC8726].

# Privacy Considerations {#privacy}

The use of the "\_for-sale" leaf node name publicly indicates the intent to sell a domain name.
Domain holders should be aware that this information is accessible to anyone querying the
DNS and may have privacy implications.

There is a risk of data scraping, such as email addresses and phone numbers.

Publishing contact information may expose domain holders to spam, or unwanted contact.

# Security Considerations {#security}

One use of the TXT record type defined in this document is to parse the content 
it contains and to automatically publish certain information from it on a 
website or elsewhere. However, there is a risk if the domain name holder 
publishes a malicious URI or one that points to improper content. 
This may result in reputational damage for the party parsing the record.

An even more serious scenario arises when the content of the TXT record 
is insufficiently validated and sanitized, potentially enabling attacks such as XSS or SQL injection.

Therefore, it is **RECOMMENDED** that any parsing and publishing is conducted with the utmost care.
Possible approaches include maintaining a list of validated URIs or applying other validation methods after parsing and before publishing.

There is also a risk that this method will be abused as a marketing tool, or to lure individuals into visiting certain sites or making contact by other
means, without there being any intention to actually sell the domain name. Therefore, this method is best suited for use by professionals.

# Ethical Considerations {#ethicalconsids}
Although not specifically designed for this purpose, the mechanisms 
described in this document may also facilitate domain name 
transactions by professional speculators, often referred to 
as domainers, and those commonly referred to as domain drop catchers. 
Some may view this as controversial.

However, by enabling domain holders to more explicitly
signal their intent to sell, the proposed approach
aims to introduce greater clarity and predictability
into the domain lifecycle. This potentially reduces the
advantage currently held by these professionals, and 
fosters a more equitable environment for all.

Furthermore, this mechanism avoids creating unnecessary 
dependencies on registries for market transactions, 
which could otherwise introduce complexities and 
potential for unintended commercial entanglements.

# Implementation Status
<!-- https://datatracker.ietf.org/doc/html/rfc7942 -->
The concept described in this document has been in use at the .nl ccTLD registry since 2022, 
when it initially started as a pilot. Since then, several hundred thousand domain names have 
been marked with the "\_for-sale" indicator. See for example:

~~~
https://www.sidn.nl/en/whois?q=example.nl
~~~

<!-- or https://api.sidn.nl/rest/whois?domain=example.nl -->

The Dutch domain name registry SIDN offers registrars the option to register a sales 
landing page via its registrar dashboard following the "fcod=" method.
When this option is used, a unique code is generated, which can be included in the "\_for-sale" record. 
If such a domain name is entered on the domain finder page of SIDN, a "for sale" 
button is displayed accordingly.



A simple demonstration of a validator is present at:

~~~
https://forsalereg.sidnlabs.nl/demo
~~~

<NOTE TO RFC EDITOR: Please remove this section before publication as per RFC7942.>

# Acknowledgements

The author would like to thank Thijs van den Hout, Caspar Schutijser, Melvin
Elderman, Ben van Hartingsveldt, Jesse Davids, Juan Stelling,
John R. Levine, and the ISE Editor for their valuable feedback.

{backmatter}

<reference anchor='PSL' target='https://publicsuffix.org/'>
 <front>
  <title>Public Suffix List</title>
  <author>
    <organization>Mozilla Foundation</organization>
  </author>
 </front>
</reference>

<reference anchor='IANA' target='https://www.iana.org/assignments/dns-parameters/dns-parameters.xml#underscored-globally-scoped-dns-node-names'>
 <front>
  <title>Underscored and Globally Scoped DNS Node Names</title>
  <author>
    <organization>IANA</organization>
  </author>
 </front>
</reference>

<reference anchor='FORSALEREG' target='https://forsalereg.sidnlabs.nl/forsale-parameters'>
 <front>
  <title>The "_for-sale" Underscored and Globally Scoped DNS Node Name</title>
  <author>
    <organization>SIDN Labs</organization>
  </author>
 </front>
</reference>

<reference anchor='ISO4217' target='https://en.wikipedia.org/wiki/ISO_4217'>
 <front>
  <title>ISO 4217</title>
  <author>
    <organization>SIX Group</organization>
  </author>
 </front>
</reference>
