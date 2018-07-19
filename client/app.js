var app = new Vue({
  el: "#app",
  data: {
    listMessages: [],
    msg: "",
    ws: null,
    events: {
      subscribe: {
        name: "subscribe",
        payload: ""
      },
      unsubscribe: {
        name: "unsubscribe",
        payload: ""
      }
    },
    rooms: ["room.$general", "room.$default", "room.$private"]
  },
  created() {
    this.wsInit();
  },
  beforeDestroy() {
    this.ws && this.ws.close();
    this.ws = null;
  },
  methods: {
    readMessage() {
      this.ws.onmessage = evt => {
        var messages = evt.data.split("\n");
        this.listMessages.push("message:" + messages[0]);
      };
    },
    closeConnection() {
      this.ws.onclose = evt => {
        this.listMessages.push("--> Connection closed.<--");
      };
    },
    openWs() {
      this.ws.onopen = evt => {
        this.join(this.rooms);
      };
    },
    wsInit() {
      if (!window["WebSocket"]) {
        this.listMessages.push(
          "--> Your browser does not support WebSockets <--"
        );
        return;
      }
      this.listMessages.push("--> Ready to start <--");
      this.ws = new WebSocket("ws://" + document.location.host + "/ws");
      this.openWs();
      this.readMessage();
      this.closeConnection();
    },
    join(rooms) {
      let { subscribe } = this.events;
      subscribe.payload = JSON.stringify(rooms);
      this.ws.send(JSON.stringify(subscribe));
    },
    left(rooms) {},

    submit() {
      if (!this.ws || !this.msg) return;
      this.ws.send(this.msg);
      this.msg = "";
    }
  }
});
