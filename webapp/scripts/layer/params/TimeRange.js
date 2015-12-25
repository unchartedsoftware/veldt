(function() {

    'use strict';

    var setTimeRange = function(timeRange) {
        var changed = false;
        if (timeRange.to !== this._params.to) {
            this._params.to = timeRange.to;
            changed = true;
        }
        if (timeRange.from !== this._params.from) {
            this._params.from = timeRange.from;
            changed = true;
        }
        if (changed) {
            this.clearExtrema();
        }
        return this;
    };

    var getTimeRange = function() {
        return {
            from: this._params.from,
            to: this._params.to
        };
    };

    var setTimeField = function(field) {
        this._params.time = field;
        return this;
    };

    var getTimeField = function() {
        return this._params.time;
    };

    module.exports = {
        setTimeRange: setTimeRange,
        getTimeRange: getTimeRange,
        setTimeField: setTimeField,
        getTimeField: getTimeField
    };

}());
