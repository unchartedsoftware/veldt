( function() {

    "use strict";

    var gulp = require('gulp'),
        concat = require('gulp-concat'),
        runSequence,
        bower,
        csso,
        filter,
        uglify;

    var projectName = 'prism';

    var output = {
        dir: 'build/',
        app: 'app',
        vendor: 'vendor'
    };

    var webappPath = 'webapp/';
    var serverPath = 'server/';

    var paths = {
        root: webappPath + '/app.js',
        server: [ serverPath + '**/*.go' ],
        scripts: [ webappPath + 'scripts/**/*.js',  webappPath + 'app.js' ],
        styles: [  webappPath + 'styles/reset.css',  webappPath + 'styles/**/*.css' ],
        index: [  webappPath + 'index.html' ],
        resources: [
            webappPath + 'index.html'
        ]
    };

    function handleError( err ) {
        console.log( err );
        this.emit( 'end' );
    }

    function bundle( bundler ) {
        var watchify = require('watchify'),
            watcher = watchify( bundler ),
            source = require('vinyl-source-stream');
        return watcher
            .on( 'update', function( ids ) {
                // When any files updates
                console.log("\nWatch detected changes to: ");
                for ( var i=0; i<ids.length; ids++ ) {
                   console.log( '\t'+ids[i] );
                }
                var updateStart = Date.now();
                watcher.bundle()
                    .on( 'error', handleError )
                    .pipe( source( projectName + ".js" ) )
                    // This is where you add uglifying etc.
                    .pipe( gulp.dest( output.dir ) );
                console.log( 'Updated in', ( Date.now() - updateStart ) + 'ms' );
            })
            .bundle() // Create the initial bundle when starting the task
            .on( 'error', handleError )
            .pipe( source( projectName + ".js" ) )
            .pipe( gulp.dest( output.dir ) );
    }

    gulp.task('clean', function () {
        var del = require('del');
        del.sync([ 'build/*']);
    });

    gulp.task('lint', function() {
        var jshint = require('gulp-jshint');
        return gulp.src( [ './webapp/**/*.js',
            '!./webapp/extern/**/*.js'] )
            .pipe( jshint() )
            .pipe( jshint('.jshintrc') )
            .pipe( jshint.reporter('jshint-stylish') );
    });

    gulp.task('build-scripts', function() {
        var browserify = require('browserify'),
            bundler = browserify( paths.root, {
                debug: true,
                standalone: projectName
            });
        return bundle( bundler );
    });

    gulp.task('build-styles', function () {
        csso = csso || require('gulp-csso');
        var concat = require('gulp-concat');
        return gulp.src( paths.styles )
            .pipe( csso() )
            .pipe( concat( projectName + '.css') )
            .pipe( gulp.dest( output.dir ) );
    });

    gulp.task('build-src', [ 'build-scripts', 'build-styles' ], function() {
    });

    gulp.task('copy-resources', function() {
        return gulp.src( paths.resources, {
                base: webappPath
            })
            .pipe( gulp.dest( output.dir ) );
    });

    gulp.task('build-vendor-js', function() {
        filter = filter || require('gulp-filter');
        bower = bower || require('main-bower-files');
        uglify = uglify || require('gulp-uglify');
        return gulp.src( bower() )
            .pipe( filter('**/*.js') ) // filter js files
            .pipe( concat('vendor.min.js') )
            .pipe( uglify() )
            .pipe( gulp.dest( output.dir ) );
    });

    gulp.task('build-vendor-css', function() {
        filter = filter || require('gulp-filter');
        bower = bower || require('main-bower-files');
        csso = csso || require('gulp-csso');
        return gulp.src( bower() )
            .pipe( filter('**/*.css') ) // filter css files
            .pipe( csso() )
            .pipe( concat('vendor.min.css') )
            .pipe( gulp.dest( output.dir ) );
    });

    gulp.task('build-vendor', [ 'build-vendor-js', 'build-vendor-css' ], function() {
    });

    gulp.task('build', function( done ) {
        runSequence = runSequence || require('run-sequence');
        runSequence(
            [ 'clean', 'lint' ],
            [ 'build-src', 'build-vendor', 'copy-resources' ],
            done );
    });

    var go;
    gulp.task("serve", function() {
        var gulpgo = require( "gulp-go" );
        go = gulpgo.run( serverPath + "main.go" );
    });

    gulp.task('watch', [ 'build' ], function( done ) {
        gulp.watch( paths.server ).on('change', function() {
            go.restart();
        });
        gulp.watch( paths.styles, [ 'build-styles' ] );
        gulp.watch( paths.resources, [ 'copy-resources' ] );
        done();
    });

    gulp.task('default', [ 'watch', 'serve' ], function() {
    });

}());
