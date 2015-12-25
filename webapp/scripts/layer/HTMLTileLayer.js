(function() {

    'use strict';

    var $ = require('jquery');
    var _ = require('lodash');
    var TileLayer = require('./TileLayer');

    var HTMLTileLayer = TileLayer.extend({

        initialize: function(meta, options) {
            L.setOptions(this, options);
        },

        onAdd: function(map) {
            var self = this;
            L.TileLayer.prototype.onAdd.call(this, map);
            map.on('zoomend', this.onZoom, this);
            map.on('click', this.onClick, this);
            $(this._container).on('mouseover', function(e) {
                self.onHover(e);
            });
        },

        onRemove: function(map) {
            map.off('zoomend', this.onZoom);
            map.off('click', this.onClick);
            $(this._container).off('mouseover');
            L.TileLayer.prototype.onRemove.call(this, map);
        },

        redraw: function() {
            if (this._map) {
                this._reset({
                    hard: true
                });
                this._update();
            }
            var self = this;
            _.forIn(this._tiles, function(tile) {
                self._redrawTile(tile);
            });
            return this;
        },

        _redrawTile: function(tile) {
            this.drawTile(tile, tile._tilePoint, this._map._zoom);
        },

        _createTile: function() {
            var tile = L.DomUtil.create('div', 'leaflet-tile leaflet-html-tile');
            tile.width = this.options.tileSize;
            tile.height = this.options.tileSize;
            tile.onselectstart = L.Util.falseFn;
            tile.onmousemove = L.Util.falseFn;
            return tile;
        },

        _loadTile: function(tile, tilePoint) {
            tile._layer = this;
            tile._tilePoint = tilePoint;
            this._adjustTilePoint(tilePoint);
            this._redrawTile(tile);
        },

        tileDrawn: function(tile) {
            this._tileOnLoad.call(tile);
        },

        render: function() {
            // override
        },

        onZoom: function() {
            // override
        },

        onHover: function() {
            // override
        },

        onClick: function() {
            // override
        }

    });

    module.exports = HTMLTileLayer;

}());
