const path = require('path')
const webpack = require('webpack')

module.exports = {
  context: path.join(__dirname, 'client'),
  entry: './index.js',
  output: {
    filename: 'bundle.js',
    path: path.join(__dirname, 'server/data/static/build')
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: [
          'babel-loader',
        ],
      },
      {
        test: [/\.vert$/, /\.frag$/],
        use: 'raw-loader',
      },
    ],
  },
  plugins: [
    new webpack.DefinePlugin({
      'CANVAS_RENDERER': JSON.stringify(true),
      'WEBGL_RENDERER': JSON.stringify(true)
    }),
  ]
}
