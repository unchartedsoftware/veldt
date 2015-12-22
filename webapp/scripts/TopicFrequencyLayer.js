(function() {

    'use strict';

    var $ = require('jquery');
    var _ = require('lodash');

    var TILE_SIZE = 256;
    var HALF_SIZE = TILE_SIZE / 2;
    var SIZE_FUNCTION = 'log';
    var MAX_NUM_WORDS = 8;
    var MIN_FONT_SIZE = 16;
    var MAX_FONT_SIZE = 22;

    var transformValue = function(value, min, max, type) {
        var clamped = Math.max(Math.min(value, max), min);
        if (type === 'log') {
            var logMin = Math.log10(min || 1);
            var logMax = Math.log10(max || 1);
            var oneOverLogRange = 1 / ((logMax - logMin) || 1);
            return (Math.log10(clamped || 1) - logMin) * oneOverLogRange;
        } else {
            var range = max - min;
            return (clamped - min) / range;
        }
    };

    var render = function(container, data, extents, highlight) {
        // convert object to array
        var values = _.keys(data).map(function(key) {
            var counts = data[key];
            return {
                topic: key,
                counts: counts,
                max: _.max(counts),
                total: _.sum(counts)
            };
        }).sort(function(a, b) {
            return b.total - a.total;
        });
        // get number of entries
        var numEntries = Math.min(values.length, MAX_NUM_WORDS);
        var $html = $('<div class="topic-frequencies" style="display:inline-block;"></div>');
        var totalHeight = 0;
        values.slice(0, numEntries).forEach(function(value) {
            var topic = value.topic;
            var counts = value.counts;
            var max = value.max;
            var total = value.total;
            var highlightClass = ( topic === highlight ) ? 'highlight' : '';
            // scale the height based on level min / max
            var percent = transformValue(total, extents.min, extents.max, SIZE_FUNCTION);
            var percentLabel = Math.round((percent * 100) / 10) * 10;
            var height = MIN_FONT_SIZE + percent * (MAX_FONT_SIZE - MIN_FONT_SIZE);
            totalHeight += height;
            // create container 'entry' for chart and hashtag
            var $entry = $('<div class="topic-frequency-entry '+highlightClass+'" ' +
                'data-word="' + topic + '"' +
                'style="' +
                'height:' + height + 'px;"></div>');
            // create chart
            var $chart = $('<div class="topic-frequency-left"' +
                'data-word="' + topic + '"' +
                '></div>');
            counts.forEach(function(count) {
                // get the percent relative to the highest count in the tile
                var relativePercent = (max !== 0) ? (count / max) * 100 : 0;
                // make invisible if zero count
                var visibility = relativePercent === 0 ? 'hidden' : '';
                // Get the style class of the bar
                var percentLabel = Math.round(relativePercent / 10) * 10;
                var barClasses = [
                    'topic-frequency-bar',
                    'topic-frequency-bar-' + percentLabel
                ].join(' ');
                var barWidth = 'calc(' + (100 / counts.length) + '% - 2px)';
                var barHeight;
                var barTop;
                // ensure there is at least a single pixel of color
                if ((relativePercent / 100) * height < 3) {
                    barHeight = '3px';
                    barTop = 'calc(100% - 3px)';
                } else {
                    barHeight = relativePercent + '%';
                    barTop = (100 - relativePercent) + '%';
                }
                // create bar
                $chart.append('<div class="' + barClasses + '" ' +
                    'data-word="' + topic + '"' +
                    'style="' +
                    'visibility:' + visibility + ';' +
                    'width:' + barWidth + ';' +
                    'height:' + barHeight + ';' +
                    'top:' + barTop + ';"></div>');
            });
            $entry.append($chart);
            var topicClasses = [
                'topic-frequency-label',
                'topic-frequency-label-' + percentLabel
            ].join(' ');
            // create tag label
            var $topic = $('<div class="topic-frequency-right">' +
                '<div class="' + topicClasses + '" ' +
                'data-word="' + topic + '"' +
                'style="' +
                'font-size:' + height + 'px;' +
                'line-height:' + height + 'px;' +
                'height:' + height + 'px">' + topic + '</div>' +
                '</div>');
            $entry.append($topic);
            $html.append($entry);
        });
        $html.css('top', HALF_SIZE - (totalHeight / 2));
        container.html($html[0].outerHTML);
    };

    var TopicFrequencyLayer = L.TileLayer.extend({

        initialize: function(url, options) {
            this._url = url;
            L.setOptions(this, options);
        },

        onAdd: function(map) {
            var self = this;
            L.TileLayer.prototype.onAdd.call(this, map);
            map.on('zoomend', this.onZoom, this);
            map.on('click', this.onClick, this);
            $(this._container).on('mouseover', function(e) {
                self.onHover(e);
            });
        },

        onRemove: function(map) {
            map.off('zoomend', this.onZoom);
            map.off('click', this.onClick);
            $(this._container).off('mouseover');
            this.highlight = null;
            L.TileLayer.prototype.onRemove.call(this, map);
        },

        onZoom: function() {
            // Clear highlighting
            $(this._container).removeClass('highlight');
            this.highlight = null;
        },

        onHover: function(e) {
            var target = $(e.originalEvent.target);
            $('.topic-frequency-entry').removeClass('hover');
            var word = target.attr('data-word');
            if (word) {
                $('.topic-frequency-entry[data-word=' + word + ']').addClass('hover');
            }
        },

        onClick: function(e) {
            // Highlight clicked-on word
            var target = $(e.originalEvent.target);
            $('.topic-frequency-entry').removeClass('highlight');
            var word = target.attr('data-word');
            if (word) {
                $(this._container).addClass('highlight');
                $('.topic-frequency-entry[data-word=' + word + ']').addClass('highlight');
                this.highlight = word;
            } else {
                $(this._container).removeClass('highlight');
                this.highlight = null;
            }
        },

        redraw: function() {
            if (this._map) {
                this._reset({
                    hard: true
                });
                this._update();
            }
            var self = this;
            _.forIn(this._tiles, function(tile) {
                self._redrawTile(tile);
            });
            return this;
        },

        _redrawTile: function(tile) {
            this.drawTile(tile, tile._tilePoint, this._map._zoom);
        },

        _createTile: function() {
            var tile = L.DomUtil.create('div', 'leaflet-tile leaflet-wordcloud');
            tile.width = tile.height = this.options.tileSize;
            tile.onselectstart = tile.onmousemove = L.Util.falseFn;
            return tile;
        },

        _loadTile: function(tile, tilePoint) {
            tile._layer = this;
            tile._tilePoint = tilePoint;
            this._adjustTilePoint(tilePoint);
            this._redrawTile(tile);
            if (!this.options.async) {
                this.tileDrawn(tile);
            }
        },

        tileDrawn: function(tile) {
            this._tileOnLoad.call(tile);
        },

        render: render,

        refresh: function() {}
    });

    module.exports = TopicFrequencyLayer;

}());
