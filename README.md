# X-Request-ID plugin for Traefik

This plugin will add the X-Request-ID header with a generated UUIDv4 value to HTTP requests and (optionally) to
responses, allowing downstream services to identify requests.

Based upon:
- github.com/mdklapwijk/traefik-plugin-request-id
- github.com/pipe01/plugin-requestid
- github.com/gamblingpro/plugin-requestid

## Disclaimer

This is not an official product of Digital H GmbH.  No support is provided.

## License

```
Copyright 2026 Digital H GmbH
Copyright 2023 M.D. Klapwijk
Copyright 2020 Vladimir Buyanov
Copyright 2020 Felipe Martínez

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
```

See [LICENSE](./LICENSE) for the full license text.
