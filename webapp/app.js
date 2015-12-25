(function() {

    'use strict';

    var $ = require('jquery');
    var _ = require('lodash');
    var moment = require('moment');

    var TileRequester = require('./scripts/TileRequester');

    var TileLayer = require('./scripts/layer/TileLayer');
    var HeatmapTileLayer = require('./scripts/layer/HeatmapTileLayer');
    var WordCloudLayer = require('./scripts/layer/WordCloudLayer');
    var TopicFrequencyLayer = require('./scripts/layer/TopicFrequencyLayer');

    var LayerMenu = require('./scripts/ui/LayerMenu');
    var Slider = require('./scripts/ui/Slider');
    var RangeSlider = require('./scripts/ui/RangeSlider');

    var AJAX_TIMEOUT = 5000;

    function createOpacitySlider(layer) {
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
            initialValue: Math.log2(layer.getResolution()),
            formatter: function(value) {
                return Math.pow(2, value);
            },
            slideStop: function(value) {
                var resolution = Math.pow(2, value);
                layer.setResolution(resolution);
                layer.redraw();
            }
        });
    }

    function createTimeSlider(layer) {
        return new RangeSlider({
            label: 'Time Range',
            min: layer.getMeta().timestamp.extrema.min,
            max: layer.getMeta().timestamp.extrema.max,
            step: 86400000,
            initialValue: [layer.getTimeRange().from, layer.getTimeRange().to],
            formatter: function(values) {
                var from = moment.unix(values[0] / 1000).format('MMM D, YYYY');
                var to = moment.unix(values[1] / 1000).format('MMM D, YYYY');
                return from + ' to ' + to;
            },
            slideStop: function(values) {
                layer.setTimeRange({
                    from: values[0],
                    to: values[1]
                });
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
        var baseLayer = new TileLayer('http://{s}.basemaps.cartocdn.com/dark_nolabels/{z}/{x}/{y}.png', {
            attribution: 'CartoDB'
        });
        baseLayer.addTo(map);

        // Label layer
        var labelLayer = new TileLayer('http://{s}.basemaps.cartocdn.com/dark_only_labels/{z}/{x}/{y}.png', {
            attribution: 'CartoDB'
        });
        labelLayer.addTo(map);

        var requester = new TileRequester('batch', function() {

            // Request meta data for the index
            $.ajax({
                url: ES_ENDPOINT + '/' + ES_INDEX + '/default'
            }).done(function(meta) {

                var heatmapLayer = new HeatmapTileLayer(meta, {
                    unloadInvisibleTiles: true
                });
                heatmapLayer.setResolution(64);
                heatmapLayer.setTimeRange({
                    from: meta.timestamp.extrema.min,
                    to: meta.timestamp.extrema.max
                });
                heatmapLayer.drawTile = function(canvas, index, zoom) {
                    $(canvas).addClass('blinking');
                    heatmapLayer.tileDrawn(canvas);

                    var params = heatmapLayer.getParams();

                    requester
                        .get(ES_ENDPOINT, ES_INDEX, 'heatmap', index.x, index.y, zoom, params)
                        .done(function() {
                            var url = ES_ENDPOINT + '/' + ES_INDEX + '/heatmap/' + zoom + '/' + index.x + '/' + index.y;
                            $.ajax({
                                url: url + '?' + $.param(params),
                                dataType: 'arraybuffer',
                                timeout: AJAX_TIMEOUT
                            }).done(function(buffer) {
                                heatmapLayer.render(canvas, buffer);
                                heatmapLayer.tileDrawn(canvas);
                            }).always(function() {
                                $(canvas).removeClass('blinking');
                            });
                        })
                        .fail(function() {
                            console.error('Could not generate tile.');
                        });
                };

                var wordCloudLayer = new WordCloudLayer(meta, {
                    unloadInvisibleTiles: true
                });
                wordCloudLayer.setTimeRange({
                    from: meta.timestamp.extrema.min,
                    to: meta.timestamp.extrema.max
                });
                wordCloudLayer.setTopics(TOPICS);
                wordCloudLayer.drawTile = function(tile, index, zoom) {
                    $(tile).addClass('blinking');
                    wordCloudLayer.tileDrawn(tile);

                    var params = wordCloudLayer.getParams();

                    requester
                        .get(ES_ENDPOINT, ES_INDEX, 'topiccount', index.x, index.y, zoom, params)
                        .done(function() {
                            var url = ES_ENDPOINT + '/' + ES_INDEX + '/topiccount/' + zoom + '/' + index.x + '/' + index.y;
                            $.ajax({
                                url: url + '?' + $.param(params),
                                dataType: 'json',
                                timeout: AJAX_TIMEOUT
                            }).done(function(wordCounts) {
                                wordCloudLayer.render(tile, wordCounts);
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
                var topicFrequencyLayer = new TopicFrequencyLayer(meta, {
                    unloadInvisibleTiles: true
                });
                topicFrequencyLayer.setTimeRange({
                    from: meta.timestamp.extrema.min,
                    to: meta.timestamp.extrema.max,
                    interval: 'month'
                });
                topicFrequencyLayer.setTimeInterval('month');
                topicFrequencyLayer.setTopics(TOPICS);
                topicFrequencyLayer.drawTile = function(tile, index, zoom) {
                    $(tile).addClass('blinking');
                    topicFrequencyLayer.tileDrawn(tile);

                    var params = topicFrequencyLayer.getParams();

                    requester
                        .get(ES_ENDPOINT, ES_INDEX, 'topicfrequency', index.x, index.y, zoom, params)
                        .done(function() {
                            var url = ES_ENDPOINT + '/' + ES_INDEX + '/topicfrequency/' + zoom + '/' + index.x + '/' + index.y;
                            $.ajax({
                                url: url + '?' + $.param(params),
                                dataType: 'json',
                                timeout: AJAX_TIMEOUT
                            }).done(function(topicCounts) {
                                topicFrequencyLayer.render(tile, topicCounts);
                                topicFrequencyLayer.tileDrawn(tile);
                            }).always(function() {
                                $(tile).removeClass('blinking');
                            });
                        })
                        .fail(function() {
                            console.error('Could not generate tile.');
                        });
                };

                $('.controls').append(createLayerControls('Word Cloud', wordCloudLayer, ['opacity', 'time']));
                //$('.controls').append(createLayerControls('Topic Trends', topicFrequencyLayer, ['opacity', 'time']));

                $('.controls').append(createLayerControls('Labels', labelLayer, ['opacity']));
                $('.controls').append(createLayerControls('Heatmap', heatmapLayer, ['opacity', 'resolution', 'time']));
                $('.controls').append(createLayerControls('Basemap', baseLayer, ['opacity']));

                // Add layers to the map
                heatmapLayer.addTo(map);
                wordCloudLayer.addTo(map);
                //topicFrequencyLayer.addTo(map);
            });

        });
    };

}());
