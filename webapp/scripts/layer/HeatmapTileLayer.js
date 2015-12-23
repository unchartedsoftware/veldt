(function() {

    'use strict';

    var _ = require('lodash');
    var CanvasTileLayer = require('./CanvasTileLayer');
    var FROM_DEFAULT = {
        r: 150,
        g: 0,
        b: 0,
        a: 150
    };
    var TO_DEFAULT = {
        r: 255,
        g: 255,
        b: 50,
        a: 255
    };

    var logTransform = function(value, min, max) {
        var logMin = Math.log(Math.max(1, min));
        var logMax = Math.log(Math.max(1, max));
        var oneOverLogRange = 1 / (logMax - logMin);
        return Math.log(value - logMin) * oneOverLogRange;
    };

    var interpolateColor = function(value, extrema, ramp) {
        var alpha = logTransform(value, extrema.min, extrema.max);
        if (value === 0) {
            return {
                r: 255,
                g: 255,
                b: 255,
                a: 0
            };
        } else {
            return {
                r: ramp.to.r * alpha + ramp.from.r * (1 - alpha),
                g: ramp.to.g * alpha + ramp.from.g * (1 - alpha),
                b: ramp.to.b * alpha + ramp.from.b * (1 - alpha),
                a: ramp.to.a * alpha + ramp.from.a * (1 - alpha)
            };
        }
    };

    var getMinMax = function(bins) {
        return {
            min: _.min(bins),
            max: _.max(bins)
        };
    };

    var blitCanvas = function(bins, resolution, extrema, ramp) {
        var canvas = document.createElement('canvas');
        canvas.height = resolution;
        canvas.width = resolution;
        var ctx = canvas.getContext('2d');
        var imageData = ctx.getImageData(0, 0, resolution, resolution);
        var data = imageData.data;
        bins.forEach(function(bin, index) {
            var rgba = interpolateColor(bin, extrema, ramp);
            data[index * 4] = rgba.r;
            data[index * 4 + 1] = rgba.g;
            data[index * 4 + 2] = rgba.b;
            data[index * 4 + 3] = rgba.a;
        });
        ctx.putImageData(imageData, 0, 0);
        return canvas;
    };

    var HeatmapTileLayer = CanvasTileLayer.extend({

        initialize: function(meta, options) {
            this._meta = meta;
            this._params = {};
            this._colorRamp = {
                from: FROM_DEFAULT,
                to: TO_DEFAULT
            };
            this.clearExtrema();
            L.setOptions(this, options);
        },

        getMeta: function() {
            return this._meta;
        },

        getParams: function() {
            return this._params;
        },

        setResolution: function(resolution) {
            if (resolution !== this._params.resolution) {
                this._params.resolution = resolution;
                this.clearExtrema();
            }
            return this;
        },

        getResolution: function() {
            return this._params.resolution;
        },

        setColorRamp: function(colorRamp) {
            if (colorRamp.to !== this._colorRamp.to ||
                colorRamp.from !== this._colorRamp.from) {
                this._colorRamp = colorRamp;
            }
            return this;
        },

        getColorRamp: function() {
            return this._colorRamp;
        },

        setTimeRange: function(timeRange) {
            if (timeRange.to !== this._params.to ||
                timeRange.from !== this._params.from) {
                this._params.to = timeRange.to;
                this._params.from = timeRange.from;
                this.clearExtrema();
            }
            return this;
        },

        getTimeRange: function() {
            return {
                from: this._params.from,
                to: this._params.to
            };
        },

        setTopics: function(topics) {
            topics = _.map(topics, function(t) {
                return t.toLowerCase();
            });
            topics.sort(function(a, b) {
                if (a < b) {
                    return -1;
                }
                if (a > b) {
                    return 1;
                }
                return 0;
            });
            var topicStr = topics.join(',');
            if (topicStr !== this._params.topics) {
                this._params.topcs = topicStr;
                this.clearExtrema();
            }
            return this;
        },

        clearExtrema: function() {
            this._extrema = {
                min: Number.MAX_VALUE,
                max: 0
            };
        },

        _updateExtrema: function(extrema) {
            var changed = false;
            if (extrema.min < this._extrema.min) {
                changed = true;
                this._extrema.min = extrema.min;
            }
            if (extrema.max > this._extrema.max) {
                changed = true;
                this._extrema.max = extrema.max;
            }
            return changed;
        },

        render: function(canvas, data) {
            var bins = new Float64Array(data);
            var minMax = getMinMax(bins);
            if (this._updateExtrema(minMax)) {
                //this.redraw();
                // FIX THIS!
                return;
            }
            var extrema = this._extrema;
            var resolution = this._params.resolution;
            var ramp = this._colorRamp;
            var tileCanvas = blitCanvas(bins, resolution, extrema, ramp);
            var ctx = canvas.getContext('2d');
            ctx.imageSmoothingEnabled = false;
            ctx.drawImage(
                tileCanvas,
                0, 0,
                resolution, resolution,
                0, 0,
                canvas.width, canvas.height);
        }

    });

    module.exports = HeatmapTileLayer;

}());
