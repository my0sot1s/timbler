interface EventSystem {
  name: string;
  payload: string;
}
interface Message {
  text: string;
  created: number;
  by: string;
  to: string;
  type: string;
}

class Client {
  connectionId: string;
  ws: WebSocket;
  rooms: Array<string>;
  wsInit(path) {
    this.ws = new WebSocket(path);
  }
  get ConnectionId() {
    return this.connectionId;
  }
  readMessage(cb) {
    this.ws.onmessage = evt => {
      var messages = evt.data.split("\n");
      cb(messages[0]);
    };
  }
  openWs() {
    this.ws.onopen = evt => {
      // sub default room
      this.subscribe(this.rooms);
    };
  }
  subscribe(rooms) {
    let subscribe: EventSystem = {
      name: "subscribe",
      payload: JSON.stringify(rooms)
    };
    this.ws.send(JSON.stringify(subscribe));
  }
  unsubscribe(rooms) {
    let unsubscribe: EventSystem = {
      name: "unsubscribe",
      payload: JSON.stringify(rooms)
    };
    this.ws.send(JSON.stringify(unsubscribe));
  }
  dead() {
    this.ws && this.ws.close();
    this.ws = null;
  }
}
