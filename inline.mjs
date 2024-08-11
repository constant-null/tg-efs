import { inlineSource } from 'inline-source';
import fs from 'node:fs';
import path from 'node:path';

const htmlpath = path.resolve('sheet/sheet.html');


try {
    const html = await inlineSource(htmlpath, {
        compress: true,
        rootpath: path.resolve('sheet'),
        // Skip all css types and png formats
        ignore: ['png'],
    });
    fs.writeFile("sheet/sheet.min.html", html, err => {
        if (err) {
            console.error(err);
        } else {
            console.log("sheet/sheet.min.html exported");
        }
    })
} catch (err) {
    // Handle error
}