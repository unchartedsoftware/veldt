(function() {

    'use strict';

    var $ = require('jquery');
    var _ = require('lodash');
    var moment = require('moment');
    //var d3 = require('d3');
    var TileRequester = require('./scripts/TileRequester');
    var WordCloudLayer = require('./scripts/WordCloudLayer');
    var TopicFrequencyLayer = require('./scripts/TopicFrequencyLayer');

    var LayerMenu = require('./scripts/LayerMenu');
    var Slider = require('./scripts/Slider');
    var RangeSlider = require('./scripts/RangeSlider');
    var AJAX_TIMEOUT = 5000;

    function mod(m, n) {
        return ((m % n) + n) % n;
    }

    function createOpacitySlider(layer) {

        layer.getOpacity = function() {
            return this.options.opacity;
        };

        return new Slider({
            label: 'Opacity',
            min: 0,
            max: 1,
            step: 0.01,
            initialValue: layer.getOpacity(),
            change: function(event) {
                layer.setOpacity(event.newValue);
            }
        });
    }

    function createResolutionSlider(layer) {
        return new Slider({
            label: 'Resolution',
            min: 0,
            max: 8,
            step: 1,
            initialValue: Math.log2(layer.resolution),
            formatter: function(value) {
                return Math.pow(2, value);
            },
            slideStop: function(value) {
                layer.resolution = Math.pow(2, value);
                layer.redraw();
            }
        });
    }

    function createTimeSlider(layer) {
        return new RangeSlider({
            label: 'Time Range',
            min: layer.meta.timestamp.extrema.min,
            max: layer.meta.timestamp.extrema.max,
            step: 86400000,
            initialValue: [layer.timeRange.from, layer.timeRange.to],
            formatter: function(values) {
                var from = moment.unix(values[0] / 1000).format('MMM D, YYYY');
                var to = moment.unix(values[1] / 1000).format('MMM D, YYYY');
                return from + ' to ' + to;
            },
            slideStop: function(values) {
                layer.timeRange.from = values[0];
                layer.timeRange.to = values[1];
                layer.redraw();
            }
        });
    }

    function getControlByType(type, layer) {
        switch (type) {
            case 'opacity':
                return createOpacitySlider(layer);
            case 'resolution':
                return createResolutionSlider(layer);
            case 'time':
                return createTimeSlider(layer);
        }
    }

    function createLayerControls(layerName, layer, controls) {
        // add visibility controls to the layer
        layer.show = function() {
            this._hidden = false;
            this._prevMap.addLayer(this);
        };
        layer.hide = function() {
            this._hidden = true;
            this._prevMap = this._map;
            this._map.removeLayer(this);
        };
        layer.isHidden = function() {
            return !!this._hidden;
        };
        var layerMenu = new LayerMenu({
            layer: layer,
            label: layerName
        });
        controls.forEach(function(control) {
            layerMenu.getBody()
                .append(getControlByType(control, layer).getElement())
                .append('<div style="clear:both;"></div>');
        });
        return layerMenu.getElement();
    }

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

        var ES_ENDPOINT = 'openstack';
        var ES_INDEX = 'isil_twitter_weekly';
        var TOPICS = [
            'isil',
            'isis',
            'terror',
            'bomb',
            'death',
            'kill',
            'cool',
            'awesome',
            'amazing',
            'badass',
            'killer',
            'dank',
            'superfly',
            'smooth',
            'radical',
            'wicked',
            'neato',
            'nifty',
            'primo',
            'gnarly',
            'crazy',
            'insane',
            'sick',
            'mint',
            'nice',
            'nasty',
            'classic',
            'tight',
            'rancid',
        ];

        var meta = null;

        // Map control
        var map = new L.Map('map', {
            zoomControl: true,
            zoom: 4,
            center: [0, 0],
            // zoom: 14,
            // center: [40.7, -73.9],
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

            var params = {
                resolution: heatmapLayer.resolution,
                from: heatmapLayer.timeRange.from,
                to: heatmapLayer.timeRange.to
            };

            index.x = mod(index.x, Math.pow(2, zoom));
            index.y = mod(index.y, Math.pow(2, zoom));

            requester
                .get(ES_ENDPOINT, ES_INDEX, 'heatmap', index.x, index.y, zoom, params)
                .done(function() {
                    var url = ES_ENDPOINT + '/' + ES_INDEX + '/heatmap/' + zoom + '/' + index.x + '/' + index.y;
                    $.ajax({
                        url: url + '?' + $.param(params),
                        dataType: 'arraybuffer',
                        timeout: AJAX_TIMEOUT
                    }).done(function(buffer) {
                        var bins = new Float64Array(buffer);

                        var tileCanvas = document.createElement('canvas');
                        tileCanvas.height = params.resolution;
                        tileCanvas.width = params.resolution;
                        var tileCtx = tileCanvas.getContext('2d');

                        var imageData = tileCtx.getImageData(0, 0, params.resolution, params.resolution);
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
                        tileCtx.putImageData(imageData, 0, 0);

                        var ctx = canvas.getContext('2d');
                        ctx.imageSmoothingEnabled = false;
                        ctx.drawImage(
                            tileCanvas,
                            0, 0,
                            params.resolution, params.resolution,
                            0, 0,
                            canvas.width, canvas.height);

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

            var params = {
                topics: TOPICS.join(','),
                from: wordCloudLayer.timeRange.from,
                to: wordCloudLayer.timeRange.to
            };

            index.x = mod(index.x, Math.pow(2, zoom));
            index.y = mod(index.y, Math.pow(2, zoom));

            requester
                .get(ES_ENDPOINT, ES_INDEX, 'topiccount', index.x, index.y, zoom, params)
                .done(function() {
                    var url = ES_ENDPOINT + '/' + ES_INDEX + '/topiccount/' + zoom + '/' + index.x + '/' + index.y;
                    $.ajax({
                        url: url + '?' + $.param(params),
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

        // Create the canvas tile layer
        var topicFrequencyLayer = new TopicFrequencyLayer(null, {
            unloadInvisibleTiles: true
        });
        topicFrequencyLayer.drawTile = function(tile, index, zoom) {
            $(tile).addClass('blinking');

            var params = {
                topics: TOPICS.join(','),
                from: topicFrequencyLayer.timeRange.from,
                to: topicFrequencyLayer.timeRange.to,
                interval: 'month'
            };

            index.x = mod(index.x, Math.pow(2, zoom));
            index.y = mod(index.y, Math.pow(2, zoom));

            requester
                .get(ES_ENDPOINT, ES_INDEX, 'topicfrequency', index.x, index.y, zoom, params)
                .done(function() {
                    var url = ES_ENDPOINT + '/' + ES_INDEX + '/topicfrequency/' + zoom + '/' + index.x + '/' + index.y;
                    $.ajax({
                        url: url + '?' + $.param(params),
                        dataType: 'json',
                        timeout: AJAX_TIMEOUT
                    }).done(function(topicCounts) {
                        topicFrequencyLayer.render(
                            $(tile),
                            topicCounts,
                            {
                                min: 0,
                                max: 1000
                            },
                            topicFrequencyLayer.highlight);
                        topicFrequencyLayer.tileDrawn(tile);
                    }).always(function() {
                        $(tile).removeClass('blinking');
                    });
                })
                .fail(function() {
                    console.error('Could not generate tile.');
                });
        };

        var requester = new TileRequester('batch', function() {

            $.ajax({
                url: ES_ENDPOINT + '/' + ES_INDEX + '/default'
            }).done(function(res) {

                meta = res;

                heatmapLayer.meta = meta;
                heatmapLayer.resolution = 32;
                heatmapLayer.timeRange = {
                    from: meta.timestamp.extrema.min,
                    to: meta.timestamp.extrema.max
                };

                wordCloudLayer.meta = meta;
                wordCloudLayer.timeRange = {
                    from: meta.timestamp.extrema.min,
                    to: meta.timestamp.extrema.max
                };

                topicFrequencyLayer.meta = meta;
                topicFrequencyLayer.timeRange = {
                    from: meta.timestamp.extrema.min,
                    to: meta.timestamp.extrema.max
                };

                //$('.controls').append(createLayerControls('Heatmap', heatmapLayer, ['opacity', 'resolution', 'time']));
                //$('.controls').append(createLayerControls('Word Cloud', wordCloudLayer, ['opacity', 'time']));
                $('.controls').append(createLayerControls('Topic Trends', topicFrequencyLayer, ['opacity', 'time']));

                // Add layer to the map
                //heatmapLayer.addTo(map);
                // Add layer to the map
                //wordCloudLayer.addTo(map);
                // Add layer to the map
                topicFrequencyLayer.addTo(map);
            });

        });
    };

}());
