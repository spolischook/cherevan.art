const sharp = require('sharp');
const fs = require('fs');
const path = require('path');

processDirectory(path.join(__dirname, 'content/calendar'));
processDirectory(path.join(__dirname, 'content/exhibitions'));
processDirectory(path.join(__dirname, 'content/press'));
processDirectory(path.join(__dirname, 'content/side-projects'));
processDirectory(path.join(__dirname, 'content/art-works'));
processDirectory(path.join(__dirname, 'static'));
console.log('Done');

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
                if (isImage(file)) {
                    try {
                        const image = sharp(absolutePath);
                        const formats = ['webp', 'avif'];
                        const sizes = [360, 375, 768];

                        formats.forEach(format => {
                            sizes.forEach(size => {
                                const resizedImagePath = path.join(directory, path.basename(file, path.extname(file)) + `_${size}.${format}`);
                                if (!fs.existsSync(resizedImagePath)) {
                                    image
                                        .resize(size)
                                        .toFormat(format)
                                        .toFile(resizedImagePath);
                                    console.log('Processed: ' + resizedImagePath);
                                }
                            });
                        });

                    } catch (err) {
                        console.error(`Error processing file ${absolutePath}: ${err.message}`);
                    }
                }
            }
        });
    });
}
function isImage(file) {
    const ext = path.extname(file).toLowerCase();
    return ext === ".jpg" || ext === ".png" || ext === ".jpeg";
}
