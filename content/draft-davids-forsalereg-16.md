%%%
# This is a comment - but only in this block
title = "The \"_for-sale\" Underscored and Globally Scoped DNS Node Name"
abbrev = "_for-sale DNS"
ipr = "trust200902" 
area = "Operations and Management"
#workgroup = 
submissiontype = "IETF"
keyword = [""]
# https://www.rfc-editor.org/rfc/rfc7991#section-2.45.14
tocdepth = 3
# date = 2022-12-22T00:00:00Z

# See FAQ: "How Do I Create an Independent IETF Document?"
# https://mmark.miek.nl/post/faq/
[seriesInfo]
name = "Internet-Draft"
value = "draft-davids-forsalereg-16"
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

This document defines an operational convention that uses the reserved 
underscored DNS leaf node name "\_for-sale" to indicate that the 
parent domain name is available for purchase. The convention can be 
deployed without disrupting existing operations, and it may be 
applied even when the domain name is still actively in use.

{removeInRFC="true"}
.# About This Document

This note is to be removed before publishing as an RFC.

This document contains a "Note to the RFC Editor" requesting removal 
of (#implementation) prior to publication. Please also review the Status of 
This Memo section and other relevant parts before publication, 
particularly (#ianaconsid).


{mainmatter}

# Introduction {#introsect}

Well-established services [@RFC3912; @RFC9083] exist to determine whether a
DNS domain name is registered. However, the fact that a domain name exists does not necessarily mean it
is unavailable; it may still be for sale.

Some registrars and other parties offer brokerage services between domain name holders and interested buyers.
Such services are of limited value when the domain name is not available for purchase, but they may be 
beneficial for domain names that are clearly being offered for sale.

This specification defines a simple method to ascertain whether a domain name, 
although registered, is available for purchase. It enables a domain name holder to 
add a reserved underscored leaf node name [@!RFC8552] in the zone, indicating that the 
domain name is for sale. The indicator can be turned on and off at will and, moreover, 
it is immediately deployable and does not require significant changes in existing 
services, allowing for a smooth introduction of the concept.

The TXT RR type [@!RFC1035] created for this purpose must follow the formal definition of
(#conventions). Its content may contain a pointer, such as a Uniform Resource Identifier (URI) 
[@!RFC3986], an Internationalized Resource Identifier (IRI) [@!RFC3987] or another string, 
allowing interested parties to obtain information or contact the domain name holder for further negotiations.
Details about whether and how such negotiations occur are out of scope.

With due caution, such information can also be incorporated into automated availability services. 
When checking a domain name for availability, the service may indicate whether it is for 
sale and provide a pointer to the seller's information.

The operational convention described in this document does not require any protocol change.

Furthermore, (#ethicalconsids) discusses some ethical considerations. In particular, 
the approach in this document aims to promote a more equitable domain aftermarket and 
minimizing potential for unintended commercial entanglements by registries.

Examples are provided in (#examples).

## Terminology

The key words "**MUST**", "**MUST NOT**", "**REQUIRED**", "**SHALL**", "**SHALL NOT**",
"**SHOULD**", "**SHOULD NOT**", "**RECOMMENDED**", "**NOT RECOMMENDED**", "**MAY**", and
"**OPTIONAL**" in this document are to be interpreted as described in BCP 14 [@!RFC2119] [@!RFC8174]
when, and only when, they appear in all capitals, as shown here.

Although the document defines an operation convention not any protocol extension, the use of normative language is used
to promote consistent and umabiguous behaviors among entities that adopt the convention.

The term "Processor" refers to an entity (person, system, or service) 
that reads, interprets, and takes appropriate actions based on "\_for-sale" DNS labels, 
whether manually or automatically.

The term "for sale" is used in a broad sense and may also refer to cases 
where the domain name is available for lease, or where the contractual right to 
use the domain name is offered to another party.

# Conventions {#conventions}

## General Record Format {#abnf}

Each "\_for-sale" TXT record **MUST** begin with a version tag, optionally followed by a string containing content that follows a simple "tag=value" syntax.

The formal definition of the record format, using ABNF [@!RFC5234; @!RFC7405], is as follows:

~~~
forsale-record  = forsale-version [forsale-content]
                  ; referred to as 'content' or RDATA
                  ; in a single character-string

forsale-version = %s"v=FORSALE1;"
                  ; %x76.3D.46.4F.52.53.41.4C.45.31.3B
                  ; version tag, case-sensitive, no spaces

forsale-content = fcod-pair / ftxt-pair / furi-pair / fval-pair
                  ; referred to as 'tag-value pairs'
                  ; only one tag-value pair per record

fcod-pair       = fcod-tag fcod-value
ftxt-pair       = ftxt-tag ftxt-value
furi-pair       = furi-tag furi-value
fval-pair       = fval-tag fval-value
                  ; the tags are referred to as 'content tags'
                  ; the values are referred to as 'content values'

fcod-tag        = %s"fcod="
ftxt-tag        = %s"ftxt="
furi-tag        = %s"furi="
fval-tag        = %s"fval="
                  ; all content tags case-sensitive lowercase

fcod-value      = 1*239OCTET

ftxt-value      = 1*239OCTET

furi-value      = URI / IRI
                  ; http, https, mailto and tel URI schemes
                  ; exactly one URI or IRI

URI             = <as defined in RFC3986, Appendix A>
IRI		= <as defined in RFC3987, Section 2.2>

fval-value      = fval-currency fval-amount
                  ; total length: 2 to 239 characters 
fval-currency   = 1*%x41-5A
                  ; one or more uppercase letters (A-Z)
                  ; indicating (crypto)currency
                  ; e.g., USD, EUR, BTC, ETH
                  ; standard three-letter fiat currencies recommended
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
is present, processors **SHOULD** assume that the domain is for sale unless a local 
policy indicates otherwise. For example:

```
_for-sale.example.com. IN TXT "v=FORSALE1;"
_for-sale.example.com. IN TXT "v=FORSALE1;fcod="
_for-sale.example.com. IN TXT "v=FORSALE1;foo=bar"
```

In such cases, processors determine how to proceed. 
An approach might be to signal that the domain is for sale and 
to rely on conventional mechanisms (e.g., WHOIS or Registration Data Access 
Protocol (RDAP)) to retrieve and present contact information.

TXT records in the same RRset, but without a version tag, **MUST NOT** be interpreted or processed as a valid "\_for-sale" indicator. 
However, they may still offer some additional information for humans when considered alongside a valid
record. For example:

```
_for-sale.example.com. IN TXT "I am for sale"
_for-sale.example.com. IN TXT "v=FORSALE1;fcod=XX-NGYyYjEyZWY"
```

If no TXT records at a leaf node contain a valid version tag, processors 
**MUST** consider the node name invalid and **MUST** discard it.

See (#contentlimits) for additional content limitations.

## Content Tag Type Definitions {#tagdefs}

The following content tags are defined as valid content tags.

Tags beginning with the character "x" may be used for development and testing.

Content tags are optional. Providing a minimum set to allow interested parties 
to engage is **RECOMMENDED**.

<!-- author tip: there are two spaces behind the content tag, to enforce a new line -->
### fcod {#fcoddef}  
This content tag is intended to contain a code that is meaningful only to processors 
that understand its semantics. The content value **MUST** consist of at least one octet. 

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

The following example shows a string encoded using Base64 [@?RFC4648] 
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
content tag and interpret it as intended. This example uses Base64 encoding 
to avoid escaping and ensure printable characters, though this is 
**OPTIONAL** and not required.

### ftxt  
This content tag is intended to contain human-readable text that conveys additional information to interested parties. For example:

~~~
_for-sale IN TXT "v=FORSALE1;ftxt=Call for info."
~~~

While a single octet is the minimum, it is **RECOMMENDED** to provide more context.

While a URI in this field is not syntactically prohibited, its 
interpretation as a URI is not guaranteed. Use of URIs in this 
field **SHOULD** be avoided in favor of the "furi=" content tag.

See (#fvalpar) for a way to explicitly indicate an asking price for easier machine parsing.

See (#handlerdata) for considerations regarding the representation of non-ASCII data in the content value.

<!-- TODO https://www.rfc-editor.org/rfc/rfc7553.html noemen, of zelfs opnemen? -->
### furi  
This content tag is intended to contain a human-readable and machine-parseable URI that can be used by interested parties to retrieve further information.

While the syntax allows any URI scheme, only the following schemes are **RECOMMENDED** 
for use: `http` and `https` [@RFC9110], `mailto` [@RFC6068; @RFC6530, (see) section 11.1], and `tel` [@RFC3966].

The content value **MUST** contain exactly one URI. For example:

~~~
_for-sale IN TXT "v=FORSALE1;furi=https://example.com/foo%20bar"
~~~

URIs **MUST** conform to the syntax and encoding requirements specified in 
[@!RFC3986, section 2.1], including the percent-encoding of characters 
not allowed unencoded (e.g., spaces must be encoded as `%20` in a URI).

(#handlerdata) provides additional guidelines on character encoding.

See the (#security, use title) section for possible risks.

Note: References to a URI in this document also encompass IRIs [@!RFC3987].

### fval {#fvalpar}
This content tag is intended to contain human-readable and machine-parseable 
text that explicitly indicates an asking price in a certain currency, as opposed to 
the price being loosely incorporated in an "ftxt=" content tag. For example:

~~~
_for-sale IN TXT "v=FORSALE1;fval=EUR999"
~~~

See (#currency) for additional operational guidelines.

### Future Tags
Future tags may be defined to accommodate operational needs. Future content 
tags **MUST NOT** alter the semantics of existing content tags.

A tag name length of 4 characters is **RECOMMENDED** for consistency with the initial tag
set and to maintain compact record formats.

Tags beginning with the character "x" are reserved for development and testing.

## Content Limitations {#contentlimits}

The "\_for-sale" TXT record [@RFC8553, (see) section 2.1] **MUST** contain content deemed valid under this specification.

Any text suggesting that a domain is not for sale is invalid content. If a domain name is not or no longer for sale, 
a "\_for-sale" indicator **SHOULD NOT** exist. The presence of a valid "_for-sale" TXT record
**SHOULD** therefore be regarded as an indication that the domain name is for sale.

The existence of a "\_for-sale" leaf node does not obligate the holder to sell the domain name; 
it may have been published in error, or withdrawn later for other reasons.

This specification does not dictate the exact use of any content values in the "\_for-sale" TXT record.
Parties may use it in their tools, perhaps even by defining specific requirements that the content
value must meet. Content values can also be represented in a human-readable format for individuals to
interpret. See the (#examples, use title) section for clarification.

See (#operationalcons) for additional guidelines.

## RRset Limitations {#rrsetlimits}

This document does not define restrictions on the number of TXT records in the
RRset of "\_for-sale" TXT records.

When multiple "\_for-sale" TXT records are present in an RRset, the 
processor **MAY** select one or more of them.

For example, a domain name registry might extract content from an RRset that includes 
a recognizable "fcod=" content tag and use it to direct visitors to a sales page as 
part of its services. An individual, on the other hand, might extract a 
phone number (if present) from a "furi=" tag in the same RRset and use it to contact a potential seller.

An example of such a combined record is provided in (#combiexample).

The RDATA [@RFC9499] of each "\_for-sale" TXT record **MUST** consist of a single character-string
[@RFC1035] with a maximum length of 255 octets, to avoid the need to concatenate multiple
character-strings during processing.

The following example illustrates an invalid "\_for-sale" TXT record due to the presence of multiple
character-strings:

~~~
_for-sale IN TXT "v=FORSALE1;" "ftxt=foo" "bar" "invalid"
~~~

## Wildcard Limitation

Wildcards are only interpreted as leaf names, so "\_for-sale.*.example." is not a valid wildcard
[@RFC4592] and is non-conformant. Hence, it is not possible to put all domains under a TLD for 
sale with just one "\_for-sale" TXT record.

The example below, however, shows a common use case where a "\_for-sale" leaf node exists alongside a
wildcard:

~~~
*         IN A    198.51.100.80
          IN AAAA 2001:db8::80
_for-sale IN TXT  "v=FORSALE1;ftxt=Only $99 at ACME"
~~~

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
there may not be a public domain name registry [@?RFC9499] capable of properly recording the rights associated with that label. 
Nevertheless, this does not constitute a violation of this document. 
One possible approach is for the involved parties to establish a mutual agreement to formalize these rights.

Note 2:
If a "\_for-sale" leaf node were to appear under the .arpa infrastructure top-level 
domain, it might be interpreted as an offer to sell IP address space. 
However, such use is explicitly out of scope for this document, and processors
**MUST** ignore any such records.

# Operational Considerations {#operationalcons}
## DNS Wildcards

DNS wildcards interact poorly with underscored names [@RFC8552, (see) section 1.4],
but they may still be encountered in practice, especially with operators who 
are not implementing this mechanism. This is why the version 
tag is a mandatory element: it allows processors to distinguish 
valid "\_for-sale" records from unrelated TXT records.

Nonetheless, any assumptions about the content of "\_for-sale" TXT 
records should be made with caution, particularly in edge 
cases where wildcard expansion - possibly combined with DNS aliases 
(e.g., CNAMEs) or redirections (e.g., DNAMEs [@?RFC6672]) - might 
result in misleading listings or unintended references to third-party domains.

## Handling of RDATA {#handlerdata}

Since this method relies on DNS TXT records, standard content rules apply as 
defined in [@RFC1035, (see) section 5.1]. This includes the possibility of 
representing non-ASCII data in the content value, by using escape sequences 
(e.g., \DDD or \X notation).

All text in content values exchanged between systems that are not part of a 
closed ecosystem **SHALL** be encoded in and interpreted as UTF-8 [@!RFC3629] and conform 
to the Network Unicode format [@?RFC5198]. The allowed subset of Unicode code points
**SHOULD** conform to [@!RFC9839, (see) section 4.3], with the exception of `%x09`, `%x0A`
and `%x0D` which **SHOULD NOT** be used.

See (#robustness) for additional guidelines and the (#security, use title)
section for possible risks.

Internationalized Domain Names (IDN) **MAY** be expressed as A-labels as well as
U-labels [@!RFC5890].

Processors **MUST** be capable of handling such encodings to ensure that 
non-ASCII content values are correctly interpreted and represented.

Note: When non-ASCII data is used, the ABNF octet limit applies to the encoded
byte sequence, not the number of visible characters. Multi-byte
characters will reduce the available character count.

For example, this TXT record contains two Japanese Kanji characters 
encoded in UTF-8 using octal escapes, resulting in a content value 
of six octets.

~~~
_for-sale IN TXT "v=FORSALE1;ftxt=\229\163\178\229\135\186"
~~~

## Currency {#currency}

The ABNF in (#abnf) allows currency codes consisting of one or 
more uppercase letters, providing flexibility to 
accommodate both standard fiat currencies and other widely 
recognized abbreviations, such as cryptocurrencies.

The use of standard fiat currencies is **RECOMMENDED**. When used, 
they **MUST** be represented by three-letter uppercase currency 
codes as specified in [@!ISO4217] (e.g., USD, EUR, GBP and JPY).

The amount component consists of an integer part, optionally 
followed by a fractional part separated by a decimal point (`%x2E`, ".").

## TTLs

Long TTLs [@!RFC1035, (see) section 3.2.1] increase the risk of outdated data misleading buyers into thinking the domain is still
available. 

## Ambiguous Constructs

Ambiguous constructs in content values **SHOULD** be avoided, as illustrated by the following
example:

~~~
_for-sale IN TXT "v=FORSALE1;fcod=TRIP-confusing;ftxt=dont_do_this"
~~~

The above example is a valid "fcod=" content tag that includes the 
string ";ftxt=" in the content value, which may be confusing, 
as it does not actually represent an "ftxt=" content tag.

## Robustness {#robustness}

Because the format of the content part is not strictly defined in this 
document, processors **MAY** apply the robustness principle of being 
liberal in what they accept. This also applies to space 
characters (`%x20`) immediately following the version tag.

Alternatively, parties may agree on a more strictly defined proprietary format 
for the content value to reduce ambiguity. However, it is out of scope to discuss
which mechanisms are put in place for such agreements. Processors **MAY** also convert 
control characters (e.g., `%x09`, `%x0A`, `%x0D`) in "ftx=" content to 
spaces (`%x20`) for correct representation.

## Scope of Application

The "\_for-sale" mechanism  relies upon the domain name being resolvable in the DNS.
This is not guaranteed, for example, during a redemption period, in
pendingDelete status [@?STD69], or when the domain is DNSSEC-signed but fails 
validation (i.e., has a bogus state).

# Security Considerations {#security}

One use of the TXT record type defined in this document is to parse the content 
it contains and to automatically publish certain information from it on a 
website or elsewhere. However, there is a risk if the domain name holder 
publishes a malicious URI or one that points to improper content. 
This may result in reputational damage to the party parsing the record.

An even more serious scenario arises when the content of the TXT record is not 
properly validated and sanitized, potentially enabling attacks such as XSS or SQL 
injection, as well as spoofing techniques based on Unicode manipulation, 
including bidirectional text attacks and homograph attacks.

Therefore, it is **RECOMMENDED** that any parsing and publishing is conducted with the utmost care.
Possible approaches include maintaining a list of validated URIs or applying other validation methods after parsing and before publishing.

There is also a risk that this method will be abused as a marketing tool, or to lure individuals into visiting certain sites or making contact by other
means, without there being any intention to actually sell the domain name. Therefore, this method is best suited for use by professionals.

# Privacy Considerations {#privacy}

The use of the "\_for-sale" leaf node name publicly indicates the intent to sell a domain name.
Domain holders should be aware that this information is accessible to anyone querying the
DNS and may have privacy implications.

There is a risk of data scraping, such as email addresses and phone numbers.

Publishing contact information may expose domain name holders to spam, or unwanted contact.

# Ethical Considerations {#ethicalconsids}
Although not specifically designed for this purpose, the mechanism 
described in this document may also facilitate domain name 
transactions by professional speculators, often referred to 
as domainers, and those commonly referred to as domain drop catchers. 
Some may view this as controversial.

However, by enabling domain name holders to more explicitly
signal their intent to sell, the "\_for-sale" approach
aims to introduce greater clarity and predictability
into the domain lifecycle. This potentially reduces the
advantage currently held by these professionals, and 
fosters a more equitable environment for all.

Furthermore, this mechanism avoids creating unnecessary 
dependencies on registries for market transactions, 
which could otherwise introduce complexities and 
potential for unintended commercial entanglements.

# Implementation Status {#implementation}
<!-- https://datatracker.ietf.org/doc/html/rfc7942 -->
The concept described in this document has been in use at the .nl ccTLD registry since 2022, 
when it initially started as a pilot. Since then, hundreds of thousands of domain names have 
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

# IANA Considerations {#ianaconsid} <!-- See RFC8126 -->

IANA is requested to add the following entry to the "Underscored and Globally Scoped DNS Node Names" 
registry [@RFC8552] :

RR Type | _NODE NAME | Reference
--------|------------|-----------
TXT | \_for-sale | <this memo>
Table: Entry for the "Underscored and Globally Scoped DNS Node Names" registry

{backmatter}

# Additional Examples {#examples}

## Example 1: Code Format

A proprietary format, defined and used by agreement between parties - for example, 
a domain name registry and its registrars - without a clearly specified meaning for third parties.
For example, it may be used to automatically redirect visitors to a web page, as described in
(#fcoddef):

~~~
_for-sale IN TXT "v=FORSALE1;fcod=XX-aHR0cHM...wbGUuY29t"
~~~

Note: the content value in the above example is truncated for readability.

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
numeric amount. See (#currency) for additional guidelines.

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

{numbered="false"}
# Acknowledgements

The author would like to thank Thijs van den Hout, Caspar Schutijser, Melvin
Elderman, Ben van Hartingsveldt, Jesse Davids, Juan Stelling, John R.&#xa0;Levine, 
Dave Lawrence, Andrew Sullivan, Paul Hoffman, Eliot Lear (ISE) and 
Mohamed 'Med' Boucadair for their valuable feedback.

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

<reference anchor='ISO4217' target='https://www.iso.org/iso-4217-currency-codes.html'>
 <front>
  <title>ISO 4217 Currency Codes</title>
  <author>
    <organization>SIX Group</organization>
  </author>
 </front>
</reference>

