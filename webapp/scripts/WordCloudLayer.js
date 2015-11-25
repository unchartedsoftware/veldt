(function() {

    'use strict';

    var $ = require('jquery');
    var _ = require('lodash');

    var HORIZONTAL_OFFSET = 10;
    var TILE_SIZE = 256;
    var HALF_SIZE = TILE_SIZE / 2;
    var VERTICAL_OFFSET = 24;
    var SIZE_FUNCTION = 'log';
    var MAX_NUM_WORDS = 15;
    var MIN_FONT_SIZE = 10;
    var MAX_FONT_SIZE = 20;
    var NUM_ATTEMPTS = 1;

    /**
     * Given an initial position, return a new position, incrementally spiralled
     * outwards.
     */
    var spiralPosition = function(pos) {
        var pi2 = 2 * Math.PI;
        var circ = pi2 * pos.radius;
        var inc = (pos.arcLength > circ / 10) ? circ / 10 : pos.arcLength;
        var da = inc / pos.radius;
        var nt = (pos.t + da);
        if (nt > pi2) {
            nt = nt % pi2;
            pos.radius = pos.radius + pos.radiusInc;
        }
        pos.t = nt;
        pos.x = pos.radius * Math.cos(nt);
        pos.y = pos.radius * Math.sin(nt);
        return pos;
    };

    /**
     *  Returns true if bounding box a intersects bounding box b
     */
    var intersectTest = function(a, b) {
        return (Math.abs(a.x - b.x) * 2 < (a.width + b.width)) &&
            (Math.abs(a.y - b.y) * 2 < (a.height + b.height));
    };

    /**
     *  Returns true if bounding box a is not fully contained inside bounding box b
     */
    var overlapTest = function(a, b) {
        return (a.x + a.width / 2 > b.x + b.width / 2 ||
            a.x - a.width / 2 < b.x - b.width / 2 ||
            a.y + a.height / 2 > b.y + b.height / 2 ||
            a.y - a.height / 2 < b.y - b.height / 2);
    };

    /**
     * Check if a word intersects another word, or is not fully contained in the
     * tile bounding box
     */
    var intersectWord = function(position, word, cloud, bb) {
        var box = {
            x: position.x,
            y: position.y,
            height: word.height,
            width: word.width
        };
        var i;
        for (i = 0; i < cloud.length; i++) {
            if (intersectTest(box, cloud[i])) {
                return true;
            }
        }
        // make sure it doesn't intersect the border;
        if (overlapTest(box, bb)) {
            // if it hits a border, increment collision count
            // and extend arc length
            position.collisions++;
            position.arcLength = position.radius;
            return true;
        }
        return false;
    };

    /**
     * Returns the word counts along with font-size, scale percent, height, and width.
     */
    var measureWords = function(wordCounts, min, max, sizeFunction) {
        // sort words by frequency
        wordCounts = wordCounts.sort(function(a, b) {
            return b.count - a.count;
        }).slice(0, MAX_NUM_WORDS);
        // build measurement html
        var html = '<div style="height:256px; width:256px;">';
        wordCounts.forEach(function(word) {
            word.percent = transformValue(word.count, min, max, sizeFunction);
            word.fontSize = MIN_FONT_SIZE + word.percent * (MAX_FONT_SIZE - MIN_FONT_SIZE);
            html += '<div class="word-cloud-label" style="visibility:hidden; font-size:' + word.fontSize + 'px;">' + word.text + '</div>';
        });
        html += '</div>';
        // append measurements
        var $temp = $(html);
        $('body').append($temp);
        $temp.children().each(function(index) {
            wordCounts[index].width = this.offsetWidth;
            wordCounts[index].height = this.offsetHeight;
        });
        $temp.remove();
        return wordCounts;
    };

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

    /**
     * Returns the word cloud words containing font size and x and y coordinates
     */
    var createWordCloud = function(wordCounts, min, max) {
        var boundingBox = {
            width: TILE_SIZE - HORIZONTAL_OFFSET * 2,
            height: TILE_SIZE - VERTICAL_OFFSET * 2,
            x: 0,
            y: 0
        };
        var cloud = [];
        // sort words by frequency
        wordCounts = measureWords(wordCounts, min, max, SIZE_FUNCTION);
        // assemble word cloud
        wordCounts.forEach(function(wordCount) {
            // starting spiral position
            var pos = {
                radius: 1,
                radiusInc: 5,
                arcLength: 10,
                x: 0,
                y: 0,
                t: 0,
                collisions: 0
            };
            // spiral outwards to find position
            while (pos.collisions < NUM_ATTEMPTS) {
                // increment position in a spiral
                pos = spiralPosition(pos);
                // test for intersection
                if (!intersectWord(pos, wordCount, cloud, boundingBox)) {
                    cloud.push({
                        text: wordCount.text,
                        fontSize: wordCount.fontSize,
                        percent: Math.round((wordCount.percent * 100) / 10) * 10, // round to nearest 10
                        x: pos.x,
                        y: pos.y,
                        width: wordCount.width,
                        height: wordCount.height
                    });
                    break;
                }
            }
        });
        return cloud;
    };

    var render = function(container, data, extents, highlight) {
        var wordCounts = _.map(data, function(count, key) {
            return {
                count: count,
                text: key
            };
        });
        // exit early if no words
        if (wordCounts.length === 0) {
            return;
        }
        // genereate the cloud
        var cloud = createWordCloud(wordCounts, extents.min, extents.max);
        // build html elements
        var html = '';
        cloud.forEach(function(word) {
            // create classes
            var classNames = [
                'word-cloud-label',
                'word-cloud-label-' + word.percent,
                word.text === highlight ? 'highlight' : ''
            ].join(' ');
            // create styles
            var styles = [
                'font-size:' + word.fontSize + 'px',
                'left:' + (HALF_SIZE + word.x - (word.width / 2)) + 'px',
                'top:' + (HALF_SIZE + word.y - (word.height / 2)) + 'px',
                'width:' + word.width + 'px',
                'height:' + word.height + 'px',
            ].join(';');
            // create html for entry
            html += '<div class="' + classNames + '"' +
                'style="' + styles + '"' +
                'data-word="' + word.text + '">' +
                word.text +
                '</div>';
        });
        container.html(html);
    };

    var WordCloudLayer = L.TileLayer.extend({

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
            $('.word-cloud-label').removeClass('hover');
            var word = target.attr('data-word');
            if (word) {
                $('.word-cloud-label[data-word=' + word + ']').addClass('hover');
            }
        },

        onClick: function(e) {
            // Highlight clicked-on word
            var target = $(e.originalEvent.target);
            $('.word-cloud-label').removeClass('highlight');
            var word = target.attr('data-word');
            if (word) {
                $(this._container).addClass('highlight');
                $('.word-cloud-label[data-word=' + word + ']').addClass('highlight');
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
            _.forIn(this._tiles, function(tile) {
                this._redrawTile(tile);
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

    module.exports = WordCloudLayer;

}());
