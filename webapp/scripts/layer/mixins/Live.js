(function() {

    'use strict';

    var MIN = Number.MAX_VALUE;
    var MAX = 0;

    var clearExtrema = function() {
        this._extrema = {
            min: MIN,
            max: MAX
        };
    };

    var getExtrema = function() {
        return this._extrema;
    };

    var updateExtrema = function(extrema) {
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
    };

    var setMeta = function(meta) {
        this._meta = meta;
        return this;
    };

    var getMeta = function() {
        return this._meta;
    };

    var setParams = function(params) {
        this._params = params;
    };

    var getParams = function() {
        return this._params;
    };

    var initialize = function(meta) {
        this._meta = meta;
        this._params = {};
        this.clearExtrema();
    };

    module.exports = {
        initialize: initialize,
        // meta
        setMeta: setMeta,
        getMeta: getMeta,
        // params
        setParams: setParams,
        getParams: getParams,
        // extrema
        clearExtrema: clearExtrema,
        getExtrema: getExtrema,
        updateExtrema: updateExtrema
    };

}());
