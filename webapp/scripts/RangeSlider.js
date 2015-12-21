(function() {
    'use strict';

    var $ = require('jquery');

    var RangeSlider = function(spec) {
        var self = this;
        // parse inputs
        spec = spec || {};
        spec.label = spec.label || 'missing-label';
        spec.min = spec.min !== undefined ? spec.min : 0.0;
        spec.max = spec.max !== undefined ? spec.max : 1.0;
        spec.step = spec.step !== undefined ? spec.step : 0.1;
        spec.range = true;
        this._formatter = spec.formatter || function(values) {
                return values[0].toFixed(2) + ' to ' + values[1].toFixed(2);
        };
        spec.initialFormattedValue = this._formatter(spec.initialValue);
        if (spec.initialValue !== undefined) {
            spec.initialValue = '[' + spec.initialValue[0] + ',' + spec.initialValue[1] + ']';
        } else {
            spec.initialValue = '[' + spec.min + ',' + spec.max + ']';
        }
        // create container element
        this._$container = $(Templates.slider(spec));
        this._$label = this._$container.find('.control-value-label');
        // create slider and attach callbacks
        this._$slider = new Slider(this._$container.find('.slider')[0], {
            tooltip: 'hide',
        });
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

    RangeSlider.prototype.getElement = function() {
        return this._$container;
    };

    module.exports = RangeSlider;

}());
