/*
Hooks for clipboard.go
 */
;(function () {

    document.addEventListener("paste", async function (e) {

        for (let item of e.clipboardData.items) {
            console.log("item:", item.kind, item.type)

            const promise = new Promise((resolve => {
                item.getAsString((data) => {
                    resolve(data)
                })
            }))


            item.getAsString(console.log)
        }
        for (let file of e.clipboardData.files) {
            console.log("file:", file.name, file.lastModified, file.type, file.size)
        }
    })

})()