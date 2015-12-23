(function() {
    'use strict';

    var $ = require('jquery');

    var ToggleBox = function(spec) {
        var self = this;
        // parse inputs
        spec = spec || {};
        spec.label = spec.label !== undefined ? spec.label : 'missing-label';
        spec.initialValue = spec.initialValue !== undefined ? spec.initialValue : true;
        this.onEnabled = spec.enabled || null;
        this.onDisabled = spec.disabled || null;
        this.onToggle = spec.toggled || null;
        // set initial state
        this._enabled = spec.initialValue;
        // create elements
        this._$container = $(Templates.togglebox(spec));
        this._$toggleIcon = this._$container.find('i.fa');
        this._$toggleBox = this._$container.find('.layer-toggle');
        // create and attach callbacks
        this._enable = function() {
            self._$toggleIcon.removeClass('fa-square-o');
            self._$toggleIcon.addClass('fa-check-square-o');
        };
        this._disable = function() {
            self._$toggleIcon.removeClass('fa-check-square-o');
            self._$toggleIcon.addClass('fa-square-o');
        };
        this._$toggleBox.click(function() {
            self._enabled = !self._enabled;
            if (self._enabled) {
                self._enable();
                if ($.isFunction(self.onEnabled)) {
                    self.onEnabled();
                }
            } else {
                self._disable();
                if ($.isFunction(self.onDisabled)) {
                    self.onDisabled();
                }
            }
            if ($.isFunction(self.onToggle)) {
                self.onToggle(self._enabled);
            }
        });
    };

    ToggleBox.prototype.getElement = function() {
        return this._$container;
    };

    module.exports = ToggleBox;

}());
