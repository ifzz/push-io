'use strict';

var sysUtil = require('util');
var _ = require('underscore');

exports.formatTable = function (rows) {
    function padString(o, size) {
        var padSize = size - o.length;
        var s = o.value;
        for (var i = 0; i < padSize; i++) {
            s += ' ';
        }
        return s;
    }

    function convertString(s) {
        s = s.toString();
        var matches = s.match(/[\u3400-\u9FBF]/g);
        var actualLength = s.length + (matches ? matches.length : 0);
        return {
            value: s,
            length: actualLength
        };
    }

    if (rows.length) {
        for (var i = 0; i < rows[0].length - 1; i++) {
            var columnSize = 0;
            for (var j = 0; j < rows.length; j++) {
                rows[j][i] = convertString(rows[j][i].toString());
                columnSize = Math.max(columnSize, rows[j][i].length);
            }
            for (j = 0; j < rows.length; j++) {
                rows[j][i] = padString(rows[j][i], columnSize);
            }
        }
        return _.map(rows, function (row) {
            return row.join('  ');
        }).join('\n');
    }
    return '';
};

exports.format = function() {
    return sysUtil.format.apply(null, arguments);
};
