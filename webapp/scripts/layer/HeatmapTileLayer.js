(function() {

    'use strict';

    var _ = require('lodash');
    var CanvasTileLayer = require('./CanvasTileLayer');
    var Live = require('./mixins/Live');
    var ColorRamp = require('./mixins/ColorRamp');
    var Binning = require('./params/Binning');
    var Topic = require('./params/Topic');
    var TimeRange = require('./params/TimeRange');

    var HeatmapTileLayer = CanvasTileLayer
        // mixins
        .extend(Live)
        .extend(ColorRamp)
        // params
        .extend(Binning)
        .extend(Topic)
        .extend(TimeRange)
        // impl
        .extend({

            initialize: function(meta, options) {
                // mixins
                Live.initialize.apply(this, arguments);
                ColorRamp.initialize.apply(this, arguments);
                // base
                L.setOptions(this, options);
            },

            render: function(canvas, data) {
                var bins = new Float64Array(data);
                var minMax = {
                    min: _.min(bins),
                    max: _.max(bins)
                };
                if (this.updateExtrema(minMax)) {
                    // if min or max has changed we need to re-render all the tiles
                    //this.redraw();
                    return;
                }
                var extrema = this.getExtrema();
                var resolution = Math.sqrt(bins.length);
                var ramp = this.getColorRamp();
                var tileCanvas = this.renderCanvas(bins, resolution, extrema, ramp, 'log');
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
