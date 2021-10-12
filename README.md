# Quick Start

You can use the demo site hosted at https://appmetadataapi.azurewebsites.net to play around.
A quick example to get metadata for app "mock app" with version 1.0, simply visit:
[https://appmetadataapi.azurewebsites.net/v1/metadata/mock%20app/1.0](https://appmetadataapi.azurewebsites.net/v1/metadata/mock%20app/1.0)

Or 
```zsh
curl https://appmetadataapi.azurewebsites.net/v1/metadata/mock%20app/1.0
```

---
# To create a new App Metadata:
Endpoint: `v1/metadata`

Method: `POST`

Content-type: `text/plain`

Body: the body should be a valid yaml format with schema as following example:
```yaml
title: Valid App 1
version: 0.0.1
maintainers:
  - name: firstmaintainer app1
    email: firstmaintainer@hotmail.com
  - name: secondmaintainer app1
    email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
  ### Interesting Title
  Some application content, and description
```

Payloads with same title and version are considered same metadata, and will be rejected if it is already created.

Example to post a new metadata using the demo site:
```zsh
curl -H "Content-Type: text/plain" -X POST https://appmetadataapi.azurewebsites.net/v1/metadata -d \
'
title: mock app
version: 1.1
maintainers:
  - name: firstmaintainer app1
    email: firstmaintainer@hotmail.com
  - name: secondmaintainer app1
    email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
  ### Interesting Title
  Some application content, and description
'
```
---
# To query for App Metadata:
If you know the app title and version, you can get the metadata directly by visiting endpoint:
```azure
/v1/metadata/:title/:version
```
where `:title` and `:version` is the url encoded value of app's title and version.
For example, to get metadata for app "mock app" with version "1.0" from the demo site:
```zsh
curl https://appmetadataapi.azurewebsites.net/v1/metadata/mock%20app/1.0
```
All queries are case-sensitive.

If you do not know app title or version, or simply want to search and explore. You can use following query parameters to do a search.
Note that all searches are for exact match and case-sensitive.

Supported query parameters:
1. title
2. version
3. maintainerName
4. maintainerEmail
5. company
6. website
7. source
8. license

For example, to find all apps' metadata that is maintained by Kai and is of MIT license:
```zsh
curl https://appmetadataapi.azurewebsites.net/v1/metadata?maintainerName=kai&license=MIT
```