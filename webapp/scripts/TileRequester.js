(function() {

    'use strict';

    var $ = require('jquery');

    function getTileHash(endpoint, index, type, x, y, z) {
        return endpoint + ':' + index + ':' + type + ':' + x + '-' + y + '-' + z;
    }

    function TileRequester(url, callback) {
        var requests = this.requests = {};
        this.socket = new WebSocket(url);
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
                tileCoord.z);
            var request = requests[tileHash];
            delete requests[tileHash];
            if (tileRes.success) {
                request.resolve(tileRes);
            } else {
                request.reject(tileRes);
            }
        };
    }

    TileRequester.prototype.get = function(endpoint, index, type, x, y, z) {
        var tileHash = getTileHash(endpoint, index, type, x, y, z);
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
            type: type
        };
        this.socket.send(JSON.stringify(msg));
        return request.promise();
    };

    module.exports = TileRequester;

}());
