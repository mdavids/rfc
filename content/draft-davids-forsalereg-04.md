%%%
# This is a comment - but only in this block
title = "Registration of Underscored and Globally Scoped 'for sale' DNS Node Name"
abbrev = "forsalereg"
ipr = "trust200902"
# area = "Internet"
# workgroup = "Internet Engineering Task Force (IETF)"
submissiontype = "IETF"
keyword = [""]
tocdepth = 5
# date = 2022-12-22T00:00:00Z

[seriesInfo]
name = "Internet-Draft"
value = "draft-davids-forsalereg-04"
stream = "IETF"
status = "bcp"  # or "informational" or "experimental" ?

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

This document defines a simple operational convention of using a reserved underscored node name ("\_for-sale") to indicate that the parent domain name above is for sale. This approach offers the advantage of easy deployment without affecting ongoing operations. As such, the method can be applied to a domain name that is still in full use.

{mainmatter}

# Introduction

Well-established services [@RFC3912; @RFC9083] exist to determine whether a domain name is registered. However, the fact that a domain name exists does not necessarily mean it
is unavailable; it may still be for sale.

Some registrars and other entities offer mediation services between domain name holders and interested parties; however, for domain names not for sale, such services may be unnecessary.

This specification defines a simple and universal method to ascertain whether a domain name, although registered, is available for purchase. It enables a domain name holder to add a reserved underscored node name [@!RFC8552] in the zone, indicating that the domain name is for sale.

The TXT record RR type [@!RFC1035] that is created for that purpose **MAY** contain a pointer, such as a
Uniform Resource Identifier (URI) [@RFC8820], allowing interested parties to obtain information or contact the domain name holder for further negotiations.

With due caution, such information can also be incorporated into automated availability services. When a domain name is checked for availability, the service can indicate whether it is for sale and provide a pointer to the seller's information.

## Terminology

The key words "**MUST**", "**MUST NOT**", "**REQUIRED**", "**SHALL**", "**SHALL NOT**",
"**SHOULD**", "**SHOULD NOT**", "**RECOMMENDED**", "**NOT RECOMMENDED**", "**MAY**", and
"**OPTIONAL**" in this document are to be interpreted as described in BCP 14 [@!RFC2119] [@!RFC8174]
when, and only when, they appear in all capitals, as shown here.

# Rationale

There are undoubtedly more ways to address this problem space. The reasons for the approach defined in this document are primarily accessibility and simplicity. The indicator can be easily turned on and off at will and moreover, it is immediately deployable and does not require significant changes in existing services. This allows for a smooth introduction of the concept.

# Conventions

## Content limitations

The TXT [@RFC8553, (see) section 2.1] record **MUST** contain any valid content, ranging from an empty string to meaningful text or URIs. However, it **SHALL NOT** contain any text that suggests that the domain is not for sale. If a domain name is not for sale, the "\_for-sale" indicator MUST NOT be used. Any existence of a "\_for-sale" TXT record **MUST** therefore be regarded as an indication that the domain name is for sale.

This specification does not dictate the exact use of any content in the "\_for-sale" TXT record, or the lack of any such content. Parties - such as Registries and registrars - may use it in their tools, perhaps even by defining additional requirements that the content must meet. Alternatively, an individual can use it in combination with existing tools to make contact with the seller.

