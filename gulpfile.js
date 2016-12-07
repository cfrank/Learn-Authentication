'use strict';

let gulp = require('gulp'),
    flow = require('gulp-flowtype'),
    browserify = require('browserify'),
    babelify = require('babelify'),
    source = require('vinyl-source-stream'),
    buffer = require('vinyl-buffer'),
    uglify = require('gulp-uglify'),
    babel = require('gulp-babel'),
    sourcemaps = require('gulp-sourcemaps'),
    sass = require('gulp-sass'),
    concat = require('gulp-concat');

const paths = {
    scripts: {
        flow: 'static/scripts/flow/**/*.js',
        entry: 'static/scripts/flow/index.js',
        js: 'static/scripts/out/'
    },

    styles: {
        sass: 'static/styles/sass/**/*.scss',
        css: 'static/styles/out'
    }
};

const SCRIPT_BUILD_NAME = 'build.js';
const STYLE_BUILD_NAME = 'build.css';

function handleError(error){
    console.log(error);
    this.emit('end');
}

/* Static type checking with flow */
gulp.task('flowcheck', () => {
    return gulp.src(paths.scripts.flow)
        .pipe(flow({
            "weak": false,
            "killFlow": false,
            "abort": false
        }).on('error', handleError));
});

/* Compile to browser complient js */
gulp.task('javascript', ['flowcheck'], () => {
    return browserify({entries: paths.scripts.entry, debug: false})
        .transform("babelify")
        .bundle()
        .on('error', handleError)
        .pipe(source(SCRIPT_BUILD_NAME))
        .pipe(new buffer())
        .pipe(sourcemaps.init())
        .pipe(uglify())
        .pipe(sourcemaps.write('.'))
        .pipe(gulp.dest(paths.scripts.js))
});

/* Sass */
gulp.task('sass', () => {
    return gulp.src(paths.styles.sass)
        .pipe(sass().on('error', sass.logError))
        .pipe(concat(STYLE_BUILD_NAME))
        .pipe(gulp.dest(paths.styles.css));
});

gulp.task('watch', () => {
    gulp.watch(paths.scripts.flow, ['javascript']);
    gulp.watch(paths.styles.sass, ['sass']);
});

gulp.task('default', ['javascript', 'sass', 'watch']);
