/**
 * Created by zhibinpan on 3/1/2017.
 */

'use strict';

var express = require('express');
var router = express.Router();
var redis = require('redis');
var bluebird = require('bluebird');
var client = redis.createClient();
var auth = require('basic-auth');

const BROKER_NODES = 'io.gf.com.cn:nodes';
const PREFIX_STATS = 'io.gf.com.cn:stats:';
const FIELD_STATS = 'clients';
const FIELD_ADDRESS = 'address';

bluebird.promisifyAll(redis.RedisClient.prototype);
bluebird.promisifyAll(redis.Multi.prototype);

/* GET home page. */
router.get('/server', function(req, res) {
    var credentials = auth(req);
    console.log(credentials);
    if (!credentials || credentials.name !== 'gftrader' || credentials.pass !== 'A98D8B1134D34F6E161463F757139') {
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
            var count = 65536;
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
            console.log(err);
            res.status(500).send(JSON.stringify(err));
        });
});

module.exports = router;
