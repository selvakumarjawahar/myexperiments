import {Component} from '@angular/core';

@Component({
  selector: 'app-servers',
  templateUrl: './servers.component.html',
})
export class ServersComponent {
allowServer = false;
serverCreationStatus = 'Server not Created!';
serverName = 'DefaultValue';
serverCreated = false;
servers = ['TestServer', 'ProductionServer'];
constructor() {
setTimeout(() => {
  this.allowServer = true;
}, 2000);
}
onCreateServer() {
  this.servers.push(this.serverName);
  this.serverCreated = true;
  this.serverCreationStatus = 'Server is Created successfully and server name is ' + this.serverName;
}
onInputEntered(event: Event) {
this.serverName =  (event.target as HTMLInputElement).value;
}
}
