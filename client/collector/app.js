/**
 * Created by zhibinpan on 20/12/2016.
 */

'use strict';

var Ack = require('./model/ack');
var utils = require('./utils');
var config = utils.config;
var _ = require('underscore');
var mqtt = require('mqtt');
var options = {
    'clientId': 'monitor',
    'clean': false,
    'username': 'monitor',
    'password': '7C7DC73CDFAB3838C5E2CE82E1BFC'
};
var client = mqtt.connect('ws://54.223.124.84:80/mqtt', options);//mqtt.connect('mqtt://54.223.124.84:1883', options);
var mongoose = require('mongoose');

mongoose.connection.on('connected', function () {
    utils.log('Mongoose default connection connected.');
});

// If the connection throws an error
mongoose.connection.on('error',function (err) {
    utils.error('Mongoose default connection error: ', err);
});

// When the connection is disconnected
mongoose.connection.on('disconnected', function () {
    utils.error('Mongoose default connection disconnected.');
});

var options = _.clone(config.DB_CONNECTION_OPTS);
//options.server = { reconnectTries: Number.MAX_VALUE, socketOptions: { connectTimeoutMS: 1000 } };
mongoose.connect(config.MONGO_DB_CONNECTION_STRING, options, function (error) {
    if (error) {
        utils.error(error);
    }
});

client.on('connect', function () {
    console.log('connect');
});

client.on('message', function (topic, message) {
    // message is Buffer
    //console.log([topic, message].join(": "));

    Ack.findOneAsync({id: message})
        .then(function (entity) {
            if (!entity) {
                var ack = new Ack({id: message});
                return ack.saveAsync();
            }
        })
        .catch(function (err) {
            utils.error(JSON.stringify(err));
        });
});

client.subscribe('ack/#', {'qos': 2}, function (err, granted) {
    if (err) {
        console.log(JSON.stringify(err));
    }
    console.log(JSON.stringify(granted));
});

console.log('running');
