const yargs = require("yargs");
const { Sticker } = require("./images/images");
const argv = yargs
    .option("addexif", {
        alias: "exif",
        description: "Add exif on webp sticker",
        type: "string",
        array: true
    })
    .help()
    .alias("help", "h")
    .argv;

if (argv.addexif) {
    const sticker = new Sticker(argv.addexif[0])
    sticker.build({
        packname: argv.addexif[1],
        author: argv.addexif[2]
    }).then(result => {
        console.log(JSON.stringify({
            success: true,
            imagePath: result
        }))
    }).catch(JSON.stringify({
      success: false
  }))
}