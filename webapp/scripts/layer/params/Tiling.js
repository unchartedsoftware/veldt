(function() {

    'use strict';

    var setXField = function(field) {
        this._params.x = field;
        return this;
    };

    var getXField = function() {
        return this._params.x;
    };

    var setYField = function(field) {
        this._params.y = field;
        return this;
    };

    var getYField = function() {
        return this._params.y;
    };

    var setExtents = function(extents) {
        if (extents.top) {
            this._params.top = extents.top;
        }
        if (extents.bottom) {
            this._params.bottom = extents.bottom;
        }
        if (extents.left) {
            this._params.left = extents.left;
        }
        if (extents.right) {
            this._params.right = extents.right;
        }
        return this;
    };

    var getExtents = function() {
        return {
            top: this._params.top,
            bottom: this._params.bottom,
            left: this._params.left,
            right: this._params.right
        };
    };

    module.exports = {
        setXField: setXField,
        getXField: getXField,
        setYField: setYField,
        getYField: getYField,
        setExtents: setExtents,
        getExtents: getExtents
    };

}());
