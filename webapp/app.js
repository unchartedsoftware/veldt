(function() {

    'use strict';

    var $ = require('jquery');
    var _ = require('lodash');
    //var d3 = require('d3');
    var TileRequester = require('./scripts/TileRequester');
    var WordCloudLayer = require('./scripts/WordCloudLayer');
    var AJAX_TIMEOUT = 5000;

    $.ajaxTransport('+arraybuffer', function(options) {
        var xhr;
        return {
            send: function(headers, callback) {
                // setup all variables
                var url = options.url;
                var type = options.type;
                var async = options.async || true;
                var dataType = 'arraybuffer';
                var data = options.data || null;
                var username = options.username || null;
                var password = options.password || null;
                // create new XMLHttpRequest
                xhr = new XMLHttpRequest();
                xhr.addEventListener('load', function() {
                    var data = {};
                    data[options.dataType] = xhr.response;
                    // make callback and send data
                    callback(xhr.status, xhr.statusText, data, xhr.getAllResponseHeaders());
                });
                xhr.open(type, url, async, username, password);
                // setup custom headers
                _.forIn(headers, function(header, key) {
                    xhr.setRequestHeader(key, header);
                });
                xhr.responseType = dataType;
                xhr.send(data);
            },
            abort: function() {
                xhr.abort();
            }
        };
    });

    window.startApp = function() {

        // Map control
        var map = new L.Map('map', {
            zoomControl: true,
            center: [40.7, -73.9],
            zoom: 14,
            maxZoom: 18,
            fadeAnimation: true,
            zoomAnimation: true
        });

        // Base map
        var baseLayer = L.tileLayer('http://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png', {
            attribution: 'CartoDB'
        });
        baseLayer.addTo(map);

        // Defines the two color values to interpolate between
        var fromColor = {
            r: 150,
            g: 0,
            b: 0,
            a: 150
        };
        var toColor = {
            r: 255,
            g: 255,
            b: 50,
            a: 255
        };

        // Due to the distribution of values, a logarithmic transform is applied
        // to give a more 'gradual' gradient
        var logTransform = function(value, min, max) {
            var logMin = Math.log(Math.max(1, min));
            var logMax = Math.log(Math.max(1, max));
            var oneOverLogRange = 1 / (logMax - logMin);
            return Math.log(value - logMin) * oneOverLogRange;
        };

        // Interpolates the color value between the minimum and maximum values provided
        var interpolateColor = function(value, min, max) {
            var alpha = logTransform(value, min, max);
            if (value === 0) {
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
        var heatmapLayer = new L.tileLayer.canvas({
            unloadInvisibleTiles: true
        });
        // Override 'drawTile' method. Requests the bin data for the tile, and
        // if it exists, renders to the canvas element for the repsecive tile.
        heatmapLayer.drawTile = function(canvas, index, zoom) {
            $(canvas).addClass('blinking');
            requester
                .get('heatmap', index.x, index.y, zoom)
                .done(function() {
                    var url = './heatmap/' + zoom + '/' + index.x + '/' + index.y;
                    $.ajax({
                        url: url,
                        dataType: 'arraybuffer',
                        timeout: AJAX_TIMEOUT
                    }).done(function(buffer) {
                        var bins = new Float64Array(buffer);
                        var ctx = canvas.getContext('2d');
                        var imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
                        var data = imageData.data;
                        var minMax = {
                            min: 0,
                            max: 1000
                        };
                        bins.forEach(function(bin, index) {
                            // Interpolate bin value to get rgba
                            var rgba = interpolateColor(bin, minMax.min, minMax.max);
                            data[index * 4] = rgba.r;
                            data[index * 4 + 1] = rgba.g;
                            data[index * 4 + 2] = rgba.b;
                            data[index * 4 + 3] = rgba.a;
                        });
                        // Overwrite original image
                        ctx.putImageData(imageData, 0, 0);
                    }).always(function() {
                        $(canvas).removeClass('blinking');
                    });
                })
                .fail(function() {
                    console.error('Could not generate tile.');
                });
        };

        // Create the canvas tile layer
        var wordCloudLayer = new WordCloudLayer(null, {
            unloadInvisibleTiles: true
        });
        wordCloudLayer.drawTile = function(tile, index, zoom) {
            $(tile).addClass('blinking');
            requester
                .get('topiccount', index.x, index.y, zoom)
                .done(function() {
                    var url = './topiccount/' + zoom + '/' + index.x + '/' + index.y;
                    $.ajax({
                        url: url,
                        dataType: 'json',
                        timeout: AJAX_TIMEOUT
                    }).done(function(wordCounts) {
                        wordCloudLayer.render(
                            $(tile),
                            wordCounts,
                            {
                                min: 0,
                                max: 1000
                            },
                            wordCloudLayer.highlight);
                        wordCloudLayer.tileDrawn(tile);
                    }).always(function() {
                        $(tile).removeClass('blinking');
                    });
                })
                .fail(function() {
                    console.error('Could not generate tile.');
                });
        };

        // // Create the canvas tile layer
        // var wordCloudLayer = new L.tileLayer.canvas();
        // // Override 'drawTile' method. Requests the bin data for the tile, and
        // // if it exists, renders to the canvas element for the repsective tile.
        // var layout = function(words, renderInfo, width, height) {
        //     var that = this;
        //     var startX = width / 2;
        //     var startY = height / 2;
        //     var maxRadius = Math.min(width, height) * 0.5;
        //     words.forEach(function(word) {
        //         var placed = false;
        //         for (var i = 10; i > 0 && !placed; i--) {
        //             var minRadius = maxRadius * (1 / i);
        //             for (var t = 0; t <= 1 && !placed; t += (1 / 50)) {
        //                 var searchRadius = minRadius + (maxRadius - minRadius) * t;
        //                 var theta = Math.random() * Math.PI * 2;
        //                 var r = Math.random() * searchRadius;
        //                 // Apply a shift no more than the width of word to avoid bunching on the right
        //                 var shiftX = -(Math.random() * renderInfo[word].bb.width);
        //                 var shiftY = -(Math.random() * renderInfo[word].bb.height);
        //                 renderInfo[word].x = startX + r * Math.sin(theta) + shiftX;
        //                 renderInfo[word].y = startY + r * Math.cos(theta) + shiftY;
        //                 if (that.fits(renderInfo[word])) {
        //                     placed = true;
        //                     that.place(word, renderInfo[word]);
        //                 }
        //             }
        //         }
        //         if (!placed) {
        //             renderInfo[word].x = -1;
        //             renderInfo[word].y = -1;
        //         }
        //     });
        // };
        //
        // var colorRamp = d3.scale.sqrt()
        //     .range(['#c77c11', '#f8d9ac'])
        //     .domain([0, 1]);
        //
        // wordCloudLayer.drawTile = function(canvas, index, zoom) {
        //     $(canvas).addClass('blinking');
        //     requester
        //         .get('topiccount', index.x, index.y, zoom)
        //         .done(function() {
        //             var url = './topiccount/' + zoom + '/' + index.x + '/' + index.y;
        //             $.getJSON(url, function(wordCounts) {
        //                 if (_.isEmpty(wordCounts)) {
        //                     // Exit early if no data
        //                     $(canvas).removeClass('blinking');
        //                     return;
        //                 }
        //                 new Cloud5()
        //                     .canvas(canvas)
        //                     .width(256)
        //                     .height(256)
        //                     .font('Helvetica')
        //                     .minFontSize(10)
        //                     .maxFontSize(30)
        //                     .wordMap(wordCounts)
        //                     .color(function(info) {
        //                         return colorRamp((info.count - info.minCount) / (info.maxCount - info.minCount));
        //                     })
        //                     .onWordOver(function() {
        //                         // handle mousing over a word
        //                         $(canvas).css('cursor', 'pointer');
        //                     //cloud.highlight( word, 'ff0000' );
        //                     })
        //                     .onWordOut(function() {
        //                         // handle mousing out of a word
        //                         $(canvas).css('cursor', 'auto');
        //                     //cloud.unhighlight( word );
        //                     })
        //                     .onWordClick(function(word) {
        //                         // handle clicking a word
        //                         console.log(word);
        //                     })
        //                     .layout(layout)
        //                     .generate();
        //                 $(canvas).removeClass('blinking');
        //             });
        //         })
        //         .fail(function() {
        //             console.error('Could not generate tile.');
        //         });
        // };

        var requester = new TileRequester('ws://localhost:8080/batch', function() {
            // Add layer to the map
            heatmapLayer.addTo(map);
            // Add layer to the map
            wordCloudLayer.addTo(map);
        });
    };

}());
