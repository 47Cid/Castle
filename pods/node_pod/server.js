// server.js
const express = require('express');
const morgan = require('morgan');
const winston = require('winston');
const xss = require('xss');
const app = express();
const port = 3030;

// Create a new winston logger that writes to a file
const logger = winston.createLogger({
  level: 'info',
  format: winston.format.json(),
  transports: [
    new winston.transports.File({ filename: 'logs/http.log' })
  ]
});

// Use morgan for logging, with a custom function to write logs using winston
app.use(morgan((tokens, req, res) => {
  return [
    tokens.method(req, res),
    tokens.url(req, res),
    tokens.status(req, res),
    tokens.res(req, res, 'content-length'), '-',
    tokens['response-time'](req, res), 'ms'
  ].join(' ')
}, { stream: { write: (message) => logger.info(message.trim()) } }));

// Basic XSS filtering
app.get('/verify', (req, res) => {
  const sanitizedInput = xss(req.query.input);
  if (req.query.input !== sanitizedInput) {
    return res.json({ isValid: "False", time: new Date().toISOString() });
  }
  res.json({ isValid: "True", time: new Date().toISOString() });
});

app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});