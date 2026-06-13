module.exports = {
  allowedDevOrigins: ['2afa05627ab850.lhr.life'],
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://192.168.0.192:8080/api/:path*',
      },
      {
        source: '/uploads/:path*',
        destination: 'http://192.168.0.192:8080/uploads/:path*',
      },
    ];
  },
};