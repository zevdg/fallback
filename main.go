// fallback project main.go
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/agl/xmpp"
	"os"
)

func main() {

	user := "fallback2"

	if len(os.Args) > 1 && os.Args[1] == "1" {
		user = "fallback"
	}

	c, err := xmpp.Dial("wtfismyip.com:5222", user, "wtfismyip.com", "password", new(xmpp.Config))
	if err != nil {
		panic(err)
	}

	fmt.Println("Sending Presence")
	if err = c.SendStanza(xmpp.ClientPresence{XMLName: xml.Name{"jabber:client", "presence"},
		Caps: new(xmpp.ClientCaps)}); err != nil {
		panic(err)
	}

	//c.Send("fallback@wtfismyip.com", "sent from go")

	go requestRoster(c)

	s, err := c.Next()
	for ; err == nil; s, err = c.Next() {

		switch val := s.Value.(type) {
		case *xmpp.ClientMessage:
			fmt.Printf("Client Message:%#v\n", val)

		case *xmpp.ClientPresence:
			fmt.Printf("Client Presence:%#v\n", val)

		default:
			fmt.Printf("unexpected type %T\n", val)

		}

	}

	if err != nil {
		panic(err)
	}
}

func requestRoster(c *xmpp.Conn) {

	fmt.Println("requesting roster")
	rosterChan, cookie, err := c.RequestRoster()
	if err != nil {
		panic(err)
	}

	var _ = cookie //TODO should probably keep this around or something

	s, ok := <-rosterChan
	if !ok {
		panic("No Roster Receieved")
	}

	roster, err := xmpp.ParseRoster(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Roster is:")
	for _, r := range roster {
		fmt.Println(r)
	}
}
