const path = require('path');
const CopyPlugin = require('copy-webpack-plugin');

module.exports = {
  devServer: {
    contentBase: path.join(__dirname, 'dist'),
    compress: true,
    port: 9000,
  },
  plugins: [
    new CopyPlugin([
      { from: 'build' }
    ])
  ]
};
