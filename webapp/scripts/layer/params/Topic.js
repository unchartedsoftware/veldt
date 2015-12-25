(function() {

    'use strict';

    var _ = require('lodash');

    var setTopics = function(topics) {
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
            this._params.topics = topicStr;
            this.clearExtrema();
        }
        return this;
    };

    var getTopics = function() {
        return (this._params.topics || '').split(',');
    };

    var setTopicField = function(field) {
        this._params.text = field;
        return this;
    };

    var getTopicField = function() {
        return this._params.text;
    };

    module.exports = {
        setTopics: setTopics,
        getTopics: getTopics,
        setTopicField: setTopicField,
        getTopicField: getTopicField
    };

}());
