module.exports = {
  // example configuration for PM2 or other supporting systems
  apps: [
    {
      name: 'discord-smtp-relay',
      // in production, point at the built binary:
      script: './discord-smtp-relay',
      // in development, you might instead use:
      // script: 'go',
      // args: 'run main.go',
      cwd: __dirname,
      // set to true if you want auto-restart on source changes
      watch: false,
      env: {
        // Discord webhook URL
        WEBHOOK: 'https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN',
        // Path to your JSON embed template
        TEMPLATE: 'etc/template.example.json',
        // Optional SMTP AUTH credentials (default “discord”/“discord”)
        USERNAME: 'discord',
        PASSWORD: 'discord',
        // SMTP server listener settings
        HOST: '0.0.0.0',
        PORT: 1025,
        DOMAIN: 'localhost',
        // Timeouts (seconds) and max message size (KB)
        READ: 10,
        WRITE: 10,
        SIZE: 1024
      }
    }
  ]
}
