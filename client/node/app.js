/**
 * Created by zhibinpan on 20/12/2016.
 */

var mqtt = require('mqtt');
var options = {
    'clientId': 'jacky',
    'clean': false,
    'username': 'gftrader',
    'password': 'A98D8B1134D34F6E161463F757139'
};
var client = mqtt.connect('mqtt://127.0.0.1:1883', options);

client.on('connect', function () {
    console.log('connect');
    //client.subscribe('test/topic');
});

client.on('message', function (topic, message) {
    // message is Buffer
    console.log(message.toString());
    //client.end();
});

client.subscribe('$SYS/brokers/emqttd@10.71.128.85/stats/', {'qos': 2}, function (err, granted) {
    if (err) {
        console.log(JSON.stringify(err));
    }
    console.log(JSON.stringify(granted));
});

//client.publish('presence', 'Hello mqtt');

console.log('done');
