
Self-Hosting an AnonVPN Service ([home](/))
===========================================

For Yourself
------------

Hosting AnonVPN for yourself is a great way to connect to private services on
your own network, or to protect your traffic while travelling by routing it to
a trusted computer on a trusted network before a malicious wi-fi access point
can be used to attack the traffic. Tampered traffic will simply be rejected by
the VPN, and most forms of tampering will be impossible anyway.


In the interest of others
-------------------------

While you should always use the [Tor Browser Bundle](https://torproject.org) for
the maximum level of anonymity when browsing the web, AnonVPN is itself a
powerful device for circumventing censorship and protecting yourself from many
meaningful classes of attacker. Providing this as a service to people who need
it is a good we encourage. If one wishes, one may combine these authentication
types on the same service as a financial-auth service, possibly as a form of
donation-collecting.

### Password-based account generation, account-based authentication

There are two "Flavors" of password-based authentication. In "Simple" mode,
the software generates a password for you, which you share with one person.
When that person connects, it generates a new password and shows it to the
systems administrator to give to the next person. A list of base64 addresses
authorized in this way, the now-useless password they created their account
with, and the time of their account creation may be recorded for administrative
purposes such as revoking keys that are compromised. This does not scale well,
but it's hopefully very easy to understand and set up, and very hard to do
unsafely, and may be appealing to human-rights organizations with non-technical
members.

In "Normal" mode passwords are long-term and part of a username-password pair.
A user with a password may use it to authorize multiple Base64 Addresses for his
that user account. This is more flexible, but in this case passwords are
sensitive and thus they are stored encrypted at all times, and revocation of
keys may need to also include revocation of passwords. This mode may be
appealing to a wide span of users.



As a commercial service
-----------------------

With it's native support for cryptocurrency payments, it's easy to set up an
AnonVPN service for commercial purposes.

### Payment-based account generation, account-based authentication


#### Small-Scale


#### Large-Scale




