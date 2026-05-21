module.exports = {
  allowedDevOrigins: ['21367f94fad27c.lhr.life'],
  async rewrites() {
    return [
      {
        source: '/api/:path*', // любой запрос к /api/...
        destination: 'http://192.168.0.192:8080/api/:path*' // ...проксируется на реальный бэкенд
      }
    ]
  }
}