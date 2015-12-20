(function() {
    'use strict';

    //var Slider = require("bootstrap-slider");
    var $ = require('jquery');

    var S = function(spec) {
        var self = this;
        // parse inputs
        spec = spec || {};
        spec.label = spec.label || 'missing-label';
        spec.min = spec.min !== undefined ? spec.min : 0.0;
        spec.max = spec.max !== undefined ? spec.max : 1.0;
        spec.step = spec.step !== undefined ? spec.step : 0.1;
        spec.initialValue = spec.initialValue !== undefined ? spec.initialValue : 0.1;
        this._formatter = spec.formatter || function(value) {
                return value.toFixed(2);
        };
        spec.initialFormattedValue = this._formatter(spec.initialValue);
        // create container element
        this._$container = $(Templates.slider(spec));
        this._$label = this._$container.find('.control-value-label');
        // create slider and attach callbacks
        //this._$slider = this._$container.find('.slider');
        this._$slider = new Slider(this._$container.find('.slider')[0], {
            tooltip: 'hide',
        });
        //this._$slider.slider({ tooltip: 'hide' });
        if ($.isFunction(spec.slideStart)) {
            this._$slider.on('slideStart', spec.slideStart);
        }
        if ($.isFunction(spec.slideStop)) {
            this._$slider.on('slideStop', spec.slideStop);
        }
        if ($.isFunction(spec.change)) {
            this._$slider.on('change', spec.change);
        }
        if ($.isFunction(spec.slide)) {
            this._$slider.on('slide', spec.slide);
        }
        this._$slider.on('change', function(event) {
            self._$label.text(self._formatter(event.newValue));
        });
    };

    S.prototype.getElement = function() {
        return this._$container;
    };

    module.exports = S;

}());
