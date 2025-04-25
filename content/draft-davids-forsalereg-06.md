%%%
# This is a comment - but only in this block
title = "Registration of the \"_for-sale\" Underscored and Globally Scoped DNS Node Name"
abbrev = "forsalereg"
ipr = "trust200902" 
# area = "Internet"
# workgroup = "Internet Engineering Task Force (IETF)"
workgroup = ""
submissiontype = "IETF"
keyword = [""]
# https://www.rfc-editor.org/rfc/rfc7991#section-2.45.14
tocdepth = 3
# date = 2022-12-22T00:00:00Z

# See FAQ: "How Do I Create an Independent IETF Document?"
# https://mmark.miek.nl/post/faq/
[seriesInfo]
name = "Internet-Draft"
value = "draft-davids-forsalereg-06"
stream = "IETF"
status = "bcp"

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

.# Abstract

This document defines an operational convention for using the reserved DNS leaf node name
"\_for-sale" to indicate that the parent domain name is available for purchase. 
This approach offers the advantage of easy deployment without affecting ongoing operations. As such, the method can be applied to a domain name that is still in full use.

{mainmatter}

# Introduction

Well-established services [@RFC3912; @RFC9083] exist to determine whether a domain name is registered. However, the fact that a domain name exists does not necessarily mean it
is unavailable; it may still be for sale.

Some registrars and other entities offer mediation services between domain name holders and interested parties. For domain names that are not for sale, such services may be
of limited value, whereas they may be beneficial for domain names that are clearly being offered for sale.

This specification defines a simple and universal method to ascertain whether a domain name, although registered, is available for purchase. It enables a domain name holder to add a reserved underscored
leaf node name [@!RFC8552] in the zone, indicating that the domain name is for sale.

