# will move to chimera-im/server soon, expect shit to break

# server
Chimera IM's rewritten backend server

### Table Of Contents:
* Community
* Features
* Contributing
* Credits

## Community
* [Homepage](https://chimera.im)
* [Discord](https://discord.gg/UPHUsumXVM)

## Features

MySpaceIM (MSIM Protocol):
* MSIMv1 - MSIMv7 Supported
* Authentication (RC4)
* Contacts (Add and Remove)
* Instant Messages
* Status Messages
* Uploading Profile Pictures
* Offline Messages
* "Zaps" (Action Messages)
* Typing Indicators
* Profile Information
* Advertisement Server (now randomized)

MSN Messenger (MSNP Protocol):
* MSNP2, MSNP3 Supported
* Dispatching to NS
* Dispatching to SB
* NS Authentication (CTP, MD5)
* SB Authentication

Please check the Issues for broken and/or missing features!

## Contributing

If you'd like to contribute, please join our [Discord](https://discord.gg/UPHUsumXVM) and message one the Developers directly.

## Credits
* EthernalRaine (Lu): Creating the original MMS-Ghost, adding beta MSIM support, adding MSNP2 support and archiving MSIM Clients
* EinTim23 (Tim): Rewriting MMS-Ghost as Phantom-IM in Golang, reversing MSIM Clients and adding almost full MSIM support
* henpett1 (Henpett): Adding rotating ads system into MSIM AdServer API
* pinksub (Tox): Implementing fixes for the dreaded 4096 byte bug
