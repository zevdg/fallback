Fallback Messenger
========

The smart multi-protocol IM and SMS client for Ubuntu Touch

At it's core, Fallback will be a multi-protocol IM client (like pidgin), but Fallback takes it a step further.  Instead of sending a message to jon@gmail.com, or 555-123-4567 (Jon's number).  You will just send a message to Jon, and the app will figure out which service to use based on where Jon is online, which service was used to talk to you last, and how long ago that was.

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
