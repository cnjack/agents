const path = require('path');

// Get the project root directory (where this config file is located)
const projectRoot = __dirname;

module.exports = {
  apps: [
    // Backend Go server
    {
      name: 'stardew-backend',
      cwd: path.join(projectRoot, 'backend'),
      script: './server',
      interpreter: 'none',
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: '1G',
      env: {
        NODE_ENV: 'production',
        PORT: 8080
      },
      env_development: {
        NODE_ENV: 'development',
        PORT: 8080
      },
      error_file: path.join(projectRoot, 'logs', 'backend-error.log'),
      out_file: path.join(projectRoot, 'logs', 'backend-out.log'),
      log_file: path.join(projectRoot, 'logs', 'backend-combined.log'),
      time: true,
      merge_logs: true,
      // Pre-start hook to build the Go binary
      pre_start: 'go build -o server ./cmd/server'
    },
    // Frontend Vue/Vite dev server
    {
      name: 'stardew-frontend',
      cwd: path.join(projectRoot, 'frontend'),
      script: 'npm',
      args: 'run dev',
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: '500M',
      env: {
        NODE_ENV: 'development',
        VITE_PORT: 5173
      },
      env_production: {
        NODE_ENV: 'production'
      },
      error_file: path.join(projectRoot, 'logs', 'frontend-error.log'),
      out_file: path.join(projectRoot, 'logs', 'frontend-out.log'),
      log_file: path.join(projectRoot, 'logs', 'frontend-combined.log'),
      time: true,
      merge_logs: true
    }
  ]
};
