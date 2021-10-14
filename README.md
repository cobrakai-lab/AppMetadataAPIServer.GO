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
```
Endpoint: v1/metadata

Method: POST

Content-type: text/plain
```
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
Note that all searches are for exact match and case-sensitive. Results are sorted ascending lexicographically by title then by version.

Supported query parameters:
```
1. title
2. version
3. maintainerName
4. maintainerEmail
5. company
6. website
7. source
8. license
```

Two more special query parameters for paginations:
1. __pageSize__: The max items per page. Default is 20. You can override this to set how many items per page.
2. __page__: The page number. Default is 1. Items displayed per page depends on pageSize. 

If either pageSize or page is invalid, default values will be used.

For example, to find all apps' metadata that is maintained by Kai and is of MIT license, with pageSize set to 1, and page number to 2:
```zsh
curl https://appmetadataapi.azurewebsites.net/v1/metadata?maintainerName=kai&license=MIT&pageSize=1&page=2
```
---
# To get all metadata
A simple API for you to get all metadata stored
```
Endpoint: v1/metadata

Method: GET
```
Pagination is supported on this as well. Results are sorted ascending lexicographically by title then by version. 

Example:
```zsh
curl https://appmetadataapi.azurewebsites.net/v1/metadata
```



---
# To delete an App Metadata:
```
Endpoint: v1/metadata/:title/:version

Method: DELETE
```
where `:title` and `:version` is the url encoded value of app's title and version. If the metadata exists, this API will return the deleted metadata content in response. If the metadata does not exist, the response will simply state so.

Exmaple to delete a metadata with title=app1 and verison=1.0:
```zsh
curl -X "DELETE" https://appmetadataapi.azurewebsites.net/v1/metadata/app1/1.0
```
---
# To get some fun statistics for the service:
You can get some fun statistics about the service such as when the service started last time, how many times each API has been called since the service restarted, etc by calling:
```
Endpoint: v1/metadata/_stats

Method: GET
```
Exmaple:
```zsh
curl https://appmetadataapi.azurewebsites.net/v1/metadata/_stats
```
Sample output:
```json
{
    "delete_metadata_api_called_counter": 0,
    "get_metadata_api_called_counter": 2,
    "post_metadata_api_called_counter": 0,
    "query_metadata_api_called_counter": 1,
    "server_start_time": "2021-10-13 05:01:50.670403762 +0000 UTC m=+0.113754108",
    "stats_viewed": 17
}
```