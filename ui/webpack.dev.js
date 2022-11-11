const { merge } = require('webpack-merge')
const common = require('./webpack.common.js')

module.exports = merge(common, {
  mode: 'development',
  devtool: 'inline-source-map',
  devServer: {
    historyApiFallback: true,
    proxy: [
      {
        context: ['/api'],
        target: 'http://localhost:80',
        changeOrigin: true
      }
    ]
  }
})
