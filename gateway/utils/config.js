'use strict';

var _ = require('underscore');
var fs = require('fs');
var path = require('path');

var argv = require('minimist')(process.argv.slice(2));
var configName = 'config.js';
var configDir = path.normalize(path.join(__dirname, '..'));

var configPath, filePath;
if (argv.c && _.isString(argv.c)) {
    filePath = argv.c;
    filePath = path.resolve(process.cwd(), filePath);
    if (fs.existsSync(filePath)) {
        configPath = filePath;
    }
} else {
    while(configDir) {
        //require('console').log(configDir);
        filePath = path.join(configDir, configName);
        if (fs.existsSync(filePath)) {
            configPath = filePath;
            break;
        }
        if (configDir === '/') {
            configDir = null;
        } else {
            configDir = path.dirname(configDir);
        }
    }
}

var config;

if (configPath) {
    config = require(configPath);
    config.$name = path.basename(configPath);
    config.$dir = path.dirname(configPath);
    config.$path = configPath;
}

module.exports = config;
