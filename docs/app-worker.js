const cacheName = "app-" + "b449a7f679f3b378ba7e108cd812e83fe239f6a8";
const resourcesToCache = ["/imgtofactbp","/imgtofactbp/app.css","/imgtofactbp/app.js","/imgtofactbp/manifest.webmanifest","/imgtofactbp/wasm_exec.js","/imgtofactbp/web/app.wasm","/imgtofactbp/web/logo-192.png","/imgtofactbp/web/logo-512.png","/imgtofactbp/web/style.css","https://cdnjs.cloudflare.com/ajax/libs/material-components-web/13.0.0/material-components-web.min.css","https://cdnjs.cloudflare.com/ajax/libs/material-components-web/13.0.0/material-components-web.min.js","https://fonts.googleapis.com/css2?family=Roboto\u0026display=swap","https://fonts.googleapis.com/icon?family=Material+Icons"];

self.addEventListener("install", (event) => {
  console.log("installing app worker b449a7f679f3b378ba7e108cd812e83fe239f6a8");

  event.waitUntil(
    caches
      .open(cacheName)
      .then((cache) => {
        return cache.addAll(resourcesToCache);
      })
      .then(() => {
        self.skipWaiting();
      })
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches.keys().then((keyList) => {
      return Promise.all(
        keyList.map((key) => {
          if (key !== cacheName) {
            return caches.delete(key);
          }
        })
      );
    })
  );
  console.log("app worker b449a7f679f3b378ba7e108cd812e83fe239f6a8 is activated");
});

self.addEventListener("fetch", (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      return response || fetch(event.request);
    })
  );
});

self.addEventListener("push", (event) => {
  if (!event.data || !event.data.text()) {
    return;
  }

  const notification = JSON.parse(event.data.text());
  if (!notification) {
    return;
  }

  const title = notification.title;
  delete notification.title;

  if (!notification.data) {
    notification.data = {};
  }
  let actions = [];
  for (let i in notification.actions) {
    const action = notification.actions[i];

    actions.push({
      action: action.action,
      path: action.path,
    });

    delete action.path;
  }
  notification.data.goapp = {
    path: notification.path,
    actions: actions,
  };
  delete notification.path;

  event.waitUntil(self.registration.showNotification(title, notification));
});

self.addEventListener("notificationclick", (event) => {
  event.notification.close();

  const notification = event.notification;
  let path = notification.data.goapp.path;

  for (let i in notification.data.goapp.actions) {
    const action = notification.data.goapp.actions[i];
    if (action.action === event.action) {
      path = action.path;
      break;
    }
  }

  event.waitUntil(
    clients
      .matchAll({
        type: "window",
      })
      .then((clientList) => {
        for (var i = 0; i < clientList.length; i++) {
          let client = clientList[i];
          if ("focus" in client) {
            client.focus();
            client.postMessage({
              goapp: {
                type: "notification",
                path: path,
              },
            });
            return;
          }
        }

        if (clients.openWindow) {
          return clients.openWindow(path);
        }
      })
  );
});
