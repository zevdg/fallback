// fallback project main.go
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/agl/xmpp"
	"gopkg.in/v0/qml"
	"os"
)

var userName string

func main() {
	os.Setenv("APP_ID", "fallback")
	userName = "fallback2"
	password := "password"

	if len(os.Args) > 1 && os.Args[1] == "1" {
		userName = "fallback"
	}

	conn, err := xmpp.Dial("wtfismyip.com:5222", userName, "wtfismyip.com", password, new(xmpp.Config))
	if err != nil {
		panic(err)
	}

	fmt.Println("Sending Presence")
	if err = conn.SendStanza(xmpp.ClientPresence{XMLName: xml.Name{"jabber:client", "presence"},
		Caps: new(xmpp.ClientCaps)}); err != nil {
		panic(err)
	}

	contacts := NewContacts()
	convos := NewConversations(conn, contacts)

	qml.Init(nil)

	qml.RegisterTypes("fallback", 1, 0, []qml.TypeSpec{
		{Init: func(value *Contact, object qml.Object) {}},
		{Init: func(value *Conversation, object qml.Object) {}},
	})

	engine := qml.NewEngine()

	engine.Context().SetVar("contactModel", contacts)

	engine.Context().SetVar("convos", convos)

	go requestRoster(conn, contacts)

	go runXmpp(conn)

	if err := runQml(engine); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type Conversations struct {
	contacts *Contacts
	conn     *xmpp.Conn
	ConvoMap map[string]*Conversation
}

func NewConversations(connection *xmpp.Conn, contacts *Contacts) *Conversations {
	return &Conversations{contacts: contacts,
		conn:     connection,
		ConvoMap: make(map[string]*Conversation)}
}

func (c *Conversations) Get(id string) *Conversation {
	convo, ok := c.ConvoMap[id]
	if !ok {

		convo = NewConversation(c.contacts.GetById(id), c.conn)
		c.ConvoMap[id] = convo
	}
	return convo
}

func (c *Conversations) remove(id string) {
	delete(c.ConvoMap, id)
}

type Conversation struct {
	With    *Contact
	conn    *xmpp.Conn
	History []Message
	//Len     int
}

func NewConversation(contact *Contact, conn *xmpp.Conn) *Conversation {
	return &Conversation{With: contact, conn: conn, History: make([]Message, 0, 10)}
}

type Message struct {
	Sender string
	Msg    string
}

func (c Conversation) Send(message string) {
	c.History = append(c.History, Message{Sender: userName, Msg: message})
	c.conn.Send(c.With.Id, message)
	//c.Len++
	//qml.Changed(c, &c.Len)
}

func (c Conversation) GetMessageByIndex(index int) Message {
	return c.History[index]
}

type Contacts struct {
	contactMap  map[string]Contact
	contactList []*Contact
	Len         int
}

func NewContacts() *Contacts {
	return &Contacts{contactMap: make(map[string]Contact)}
}

type Contact struct {
	Id string
}

func (c *Contacts) add(id string) {
	contact := Contact{Id: id}
	c.contactMap[id] = contact
	c.contactList = append(c.contactList, &contact)
	c.Len++
	qml.Changed(c, &c.Len)
}

func (c *Contacts) GetByIndex(index int) *Contact {
	return c.contactList[index]
}

func (c *Contacts) GetById(id string) *Contact {
	contact, ok := c.contactMap[id]
	if !ok {
		panic("contact " + id + " doesn't exist")
	}
	return &contact
}

func runQml(engine *qml.Engine) error {

	engine.On("quit", func() { os.Exit(0) })

	component, err := engine.LoadFile("fallback.qml")
	if err != nil {
		return err
	}
	window := component.CreateWindow(nil)
	window.Show()
	window.Wait()

	return nil
}

func runXmpp(conn *xmpp.Conn) {
	s, err := conn.Next()
	for ; err == nil; s, err = conn.Next() {

		switch val := s.Value.(type) {
		case *xmpp.ClientMessage:
			//fmt.Printf("Client Message: %#v\n", val)
			fmt.Printf("expected type %T\n", val)

		case *xmpp.ClientPresence:
			//fmt.Printf("Client Presence: %#v\n", val)
			fmt.Printf("expected type %T\n", val)

		case *xmpp.ClientIQ:
			fmt.Printf("ClientIQ: %#v\n", val)

		default:
			fmt.Printf("unexpected type %T\n", val)

		}

	}

	if err != nil {
		panic(err)
	}
}

func requestRoster(conn *xmpp.Conn, model *Contacts) {

	fmt.Println("requesting roster")
	rosterChan, cookie, err := conn.RequestRoster()
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
		fmt.Printf("%#v\n", r)
		name := r.Name
		if name == "" {
			name = r.Jid
		}
		model.add(name)
	}
}
