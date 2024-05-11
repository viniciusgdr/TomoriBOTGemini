const webp = require('node-webpmux');
const fs = require('fs');

class Sticker {
    media
    fileNameTemp
    defaultDir
    constructor(media) {
        this.media = media
        this.fileNameTemp = Date.now()
        this.defaultDir = './assets/temp/'
    }
    build({
        packname,
        author
    }) {
        return new Promise(async (resolve, reject) => {
            let fileWithExif = this.media + '_exif.webp'
            const img = new webp.Image()
            const json = {
                "sticker-pack-id": 'Feito usando o bot de stickers do WhatsApp',
                "sticker-pack-name": packname,
                "sticker-pack-publisher": author,
                "emojis": [""]
            }

            const exifAttr = Buffer.from([0x49, 0x49, 0x2A, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x41, 0x57, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00])
            const jsonBuff = Buffer.from(JSON.stringify(json), "utf-8")
            const exif = Buffer.concat([exifAttr, jsonBuff])
            exif.writeUIntLE(jsonBuff.length, 14, 4)

            await img.load(this.media)
            img.exif = exif
            await img.save(fileWithExif)
            
            fs.unlinkSync(this.media)
            resolve(fileWithExif)
        })
    }
}

module.exports = {
    Sticker
}