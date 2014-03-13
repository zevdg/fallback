import QtQuick 2.0
import Ubuntu.Components 0.1

MainView {
    width: units.gu(48)
    height: units.gu(60)

    Page {
        title: "Contact List"

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
                    Text { text: contactModel.contact(index).name }
                }
                Column {
                    id: status
                    width: units.gu(5)
                    anchors.right : parent.right
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