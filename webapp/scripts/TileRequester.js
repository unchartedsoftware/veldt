(function() {

    'use strict';

    var $ = require('jquery');
    var _ = require('lodash');

    function getTileHash(endpoint, index, type, x, y, z, params) {
        var ps = _.map(params, function(val, key) {
            return key.toLowerCase() + '=' + val;
        });
        ps.sort(function(a, b) {
            return (a < b) ? -1 : (a > b) ? 1 : 0;
        });
        return endpoint + '-' + index + '-' + type + '-' + x + '-' + y + '-' + z + '-' + ps.join('-');
    }

    function getHost() {
        var loc = window.location;
        var new_uri;
        if (loc.protocol === 'https:') {
            new_uri = 'wss:';
        } else {
            new_uri = 'ws:';
        }
        return new_uri + '//' + loc.host + loc.pathname;
    }

    function TileRequester(url, callback) {
        var that = this;
        this.requests = {};
        this.socket = new WebSocket(getHost() + url);
        this.socket.onopen = callback;
        this.socket.onmessage = function(event) {
            var tileRes = JSON.parse(event.data);
            var tileCoord = tileRes.tilecoord;
            var tileHash = getTileHash(
                tileRes.endpoint,
                tileRes.index,
                tileRes.type,
                tileCoord.x,
                tileCoord.y,
                tileCoord.z,
                tileRes.params);
            var request = that.requests[tileHash];
            delete that.requests[tileHash];
            if (tileRes.success) {
                request.resolve(tileRes);
            } else {
                request.reject(tileRes);
            }
        };
        this.socket.onclose = function() {
            console.warn('Websocket connection closed.');
        };
    }

    TileRequester.prototype.get = function(endpoint, index, type, x, y, z, params) {
        var tileHash = getTileHash(endpoint, index, type, x, y, z, params);
        var request = this.requests[tileHash];
        if (request) {
            return request.promise();
        }
        request = this.requests[tileHash] = $.Deferred();
        var msg = {
            tilecoord: {
                x: x,
                y: y,
                z: z
            },
            endpoint: endpoint,
            index: index,
            type: type,
            params: params
        };
        this.socket.send(JSON.stringify(msg));
        return request.promise();
    };

    module.exports = TileRequester;

}());
