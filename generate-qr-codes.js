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
const LOGO_SIZE_PERCENTAGE = 0.13; // Logo size as a percentage of QR code size
const USE_LOGO = true; // Set to false to disable the logo in the center of QR codes
const USE_ROUNDED_DOTS = false; // Set to true for rounded dots, false for traditional square dots
const QR_CODES_JSON = path.join(QR_CODES_DIR, 'qr-codes.json');
const BORDER_WIDTH = 5; // Width of the border in pixels
const BORDER_PADDING = 10; // Padding between QR code and border in pixels
const BORDER_RADIUS = 40; // Radius of the rounded corners in pixels
const BORDER_COLOR = '#000000'; // Border color (black)
const QR_MARGIN = 2; // Margin size for QR code (smaller value = less white space)
const QR_ERROR_CORRECTION = 'H'; // Error correction level (L, M, Q, H)
const TITLE_FONT_SIZE = 35; // Font size for title text
const TITLE_MARGIN_TOP = 40; // Margin between QR code border and title
const TITLE_COLOR = '#000000'; // Title text color

// Load QR code data from JSON file if it exists
let qrCodesData = [];
if (fs.existsSync(QR_CODES_JSON)) {
  try {
    qrCodesData = JSON.parse(fs.readFileSync(QR_CODES_JSON, 'utf8'));
  } catch (error) {
    console.error('Error reading QR codes JSON file:', error);
  }
}

// Define links if no JSON data is available
const DEFAULT_LINKS = [
  { name: 'tetiana-cherevan-website', url: 'https://www.cherevan.art' },
  { name: 'tetiana-cherevan-instagram', url: 'https://www.instagram.com/tetianacherevan/' },
  { name: 'tetiana-cherevan-shibari-artkingdom-instagram', url: 'https://www.instagram.com/shibari.artkingdom' },
  { name: 'tetiana-cherevan-art-finder', url: 'https://www.artfinder.com/artist/tetiana-cherevan/' },
  { name: 'tetiana-cherevan-etsy', url: 'https://www.etsy.com/shop/CherevanArtGallery' },
  { name: 'tetiana-cherevan-artsy', url: 'https://www.artsy.net/artist/tetiana-cherevan' }
];

// Use data from JSON file if available, otherwise use default links
const LINKS = qrCodesData.length > 0 ? qrCodesData : DEFAULT_LINKS;

// Ensure the QR codes directory exists
if (!fs.existsSync(QR_CODES_DIR)) {
  fs.mkdirSync(QR_CODES_DIR, { recursive: true });
  console.log(`Created directory: ${QR_CODES_DIR}`);
}

/**
 * Generate an artistic QR code with rounded dots
 * @param {string} url - The URL to encode in the QR code
 * @param {string} outputPath - Path to save the QR code image
 * @param {string} logoSvg - Optional SVG logo to place in the center of the QR code
 * @param {string} title - Optional title to display below the QR code
 */
