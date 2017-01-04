'use strict';

var _ = require('underscore');
var Promise = require('bluebird');

var logger = require('./logger');

exports.applyThreshold = function (innerFunction, concurrentThreshold, tolerantThreshold) {
    var waitingList = [];
    var errorCount = 0;

    function acquire() {
        //logger.info('acquire ' + concurrentThreshold);
        if (concurrentThreshold > 0) {
            concurrentThreshold--;
            return Promise.resolve();
        } else {
            return new Promise(function (resolve) {
                waitingList.push(resolve);
            });
        }
    }

    function release() {
        //logger.info('release ' + concurrentThreshold);
        if (!_.isEmpty(waitingList)) {
            waitingList.splice(0, 1)[0]();
        } else {
            concurrentThreshold++;
        }
    }

    return function () {
        var parameters = _.toArray(arguments);
        return acquire().then(function () {
            return innerFunction.apply(null, parameters);
        }).catch(function (e) {
            logger.error(e);
            errorCount++;
            if (!tolerantThreshold || errorCount >= tolerantThreshold) {
                return Promise.reject('Too many errors...');
            }
        }).finally(function () {
            release();
        });
    };
};

exports.retry = function (method, counter, delay, retryCallback) {
    if (_.isFunction(delay)) {
        // make delay optional
        retryCallback = delay;
        delay = 0;
    }
    return function () {
        var args = arguments;
        var count = counter;
        var invoke = function (error) {
            if (count > 0) {
                if (count < counter) {
                    if (retryCallback) {
                        retryCallback(counter - count);
                    }
                }
                count--;
                return method.apply(null, args).catch(function(error) {
                    if (error) {
                        logger.error(error);
                    }
                    if (delay) {
                        return Promise
                            .delay(delay)
                            .then(function() {
                                return invoke(error);
                            });
                    } else {
                        return invoke(error);
                    }
                });
            } else {
                if (retryCallback) {
                    retryCallback(0);
                }
                return Promise.reject(error);
            }
        };
        return invoke();
    };
};
