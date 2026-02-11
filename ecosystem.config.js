module.exports = {
  apps: [
    {
      name: 'stardew-backend',
      cwd: '/root/workpath/claude/stardew-agent/backend',
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
      error_file: '/root/workpath/claude/stardew-agent/logs/backend-error.log',
      out_file: '/root/workpath/claude/stardew-agent/logs/backend-out.log',
      log_file: '/root/workpath/claude/stardew-agent/logs/backend-combined.log',
      time: true,
      merge_logs: true,
      // Pre-start hook to build the Go binary
      pre_start: 'go build -o server ./cmd/server'
    },
    {
      name: 'stardew-frontend',
      cwd: '/root/workpath/claude/stardew-agent/frontend',
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
      error_file: '/root/workpath/claude/stardew-agent/logs/frontend-error.log',
      out_file: '/root/workpath/claude/stardew-agent/logs/frontend-out.log',
      log_file: '/root/workpath/claude/stardew-agent/logs/frontend-combined.log',
      time: true,
      merge_logs: true
    }
  ]
};