The content of the TXT record is "as is" and characters such as ";" between two URIs for example, have no defined meaning. It is up to the processor of the content to decide how to handle it. See (#security)
for additional guidelines.

## RRset limitations

This specification does not define any restrictions on the number of TXT records in the RRset, although it is recommended to limit it to one. It is also recommended that the length of the RDATA [@RFC8499] does not exceed 255 bytes. If the RRset contains multiple records or the total size exceeds 255 bytes, it is up to the processor to determine which data to use.. For example, a
registry might pick a mandatory URI from the RRset to display on a website as part of its service, while an individual might just pick a phone number (if present) and dial it to make contact with a potential seller.

## RR Type limitation

Adding any other RR types under the "\_for-sale" leaf but TXT is **NOT RECOMMENDED** and they **MUST** be ignored for the purpose of this document.

## TTL limitation

A TTL longer than 86400 is **NOT RECOMMENDED**. Long TTLs increase the risk of outdated information persisting, potentially misleading buyers into believing the domain is still available for purchase.

## Wildcard limitation

The "\_for-sale" leaf **SHOULD NOT** be a wildcard.

## CNAME limitation

The "\_for-sale" leaf **MAY** be a CNAME pointing to a TXT RR type.

## Placement of node name

The "\_for-sale" leaf node name **MAY** be placed on the top level domain, or any domain directly below. It **MAY** also be placed at a lower level, but only when that level is mentioned in the Public Suffix List [@PSL]. 

Any other placement of the record **MUST NOT** be regarded as a signal that the domain above it is for sale. 

See (#placements) for further explanation.

Name | Situation | Verdict
-----|-----------|--------
\_for-sale.example | root zone | For sale
\_for-sale.aaa.example | Second level | For sale
\_for-sale.acme.bbb.example | bbb.example in PSL | For sale
\_for-sale.www.ccc.example | Other | Invalid
Table: Allowed placements of TXT record {#placements}

# Examples {#examples}

## Example 1: A URI

The holder of 'example.com' wishes to signal that the domain is for sale and adds this record to the 'example.com' zone:

~~~
_for-sale.example.com. IN TXT "https://example.com/forsale.html"
~~~

An interested party notices this signal and can visit the URI mentioned for further information. The TXT record can also be processed by automated tools, but see the (#security, use title) section for possible risks. 

As an alternative, a mailto: URI could also be used:

~~~
_for-sale.example.com. IN TXT "mailto:owner@example.com"
~~~

Or a telephone URI:

~~~
_for-sale.example.com. IN TXT "tel:+1-201-555-0123"
~~~

There can be a use case for a telephone URI, especially since WHOIS (or RDAP) often has privacy restrictions.

## Example 2: Various other approaches

Free format text:

~~~
_for-sale.example.com. IN TXT "I'm for sale: info [at] example.com"
~~~

Proprietary format, used by a registry or registrar to automatically redirect visitors to a web page, and which has no well-defined meaning to third parties.

~~~
_for-sale.example.com. IN TXT "fscode=aHR0cHM...V4YW1wbGUuY29t"
~~~

The content in the following example could be malicious, but it is not in violation of this specification (see (#security)):

~~~
_for-sale.example.com. IN TXT "<script>alert('H4x0r')</script>"
~~~

# Operational Guidelines {#guidelines}
DNS wildcards interact poorly with underscored names. And even though wildcards
are NOT RECOMMENDED, they can still occur. As such, no assumptions SHOULD be made
about the content of "\_for-sale" TXT records. 

For example, some operators use wildcards to enforce a "v=spf1 -all"
response for every subdomain. But obvisously, there is a reasonable change
that the existence of a "\_for-sale" TXT record with such content is not for
sale. It's possible to circumvent this by adding a "\_for-sale" record of a
different RR type, but processors SHOULD NOT expect this to be the case.

For example:

~~~
_for-sale.example.com. IN NULL \# 1 FF
~~~

Hence, it is RECOMMENDED to work with with content that is recognizable,
either for humans or automated processes. Such as the "fscode="-string in
the (#examples, use title) section, or a descriptive string that humans can can easily
interpret.

# IANA Considerations

IANA has established the "Underscored and Globally Scoped DNS Node Names" registry [@!RFC8552; @IANA]. The underscored node name defined in this specification should be added as follows:

~~~ ascii-art
             +-----------+--------------+-------------+
             | RR Type   | _NODE NAME   | Reference   |
             +-----------+--------------+-------------+
             | TXT       | _for-sale    | TBD         |
             +-----------+--------------+-------------+
~~~
Figure: Entry for the "Underscored and Globally Scoped DNS Node Names" registry


# Privacy Considerations

There is a risk of data scraping, such as email addresses and phone numbers.

# Security Considerations {#security}

One use of the TXT record type defined in this document is to parse the content it contains and to automatically publish certain information from it on a website or elsewhere. However, there is a risk involved in this if the domain name holder  publishes a malicious URI or one that points to improper content. This may result in reputational damage for the party parsing the record.

Even worse is a scenario in which the content of the TXT record is not validated and sanitized sufficiently, opening doors to - for example - XSS attacks among other things. 

Therefore, it is **RECOMMENDED** that any parsing and publishing is conducted with the utmost care.

There is also a risk that this method will be abused as a marketing tool, or to otherwise lure individuals into visiting certain sites or attempting other forms of contact, without there being any intention to actually sell the particular domain name. Therefore, it is recommended that this method is primarily used by professionals.

# Implementation Status

The concept described in this document is in use with the .nl ccTLD registry.

[note to editor: please remove this section before publication]

# Acknowledgements

The author would like to thank Thijs van den Hout, Caspar Schutijser, Melvin Elderman and Paul Bakker for their valuable feedback.

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


