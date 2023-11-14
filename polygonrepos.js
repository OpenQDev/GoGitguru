const fs = require('fs');
const path = require('path');

// Read the file
const filePath = path.join(__dirname, 'polygonrepos.txt');
const fileContent = fs.readFileSync(filePath, 'utf-8');

// Split the content by newlines, convert to lowercase, and trim
const lines = fileContent.split('\n').map(line => line.toLowerCase().trim());

// Output as a JSON array
const jsonArray = JSON.stringify(lines, null, 2);

// Write to a file
const outputFilePath = path.join(__dirname, 'output.json');
fs.writeFileSync(outputFilePath, jsonArray);