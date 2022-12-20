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

This document defines the operational convention of using a reserved underscored node name TXT RRset in DNS ("\_for-sale") to indicate that the parent domain name above it is for sale. The TXT record type that is created **MAY** contain pointers, such as a URI that allows an interested party to find more information or to engage in further arrangements.

# Introduction

## Terminology

The key words "**MUST**", "**MUST NOT**", "**REQUIRED**", "**SHALL**", "**SHALL NOT**",
"**SHOULD**", "**SHOULD NOT**", "**RECOMMENDED**", "**NOT RECOMMENDED**", "**MAY**", and
"**OPTIONAL**" in this document are to be interpreted as described in BCP 14 [@!RFC2119] [@!RFC8174]
when, and only when, they appear in all capitals, as shown here.

# Rationale

# Convention

## RRset

## Content of TXT record

## TTL

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


