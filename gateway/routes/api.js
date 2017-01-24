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
const FIELD_CLIENTS = 'clients';
const FIELD_ADDRESS = 'address';
const FIELD_STATUS = 'status';
const FIELD_TOTAL_MEMORY = 'total_memory';
const FIELD_USED_MEMORY = 'used_memory';
const FIELD_LOAD_1 = 'load1';
const FIELD_LOAD_5 = 'load5';
const FIELD_LOAD_15 = 'load15';

bluebird.promisifyAll(redis.RedisClient.prototype);
bluebird.promisifyAll(redis.Multi.prototype);

/* GET home page. */
// query which server is available
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
                stats.push(client.hmgetAsync(PREFIX_STATS + nodes[i], FIELD_ADDRESS, FIELD_CLIENTS));
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

router.get('/nodes', function (req, res) {
    client.smembersAsync(BROKER_NODES)
        .then(function (nodes) {
            var stats = [];
            for (var i = 0; i < nodes.length; i++) {
                stats.push(client.hmgetAsync(PREFIX_STATS + nodes[i], FIELD_ADDRESS, FIELD_CLIENTS, FIELD_STATUS,
                    FIELD_TOTAL_MEMORY, FIELD_USED_MEMORY, FIELD_LOAD_1, FIELD_LOAD_5, FIELD_LOAD_15));
            }
            return Promise.all(stats);
        })
        .then(function (nodes) {
            var result = [];
            _.forEach(nodes, function (node) {
                result.push({
                    'address': node[0],
                    'clients': node[1],
                    'status': node[2],
                    'total_memory': node[3],
                    'used_memory': node[4],
                    'load1': node[5],
                    'load5': node[6],
                    'load6': node[7]
                });
            });
            res.json({nodes: result});
        })
        .catch(function (err) {
            logger.error(JSON.stringify(err));
            res.status(500).send(JSON.stringify(err));
        });
});

module.exports = router;
