var app = new Vue({
  el: "#app",
  data: {
    listMessages: [],
    msg: "",
    ws: null,
    ID: "",
    events: {
      subscribe: {
        name: "subscribe",
        payload: ""
      },
      unsubscribe: {
        name: "unsubscribe",
        payload: ""
      },
      message: {
        text: "",
        created: Date.now(),
        by: "",
        to: "",
        type: "text"
      }
    },
    rooms: ["room.$default"]
  },
  created() {
    this.wsInit();
  },
  beforeDestroy() {
    this.ws && this.ws.close();
    this.ws = null;
  },
  methods: {
    messageFlow(message) {
      message = JSON.parse(message);
      switch (message.type) {
        case "system":
          this.ID = message.text;
        case "text":
          this.listMessages.push(message);
      }
    },
    readMessage() {
      this.ws.onmessage = evt => {
        var messages = evt.data.split("\n");
        this.messageFlow(messages[0]);
        // this.listMessages.push("message:" + messages[0]);
      };
    },
    closeConnection() {
      this.ws.onclose = evt => {
        this.listMessages.push("--> Connection closed.<--");
      };
    },
    openWs() {
      this.ws.onopen = evt => {
        this.subscribe(this.rooms);
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
    subscribe(rooms) {
      let { subscribe } = this.events;
      subscribe.payload = JSON.stringify(rooms);
      this.ws.send(JSON.stringify(subscribe));
    },
    unsubscribe(rooms) {
      let { unsubscribe } = this.events;
      unsubscribe.payload = JSON.stringify(rooms);
      this.ws.send(JSON.stringify(unsubscribe));
    },

    submit() {
      if (!this.ws || !this.msg) return;
      let { message } = this.events;
      this.ws.send(
        JSON.stringify({
          name: "message",
          payload: JSON.stringify({
            ...message,
            text: this.msg,
            to: "room.$pim",
            by: this.ID
          })
        })
      );
      this.msg = "";
    },
    join(name) {
      this.subscribe([name]);
    },
    left(name) {
      this.unsubscribe([name]);
    }
  }
});

