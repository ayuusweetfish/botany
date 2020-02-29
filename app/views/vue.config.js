module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  productionSourceMap: false,
  devServer: {
    proxy: {
      '/api': {
        target: 'http://140.143.120.174:3389/',
        changeOrigin: true,
        pathRewrite: {
          '^/api': '/api'
        }
      }
    }
  }
}
