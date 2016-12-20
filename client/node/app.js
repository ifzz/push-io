/**
 * Created by zhibinpan on 20/12/2016.
 */

var mqtt = require('mqtt');
var client = mqtt.connect('mqtt://54.222.243.29:1883');

client.on('connect', function () {
    console.log('connect');
    client.subscribe('test/topic');
});

client.on('message', function (topic, message) {
    // message is Buffer
    console.log(message.toString());
    //client.end();
});

client.publish('presence', 'Hello mqtt');

console.log('done');
