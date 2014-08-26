var gulp            = require('gulp'),
    component       = require('gulp-component'),
    uglify          = require('gulp-uglify'),
    rename          = require('gulp-rename'),
    less            = require('gulp-less'),
    jshint          = require('gulp-jshint'),
    util            = require('gulp-util'),
    minifyCSS       = require('gulp-minify-css'),
    shell           = require('gulp-shell')
    //mocha           = require('gulp-mocha'),
    //mochaPhantomJS  = require('gulp-mocha-phantomjs'),
    //server          = require('./server'),
    //serverPort      = 5000

var paths = {
  scripts: ['src/**/*.js'],
  tests: 'test/**/*.js'
}

gulp.task('default', function () {
})

gulp.task('server', function () {
    gulp.watch(['component.json', 'src/**/*'], ['build'])
    server.listen(serverPort)
})

gulp.task('go-run', shell.task([
   'run'
]))

gulp.task('less', function () {
    gulp.src('assets/css/admin/style.less')
        .pipe(less())
        .pipe(gulp.dest('assets/css/admin/'))
        .pipe(minifyCSS({keepBreaks:false}))
        .pipe(gulp.dest('assets/css/admin/'))
})

gulp.task('watch', function () {
    //server.listen(serverPort)
    //gulp.watch(['**/*.go'], ['go-run'])
    gulp.watch(['*.exe'], ['go-run'])
    gulp.watch(['assets/css/**/*.less'], ['less'])
    //gulp.watch(['test/**/*.*'], ['test'])
})

