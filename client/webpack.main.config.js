const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    entry: {
        main: './src/index.tsx', // Main entry point
    },
    output: {
        filename: 'bundle.js', // Output bundle file in separate directories
        path: path.resolve(__dirname, 'dist/main'),
        publicPath: './', // Set publicPath to './' to ensure correct script file path
        clean: true, // Clean the output directory before emit
    },
    resolve: {
        extensions: ['.ts', '.tsx', '.js', '.jsx'], // Resolve these extensions
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/, // Match .ts and .tsx files
                use: 'ts-loader', // Use ts-loader to transpile TypeScript
                exclude: /node_modules/, // Exclude node_modules
            },
            {
                test: /\.css$/, // Match .css files
                use: ['style-loader', 'css-loader'], // Use style-loader and css-loader
            },
        ],
    },
    devtool: 'source-map', // Enable source maps for debugging
    devServer: {
        static: path.join(__dirname, 'dist'), // Serve static files from dist
        compress: true, // Enable gzip compression
        port: 3000, // Serve on port 3000
        hot: true, // Enable Hot Module Replacement
    },
    plugins: [
        new HtmlWebpackPlugin({
            filename: 'index.html', // Output HTML file for the main entry
            template: './src/index.html', // Template for the main entry HTML file
            inject: true,
            scriptLoading: 'defer', // Add the defer attribute to the generated script tag
            publicPath: '', // Set publicPath to empty to remove any path from the script src
        }),
    ],
};