/**
 * Created by zhibinpan on 10/1/2017.
 */

'use strict';

var Ack = require('../model/ack');
var express = require('express');
var router = express.Router();

/* GET users listing. */
router.get('/message/:id', function(req, res, next) {
    var id = req.params.id;



    res.send('respond with a resource');
});

module.exports = router;
