#!/bin/bash

#verify script is beig sourced
if [[ $_ == $0 ]]
then
	echo "script must be sourced.  run command as"
	echo "source $0"
	exit 1
fi

#add repositories
sudo add-apt-repository --yes ppa:ubuntu-sdk-team/ppa
sudo add-apt-repository --yes ppa:webupd8team/sublime-text-3


#udate repo list
sudo apt-get update


#install packages
sudo apt-get install -y ubuntu-sdk qtbase5-private-dev qtdeclarative5-private-dev libqt5opengl5-dev golang sublime-text-installer git


#set gopath
export GOPATH=$HOME/go
mkdir $GOPATH
echo "export GOPATH=$GOPATH" >> $HOME/.bashrc


#get fallback
go get github.com/psywolf/fallback
