module.exports = {
  allowedDevOrigins: ['decb075dc23ac9.lhr.life'],
  async rewrites() {
    return [
      {
        source: '/api/:path*', // любой запрос к /api/...
        destination: 'http://192.168.0.192:8080/api/:path*' // ...проксируется на реальный бэкенд
      }
    ]
  }
}