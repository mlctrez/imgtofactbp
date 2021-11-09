const cacheName = "app-" + "7889294808d9f4049514cebd0ce178c11623715e";

self.addEventListener("install", event => {
  console.log("installing app worker 7889294808d9f4049514cebd0ce178c11623715e");

  event.waitUntil(
    caches.open(cacheName).
      then(cache => {
        return cache.addAll([
          "/imgtofactbp",
          "/imgtofactbp/app.css",
          "/imgtofactbp/app.js",
          "/imgtofactbp/manifest.webmanifest",
          "/imgtofactbp/wasm_exec.js",
          "/imgtofactbp/web/app.wasm",
          "/imgtofactbp/web/logo-192.png",
          "/imgtofactbp/web/logo-512.png",
          "/imgtofactbp/web/style.css",
          "https://cdnjs.cloudflare.com/ajax/libs/material-components-web/13.0.0/material-components-web.min.css",
          "https://cdnjs.cloudflare.com/ajax/libs/material-components-web/13.0.0/material-components-web.min.js",
          "https://fonts.googleapis.com/css2?family=Roboto&display=swap",
          "https://fonts.googleapis.com/icon?family=Material+Icons",
          
        ]);
      }).
      then(() => {
        self.skipWaiting();
      })
  );
});

self.addEventListener("activate", event => {
  event.waitUntil(
    caches.keys().then(keyList => {
      return Promise.all(
        keyList.map(key => {
          if (key !== cacheName) {
            return caches.delete(key);
          }
        })
      );
    })
  );
  console.log("app worker 7889294808d9f4049514cebd0ce178c11623715e is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
