# IM Server
Phantom IM's backend server

### Table Of Contents:
* Overview
* Features
* Setup
* Contributing
* Credits

## Overview
Phantom IM is a fully featured backend re-implementation for old IMs. Supported protocols listed below.

Useful Links:
* [Wiki](https://wiki.phantom-im.xyz)

*Registration and Homepage coming soon*

## Features

MySpaceIM (MSIM Protocol):
* MSIMv1 - MSIMv6 Supported
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

## Setup

0. Clone Repo
1. Create a Database with the following included .sql file
2. Run "go build"
3. Rename config.example to config and configuration file
4. Start server
5. Login using user test and password test (currently stored plaintext) **(please change the details of this user before any public use)**

## Contributing

If you'd like to contribute, please join our [Discord](https://discord.gg/UPHUsumXVM) and message one the Developers directly.

## Credits
* EthernalRaine (Lu): Creating the original MMS-Ghost, adding beta MSIM support, adding MSNP2 support and archiving MSIM Clients
* EinTim23 (Tim): Rewriting MMS-Ghost as Phantom-IM in Golang, reversing MSIM Clients and adding almost full MSIM support
