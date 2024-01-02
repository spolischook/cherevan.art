const sharp = require('sharp');
const fs = require('fs');
const path = require('path');

const directoryPath = path.join(__dirname, 'content');

function processDirectory(directory) {
    fs.readdir(directory, function (err, files) {
        if (err) {
            return console.log('Unable to scan directory: ' + err);
        }
        files.forEach(function (file) {
            const absolutePath = path.join(directory, file);
            if (fs.statSync(absolutePath).isDirectory()) {
                processDirectory(absolutePath);
            } else {
                const ext = path.extname(file).toLowerCase();
                if (ext === ".jpg" || ext === ".png" || ext === ".jpeg") {
                    try {
                        const image = sharp(absolutePath);
                        const webpPath = path.join(directory, path.basename(file, path.extname(file)) + '.webp');
                        const avifPath = path.join(directory, path.basename(file, path.extname(file)) + '.avif');

                        if (!fs.existsSync(webpPath)) {
                            image
                                .toFormat('webp')
                                .toFile(webpPath);
                            console.log('Processed: ' + webpPath);
                        }

                        if (!fs.existsSync(avifPath)) {
                            image
                                .toFormat('avif')
                                .toFile(avifPath);
                            console.log('Processed: ' + avifPath);
                        }
                    } catch (err) {
                        console.error(`Error processing file ${absolutePath}: ${err}`);
                    }
                }
            }
        });
    });
}

processDirectory(directoryPath);
console.log('Done');