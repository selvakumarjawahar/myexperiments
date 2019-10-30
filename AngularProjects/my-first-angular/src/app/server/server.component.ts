import {Component} from '@angular/core';

@Component({
  selector: 'app-server',
  templateUrl: './server.component.html'
})
export class ServerComponent {
  serverID = 10;
  serverStatus = 'offline';

  constructor() {
    this.serverStatus = Math.random() < 0.5 ? 'online' : 'offline';
  }

  getServerStatus() {
    return this.serverStatus;
  }

  getcolor() {
    return this.serverStatus === 'online' ? 'green' : 'red';
  }

}
