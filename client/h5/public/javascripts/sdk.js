/**
 * Created by zhibinpan on 21/12/2016.
 */

'use strict';
var gfWebsocket = (function () {
    var clientId = localStorage.getItem(MQTT_CLIENT_ID_KEY);
    if (!clientId) {
        clientId = guid();
        localStorage.setItem(MQTT_CLIENT_ID_KEY, clientId);
    }

    var options = {
        'clientId': clientId,
        'clean': false,
        'username': PUSH_APP_ID,
        'password': PUSH_APP_KEY
    };

    function guid() {
        return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
            s4() + '-' + s4() + s4() + s4();
    }

    function s4() {
        return Math.floor((1 + Math.random()) * 0x10000)
            .toString(16)
            .substring(1);
    }

    var _listener = null;

    var client = mqtt.connect(MQTT_BROKER_URL, options);

    client.on('message', function (topic, payload) {
        console.log([topic, payload].join(': '));
        if (_listener) {
            _listener([payload].join(''));
        }
    });

    client.on('connect', function (connack) {
        console.log('connect: ' + JSON.stringify(connack));
        //client.publish(topic, null, {retain: true});
    });

    client.on('reconnect', function () {
        console.log('reconnect');
    });

    client.on('close', function () {
        console.log('close');
    });

    client.on('offline', function () {
        console.log('offline');
    });

    client.on('error', function (error) {
        console.log('error: ' + JSON.stringify(error));
    });

    client.on('packetreceive', function (packet) {
        //console.log('packetreceive: ' + JSON.stringify(packet));
    });

    function Ws() {
        this.ver = '2.0.0';
    }

    Ws.prototype.watch = function (appId, regId, fn) {
        if (fn && typeof fn === 'function') {
            _listener = fn;
        }

        var topic = localStorage.getItem(MQTT_SUB_TOPIC_KEY);
        if (topic) {
            if (topic === appId + '/' + regId) {
                return;
            }
        }
        topic = appId + '/' + regId;
        localStorage.setItem(MQTT_SUB_TOPIC_KEY, topic);
        client.subscribe(topic, {'qos': 2}, function (err, granted) {
            if (err) {
                console.log(JSON.stringify(err));
            }
            console.log(JSON.stringify(granted));
        });
    };

    var instance = new Ws();

    return instance;
}());
console.log(gfWebsocket.ver);
