(function() {

    "use strict";

    var $ = require('jquery');
    //var _ = require('lodash');

    window.startApp = function() {

        // Map control
        var map = new L.Map('map', {
            zoomControl: true,
            center: [40.7, -73.9],
            zoom: 11,
            maxZoom: 18
        });

        // Base map
        L.tileLayer('http://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png', {
            attribution: 'CartoDB'
        }).addTo(map);

        // GET request that executes a callback with either an Float32Array
        // containing bin values or null if no data exists
        var getArrayBuffer = function( url, callback ) {
            var xhr = new XMLHttpRequest();
            xhr.open( 'GET', url, true );
            xhr.responseType = 'arraybuffer';
            xhr.onload = function() {
                if ( this.status === 200 ) {
                    callback( new Float64Array( this.response ) );
                } else {
                    callback( null );
                }
            };
            xhr.send();
        };

        // Defines the two color values to interpolate between
        var fromColor = { r: 150, g: 0, b: 0, a: 150 };
        var toColor = { r: 255, g: 255, b: 50, a: 255 };

        // Due to the distribution of values, a logarithmic transform is applied
        // to give a more 'gradual' gradient
        var logTransform = function(value, min, max) {
            var logMin = Math.log( Math.max( 1, min ) );
            var logMax = Math.log( Math.max( 1, max ) );
            var oneOverLogRange = 1 / ( logMax - logMin );
            return Math.log(value - logMin) * oneOverLogRange;
        };

        // Interpolates the color value between the minimum and maximum values provided
        var interpolateColor = function(value, min, max) {
            var alpha = logTransform( value, min, max );
            if ( value === 0 ) {
                return {
                    r: 255,
                    g: 255,
                    b: 255,
                    a: 0
                };
            } else {
                return {
                    r: toColor.r * alpha + fromColor.r * (1 - alpha),
                    g: toColor.g * alpha + fromColor.g * (1 - alpha),
                    b: toColor.b * alpha + fromColor.b * (1 - alpha),
                    a: toColor.a * alpha + fromColor.a * (1 - alpha)
                };
            }
        };

        // Create the canvas tile layer
        var heatmapLayer = new L.tileLayer.canvas();

        // Override 'drawTile' method. Requests the bin data for the tile, and
        // if it exists, renders to the canvas element for the repsecive tile.
        heatmapLayer.drawTile = function(canvas, index, zoom) {
            var url = './tile/'+zoom+'/'+index.x+'/'+index.y;
            $( canvas ).addClass( 'blinking' );
            getArrayBuffer( url, function( bins ) {
                if (!bins) {
                    // Exit early if no data
                    return;
                }
                var ctx = canvas.getContext("2d");
                var imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
                var data = imageData.data;
                var minMax =  {
                    min: 0,
                    max: 1000
                };
                bins.forEach(function(bin,index) {
                    // Interpolate bin value to get rgba
                    var rgba = interpolateColor(bin, minMax.min, minMax.max);
                    data[index*4] = rgba.r;
                    data[index*4+1] = rgba.g;
                    data[index*4+2] = rgba.b;
                    data[index*4+3] = rgba.a;
                });
                // Overwrite original image
                ctx.putImageData(imageData, 0, 0);
                $( canvas ).removeClass( 'blinking' );
            });
        };
        // Add layer to the map
        heatmapLayer.addTo(map);
    };

}());
