module.exports = {
  allowedDevOrigins: ['5e940d184bf634.lhr.life'],
  allowedDevOrigins: ['0f087b8ec4a99c.lhr.life'],
  async rewrites() {
    return [
      {
        source: '/api/:path*', // любой запрос к /api/...
        destination: 'http://192.168.0.192:8080/api/:path*' // ...проксируется на реальный бэкенд
      }
    ]
  }
}