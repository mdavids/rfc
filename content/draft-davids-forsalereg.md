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
# date = 2022-12-20T00:00:00Z

[seriesInfo]
name = "Internet-Draft"
value = "draft-davids-forsale-00"
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
  city = "Arnhem"
  code = "6825 MD"
  pobox = "Meander 501"
  cityarea = "Gld"
%%%

{mainmatter}

# Abstract

This document defines a simple operational convention of using a reserved underscored node name ("\_for-sale") to indicate that the parent domain name above, is for sale. It has the advantage that it can be easily deployed, without affecting any running operations. As such, the method can be applied to a domain name that is still in full use.

# Introduction

Well established services [@RFC3912; @RFC7483] exist, to find out if a domain name is registered or not. But the fact that a domain name exists does not exclude the possibility to obtain it, because it may be up for sale.

Some registrars and various other parties offer (payed) mediation services between domain name holders and interested parties, but for a domain name that is not for sale, such services are a waste of money and time.

This specification defines a simple universal way to find out if a domain name, even thouh it is taken, might be obtained nevertheless. It enables a domain name holder to add a reserved underscored node name [@!RFC8552] in the zone, indicating that the domain name is actually for sale.

The TXT record RRtype [@!RFC1035] that is created for that purpose **MAY** contain a pointer, such as a URI [@RFC8820], to allow an interested party to find information or to get in touch and engage in further arrangements.

With due caution, this information can also be incorporated in the automated availability services, so that when the domain name is checked for availabilty, the service can also indicate it is for sale, including a referral to the selling party's information.

## Terminology

The key words "**MUST**", "**MUST NOT**", "**REQUIRED**", "**SHALL**", "**SHALL NOT**",
"**SHOULD**", "**SHOULD NOT**", "**RECOMMENDED**", "**NOT RECOMMENDED**", "**MAY**", and
"**OPTIONAL**" in this document are to be interpreted as described in BCP 14 [@!RFC2119] [@!RFC8174]
when, and only when, they appear in all capitals, as shown here.

# Rationale

There are undoubtedly more ways to address this problem space. The reasons for the approach defined in this document are primarly accessibility and simplicity. The indicator can be easilty turned on and off at will and more over, it is available right away and does not require major changes in existing services. This allows for a smooth introduction of the concept.

# Convention

## Content limitations

The TXT record **MUST** contain any valid content, ranging from an empty string to sensible text or URI's. However, it **SHALL NOT** contain any text that is suggesting that the domain is not for sale. In the case a domain name is not for sale, the "\_for-sale" indicator MUST NOT be used. Any existence of a "\_for-sale" TXT record **MUST** therefore be regarded as an indication that the domain name is for sale.

This specification does not dictate the exact use of any content in the "\_for-sale" TXT record, or the lack of any such content. Parties, such as Registries and Registrars may use it in their tools, perhaps even by defining additional requirements that the content must meet. Or an individual can use it in combination with existing tools to get in touch with the seller.

The content of the TXT record is "as is" and characters such as ";" between two URIs for example, have no defined meaning. It is up to the processor of the content to decide how to handle them.

## RRset limitations

This specification does not define any restrictions to the number of TXT records in the RRset, although it is recommended to limit it to one. When the RRset contains multiple records, it is at the discretion of the processor to make a selection. For example, a registry might pick a mandatory URI from the RRset, to display on a website as part of their service, whilst and indivual might just pick a possibly present phone number and dial it to get in touch.

## RRtype limitation

Adding any other RRtypes under the "\_for-sale" leaf but TXT is **NOT RECOMMENDED** and they **MUST** be ignored for the purpose of this document.

## TTL limitation

A TTL longer than 86400 is **NOT RECOMMENDED**.

## Wildcard limitation

The "\_for-sale" leaf **MUST NOT** be a wildcard.

## CNAME limitation

The "\_for-sale" leaf **MAY** be a CNAME pointing to a TXT RRtype.

## Placement of node name

The "\_for-sale" leaf node name **MAY** be placed on the top level domain, or any domain directly below. It **MAY** also be placed at a lower level, but only when that level is mentioned in the Public Suffix List [@PSL]. 

Any other placement of the record **MUST NOT** be regarded as a signal that the domain above it is for sale.

See (#placements) for further explanation.

\_for-sale.domain | Situation | Verdict
-------|---------------------|--------
\_for-sale.example | root zone | For sale
\_for-sale.aaa.example | Second level | For sale
\_for-sale.co.bbb.example | bbb.example in PSL | For sale
\_for-sale.www.ccc.example | Other | Invalid
Table: Allowed placements {#placements}

# Examples

## Example 1: a URI

The owner of 'example.com' wishes to signal that the domain is for sale and adds this record to the 'example.com' zone:

~~~
_for-sale.example.com IN TXT  "https://example.com/forsale.html"
~~~

An interested party notices this signal and can visit the URI mentioned for further information. The TXT record can also be processed by automated tools. See the (#security, use title) section for possible risks.

As an alternative, a mailto: URI could also be used:

~~~
_for-sale.example.com IN TXT "mailto:owner@example.com"
~~~

Or a telephone URI:

~~~
_for-sale.example.com IN TXT "tel:+1-201-555-0123"
~~~

There can be a use case for this, especially since WHOIS (or RDAP) often has privacy restrictions.

## Example 2: Various other possibilities

Free format text:

~~~
_for-sale.example.com IN TXT  "I'm for sale: info [at] example.com"
~~~

The content in the next example could be malicious, but it is not in violation of this specification (see (#security)):

~~~
_for-sale.example.com IN TXT  "<script>alert('H4x0r')</script>"
~~~


# IANA Considerations

IANA has established the "Underscored and Globally Scoped DNS Node Names" registry [@!RFC8552; @IANA]. The underscored node name defined in this specification should be added as follows:

~~~ ascii-art
             +-----------+--------------+-------------+
             | RR Type   | _NODE NAME   | Reference   |
             +-----------+--------------+-------------+
             | TXT       | _for-sale    | TBD         |
             +-----------+--------------+-------------+
~~~
Figure: Entry for the "Underscored and Globally Scoped DNS Node Names" Registry


# Security Considerations {#security}

One use of the TXT record type defined in this document is to parse the content it contains and to automatically publish certain information from it on a website or otherwise. There is a risk involved in this, when the domain owner publishes a malicious URI or one that points to improper content. This may result in reputational damage for the party parsing the record.

Even worse is the scenario where the content of the TXT record is not validated sufficiently, opening doors to XSS attacks among other things.

Therefore it is **RECOMMENDED** that any parsing and publishing is done with utmost care and sufficient validation.

There is also a potential risk that this method is abused as a marketing tool, or to otherwise lure individuals into visiting certain sites or other forms of contact, without the intention of actually selling the particular domain name. It is therefore recommended that this method is primarily used by professionals who are sufficiently alert and aware.

# Acknowledgements

The author would like to thank Thijs van den Hout and Caspar Schutijser for their valuable feedback.

[@-RFC8553, section 2.1]

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


