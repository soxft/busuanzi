const path = require('path');
const webpack = require('webpack');

module.exports = {
    entry:'./busuanzi.ts',
    output:{
        filename:'busuanzi.js',
        path: path.resolve(__dirname,'./')
    },
    mode: 'production',
    module:{
        rules: [{
            test: /\.tsx?$/,
            use: 'ts-loader',
            exclude: /node_modules/
        }]
    },
    resolve: {
        extensions: ['.ts']
    },
}