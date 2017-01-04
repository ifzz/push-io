'use strict';

var crypto = require('crypto');

exports.hash = function (content, algorithm) {
    var hash = crypto.createHash(algorithm || 'sha256');
    hash.update(content, 'utf8');
    return hash.digest('hex');
};
