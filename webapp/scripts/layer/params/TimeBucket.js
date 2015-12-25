(function() {

    'use strict';

    var TimeRange = require('./TimeRange');

    var setTimeInterval = function(interval) {
        if (interval !== this._params.interval) {
            this._params.interval = interval;
            this.clearExtrema();
        }
        return this;
    };

    var getTimeInterval = function() {
        return this._params.interval;
    };

    module.exports = {
        // time range
        setTimeRange: TimeRange.setTimeRange,
        getTimeRange: TimeRange.getTimeRange,
        setTimeField: TimeRange.setTimeField,
        getTimeField: TimeRange.getTimeField,
        // time bucketing
        setTimeInterval: setTimeInterval,
        getTimeInterval: getTimeInterval
    };

}());
