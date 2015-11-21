(function() {

    'use strict';

    var $ = require('jquery');

    function getTileHash( type, x, y, z ) {
        return type + '-' + x + '-' + y + '-' + z;
    }

    function TileRequester( url, callback ) {
        var that = this;
        this.url = url;
        this.socket = new WebSocket( url );
        this.socket.onopen = callback;
        this.socket.onmessage = function( event ) {
            var tileRes = JSON.parse( event.data );
            var tileCoord = tileRes.tilecoord;
            var tileHash = getTileHash( tileRes.type, tileCoord.x, tileCoord.y, tileCoord.z );
            var request = that.requests[ tileHash ];
            delete that.requests[ tileHash ];
            if ( !request ) {
                console.log( 'Rejecting: ' + JSON.stringify( tileRes ) );
            } else {
                console.log( 'Resolving: ' + JSON.stringify( tileRes ) );
            }
            if ( tileRes.success ) {
                request.resolve( tileRes );
            } else {
                request.reject( tileRes );
            }
        };
        this.requests = {};
    }

    TileRequester.prototype.get = function( type, x, y, z ) {
        var tileHash = getTileHash( type, x, y, z );
        var request = this.requests[ tileHash ];
        if ( request ) {
            return request;
        }
        request = this.requests[ tileHash ] = $.Deferred();
        this.socket.send( JSON.stringify({
            tilecoord: {
                x: x,
                y: y,
                z: z
            },
            type: type
        }));
        return request.promise();
    };

    module.exports = TileRequester;

}());
