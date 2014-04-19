import QtQuick 2.0
import Ubuntu.Components 0.1
import QtWebKit 3.0

MainView {
    width: units.gu(48)
    height: units.gu(60)

    PageStack {
            id: pageStack
            Component.onCompleted: initialize();

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
                            width: units.gu(48)
                            height: units.gu(5)
                            Item {
                                id: icon
                                width: units.gu(5)
                            }
                            Item {
                                id: name
                                anchors.left: icon.right
                                Text { text: contactModel.getByIndex(index).name() }
                            }
                            Item {
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
                        focus: true
                    }
                }
            }
            Tab {
                id: conversationTab
                title: "Contact Name Goes Here"
                property string withId
                Page {
                    id: conversationPage
                    
                    Component {
                        id: messageDelegate
                        Column {
                            width: units.gu(48)
                            Row{
                                width: parent.width 
                                Text { 
                                    width: parent.width
                                    wrapMode: Text.Wrap
                                    textFormat: Text.StyledText
                                    text: makeMsgText(convos.current.getMessageByIndex(index))

                                    function makeMsgText(msg){
                                        var text = '<font color="'
                                        if ( msg.sender.isMe ){
                                            text += 'red';
                                        }else{
                                            text += 'blue';
                                        }
                                        text += '">' + msg.sender.name() + ": </font>" + msg.msg
                                        return text
                                    }
                                }
                            }
                        }
                    }

                    ListView {
                        id: messageList
                        anchors.top : parent.top
                        anchors.bottom : entry.top
                        anchors.left : parent.left
                        anchors.right : parent.right
                        model: convos.current.len
                        delegate: messageDelegate
                        focus: true
                    }
                    
                    TextArea {
                        id: entry
                        anchors.bottom: parent.bottom
                        anchors.left: parent.left
                        anchors.right: sendBtn.left
                        focus: true
                        Keys.onReturnPressed: parent.send()

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
                        onClicked: parent.send()
                    }

                    function send(){
                        if( entry.text == ""){
                            return;
                        }
                        convos.current.send(entry.text);
                        entry.text = "";
                    }
                }
                function init(id){
                    conversationTab.title = id;
                    tabs.selectedTabIndex = 1;
                    convos.changeCurrent(id)
                }
            }
        }

        Page {
            id: firstRun
            visible: false
            title: "Fallback Messenger"

            Button {
                anchors.centerIn: parent
                onClicked: {
                    pageStack.push(oauthPage)
                }
                text: "Sign in to Google"
            }

            tools: ToolbarItems {
                        locked: true
                        opened: false
                    }

        }

        Page {
            id: oauthPage
            visible: false
            title: "Google Authentication"
            WebView {
                id: webview
                url: oauth.getInitialRequestUrl()
                width: parent.width
                height: parent.height

                onNavigationRequested: {
                    // detectCodeResponse
                    if (oauth.getCodeFromUrl(request.url)) {
                        request.action = WebView.IgnoreRequest;
                        oauth.getAccessToken(loginWithToken)
                        while(pageStack.depth > 1){
                            pageStack.pop();
                        }
                    } else {
                        request.action = WebView.AcceptRequest;
                    }
                }
            }

        }
    }

    function initialize(){
    	console.log("initializing...");
        pageStack.push(tabs);
        if(!oauth.refreshAccessToken(loginWithToken)){
            pageStack.push(firstRun);
        }
    }

    GoogleAuthentication{
        id: oauth
    }

    function loginWithToken (accessToken){
        oauth.getEmail(function(email){
        	console.log("attempting login...");
            convos.login(email, accessToken)
            console.log("finished login")
        })
    }
}