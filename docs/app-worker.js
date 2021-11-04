const cacheName = "app-" + "737e1f7aa2433e2786af920a45357a2eed16dace";

self.addEventListener("install", event => {
  console.log("installing app worker 737e1f7aa2433e2786af920a45357a2eed16dace");

  event.waitUntil(
    caches.open(cacheName).
      then(cache => {
        return cache.addAll([
          "/",
          "/app.css",
          "/app.js",
          "/manifest.webmanifest",
          "/wasm_exec.js",
          "/web/app.wasm",
          "/web/logo-192.png",
          "/web/logo-512.png",
          "/web/style.css",
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
  console.log("app worker 737e1f7aa2433e2786af920a45357a2eed16dace is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
