<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Welcome to BugSim</title>
	</head>
	<script>
	  'use strict';

      let socket = new WebSocket("ws://".concat(window.location.host).concat("/ws/pod_list"));

      socket.onopen = function(e) {
        // alert("[open] Connection established");
        // alert("Sending to server");
        // socket.send("My name is John");
      };
      
      socket.onmessage = function(event) {
        //alert(`[message] Data received from server: ${event.data}`);
        document.getElementById("pod-tab").innerHTML = event.data;
      };
      
      socket.onclose = function(event) {
        if (event.wasClean) {
          //alert(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
        } else {
          // e.g. server process killed or network down
          // event.code is usually 1006 in this case
          alert('[close] Connection died');
        }
      };
      
      socket.onerror = function(error) {
        alert(`error:`.concat(error));
      };
	</script>
	<body>
	    <div class=welcome-msg>Willkommen zum Bugsimulator {{ .Version }}<p>Möchtest du Bug spielen? Du hast eine {{ .SuccessRate }}% Wahrscheinlichkeit, den Pod zu erschießen, den du auswählst. Klicke einfach einen der Links mit den Namen der Pods unten</div>
        <div id="pod-tab">
            <div id="pre-table-data">
                Warte auf Server Data....
            </div>
        </div>
	</body>
</html>