async function generateArtisticQRCode(url, outputPath, logoSvg, title) {
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
        // Create QR code with configured error correction
        const qrcode = qr.create(url, { errorCorrectionLevel: ErrorCorrectLevel[QR_ERROR_CORRECTION] });
        resolve(qrcode.modules);
      } catch (error) {
        reject(error);
      }
    });
    
    // Get QR code size (number of modules)
    const size = qrData.size;
    const cellSize = Math.floor(QR_SIZE / (size + QR_MARGIN)); // Use margin from configuration
    const margin = Math.floor((QR_SIZE - (size * cellSize)) / 2);
    
    // Calculate the center region to leave empty for the logo (if enabled)
    let logoSizeInModules = 0;
    let logoStartModule = 0;
    let logoEndModule = 0;
    
    if (USE_LOGO) {
      logoSizeInModules = Math.ceil(size * LOGO_SIZE_PERCENTAGE * 1.5); // Adjusted for better visibility with less whitespace
      logoStartModule = Math.floor((size - logoSizeInModules) / 2);
      logoEndModule = logoStartModule + logoSizeInModules;
    }
    
    // Create SVG with either circles (rounded dots) or rectangles (square dots) based on configuration
    let svgContent = `<svg xmlns="http://www.w3.org/2000/svg" width="${QR_SIZE}" height="${QR_SIZE}" viewBox="0 0 ${QR_SIZE} ${QR_SIZE}">
`;
    svgContent += `  <rect width="${QR_SIZE}" height="${QR_SIZE}" fill="#FFFFFF"/>
`;
    
    // Draw each dot as either a circle (rounded) or rectangle (square) based on configuration
    for (let y = 0; y < size; y++) {
      for (let x = 0; x < size; x++) {
        // Skip the center region if a logo will be placed there
        if (USE_LOGO && 
            x >= logoStartModule && x <= logoEndModule && 
            y >= logoStartModule && y <= logoEndModule) {
          continue;
        }
        
        // Check if this module is dark (1) or light (0)
        if (qrData.get(x, y) === 1) {
          if (USE_ROUNDED_DOTS) {
            // Draw as circle for rounded dots
            const cx = margin + (x * cellSize) + (cellSize / 2);
            const cy = margin + (y * cellSize) + (cellSize / 2);
            const radius = cellSize / 2 * 0.85; // Slightly smaller than half cell size for spacing
            
            svgContent += `  <circle cx="${cx}" cy="${cy}" r="${radius}" fill="#000000"/>
`;
          } else {
            // Draw as rectangle for traditional square dots
            const rectX = margin + (x * cellSize);
            const rectY = margin + (y * cellSize);
            const rectSize = cellSize * 0.9; // Slightly smaller than cell size for spacing
            const rectOffset = (cellSize - rectSize) / 2;
            
            svgContent += `  <rect x="${rectX + rectOffset}" y="${rectY + rectOffset}" width="${rectSize}" height="${rectSize}" fill="#000000"/>
`;
          }
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
    
    if (USE_LOGO && logoSvg) {
      // Calculate logo size and position
      const logoSize = Math.round(QR_SIZE * LOGO_SIZE_PERCENTAGE * 1.2);
      const logoPosition = Math.round((QR_SIZE - logoSize) / 2);
      
      // Create temporary files for the logo
      const tempLogoSvgPath = `${outputPath}.logo.svg`;
      const tempLogoPngPath = `${outputPath}.logo.png`;
      
      // Write the SVG logo to a temporary file
      fs.writeFileSync(tempLogoSvgPath, logoSvg);
      
      // Convert SVG logo to PNG and resize
      await sharp(tempLogoSvgPath)
        .resize(logoSize, logoSize)
        .toFile(tempLogoPngPath);
      
      // Create a QR code with logo
      const qrWithLogo = await sharp(tempQrPngPath)
        .composite([
          {
            input: tempLogoPngPath,
            top: logoPosition,
            left: logoPosition
          }
        ])
        .toBuffer();
        
      // Create a new image with border, padding, and title (if provided)
      const totalPadding = BORDER_PADDING * 2; // Padding on all sides
      const titleHeight = title ? TITLE_FONT_SIZE + TITLE_MARGIN_TOP : 0;
      const finalWidth = QR_SIZE + totalPadding + (BORDER_WIDTH * 2);
      const finalHeight = finalWidth + titleHeight;
      const innerSize = QR_SIZE + totalPadding;
      
      let svgContent = `<svg width="${finalWidth}" height="${finalHeight}" viewBox="0 0 ${finalWidth} ${finalHeight}" xmlns="http://www.w3.org/2000/svg">
        <rect x="0" y="0" width="${finalWidth}" height="${finalWidth}" rx="${BORDER_RADIUS}" ry="${BORDER_RADIUS}" fill="${BORDER_COLOR}"/>
        <rect x="${BORDER_WIDTH}" y="${BORDER_WIDTH}" width="${innerSize}" height="${innerSize}" rx="${BORDER_RADIUS - BORDER_WIDTH}" ry="${BORDER_RADIUS - BORDER_WIDTH}" fill="white"/>`;
      
      // Add title if provided
      if (title) {
        svgContent += `
        <text x="${finalWidth / 2}" y="${finalWidth + TITLE_MARGIN_TOP}" font-family="Arial, sans-serif" font-size="${TITLE_FONT_SIZE}" text-anchor="middle" fill="${TITLE_COLOR}">${title}</text>`;
      }
      
      svgContent += `
      </svg>`;
      
      const finalImage = Buffer.from(svgContent);
      
      // Composite the QR code with logo onto the bordered background
      await sharp(finalImage)
        .composite([
          {
            input: qrWithLogo,
            top: BORDER_WIDTH + BORDER_PADDING,
            left: BORDER_WIDTH + BORDER_PADDING
          }
        ])
        .toFile(outputPath);
        
      // Remove the temporary logo files
      fs.unlinkSync(tempLogoSvgPath);
      fs.unlinkSync(tempLogoPngPath);
    } else {
      // If no logo is used, add a border to the QR code
      // Create a new image with border, padding, and title (if provided)
      const totalPadding = BORDER_PADDING * 2; // Padding on all sides
      const titleHeight = title ? TITLE_FONT_SIZE + TITLE_MARGIN_TOP : 0;
      const finalWidth = QR_SIZE + totalPadding + (BORDER_WIDTH * 2);
      const finalHeight = finalWidth + titleHeight;
      const innerSize = QR_SIZE + totalPadding;
      
      let svgContent = `<svg width="${finalWidth}" height="${finalHeight}" viewBox="0 0 ${finalWidth} ${finalHeight}" xmlns="http://www.w3.org/2000/svg">
        <rect x="0" y="0" width="${finalWidth}" height="${finalWidth}" rx="${BORDER_RADIUS}" ry="${BORDER_RADIUS}" fill="${BORDER_COLOR}"/>
        <rect x="${BORDER_WIDTH}" y="${BORDER_WIDTH}" width="${innerSize}" height="${innerSize}" rx="${BORDER_RADIUS - BORDER_WIDTH}" ry="${BORDER_RADIUS - BORDER_WIDTH}" fill="white"/>`;
      
      // Add title if provided
      if (title) {
        svgContent += `
        <text x="${finalWidth / 2}" y="${finalWidth + TITLE_MARGIN_TOP}" font-family="Arial, sans-serif" font-size="${TITLE_FONT_SIZE}" text-anchor="middle" fill="${TITLE_COLOR}">${title}</text>`;
      }
      
      svgContent += `
      </svg>`;
      
      const finalImage = Buffer.from(svgContent);
      
      // Composite the QR code onto the bordered background
      await sharp(finalImage)
        .composite([
          {
            input: tempQrPngPath,
            top: BORDER_WIDTH + BORDER_PADDING,
            left: BORDER_WIDTH + BORDER_PADDING
          }
        ])
        .toFile(outputPath);
    }
    
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
  
  for (const link of LINKS) {
    const { name, url, logo, title } = link;
    const outputPath = path.join(QR_CODES_DIR, `${name}.png`);
    console.log(`Generating QR code for ${url} as ${outputPath}`);
    
    try {
      await generateArtisticQRCode(url, outputPath, logo, title);
    } catch (error) {
      console.error(`Failed to generate QR code for ${name}:`, error);
    }
  }
  
  console.log('QR code generation completed!');
}

// Run the main function
generateAllQRCodes().catch(console.error);
