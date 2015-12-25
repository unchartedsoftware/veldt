(function() {

    'use strict';

    var FROM = {
        r: 150,
        g: 0,
        b: 0,
        a: 150
    };

    var TO = {
        r: 255,
        g: 255,
        b: 50,
        a: 255
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

    var interpolateColor = function(value, extrema, ramp, type) {
        var alpha = transformValue(value, extrema.min, extrema.max, type);
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

    var renderCanvas = function(bins, resolution, extrema, ramp, type) {
        var canvas = document.createElement('canvas');
        canvas.height = resolution;
        canvas.width = resolution;
        var ctx = canvas.getContext('2d');
        var imageData = ctx.getImageData(0, 0, resolution, resolution);
        var data = imageData.data;
        bins.forEach(function(bin, index) {
            var rgba = interpolateColor(bin, extrema, ramp, type);
            data[index * 4] = rgba.r;
            data[index * 4 + 1] = rgba.g;
            data[index * 4 + 2] = rgba.b;
            data[index * 4 + 3] = rgba.a;
        });
        ctx.putImageData(imageData, 0, 0);
        return canvas;
    };

    var setColorRamp = function(colorRamp) {
        this._colorRamp = this._colorRamp || {};
        if (colorRamp.to) {
            this._colorRamp.to = colorRamp.to;
        }
        if (colorRamp.from) {
            this._colorRamp.from = colorRamp.from;
        }
        return this;
    };

    var getColorRamp = function() {
        return this._colorRamp;
    };

    var initialize = function() {
        this._colorRamp = {
            from: FROM,
            to: TO
        };
    };

    module.exports = {
        initialize: initialize,
        setColorRamp: setColorRamp,
        getColorRamp: getColorRamp,
        renderCanvas: renderCanvas
    };

}());
