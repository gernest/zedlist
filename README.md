zedlist [![Coverage Status](https://coveralls.io/repos/gernest/zedlist/badge.svg?branch=master&service=github)](https://coveralls.io/github/gernest/zedlist?branch=master) [![Build Status](https://drone.io/github.com/gernest/zedlist/status.png)](https://drone.io/github.com/gernest/zedlist/latest)
========
A humble job recruitment service.

# Motivation
I have never been employed. I'm in Tanzania, and God knows if there is anyone who cared about my applications.

The purpose of zedlist is to provide infastructure that bridges the gap between job seekers and employers in the african continent (initially for Tanzania).

# Features
* Job listing.
* Register/Delete/Rename account.
* Create/Delete/Update jobs via JSON API.
* JSON Web Tokens support.
* Structured Resume support.
* Builtin Job applications processing.
* Social account login (facebook, google+, github)
* Search.
* Support multiple databases (currently mysql and postgresql) 

Some of the features aren't complete yet.

# Prerequisites

You must have a database. It doesnt matter if its local or remote, only that you have a working database connection.


# Installation

For installation of zedlist please read [INSTALL.md](INSTALL.md)


# A note about translation
Only swahili and English are supported

# Acknowledgement

* Project structure is heavily inspired by [Gogs](https://github.com/gogits/gogs)
* Middlewares and routes are based on [Echo](https://github.com/labstack/echo)


# Contributing

Contributions are welcome and before you do anything please read [DEVELOPER.md](DEVELOPER.md) for more details about zedlist and [CONTRIBUTING.md](CONTRIBUTING.md) for contributions guidelines.

## Author
Geofrey Ernest <geofreyernest@live.com>

## Licence
This project is released under MIT licence see [LICENCE](LICENCE) for more details.
