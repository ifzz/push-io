/**
 * Created by zhibinpan on 4/1/2017.
 */

'use strict';

var Promise = require('bluebird');
var mongoose = require('mongoose');

var Schema = mongoose.Schema;

var MessageSchema = new Schema({
    id: String,
    timestamp: {type: Date, default: Date.now},
}, {});

// To get a message
MessageSchema.index({id: 1}, {unique: true});

module.exports = mongoose.model('Ack', MessageSchema);

Promise.promisifyAll(mongoose.Model);
Promise.promisifyAll(mongoose.Model.prototype);
Promise.promisifyAll(mongoose.Query.prototype);
