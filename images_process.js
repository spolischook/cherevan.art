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
                if(ext === ".jpg" || ext === ".png" || ext === ".jpeg"){
                    const image = sharp(absolutePath);
                    image
                    .toFormat('webp')
                    .toFile(path.join(directory, path.basename(file, path.extname(file)) + '.webp'));

                    image
                    .toFormat('avif')
                    .toFile(path.join(directory, path.basename(file, path.extname(file)) + '.avif'));
                }
            }
        });
    });
}

processDirectory(directoryPath);