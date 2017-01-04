'use strict';

var path = require('path');
var fs = require('fs');
var _ = require('underscore');
var Promise = require('bluebird');
var readline = require('readline');

var config = require('./config');

exports.fileExists = function(filePath) {
    try {
        fs.statSync(filePath);
        return true;
    }
    catch (e) {
        return false;
    }
};

exports.getAppDir = function () {
    return config.$dir;
};

exports.ensureDir = function (dir) {
    if (!fs.existsSync(dir)) {
        var parentDir = path.dirname(path.resolve(dir));
        exports.ensureDir(parentDir);
        fs.mkdirSync(dir);
    }
};

exports.clearDir = function (dir) {
    _.each(fs.readdirSync(dir), function (file) {
        var filePath = path.join(dir, file);
        if (fs.statSync(filePath).isFile()) {
            fs.unlinkSync(filePath);
        }
    });
};

var tempUploadDir;

exports.getTempUploadDir = function () {
    if (!tempUploadDir) {
        tempUploadDir = config.TEMP_UPLOAD_DIR || path.join(exports.getAppDir(), 'upload');
        exports.ensureDir(tempUploadDir);
    }
    return tempUploadDir;
};

exports.question = function (question, choices) {
    return new Promise(function (resolve) {
        var rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout
        });

        var text = question;
        if (choices) {
            text += ' (' + choices.join('/') + ')';
        }
        var lChoices = _.map(choices, function (choice) {
            return choice.toLowerCase();
        });

        function handleAnswer(answer) {
            if (choices) {
                var index = _.indexOf(lChoices, answer.toLowerCase());
                if (index >= 0) {
                    rl.close();
                    resolve(choices[index]);
                } else {
                    rl.question(text, handleAnswer);
                }
            } else {
                rl.close();
                resolve(answer);
            }
        }

        rl.question(text, handleAnswer);
    });
};

exports.readJSONFile = function (fileName) {
    return new Promise(function (resolve, reject) {
        fs.readFile(fileName, {
            encoding: 'utf8'
        }, function (err, data) {
            if (err) {
                reject(err);
            } else {
                try {
                    resolve(JSON.parse(data));
                } catch (e) {
                    reject(e);
                }
            }
        });
    });
};

exports.writeJSONFile = function (fileName, data) {
    return new Promise(function (resolve, reject) {
        var text = JSON.stringify(data, null, 2);
        fs.writeFile(fileName, text, function (err) {
            if (err) {
                reject(err);
            } else {
                resolve();
            }
        });
    });
};

exports.readJSONStream = function (stream) {
    return new Promise(function (resolve, reject) {
        var content = '';
        stream.on('data', function (chunk) {
            content += chunk;
        });
        stream.on('end', function () {
            resolve(JSON.parse(content));
        });
        stream.on('error', function (error) {
            reject(error);
        });
        stream.setEncoding('utf-8');
        stream.resume();
    });
};

exports.createStreamPromise = function (stream) {
    var completed;
    return new Promise(function (resolve, reject) {
        stream.on('error', function (error) {
            reject(error);
        });
        stream.on('end', function () {
            if (!completed) {
                resolve();
                completed = true;
            }
        });
        stream.on('finish', function () {
            if (!completed) {
                resolve();
                completed = true;
            }
        });
        stream.on('close', function () {
            if (!completed) {
                resolve();
                completed = true;
            }
        });
    });
};

