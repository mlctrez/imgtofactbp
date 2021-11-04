const cacheName = "app-" + "7adadf3915c608520570f4f55f2fa374524c46bf";

self.addEventListener("install", event => {
  console.log("installing app worker 7adadf3915c608520570f4f55f2fa374524c46bf");

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
          "https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css",
          "https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.bundle.min.js",
          "https://fonts.googleapis.com/css2?family=Roboto&display=swap",
          
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
  console.log("app worker 7adadf3915c608520570f4f55f2fa374524c46bf is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
