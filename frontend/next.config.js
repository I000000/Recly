module.exports = {
  allowedDevOrigins: ['c33e88280607ea.lhr.life'],
  async rewrites() {
    return [
      {
        source: '/api/:path*', // любой запрос к /api/...
        destination: 'http://192.168.0.192:8080/api/:path*' // ...проксируется на реальный бэкенд
      }
    ]
  }
}