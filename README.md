Fallback Messenger
========

The smart multi-protocol IM and SMS client for Ubuntu Touch

## Tentative Roadmap Milestones
- a simple xmpp client
- gchat support with [Ubuntu.OnlineAccounts API](http://developer.ubuntu.com/api/qml/sdk-14.04/Ubuntu.OnlineAccounts/)
- merge contacts between multiple accounts and phone's contact list
- SMS support
- implement "fallback" logic (automatically choose service based on availability)
- facebook chat support

## Setup
#### Prereqs
golang 1.2 or higher
Ubuntu Touch SDK

#### Checkout
go get github.com/psywolf/fallback


#### Compile
cd $GOPATH/src/github.com/psywolf/fallback

go build

#### Run
./fallback