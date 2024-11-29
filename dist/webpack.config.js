const path = require('path');
const glob = require('glob');

module.exports = {
    entry: glob.sync('./*.ts').reduce((entries, file) => {
        const name = path.basename(file, '.ts');
        entries[name] = `./${file}`;
        return entries;
    }, {}),
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, './')
    },
    mode: 'production',
    module: {
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