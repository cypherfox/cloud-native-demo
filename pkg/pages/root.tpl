<!DOCTYPE html>
<html>
	<head>
		<title>Welcome to BugSim</title>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>Bootstrap demo</title>
        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-rbsA2VBKQhggwzxH7pPCaAqO46MgnOM80zW1RWuH61DGLwZJEdK2Kadq2F9CUG65" crossorigin="anonymous">
        <link href="style.css" rel="stylesheet">
        
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
        <div class="container-fluid">
            <div class="row">

                <div class="d-flex flex-column flex-shrink-0 p-3 text-bg-dark" style="width: 280px;">
                    <a href="/" class="d-flex align-items-center mb-3 mb-md-0 me-md-auto text-white text-decoration-none">
                        <svg class="bi pe-none me-2" width="40" height="32"><use xlink:href="#bootstrap"/></svg>
                        <span class="fs-4">bugsim</span>
                    </a>
                    <hr>
                    <ul class="nav nav-pills flex-column mb-auto">
                        <li class="nav-item">
                            <a href="#" class="nav-link text-white" aria-current="page">
                                <svg class="bi pe-none me-2" width="16" height="16"><use xlink:href="#home"/></svg>
                                Introduction
                            </a>
                        </li>
                        <li>
                            <a href="#" class="nav-link active">
                                <svg class="bi pe-none me-2" width="16" height="16"><use xlink:href="#speedometer2"/></svg>
                                Desired State
                            </a>
                        </li>
                        <li>
                            <a href="#" class="nav-link text-white">
                                <svg class="bi pe-none me-2" width="16" height="16"><use xlink:href="#speedometer2"/></svg>
                                Metrics Dashboards
                            </a>
                        </li>
                        <li>
                            <a href="#" class="nav-link text-white">
                                <svg class="bi pe-none me-2" width="16" height="16"><use xlink:href="#speedometer2"/></svg>
                                Logging
                            </a>
                        </li>
                        <li>
                            <a href="#" class="nav-link text-white">
                                <svg class="bi pe-none me-2" width="16" height="16"><use xlink:href="#table"/></svg>
                                Auto-Scaling
                            </a>
                        </li>
                        <li>
                            <a href="#" class="nav-link text-white">
                                <svg class="bi pe-none me-2" width="16" height="16"><use xlink:href="#grid"/></svg>
                                Ingress
                            </a>
                        </li>
                        <li>
                            <a href="#" class="nav-link text-white">
                                <svg class="bi pe-none me-2" width="16" height="16"><use xlink:href="#people-circle"/></svg>
                                Cert-Manager
                            </a>
                        </li>
                        <li>
                            <a href="#" class="nav-link text-white">
                                <svg class="bi pe-none me-2" width="16" height="16"><use xlink:href="#people-circle"/></svg>
                                Service Mesh
                            </a>
                        </li>
                    </ul>
                </div>
                
                <div class="col">
                    <div class="row"><h1 class="row welcome-msg">Willkommen zum Bugsimulator {{ .Version }}</h1></div>
                    <div class="row">Möchtest du Bug spielen? Du hast eine {{ .SuccessRate }}% Wahrscheinlichkeit, den Pod zu erschießen, den du auswählst. 
                        <p>Klicke einfach einen der Links mit den Namen der Pods unten.
                    </div>

                    <div class="row" id="pod-tab">
                        <div id="pre-table-data">
                            Warte auf Server Data....
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-kenU1KFdBIe4zVF0s0G1M5b4hcpxyD9F7jL+jjXkk+Q2h455rYXK/7HAuoJl+0I4" crossorigin="anonymous"></script>
	</body>
</html>