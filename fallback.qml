import QtQuick 2.0
import Ubuntu.Components 0.1
import fallback 1.0

MainView {
    width: units.gu(48)
    height: units.gu(60)

    PageStack {
            id: pageStack

        Tabs {
            id: tabs


            //TabBar {
            //    selectionMode: false
            //}


            Tab {
                id: contactListTab
                title: i18n.tr("Contact List")
                Page {
                    id: contactListPage
                    Component {
                        id: contactDelegate
                        Item {
                            width: units.gu(48); height: units.gu(5)
                            Column {
                                id: icon
                                width: units.gu(5)
                            }
                            Column {
                                id: name
                                anchors.left : icon.right
                                Text { text: contactModel.getByIndex(index).id }
                            }
                            Column {
                                id: status
                                width: units.gu(5)
                                anchors.right : parent.right
                            }
                            MouseArea {
                                anchors.fill: parent
                                onClicked: { 
                                    conversationTab.init(contactModel.getByIndex(index).id)
                                }
                            }
                        }
                    }

                    ListView {
                        anchors.fill: parent
                        model: contactModel.len
                        delegate: contactDelegate
                        highlight: Rectangle { color: "lightsteelblue"; radius: 5 }
                        focus: true
                    }
                }
            }
            Tab {
                id: conversationTab
                title: i18n.tr("Contact Name")
                property string withId
                Page {
                    id: conversationPage
                    Column {
                        id: messages

                    }
                    /*
                    Component {
                        id: messageDelegate
                        Item {
                            Text { text: convo().getMessageByIndex(index).sender + ": "
                                    + convo().getMessageByIndex(index).msg }
                        }
                    }
                    ListView {
                        id: messageList
                        anchors.top : parent.top
                        anchors.bottom : entry.bottom
                        anchors.left : parent.left
                        anchors.right : parent.right
                        //model: convo().len
                        delegate: messageDelegate
                        highlight: Rectangle { color: "lightsteelblue"; radius: 5 }
                        focus: true
                    }
                    */
                    TextArea {
                        id: entry
                        anchors.bottom: parent.bottom
                        anchors.left: parent.left
                        anchors.right: sendBtn.left

                        autoSize: true
                        maximumLineCount: 3
                    }

                    Button {
                        id: sendBtn
                        width: 32
                        height: 32
                        anchors.bottom: parent.bottom
                        anchors.right: parent.right

                        iconSource: '/usr/share/icons/gnome/48x48/actions/document-send.png'
                        onClicked: {
                            parent.convo().send(entry.text)
                            entry.text = ""
                        }
                    }

                    function convo(){
                        return convos.get(conversationTab.withId)
                    }
                }
                function init(id){
                    conversationTab.title = id;
                    conversationTab.withId = id;
                    tabs.selectedTabIndex = 1;
                    //messageList.model = conversationPage.convo().len
                }
            }
        }
    }

}