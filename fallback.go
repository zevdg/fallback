// fallback project main.go
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/agl/xmpp"
	"gopkg.in/v0/qml"
	"os"
	"strings"
)

var me *Contact

func main() {
	os.Setenv("APP_ID", "fallback")
	me = &Contact{Id: "fallback@wtfismyip.com", Alias: "Me", IsMe: true}
	password := "password"

	if len(os.Args) > 1 && os.Args[1] == "1" {
		me.Id = "fallback2@wtfismyip.com"
	}
	userName := me.Id[:strings.Index(me.Id, "@")]
	conn, err := xmpp.Dial("wtfismyip.com:5222", userName, "wtfismyip.com", password, new(xmpp.Config))
	if err != nil {
		panic(err)
	}

	if err = conn.SendStanza(xmpp.ClientPresence{XMLName: xml.Name{"jabber:client", "presence"},
		Caps: new(xmpp.ClientCaps)}); err != nil {
		panic(err)
	}

	contacts := NewContacts()
	convos := NewConversations(conn, contacts)

	qml.Init(nil)

	engine := qml.NewEngine()

	engine.Context().SetVar("contactModel", contacts)

	engine.Context().SetVar("convos", convos)

	contacts.add(me)

	go requestRoster(conn, contacts)

	go runXmpp(conn, convos, contacts)

	if err := runQml(engine); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type Conversations struct {
	contacts *Contacts
	conn     *xmpp.Conn
	ConvoMap map[string]*Conversation
	Current  *Conversation
}

func NewConversations(connection *xmpp.Conn, contacts *Contacts) *Conversations {
	return &Conversations{contacts: contacts,
		conn:     connection,
		ConvoMap: make(map[string]*Conversation),
		Current:  &Conversation{}}
}

func (c *Conversations) Get(id string) *Conversation {
	convo, ok := c.ConvoMap[id]
	if !ok {
		convo = NewConversation(c.contacts.GetById(id), c.conn)
		c.ConvoMap[id] = convo
	}
	return convo
}

func (c *Conversations) ChangeCurrent(id string) {
	c.Current = c.Get(id)
	qml.Changed(c, &c.Current)
}

func (c *Conversations) remove(id string) {
	delete(c.ConvoMap, id)
}

type Conversation struct {
	With    *Contact
	conn    *xmpp.Conn
	history []*Message
	Len     int
}

func NewConversation(contact *Contact, conn *xmpp.Conn) *Conversation {
	return &Conversation{With: contact, conn: conn, history: make([]*Message, 0, 10)}
}

type Message struct {
	Sender *Contact
	Msg    string
}

//func (m Message)

func (c *Conversation) AddMsg(msg *Message) {
	c.history = append(c.history, msg)
	c.Len++
	qml.Changed(c, &c.Len)
}

func (c *Conversation) Send(message string) {
	c.AddMsg(&Message{Sender: me, Msg: message})
	fmt.Printf("msg: %#v\n", c.history[len(c.history)-1])
	c.conn.Send(c.With.Id, message)
}

func (c *Conversation) GetMessageByIndex(index int) *Message {
	return c.history[index]
}

type Contacts struct {
	contactMap  map[string]*Contact
	contactList []*Contact
	Len         int
}

func NewContacts() *Contacts {
	return &Contacts{contactMap: make(map[string]*Contact)}
}

func (c *Contacts) add(contact *Contact) {
	c.contactMap[contact.Id] = contact
	if !contact.IsMe {
		c.contactList = append(c.contactList, contact)
		c.Len++
		qml.Changed(c, &c.Len)
	}
}

func (c *Contacts) GetByIndex(index int) *Contact {
	return c.contactList[index]
}

func (c *Contacts) GetById(id string) *Contact {
	contact, ok := c.contactMap[id]
	if !ok {
		panic("contact '" + id + "' doesn't exist")
	}
	return contact
}

type Contact struct {
	Id    string
	Alias string
	IsMe  bool
}

func (c Contact) Name() string {
	name := c.Alias
	if name == "" {
		name = c.Id[:strings.LastIndex(c.Id, "@")]
	}
	return name
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

func runXmpp(conn *xmpp.Conn, convos *Conversations, contacts *Contacts) {
	s, err := conn.Next()
	for ; err == nil; s, err = conn.Next() {

		switch val := s.Value.(type) {
		case *xmpp.ClientMessage:
			fmt.Printf("Client Message: %#v\n", val)
			sender := val.From[:strings.LastIndex(val.From, "/")]
			convo := convos.Get(sender)
			if val.Body != "" {
				convo.AddMsg(&Message{Sender: contacts.GetById(sender), Msg: val.Body})
			}

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

func requestRoster(conn *xmpp.Conn, contacts *Contacts) {

	rosterChan, cookie, err := conn.RequestRoster()
	if err != nil {
		panic(err)
	}

	var _ = cookie //TODO should probably keep this around or something

	s, ok := <-rosterChan
	if !ok {
		panic("Error recieving on roster channel")
	}

	roster, err := xmpp.ParseRoster(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Roster is:")
	for _, r := range roster {
		fmt.Printf("%#v\n", r)
		contacts.add(&Contact{Id: r.Jid, Alias: r.Name})
	}
}
