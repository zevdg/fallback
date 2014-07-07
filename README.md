Fallback Messenger
========

The smart multi-protocol IM and SMS client for Ubuntu Touch

At it's core, Fallback will be a multi-protocol IM client (like pidgin), but Fallback takes it a step further.  Instead of sending a message to Jon's gmail or Jon's facebook or Jon's phone number.  You will just send a message to Jon, and the app will figure out which service to use based on where Jon is online, which service was used to talk to you last, and how recently that was.

 For updates and release announcements, [subscribe here](https://plus.google.com/u/0/communities/101599674231948077682) and turn on notifications.

## Tentative Roadmap Milestones
- ~~a simple xmpp client~~
- ~~google hangouts (a.k.a google talk, gtalk, gchat) support~~
- release as click app on Ubuntu Touch
- merge contacts between multiple accounts and phone's contact list
- SMS support
- implement "fallback" logic (automatically choose service based on availability)
- facebook chat support
- [Ubuntu.OnlineAccounts API](http://developer.ubuntu.com/api/qml/sdk-14.04/Ubuntu.OnlineAccounts/) integration

## Setup
wget https://raw.githubusercontent.com/psywolf/fallback/master/dev_env_setup.sh
source dev_env_setup.sh
tested on a fresh install of 14.04 64 bit
(warning: this script will set your gopath to $HOME/go)

#### Alternatively, install prereqs and then clone
#Prereqs
golang 1.2 or higher
Ubuntu Touch SDK
go-qml
#clone
go get github.com/psywolf/fallback


#### Compile
cd $GOPATH/src/github.com/psywolf/fallback

go install

#### Run
$GOPATH/bin/fallback
