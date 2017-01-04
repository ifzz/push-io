/**
 * Created by zhibinpan on 3/1/2017.
 */

'use strict';
var utils = require('../utils');
var config = utils.config;
var logger = utils.getLogger('[API] ');

var express = require('express');
var router = express.Router();
var redis = require('redis');
var bluebird = require('bluebird');
var client = redis.createClient(config.REDIS_PORT, config.REDIS_HOST, config.REDIS_OPTIONS);
var auth = require('basic-auth');
var _ = require('lodash');

const BROKER_NODES = 'io.gf.com.cn:nodes';
const PREFIX_STATS = 'io.gf.com.cn:stats:';
const FIELD_STATS = 'clients';
const FIELD_ADDRESS = 'address';

bluebird.promisifyAll(redis.RedisClient.prototype);
bluebird.promisifyAll(redis.Multi.prototype);

/* GET home page. */
router.get('/server', function(req, res) {
    var credentials = auth(req);
    //logger.debug(JSON.stringify(credentials));
    if (!credentials) {
        res.statusCode = 401;
        res.end('Access denied');
        return;
    }

    var found = _.find(config.ADMIN_CREDENTIALS, function (user) {
        return user.name === credentials.name && user.pass === credentials.pass;
    });
    if (!found) {
        res.statusCode = 401;
        res.end('Access denied');
        return;
    }

    client.smembersAsync(BROKER_NODES)
        .then(function (nodes) {
            var stats = [];
            for (var i = 0; i < nodes.length; i++) {
                stats.push(client.hmgetAsync(PREFIX_STATS + nodes[i], FIELD_ADDRESS, FIELD_STATS));
            }
            return Promise.all(stats);
        })
        .then(function (nodes) {
            var host = null;
            var count = Number.MAX_SAFE_INTEGER;
            for (var i = 0; i < nodes.length; i++) {
                if (nodes[i][1] < count) {
                    count = nodes[i][1];
                    host = nodes[i][0];
                }
            }

            if (host) {
                res.json({'host': host, 'count': count});
            } else {
                res.status(500).send('no connector available');
            }
        })
        .catch(function (err) {
            logger.error(JSON.stringify(err));
            res.status(500).send(JSON.stringify(err));
        });
});

module.exports = router;
