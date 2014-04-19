import QtQuick 2.0
import Fallback.Messenger.FileIO 1.0
import Ubuntu.Components 0.1
import QtWebKit 3.0

Rectangle {
    id:base
    visible:false
    property string baseUrl: "https://accounts.google.com/o/oauth2"
    property string tokenUrl: baseUrl + "/token"
    property string emailUrl: "https://www.googleapis.com/oauth2/v1/userinfo"

    property string accessToken
    property string refreshToken
    property int expiresIn
    property string tokenType
    property string id_token

    property string clientId: "166761033436.apps.googleusercontent.com"
    property string clientSecret:"pzltJu6yqEReQYpz7ejoywbt"

    property string userInfoEmailScope: "https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email"
    property string contactsScope: "https%3A%2F%2Fwww.google.com%2Fm8%2Ffeeds"
    property string gtalkScope: "https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fgoogletalk"
    property string redirectUri: "http://localhost:2376"
    property string responseType: "code"
    property string code

    function getInitialRequestUrl(){
        return "https://accounts.google.com/o/oauth2/auth?scope="+contactsScope+"+"+gtalkScope+"+"+userInfoEmailScope+"&redirect_uri="+redirectUri+"&response_type="+responseType+"&client_id="+clientId
    }

    function getEmail(funcHandler){
        httpRequest(emailUrl+"?access_token="+accessToken, function(response){funcHandler(response.email);});
    }

    function getCodeFromUrl(url){
        var schemaRE = /^http:\/\/localhost:2376\/\?code=(.+)/;
        var match = schemaRE.exec(url)
        if (match && match.length === 2) {
            code = match[1];
            return true;
        } else {
            return false;
        }
    }

    function generateCallback(funcHandler){
        function callbackFunc(response){
            updateTokens(response);
            funcHandler(base.accessToken);
        }
        return callbackFunc;
    }

    function updateTokens(response){

        base.accessToken = response.access_token;
        base.expiresIn = response.expires_in;
        base.tokenType = response.token_type;

        if(response.id_token){
            id_token = response.id_token;
        }

        if(response.refresh_token){
            base.refreshToken = response.refresh_token;
            oauthFile.write(base.refreshToken);
        }
    }

    function getAccessToken(funcHandler){
        httpRequest(tokenUrl, generateCallback(funcHandler),
                    "code="+code
                    +"&client_id="+clientId
                    +"&client_secret="+clientSecret
                    +"&redirect_uri="+redirectUri
                    +"&grant_type=authorization_code");
    }

    function refreshAccessToken(funcHandler){

        if(!refreshToken){
            refreshToken = oauthFile.read();
        }
        if(!refreshToken){
            return false;
        }

        httpRequest(tokenUrl, generateCallback(funcHandler),
                    "client_id="+clientId
                    +"&client_secret="+clientSecret
                    +"&refresh_token="+refreshToken
                    +"&grant_type=refresh_token");
        return true;
    }

    function httpRequest(url, funcHandler, postData) {


        //constants and helper functions are above
        //real code starts here

        var doc = new XMLHttpRequest();
        doc.onreadystatechange = function() {
            if (doc.readyState === XMLHttpRequest.DONE) {
                console.debug(doc.responseText);
                var response = JSON.parse(doc.responseText);
                funcHandler(response);
            }
        }

        if(postData){
            doc.open("POST", url );
            doc.setRequestHeader("Content-Type", 'application/x-www-form-urlencoded');
            doc.setRequestHeader('Content-length', postData.length);
        }else{
            doc.open("GET", url);
        }

        doc.setRequestHeader("Connection", "close");

        doc.send(postData);
    }

    FileIO {
        id: oauthFile
        source: oauthFile.homePath() + "/.fallback/oauth"
        //onError: console.log(msg)
    }
}
