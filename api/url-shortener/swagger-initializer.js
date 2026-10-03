window.onload = function () {
  window.ui = SwaggerUIBundle({
    url: "/api/url-shortener/url-shortener.openapi.yaml",
    // urls: [
    //   {
    //     url: "./api/url-shortener/url-shortener.openapi.yaml",
    //     name: "url-shortener"
    //   },
    //   {
    //     url: "./api/url-shortener/schemas/shorten-post-request.yaml",
    //     name: "shorten-post-request"
    //   },
    //   {
    //     url: "./api/url-shortener/schemas/shorten-post-201-response.yaml",
    //     name: "shorten-post-201-response"
    //   }
    // ],
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};
