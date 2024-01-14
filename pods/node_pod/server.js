// server.js
const express = require('express');
const app = express();
const port = 3030;

app.get('/verify', (req, res) => {
  res.json({ isValid: "True" });
});

app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});