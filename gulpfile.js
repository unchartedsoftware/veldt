( function() {

    'use strict';

    var gulp = require('gulp'),
        concat = require('gulp-concat'),
        runSequence,
        bower,
        csso,
        filter,
        uglify;

    var projectName = 'prism';
    var webappPath = 'webapp/';
    var serverPath = 'server/';
    var outputDir = 'build/public/';
    var paths = {
        serverRoot: serverPath + '/main.go',
        webappRoot: webappPath + '/app.js',
        go: [ './**/*.go' ],
        templates: [ webappPath + 'templates/**/*.hbs'],
        scripts: [ webappPath + 'scripts/**/*.js',  webappPath + 'app.js' ],
        styles: [  webappPath + 'styles/reset.css',  webappPath + 'styles/**/*.css' ],
        index: [  webappPath + 'index.html' ],
        resources: [
            webappPath + 'index.html',
            webappPath + 'favicons/*'
        ]
    };

    function handleError( err ) {
        console.log( err );
        this.emit( 'end' );
    }

    function bundle( bundler, watch ) {
        var source = require('vinyl-source-stream');
        if ( watch ) {
            var watchify = require('watchify');
            var watcher = watchify( bundler );
            watcher.on( 'update', function( ids ) {
                // When any files updates
                console.log('\nWatch detected changes to: ');
                for ( var i=0; i<ids.length; ids++ ) {
                   console.log( '\t'+ids[i] );
                }
                var updateStart = Date.now();
                watcher.bundle()
                    .on( 'error', handleError )
                    .pipe( source( projectName + '.js' ) )
                    // This is where you add uglifying etc.
                    .pipe( gulp.dest( outputDir ) );
                console.log( 'Updated in', ( Date.now() - updateStart ) + 'ms' );
            });
            bundler = watcher;
        }
        return bundler
            .bundle() // Create the initial bundle when starting the task
            .on( 'error', handleError )
            .pipe( source( projectName + '.js' ) )
            .pipe( gulp.dest( outputDir ) );
    }

    gulp.task('clean', function () {
        var del = require('del');
        del.sync([ outputDir + '/*']);
    });

    gulp.task('lint', function() {
        var jshint = require('gulp-jshint');
        return gulp.src( [ './webapp/**/*.js',
            '!./webapp/extern/**/*.js'] )
            .pipe( jshint() )
            .pipe( jshint('.jshintrc') )
            .pipe( jshint.reporter('jshint-stylish') );
    });

    gulp.task('build-and-watch-scripts', function() {
        var browserify = require('browserify'),
            bundler = browserify( paths.webappRoot, {
                debug: true,
                standalone: projectName
            });
        return bundle( bundler, true );
    });

    gulp.task('build-scripts', function() {
        var browserify = require('browserify'),
            bundler = browserify( paths.webappRoot, {
                debug: true,
                standalone: projectName
            });
        return bundle( bundler, false );
    });

    gulp.task('build-templates',function() {
        var handlebars = require('gulp-handlebars');
        var wrap = require('gulp-wrap');
        var declare = require('gulp-declare');
        return gulp.src( paths.templates )
            .pipe( handlebars({
                // Pass your local handlebars version
                handlebars: require('handlebars')
            }))
            .pipe( wrap('Handlebars.template(<%= contents %>)'))
            .pipe(declare({
                namespace: 'Templates',
                noRedeclare: true, // avoid duplicate declarations
            }))
            .pipe( concat( projectName + '.templates.js' ) )
            .pipe( gulp.dest( outputDir ) );
    });

    gulp.task('build-styles', function () {
        csso = csso || require('gulp-csso');
        var concat = require('gulp-concat');
        return gulp.src( paths.styles )
            .pipe( csso() )
            .pipe( concat( projectName + '.css') )
            .pipe( gulp.dest( outputDir ) );
    });

    gulp.task('copy-resources', function() {
        return gulp.src( paths.resources, {
                base: webappPath
            })
            .pipe( gulp.dest( outputDir ) );
    });

    gulp.task('build-vendor-js', function() {
        filter = filter || require('gulp-filter');
        bower = bower || require('main-bower-files');
        uglify = uglify || require('gulp-uglify');
        return gulp.src( bower() )
            .pipe( filter('**/*.js') ) // filter js files
            .pipe( concat('vendor.min.js') )
            .pipe( uglify() )
            .pipe( gulp.dest( outputDir ) );
    });

    gulp.task('build-vendor-css', function() {
        filter = filter || require('gulp-filter');
        bower = bower || require('main-bower-files');
        csso = csso || require('gulp-csso');
        return gulp.src( bower() )
            .pipe( filter('**/*.css') ) // filter css files
            .pipe( csso() )
            .pipe( concat('vendor.min.css') )
            .pipe( gulp.dest( outputDir ) );
    });

    gulp.task('build-vendor', function( done ) {
        runSequence = runSequence || require('run-sequence');
        runSequence([
            'build-vendor-js',
            'build-vendor-css' ],
            done );
    });

    gulp.task('build', function( done ) {
        runSequence = runSequence || require('run-sequence');
        runSequence(
            [ 'clean', 'lint' ],
            [ 'build-and-watch-scripts', 'build-templates', 'build-styles', 'build-vendor', 'copy-resources' ],
            done );
    });

    gulp.task('deploy', function( done ) {
        runSequence = runSequence || require('run-sequence');
        runSequence(
            [ 'clean', 'lint' ],
            [ 'build-scripts', 'build-templates', 'build-styles', 'build-vendor', 'copy-resources' ],
            done );
    });

    var go;
    gulp.task('serve', function() {
        var gulpgo = require( 'gulp-go' );
        go = gulpgo.run( paths.serverRoot, [
            '-alias', 'development=http://192.168.0.41:9200',
            '-alias', 'production=http://10.65.16.13:9200',
            '-alias', 'openstack=http://10.64.16.120:9200',
        ], {
            cwd: __dirname,
            stdio: 'inherit'
        });
    });

    gulp.task('watch', [ 'build' ], function( done ) {
        gulp.watch( paths.go ).on('change', function() {
            go.restart();
        });
        gulp.watch( paths.styles, [ 'build-styles' ] );
        gulp.watch( paths.templates, [ 'build-templates' ]);
        gulp.watch( paths.resources, [ 'copy-resources' ] );
        done();
    });

    gulp.task('default', [ 'watch', 'serve' ], function() {
    });

}());
