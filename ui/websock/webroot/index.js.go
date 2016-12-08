package webroot

func init() {
	webFiles["index.js"] = indexJS
}

const indexJS = `
var intercept = {
	EDITORCHANNEL: "proxy/intercept/editor",
	SETTINGSCHANNEL: "proxy/intercept/options",
	HISTORYCHANNEL: "proxy/httpHistory",
}
var waptyServer= new WebSocket("ws://localhost:8081/ws");
waptyServer.onopen = function(event){
	console.log("WebSocket connected");
	var msg = {
		Action: "intercept",
		Channel: intercept.SETTINGSCHANNEL
	}
	waptyServer.send(JSON.stringify(msg));

}

//scope this
var tmpHistory = null
var historyTable = document.getElementById("historyTable");

waptyServer.onmessage = function(event){
	//	console.log(event.data);
	msg = JSON.parse(event.data);
	//console.log(msg);
	switch (msg.Channel){
		case intercept.EDITORCHANNEL:
			//if ('Payload' in msg){
			document.getElementById("proxybuffer").value=atob(msg.Payload);
			controls = true;
			//}
			break;
		case intercept.SETTINGSCHANNEL:
			switch (msg.Action){
				case "intercept":
					document.getElementById("interceptToggle").checked = msg.Args[0] === "true";
			}
			break;
		case intercept.HISTORYCHANNEL:
			switch (msg.Action){
				case "metaData":
					var metaData = JSON.parse(msg.Args[0])
					if (tmpHistory===null){
						var h = document.getElementById("historyHeader")
						for (var key in metaData) {
							if (metaData.hasOwnProperty(key)) {
								var cell = h.insertCell(-1)
								cell.innerText = key
							}
						}
						tmpHistory = {}
						historyTable.style.display='block';
						//FIXME this doesn't work. Why?
						$("#historyTable").colResizable({resizeMode:'overflow'});
						$("#historyTable").tablesorter(); 
					}
					var stringID=""+metaData.id;
					console.log("Got metaData for id " + stringID);
					if (stringID in tmpHistory){
						tmp = tmpHistory[stringID]
						for (var key in metaData) {
							if (metaData.hasOwnProperty(key)) {
								tmp[key].innerText=metaData[key]
							}
						}
						delete tmpHistory[stringID]
					}else{
						var row=historyTable.insertRow(1);
						var tmp={}
						for (var key in metaData) {
							if (metaData.hasOwnProperty(key)) {
								var cell = row.insertCell(-1)
								cell.innerText = metaData[key]
								tmp[key] = cell
							}
						}
						tmpHistory[stringID]=tmp
					}
					/*case "metaData":*/
					//var problem = false;
					//var metaData = JSON.parse(msg.Args[0])
					//console.log("Metadata for request " + metaData.Id + " received:");
					//console.log(metaData)
					//if (""+metaData.Id in debugHistory){
					//if (debugHistory[""+metaData.Id]==1){
					//debugHistory[""+metaData.Id]=2
					//}else{
					//console.log("Problem with request " + metaData.Id);
					//problem=true;
					//}
					//}else{
					//debugHistory[""+metaData.Id]=1
					//}
					//document.getElementById("historyTable").innerHTML=document.getElementById("historyTable").innerHTML + "<tr"+
					//(problem?" style='color:red;' ":"")+
					//"><td>"+metaData.Id+"</td>"+
					//"<td>"+metaData.Host+"</td>"+
					//"<td>"+metaData.Path+"</td>"+
					//"</tr>";
					break;
				case "fetch":
					var pl = JSON.parse(atob(msg.Payload))
					console.log(atob(pl.RawReq))
					console.log(atob(pl.RawRes))
					console.log(atob(pl.RawEditedReq))
					console.log(atob(pl.RawEditedRes))
					break;
			}
			break;
	}
}
waptyServer.onclose=function(event){
	var value = ("Server connection lost, would you like to try to reconnect?")
	if (value){
		location.reload()
	}
}

var controls = false;

function clickhandler(){
	if (!controls){
		return;
	}
	switch (event.target.id){
		case "forwardOriginal":
			var msg = {
				Action: "forward",
				Channel: intercept.EDITORCHANNEL
			}
			controls = false;
			document.getElementById("proxybuffer").value="";
			waptyServer.send(JSON.stringify(msg));
			break;
		case "forwardModified":
			var payload = btoa(document.getElementById("proxybuffer").value);
			var msg = {
				Action: "edit",
				Channel: intercept.EDITORCHANNEL,
				Payload: payload
			}
			controls = false;
			document.getElementById("proxybuffer").value="";
			waptyServer.send(JSON.stringify(msg));
			//var xhr = new XMLHttpRequest();
			//xhr.open("POST", "/edit", true);
			//xhr.setRequestHeader('Content-Type', 'application/json');
			//xhr.send(JSON.stringify(msg));
			break;
		case "drop":
			var msg = {
				Action: "drop",
				Channel: intercept.EDITORCHANNEL,
			}
			controls = false;
			document.getElementById("proxybuffer").value="";
			waptyServer.send(JSON.stringify(msg));
			break;
		case "provideResponse":
			var payload = btoa(document.getElementById("proxybuffer").value);

			var msg = {
				Action: 	"provideResp",
				Channel: intercept.EDITORCHANNEL,
				Payload: payload
			}
			controls = false;
			document.getElementById("proxybuffer").value="";
			waptyServer.send(JSON.stringify(msg));
			//var xhr = new XMLHttpRequest();
			//xhr.open("POST", "/edit", true);
			//xhr.setRequestHeader('Content-Type', 'application/json');
			//xhr.send(JSON.stringify(msg));
			break;
		default:
			console.log("unknown event")
	}
}
function toggler(){
	var msg = {
		Action: "intercept",
		Channel: intercept.SETTINGSCHANNEL,
		Args: [""+document.getElementById("interceptToggle").checked]
	}
	waptyServer.send(JSON.stringify(msg));
}
function fetchHistory(id){
	var msg = {
		Action: "fetch",
		Channel: intercept.HISTORYCHANNEL,
		Args: [""+id]
	}
	waptyServer.send(JSON.stringify(msg));
}
`
