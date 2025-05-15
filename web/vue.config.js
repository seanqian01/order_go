const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  // 设置页面标题
  chainWebpack: config => {
    config.plugin('html').tap(args => {
      args[0].title = 'IBS量化交易管理系统'
      return args
    })
  },
  // 开发服务器代理配置
  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true
      }
    }
  }
})
