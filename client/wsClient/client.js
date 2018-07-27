var Client = /** @class */ (function () {
    function Client() {
    }
    Client.prototype.wsInit = function (path) {
        this.ws = new WebSocket(path);
    };
    Object.defineProperty(Client.prototype, "ConnectionId", {
        get: function () {
            return this.connectionId;
        },
        enumerable: true,
        configurable: true
    });
    Client.prototype.readMessage = function (cb) {
        this.ws.onmessage = function (evt) {
            var messages = evt.data.split("\n");
            cb(messages[0]);
        };
    };
    Client.prototype.openWs = function () {
        var _this = this;
        this.ws.onopen = function (evt) {
            // sub default room
            _this.subscribe(_this.rooms);
        };
    };
    Client.prototype.subscribe = function (rooms) {
        var subscribe = {
            name: "subscribe",
            payload: JSON.stringify(rooms)
        };
        this.ws.send(JSON.stringify(subscribe));
    };
    Client.prototype.unsubscribe = function (rooms) {
        var unsubscribe = {
            name: "unsubscribe",
            payload: JSON.stringify(rooms)
        };
        this.ws.send(JSON.stringify(unsubscribe));
    };
    Client.prototype.dead = function () {
        this.ws && this.ws.close();
        this.ws = null;
    };
    return Client;
}());
