%%%
# This is a comment - but only in this block
title = "Registration of underscored 'for sale' DNS Node Name"
abbrev = "forsalereg"
ipr = "trust200902"
# area = "Internet"
# workgroup = "Internet Engineering Task Force (IETF)"
submissiontype = "IETF"
keyword = [""]
tocdepth = 5
# date = 2022-12-20T00:00:00Z

[seriesInfo]
name = "Internet-Draft"
value = "draft-davids-forsale-00"
stream = "IETF"
status = "bcp"	# or "informational" or "experimental" ?

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
  city = "Arnhem"
  code = "6825 MD"
  pobox = "Meander 501"
  cityarea = "Gld"
%%%

{mainmatter}

# Abstract

This document defines a simple operational convention of using a reserved underscored node name ("\_for-sale" TXT RRset in DNS [@!RFC1035]) to indicate that the parent domain name above it is for sale.

# Introduction

Well established services exist, to find out if a domain name is registered or not. But the fact that a domain name exists does not exclude the possibility that it can be obtained.

Some registrars offer (payed) mediation services between domain name holders and interested parties, but for a domain  name that is not for sale, such services are a waste of money and   time.

This specification defines a simple universal way to find out if a domain name, even thouh it is taken, might be obtained nevertheless. It enables domain name holder to add a reserved underscored node name in the zone, indicating that the domain name is actually for sale.

The TXT record type that is created for that purpose **MAY** contain a pointer, such as a URI, to allow an interested party to find information or to get in touch and engage in further arrangements.

With due caution, this information can also be incorporated in the automated availability services, so that when the domain name is checked for availabilty, the service can also indicate it is for sale, including a referral to the selling party's information.

## Terminology

The key words "**MUST**", "**MUST NOT**", "**REQUIRED**", "**SHALL**", "**SHALL NOT**",
"**SHOULD**", "**SHOULD NOT**", "**RECOMMENDED**", "**NOT RECOMMENDED**", "**MAY**", and
"**OPTIONAL**" in this document are to be interpreted as described in BCP 14 [@!RFC2119] [@!RFC8174]
when, and only when, they appear in all capitals, as shown here.

# Rationale

[todo]

# Convention

## RRset

This specification does not define any restrictions to the number of TXT records in the RRset, although it is recommended to limit it to one.

## Content of TXT record

The TXT record may contain any valid content, ranging from an empty string to sensible text or URI's. However, it **SHALL NOT** contain any text that is suggesting that the domain is not for sale. Only the absence of the "\for-sale" record is an indication of that. Any existence of a "\_for-sale" TXT record **MUST** be regarded as an indication that it is for sale.

## TTL

There is no special obligation for a certain TTL of the "\_for-sale" TXT record, but it is recommended to keep the TTL at a reasonable level. Not too short, to avoid redundant DNS queries, and no too long, to avoid that a domain name might still show up as for sale, while in fact it no longer is. Any TTL between 3600 and 28800 seems reasonable. A TTL longer than 86400 is **NOT RECOMMENDED**.

## Placement of node name

The "\_for-sale" leaf node name **MAY** be placed on the top level domain, or any domain directly below. It **MAY** also be placed at a lower level, but only when that level is mentioned in the Public Suffix List (PSL| [todo: referentie]. 

Any other placement of the record **MUST NOT** be regarded as a signal that the domain above it is for sale.

See (#placements) for further explanation.

\_for-sale.domain | Situation | Verdict
-------|---------------------|--------
\_for-sale.example | root zone | Allowed
\_for-sale.co.aaa.example | aaa.example in PSL | Allowed
\_for-sale.bbb.example | Second level | Allowed
\_for-sale.www.ccc.example | Other | Not allowed
Table: Allowed placements {#placements}

# Examples

## Example 1: a URI

The owner of 'example.com' wishes to signal that the domain is for sale and adds this record to the 'example.com' zone:

~~~
_for-sale.example.com IN TXT  "https://example.com/forsale.html"
~~~

And interested party notices this signal and can visit the URI mentioned for further information.

As an alternative, a mailto: URI could also be used:

~~~
_for-sale.example.com IN TXT "mailto:owner@example.com"
~~~

[todo] wel/geen subject/body erbij?

There can be a use case for this, especially since WHOIS (or RDAP) often has privact restrictions.

# IANA Considerations

IANA has established the "Underscored and Globally Scoped DNS Node Names" registry. The underscored node name defined in this specification should be added as follows:

~~~ ascii-art
             +-----------+--------------+-------------+
             | RR Type   | _NODE NAME   | Reference   |
             +-----------+--------------+-------------+
             | TXT       | _for-sale    | TBD         |
             +-----------+--------------+-------------+
~~~
Figure: Entry for the "Underscored and Globally Scoped DNS Node Names" Registry


# Security Considerations

One use of the TXT record type defined in this document is to parse the content it contains and to automatically publisch certain information from it on a website or otherwise. There is a risk involved, when the domain owner publishes a malicious URI or one that points to improper content. This may result in reputational damage for the part parsing the record.

Even worse is when the content of the TXT record is not validated sufficiently, opening doors to XSS attacks. Therefore it is **RECOMMENDED** that any parsing and publishing is done with utmost care.

# Acknowledgements

The author would like to thank [todo] for their valuable feedback.

[@-RFC1035] [@-RFC8552] [@-RFC8553] [@-RFC7553]

{backmatter}