The TXT RR type [@!RFC1035] created for this purpose **MUST** follow the formal definition of
(#recformat). Its content **MAY** contain a pointer, such as a Uniform Resource Identifier (URI) [@RFC8820], or another string, 
allowing interested parties to obtain information or contact the domain name holder for further negotiations.

With due caution, such information can also be incorporated into automated availability services. When checking a domain name for availability, the service may indicate whether it is for sale and provide a pointer to the seller's information.

## Terminology

The key words "**MUST**", "**MUST NOT**", "**REQUIRED**", "**SHALL**", "**SHALL NOT**",
"**SHOULD**", "**SHOULD NOT**", "**RECOMMENDED**", "**NOT RECOMMENDED**", "**MAY**", and
"**OPTIONAL**" in this document are to be interpreted as described in BCP 14 [@!RFC2119] [@!RFC8174]
when, and only when, they appear in all capitals, as shown here.

# Rationale

There are undoubtedly more ways to address this problem space. The reasons for the approach defined in this document are primarily accessibility and simplicity. The indicator can be easily turned on and off at will and moreover, it is immediately deployable and does not require significant changes in existing services. This allows for a smooth introduction of the concept.

# Conventions

## General Record Format {#recformat}

<!-- TODO: see https://www.rfc-editor.org/rfc/rfc8461.html#section-3.1 for inspiration -->

The "\_for-sale" TXT record **MUST** start with a version tag, possibly followed by a string.

The formal definition of the record format, using ABNF [@!RFC5234; @!RFC7405], is as follows:

~~~
forsale-record  = forsale-version forsale-content
forsale-version = %s"v=FORSALE1;"
                  ; case sensitive, no spaces
forsale-content = 0*244OCTET
                  ; referred to as content or data
~~~

<!-- TODO: double check on https://author-tools.ietf.org/abnf -->

## Content limitations

The TXT [@RFC8553, (see) section 2.1] record **MUST** contain any valid content, ranging from an empty string to meaningful text or URIs.
Any text that suggests that the domain is not for sale is invalid content. If a domain name is not for sale, 
a "\_for-sale" indicator is pointless and any existence of a valid "\_for-sale" TXT record **MAY**
therefore be regarded as an indication that the domain name is for sale.

This specification does not dictate the exact use of any content in the "\_for-sale" TXT record, or the lack of any such content.
Parties - such as registries and registrars - **MAY** use it in their tools, perhaps even by defining specific requirements that the content must meet.
Content can also be represented in a human-readable format for individuals to
interpret. See the (#examples, use title) section for clarification.

Since the content of TXT record has no defined meaning, it is up to the processor of the content to decide how to handle it. 

See (#guidelines) for additional guidelines.

## RRset limitations

This specification does not define any restrictions on the number of TXT records in the RRset, but limiting it to one is **RECOMMENDED**. 
It is also **RECOMMENDED** that the length of the RDATA [@RFC9499] per TXT record does not exceed 255 octets. 
If this is not the case, the processor **SHOULD**  determine which content to use. 

For example, a registry might select content that includes a recognizable code, which can be used to direct visitors to a sales page 
as part of its services, whereas an individual might simply extract a phone number (if present) and use it to contact a potential seller.

## RR Type limitations

Adding any resource record (RR) types under the "\_for-sale" leaf other than TXT is **NOT RECOMMENDED**. 
Such records **MUST** be ignored for the purposes of this document.

## TTL limitation

A TTL longer than 86400 is **NOT RECOMMENDED**. Long TTLs increase the risk of outdated information persisting, potentially misleading buyers into believing the domain is still available for purchase.

## Wildcard limitation

The "\_for-sale" leaf node name **SHOULD NOT** be a wildcard.

## CNAME limitation

The "\_for-sale" leaf node name **MAY** be an alias, but if
that is the case, the CNAME record it is associated with it **SHOULD** also be
named "\_for-sale", for example:

~~~
_for-sale.example.com. IN CNAME _for-sale.example.org.
~~~

## Placement of leaf node name

The "\_for-sale" leaf node name **MAY** be placed on the top level domain, or any domain directly below, with the exception of the .arpa infrastructure top-level domain.

It **MAY** also be placed at a lower level, but only when that level is mentioned in the Public Suffix List [@PSL]. 

Any other placement of the record **MUST NOT** be regarded as a signal that the domain above it is for sale. 

(#placements) provides further clarification.

Name | Situation | Verdict
-----|-----------|--------
\_for-sale.example | root zone | For sale
\_for-sale.aaa.example | Second level | For sale
\_for-sale.acme.bbb.example | bbb.example in PSL | For sale
\_for-sale.www.ccc.example | Other | Invalid
\_for-sale.51.198.in-addr.arpa | infrastructure TLD | Invalid
Table: Allowed placements of TXT record {#placements}

# Examples {#examples}

## Example 1: A URI

The holder of 'example.com' wishes to signal that the domain is for sale and adds this record to the 'example.com' zone:

~~~
_for-sale.example.com. IN TXT "v=FORSALE1;https://buy.example.com/"
~~~

An interested party notices this signal and can visit the URI mentioned for further information. The TXT record
may also be processed by automated tools, but see the (#security, use title) section for possible risks. 

As an alternative, a mailto: URI could also be used:

~~~
_for-sale.example.com. IN TXT "v=FORSALE1;mailto:owner@example.com"
~~~

Or a telephone URI:

~~~
_for-sale.example.com. IN TXT "v=FORSALE1;tel:+1-201-555-0123"
~~~

There can be a use case for these URIs, especially since WHOIS (or RDAP) often has privacy restrictions.
But see the (#privacy, use title) section for possible downsides.

## Example 2: Various other approaches

Free format text, with some additional unstructured information, aimed at
being human-readable:

~~~
_for-sale.example.com. IN TXT "v=FORSALE1;$500, info[at]example.com"
~~~

A proprietary format, defined by a registry or registrar to automatically redirect visitors to a web page, 
but without a clearly defined meaning to third parties:

~~~
_for-sale.example.com. IN TXT "v=FORSALE1;fscode=aHR0cHM...wbGUuY29t"
~~~

The content in the following example could be malicious, but it is not in violation of this specification (see (#security)):

~~~
_for-sale.example.com. IN TXT "v=FORSALE1;<script>alert('')</script>"
~~~

# Operational Guidelines {#guidelines}
DNS wildcards interact poorly with underscored names. Therefore, the use of wildcards 
is **NOT RECOMMENDED** when deploying this mechanism. However, wildcards may still be encountered 
in practice, especially with operators who are not implementing this mechanism. 
This is why the version tag is a **REQUIRED** element: it helps distinguish
valid "\_for-sale" records from unrelated TXT records. Nonetheless, any assumptions about the 
content of "\_for-sale" TXT records **SHOULD** be made with caution.

It is also **RECOMMENDED** that the content string be limited to visible ASCII characters, 
excluding the double quote (") and backslash (\\). In ABNF syntax, this would be:

~~~
forsale-content     = 0*244recommended-char
recommended-char    = %x20-21 / %x23-5B / %x5D-7E
~~~

# IANA Considerations

IANA has established the "Underscored and Globally Scoped DNS Node Names" registry [@!RFC8552; @IANA]. The underscored
leaf node name defined in this specification should be added as follows:

~~~ ascii-art
             +-----------+--------------+-------------+
             | RR Type   | _NODE NAME   | Reference   |
             +-----------+--------------+-------------+
             | TXT       | _for-sale    | TBD         |
             +-----------+--------------+-------------+
~~~
Figure: Entry for the "Underscored and Globally Scoped DNS Node Names" registry

This specification does not require the creation of an IANA registry for record fields.

# Privacy Considerations {#privacy}

The use of the "\_for-sale" leaf node name publicly indicates the intent to sell a domain name.
Domain owners should be aware that this information is accessible to anyone querying the
DNS and may have privacy implications.

There is a risk of data scraping, such as email addresses and phone numbers.

# Security Considerations {#security}

One use of the TXT record type defined in this document is to parse the content it contains and to automatically publish certain information from it on a website or elsewhere. However, there is a risk if the domain name holder  publishes a malicious URI or one that points to improper content. This may result in reputational damage for the party parsing the record.

An even more serious scenario occurs when the content of the TXT record is not validated and sanitized sufficiently, opening doors to - for example - XSS attacks among other things. 

Therefore, it is **RECOMMENDED** that any parsing and publishing is conducted with the utmost care.

There is also a risk that this method will be abused as a marketing tool, or to lure individuals into visiting certain sites or making contact by other
means, without there being any intention to actually sell the particular domain name. Therefore, this method is best suited for use by professionals.


# Implementation Status

The concept described in this document is in use with the .nl ccTLD
registry. See for example:

~~~
https://www.sidn.nl/en/whois?q=example.nl
~~~

<NOTE TO RFC EDITOR: Please remove this section before publication.>

# Acknowledgements

The author would like to thank Thijs van den Hout, Caspar Schutijser, Melvin
Elderman, Paul Bakker and Ben van Hartingsveldt for their valuable feedback.

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
