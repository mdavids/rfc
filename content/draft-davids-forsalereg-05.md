%%%
# This is a comment - but only in this block
title = "Registration of Underscored and Globally Scoped DNS Node Name: \"_for-sale\""
abbrev = "forsalereg"
#ipr = "trust200902" # mmark docs say 'none' for independent submissions
ipr = "none"
# area = "Internet"
# workgroup = "Internet Engineering Task Force (IETF)"
submissiontype = "independent"
keyword = [""]
# https://www.rfc-editor.org/rfc/rfc7991#section-2.45.14
tocdepth = 3
# date = 2022-12-22T00:00:00Z

# See FAQ: "How Do I Create an Independent IETF Document?"
# https://mmark.miek.nl/post/faq/
[seriesInfo]
name = "Internet-Draft"
value = "draft-davids-forsalereg-05"
stream = "independent"
status = "informational"  # or "bcp" or "experimental" ?

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

This document defines an operational convention for using the reserved DNS node name
"\_for-sale" to indicate that the parent domain name is available for purchase. 
This approach offers the advantage of easy deployment without affecting ongoing operations. As such, the method can be applied to a domain name that is still in full use.

{mainmatter}

# Introduction

Well-established services [@RFC3912; @RFC9083] exist to determine whether a domain name is registered. However, the fact that a domain name exists does not necessarily mean it
is unavailable; it may still be for sale.

Some registrars and other entities offer mediation services between domain name holders and interested parties; however, for domain names not for sale, such services may be unnecessary.

This specification defines a simple and universal method to ascertain whether a domain name, although registered, is available for purchase. It enables a domain name holder to add a reserved underscored node name [@!RFC8552] in the zone, indicating that the domain name is for sale.

The TXT RR type [@!RFC1035] that is created for that purpose **MAY** contain a pointer, such as a
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

The TXT [@RFC8553, (see) section 2.1] record **MUST** contain any valid content, ranging from an empty string to meaningful text or URIs. However, it **SHALL NOT** contain any text that suggests that the domain is not for sale. If a domain name is not for sale, the "\_for-sale" indicator
**MUST NOT** be used. Any existence of a "\_for-sale" TXT
record, assuming it is not a wildcard,  **MAY** therefore be regarded as an indication that the domain name is for sale.

This specification does not dictate the exact use of any content in the "\_for-sale" TXT record, or the lack of any such content. Parties - such as
registries and registrars - may use it in their tools, perhaps even by defining additional requirements that the content must meet. Alternatively, an individual can use it in combination with existing tools to make contact with the seller.

The content of the TXT record is "as is" and characters such as ";" between two URIs for example, have no defined meaning. It is up to the processor of the content to decide how to handle it. See
(#guidelines) for additional guidelines.

## RRset limitations

This specification does not define any restrictions on the number of TXT records in the RRset, although it is recommended to limit it to one. It is also recommended that the length of the RDATA [@RFC8499] does not exceed 255 bytes. If the RRset contains multiple records or the total size exceeds 255 bytes, it is up to the processor to determine which data to use. For example, a
registry might pick a mandatory URI from the RRset to display on a website as part of its service, while an individual might just pick a phone number (if present) and dial it to make contact with a potential seller.

## RR Type limitations

Adding any other RR types under the "\_for-sale" leaf but TXT is not
recommended and they **MUST** be ignored for the purpose of this document.

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
_for-sale.example.com. IN TXT "https://broker.example.net/offer?id=3"
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

There can be a use case for these URIs, especially since WHOIS (or RDAP) often has privacy restrictions.
But see the (#privacy, use title) section for possible downsides.

## Example 2: Various other approaches

Free format text, to make the availability more explicit:

~~~
_for-sale.example.com. IN TXT "I'm for sale: info [at] example.com"
~~~

Proprietary format, used by a registry or registrar to automatically redirect visitors to a web page,
but which has no well-defined meaning to third parties:

~~~
_for-sale.example.com. IN TXT "fscode=aHR0cHM...V4YW1wbGUuY29t"
~~~

The content in the following example could be malicious, but it is not in violation of this specification (see (#security)):

~~~
_for-sale.example.com. IN TXT "<script>alert('H4x0r')</script>"
~~~

# Operational Guidelines {#guidelines}
DNS wildcards interact poorly with underscored names, which is why the use of wildcards
is **NOT RECOMMENDED** when deploying this mechanism. But they may still be 
encountered in practice, especially by operators who are not deploying this
mechanism. Therefore, any assumptions about the content of "\_for-sale" 
TXT records should be made with caution. 

For instance, some operators configure wildcards to return a fixed "v=spf1 -all"
TXT record for all subdomains. In such cases, the presence of a "\_for-sale" TXT record 
containing this content does not indicate that the domain is actually for sale. 

To minimize confusion, it is **RECOMMENDED** to include content that is recognizable either 
by humans or automated systems, such as the "fscode=" string or the descriptive text 
shown in the (#examples, use title) section.

As an alternative. the situation can be circumvented by adding a "\_for-sale" leaf node with a 
different RR type, anything other than TXT. Although being an exception to the
recommendations, it will prevent confusing wildcard responses to TXT queries.

For example:

~~~
_for-sale.example.com. IN HINFO "NOT A TXT" "NOT FOR SALE"
~~~

In general it is best to avoid the above wildcard situation completely.

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


# Privacy Considerations {#privacy}

The use of the "\_for-sale" node name publicly indicates the intent to sell a domain name.
Domain owners should be aware that this information is accessible to anyone querying the
DNS and may have privacy implications.

There is a risk of data scraping, such as email addresses and phone numbers.

# Security Considerations {#security}

One use of the TXT record type defined in this document is to parse the content it contains and to automatically publish certain information from it on a website or elsewhere. However, there is a risk if the domain name holder  publishes a malicious URI or one that points to improper content. This may result in reputational damage for the party parsing the record.

Even worse is a scenario in which the content of the TXT record is not validated and sanitized sufficiently, opening doors to - for example - XSS attacks among other things. 

Therefore, it is **RECOMMENDED** that any parsing and publishing is conducted with the utmost care.

There is also a risk that this method will be abused as a marketing tool, or to otherwise lure individuals into visiting certain sites or attempting other forms of contact, without there being any intention to actually sell the particular domain name. Therefore, it is recommended that this method is primarily used by professionals.

# Implementation Status

The concept described in this document is in use with the .nl ccTLD
registry. See for example:

~~~
https://www.sidn.nl/en/whois?q=example.nl
~~~

[*note to editor: please remove this section before publication*]

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


