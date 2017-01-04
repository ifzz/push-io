'use strict';

var _ = require('underscore');
var winston = require('winston');

var config = require('./config');
var io = require('./io');

var logger;
exports.log = function() { logger.info.apply(logger, arguments); };
exports.debug = function() { logger.debug.apply(logger, arguments); };
exports.info = function() { logger.info.apply(logger, arguments); };
exports.warn = function() { logger.warn.apply(logger, arguments); };
exports.error = function() { logger.error.apply(logger, arguments); };

exports.getLogger = function(prefix) {
    return {
        disabled: false,
        log: function () {
            if (this.disabled) { return; }
            logger.info.apply(logger, addPrefix(arguments));
        },
        debug: function () {
            if (this.disabled) { return; }
            logger.debug.apply(logger, addPrefix(arguments));
        },
        info: function () {
            if (this.disabled) { return; }
            logger.info.apply(logger, addPrefix(arguments));
        },
        warn: function () {
            if (this.disabled) { return; }
            logger.warn.apply(logger, addPrefix(arguments));
        },
        error: function () {
            if (this.disabled) { return; }
            logger.error.apply(logger, addPrefix(arguments));
        }
    };
    function addPrefix(args) {
        var first = _.first(args);
        if (first) {
            return [prefix + first].concat(_.rest(args));
        } else {
            return args;
        }
    }
};

(function() {
    var logDir = config.LOG_DIR;
    if (!logDir) {
        logDir = io.getAppDir() + '/log';
    }
    io.ensureDir(logDir);
    var appLogFile = logDir + '/app.log';
    var errorLogFile = logDir + '/error.log';
    var fileSize = 1024 * 1024 * 10;
    logger = new (winston.Logger)({
        transports: [
            new (winston.transports.Console)({level: 'debug', handleExceptions: true, debugStdout: true})
        ]
    });
})();
