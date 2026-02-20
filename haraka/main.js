const SMTPServer = require('smtp-server').SMTPServer;
const { Client } = require('pg');
const pgClient = new Client({
  host: process.env.DB_HOST || 'postgres.besend.svc.cluster.local',
  port: process.env.DB_PORT || 5432,
  database: process.env.DB_NAME || 'besend',
  user: process.env.DB_USER || 'besenduser',
  password: process.env.DB_PASSWORD || 'besend_secure_password_123',
});
pgClient.connect().catch(err => {
  console.error('Database connection error:', err);
  process.exit(1);
});
const server = new SMTPServer({
  secure: false,
  requireTLS: false,
  allowInsecureAuth: true,
  onAuth: async (auth, session, callback) => {
    if (!auth.username || !auth.password) {
      return callback(new Error('Username and password required'));
    }
    try {
      const result = await pgClient.query(
        'SELECT id FROM smtp_credentials WHERE username = $1 AND password = crypt($2, password)',
        [auth.username, auth.password]
      );
      if (result.rows.length === 0) {
        return callback(new Error('Invalid credentials'));
      }
      callback(null, { user: auth.username });
    } catch (err) {
      console.error('Auth error:', err);
      callback(err);
    }
  },
  onData: async (stream, session, callback) => {
    let email = '';
    stream.on('data', chunk => {
      email += chunk.toString();
    });
    stream.on('end', async () => {
      try {
        const messageId = require('crypto').randomUUID();
        console.log(`Email from ${session.user || 'unknown'}: ${messageId}`);
        await pgClient.query(
          `INSERT INTO email_audit_logs 
           (id, customer_id, postal_message_id, status, created_at) 
           VALUES ($1, $2, $3, 'delivered', NOW())`,
          [messageId, session.user || 'unknown', messageId]
        );
        callback();
      } catch (err) {
        console.error('Data error:', err);
        callback(err);
      }
    });
  }
});
const port = process.env.SMTP_PORT || 587;
server.listen(port, '0.0.0.0', () => {
  console.log(`SMTP Server listening on port ${port}`);
});
server.on('error', err => {
  console.error('SMTP Server error:', err);
});
