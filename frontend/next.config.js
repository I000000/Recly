// frontend/next.config.js
const apiBaseUrl = process.env.API_BASE_URL || 'http://localhost:8080';

module.exports = {
  allowedDevOrigins: ['a6b972552844ae.lhr.life'],
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${apiBaseUrl}/api/:path*`,
      },
      {
        source: '/uploads/:path*',
        destination: `${apiBaseUrl}/uploads/:path*`,
      },
    ];
  },
};