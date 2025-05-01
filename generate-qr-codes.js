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
const sharp = require('sharp');

// Configuration
const QR_SIZE = 500; // Size of QR code in pixels
const QR_CODES_DIR = path.join(__dirname, 'static', 'qr-codes');
const LOGO_PATH = path.join(__dirname, 'static', 'images', 'qr-logo.png');
const LOGO_SIZE_PERCENTAGE = 0.17; // Logo size as a percentage of QR code size

const LINKS = [
  { name: 'website', url: 'https://www.cherevan.art' },
  { name: 'instagram', url: 'https://www.instagram.com/tetianacherevan/' },
  { name: 'shibari-artkingdom', url: 'https://www.instagram.com/shibari.artkingdom' },
  { name: 'art-finder', url: 'https://www.artfinder.com/artist/tetiana-cherevan/' },
  { name: 'etsy', url: 'https://www.etsy.com/shop/CherevanArtGallery' },
  { name: 'artsy', url: 'https://www.artsy.net/artist/tetiana-cherevan' }
];

// Ensure the QR codes directory exists
if (!fs.existsSync(QR_CODES_DIR)) {
  fs.mkdirSync(QR_CODES_DIR, { recursive: true });
  console.log(`Created directory: ${QR_CODES_DIR}`);
}

/**
 * Generate an artistic QR code with rounded dots
 * @param {string} url - The URL to encode in the QR code
 * @param {string} outputPath - Path to save the QR code image
 */
async function generateArtisticQRCode(url, outputPath) {
  try {
    // Create temporary file paths
    const tempSvgPath = `${outputPath}.svg`;
    const tempQrPngPath = `${outputPath}.temp.png`;
    
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
    
    // Calculate the center region to leave empty for the logo
    const logoSizeInModules = Math.ceil(size * LOGO_SIZE_PERCENTAGE * 2); // Double the percentage for better visibility
    const logoStartModule = Math.floor((size - logoSizeInModules) / 2);
    const logoEndModule = logoStartModule + logoSizeInModules;
    
    // Create SVG with circles instead of rectangles for truly rounded dots
    let svgContent = `<svg xmlns="http://www.w3.org/2000/svg" width="${QR_SIZE}" height="${QR_SIZE}" viewBox="0 0 ${QR_SIZE} ${QR_SIZE}">
`;
    svgContent += `  <rect width="${QR_SIZE}" height="${QR_SIZE}" fill="#FFFFFF"/>
`;
    
    // Add circles for each dark module in the QR code
    for (let row = 0; row < size; row++) {
      for (let col = 0; col < size; col++) {
        // Skip the center area where the logo will be placed
        if (row >= logoStartModule && row < logoEndModule && 
            col >= logoStartModule && col < logoEndModule) {
          continue;
        }
        
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
    
    // Convert SVG to PNG
    await sharp(tempSvgPath)
      .resize(QR_SIZE, QR_SIZE)
      .png()
      .toFile(tempQrPngPath);
    
    // Calculate logo size and position
    const logoSize = Math.round(QR_SIZE * LOGO_SIZE_PERCENTAGE);
    const logoPosition = Math.round((QR_SIZE - logoSize) / 2);
    
    // Create a temporary resized logo
    const tempLogoPath = `${outputPath}.logo.png`;
    
    // Resize the logo first
    await sharp(LOGO_PATH)
      .resize(logoSize, logoSize)
      .toFile(tempLogoPath);
    
    // Composite the resized logo onto the QR code
    await sharp(tempQrPngPath)
      .composite([
        {
          input: tempLogoPath,
          top: logoPosition,
          left: logoPosition
        }
      ])
      .toFile(outputPath);
      
    // Remove the temporary logo file
    fs.unlinkSync(tempLogoPath);
    
    // Remove temporary files
    fs.unlinkSync(tempSvgPath);
    fs.unlinkSync(tempQrPngPath);

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
    const outputPath = path.join(QR_CODES_DIR, `${link.name}.png`);
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
  const jsonOutputPath = path.join(QR_CODES_DIR, 'qr-codes.json');
  fs.writeFileSync(jsonOutputPath, JSON.stringify(results, null, 2));
  console.log(`QR code information saved to ${jsonOutputPath}`);
  
  console.log('QR code generation completed!');
}

// Run the main function
generateAllQRCodes().catch(console.error);
