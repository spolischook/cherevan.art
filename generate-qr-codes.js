/**
 * QR Code Generator for Cherevan.art Link Tree
 * 
 * This script generates artistic QR codes for all links in the link tree
 * and saves them to the Hugo static assets folder.
 * 
 * Uses qrcode library with radius option for rounded dots.
 */

const fs = require('fs');
const path = require('path');
const QRCode = require('qrcode');

// Configuration
const QR_SIZE = 1024; // Size of QR code in pixels
const OUTPUT_DIR = path.join(__dirname, 'static', 'qr-codes');
const LINKS = [
  { name: 'website', url: 'https://www.cherevan.art' },
  { name: 'instagram', url: 'https://www.instagram.com/tetianacherevan/' },
  { name: 'shibari-artkingdom', url: 'https://www.instagram.com/shibari.artkingdom' },
  { name: 'art-finder', url: 'https://www.artfinder.com/artist/tetiana-cherevan/' },
  { name: 'etsy', url: 'https://www.etsy.com/shop/CherevanArtGallery' },
  { name: 'artsy', url: 'https://www.artsy.net/artist/tetiana-cherevan' }
];

// Create output directory if it doesn't exist
if (!fs.existsSync(OUTPUT_DIR)) {
  fs.mkdirSync(OUTPUT_DIR, { recursive: true });
  console.log(`Created directory: ${OUTPUT_DIR}`);
}

/**
 * Generate an artistic QR code with rounded dots
 * @param {string} url - The URL to encode in the QR code
 * @param {string} outputPath - Path to save the QR code image
 */
async function generateArtisticQRCode(url, outputPath) {
  try {
    // Create a temporary SVG file path
    const tempSvgPath = `${outputPath}.svg`;
    
    // First, generate the QR code matrix data
    const qrData = await new Promise((resolve, reject) => {
      // Use the QRCode.create method to get the raw QR code data
      const qr = require('qrcode/lib/core/qrcode');
      const ErrorCorrectLevel = require('qrcode/lib/core/error-correction-level');
      
      try {
        // Create QR code with high error correction
        const qrcode = qr.create(url, { errorCorrectionLevel: ErrorCorrectLevel.H });
        resolve(qrcode.modules);
      } catch (error) {
        reject(error);
      }
    });
    
    // Get QR code size (number of modules)
    const size = qrData.size;
    const cellSize = Math.floor(QR_SIZE / (size + 8)); // Add margin
    const margin = Math.floor((QR_SIZE - (size * cellSize)) / 2);
    
    // Create SVG with circles instead of rectangles for truly rounded dots
    let svgContent = `<svg xmlns="http://www.w3.org/2000/svg" width="${QR_SIZE}" height="${QR_SIZE}" viewBox="0 0 ${QR_SIZE} ${QR_SIZE}">
`;
    svgContent += `  <rect width="${QR_SIZE}" height="${QR_SIZE}" fill="#FFFFFF"/>
`;
    
    // Add circles for each dark module in the QR code
    for (let row = 0; row < size; row++) {
      for (let col = 0; col < size; col++) {
        // Check if this module is dark (true)
        if (qrData.data[row * size + col]) {
          const x = margin + col * cellSize + cellSize / 2;
          const y = margin + row * cellSize + cellSize / 2;
          const radius = cellSize / 2 * 0.9; // Slightly smaller than half cell size for spacing
          
          svgContent += `  <circle cx="${x}" cy="${y}" r="${radius}" fill="#000000"/>
`;
        }
      }
    }
    
    svgContent += '</svg>';
    
    // Write the SVG to a temporary file
    fs.writeFileSync(tempSvgPath, svgContent);
    
    // Convert SVG to PNG using sharp
    await require('sharp')(tempSvgPath)
      .resize(QR_SIZE, QR_SIZE)
      .png()
      .toFile(outputPath);
    
    // Remove the temporary SVG file
    fs.unlinkSync(tempSvgPath);

    console.log(`Generated QR code for ${url} at ${outputPath}`);
    return outputPath;
  } catch (error) {
    console.error(`Error generating QR code for ${url}:`, error);
    throw error;
  }
}

/**
 * Main function to generate all QR codes
 */
async function generateAllQRCodes() {
  console.log('Starting QR code generation...');
  
  const results = [];
  
  for (const link of LINKS) {
    const outputPath = path.join(OUTPUT_DIR, `${link.name}.png`);
    try {
      await generateArtisticQRCode(link.url, outputPath);
      results.push({
        name: link.name,
        url: link.url,
        qrPath: `/qr-codes/${link.name}.png`
      });
    } catch (error) {
      console.error(`Failed to generate QR code for ${link.name}:`, error);
    }
  }
  
  // Save a JSON file with information about all generated QR codes
  const jsonOutputPath = path.join(OUTPUT_DIR, 'qr-codes.json');
  fs.writeFileSync(jsonOutputPath, JSON.stringify(results, null, 2));
  console.log(`QR code information saved to ${jsonOutputPath}`);
  
  console.log('QR code generation completed!');
}

// Run the main function
generateAllQRCodes().catch(console.error);
