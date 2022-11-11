const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')

module.exports = {
  module: {
    rules: [
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader']
      },
      {
        test: /\.(?:jpg|png|svg)$/,
        type: 'asset/resource'
      },
      {
        test: /\.js$/,
        include: path.resolve(__dirname, './src'),
        use: 'babel-loader'
      }
    ]
  },
  output: {
    clean: true,
    filename: '[name].[contenthash].js',
    publicPath: '/'
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, 'public/index.html')
    }),
  ]
}
