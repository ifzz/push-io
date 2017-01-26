/**
 * Created by zhibinpan on 10/1/2017.
 */

module.exports.LOG_DIR = null;

// Default to 127.0.0.1
module.exports.REDIS_HOST = '10.71.2.39';

// Default to 6379
module.exports.REDIS_PORT = 6379;

module.exports.REDIS_OPTIONS = null;

module.exports.ADMIN_CREDENTIALS = [
    {
        name: 'gftrader',
        pass: '1163CFFD87155CD634CBD3DA9F53D'
    },
    {
        name: 'gfnbop',
        pass: '988733B9DDB81E626E4A84D232676'
    }
];

module.exports.STATSD_HOST = '10.71.0.8';

module.exports.STATSD_PORT = 8125;
