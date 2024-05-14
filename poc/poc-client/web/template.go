package web

import "text/template"

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

	var sendHandshakeRequest = function() {
		console.log("Gg")
		// ws.send("HANDSHAKE_REQUEST");
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "http://localhost:8081/start-handshaking", "GSGS");
		xhr.onreadystatechange = function() {
			if (xhr.readyState === XMLHttpRequest.DONE) {
				if (xhr.status === 200) {
					console.log("Handshake request sent successfully");
				} else {
					console.error("Failed to send handshake request");
				}
			}
		};
		xhr.send();
		return false;
	};

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        var connInfo = document.getElementById("serverPort").value;
		const serverPort = connInfo.split(":")[0];
		const clientID = connInfo.split(":")[1];
        ws = new WebSocket("ws://localhost:" + serverPort + "/echo/" + clientID);
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("handshake").onclick = function(evt) {
        sendHandshakeRequest();
		// ws.send("HANDSHAKE_REQUEST");
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<label for="serverPort">Connection Info(PORT:ClientID):</label>
<input type="text" id="serverPort" value="8080"><br>
<button id="open">Open</button>
<button id="close">Close</button>

<p><input id="handshake-input" type="text" value="">
<button id="handshake" type="button">Handshake</button>

<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
